//go:build linux
// +build linux

package pty

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

type (
	_CInt  int32
	_CUint uint32
)

const PtsPreffix = "/dev/pts/"
const ReadBufferSize = 1024
const InteractiveModeThreshold = 500 // in miliseconds

func ioctl(f *os.File, cmd, ptr uintptr) error {
	_, _, e := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), cmd, ptr)
	if e != 0 {
		return e
	}
	return nil
}

func disableEchoOnSlave(slave *os.File) error {
	var termios unix.Termios
	termios.Lflag &^= unix.ECHO | unix.ECHONL

	if _, _, errno := unix.Syscall6(
		unix.SYS_IOCTL,
		slave.Fd(),
		uintptr(unix.TCSETS),
		uintptr(unsafe.Pointer(&termios)),
		0, 0, 0,
	); errno != 0 {
		return fmt.Errorf("can't disable echo on slave: %w", errno)
	}
	return nil
}

func ptsName(f *os.File) (string, error) {
	var n _CUint
	err := ioctl(f, syscall.TIOCGPTN, uintptr(unsafe.Pointer(&n)))
	if err != nil {
		return "", err
	}
	return PtsPreffix + strconv.Itoa(int(n)), nil
}

func unlockPt(f *os.File) error {
	var u _CInt
	return ioctl(f, syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&u)))
}

func open() (*os.File, *os.File, error) {
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	defer func() {
		if err != nil {
			_ = master.Close()
		}
	}()
	if err != nil {
		return nil, nil, err
	}

	slaveName, err := ptsName(master)
	if err != nil {
		return nil, nil, err
	}

	if err = unlockPt(master); err != nil {
		return nil, nil, err
	}

	slave, err := os.OpenFile(slaveName, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}

	if err = disableEchoOnSlave(slave); err != nil {
		_ = unix.Close(int(master.Fd()))
		_ = unix.Close(int(slave.Fd()))
		return nil, nil, err
	}

	return master, slave, nil
}

func setSize(ptyMaster *os.File, windowSize *unix.Winsize) error {
	if err := unix.IoctlSetWinsize(int(ptyMaster.Fd()), unix.TIOCSWINSZ, windowSize); err != nil {
		return fmt.Errorf("failed to set winsize: %w", err)
	}
	return nil
}

// start command in pty
func Start(cmd *exec.Cmd, winSize *unix.Winsize) (*os.File, error) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setsid = true
	cmd.SysProcAttr.Setctty = true

	ptm, pts, err := open()
	if err != nil {
		return nil, err
	}
	defer func() { _ = pts.Close() }()

	if err = setSize(ptm, winSize); err != nil {
		return nil, err
	}
	if cmd.Stdout == nil {
		cmd.Stdout = pts
	}
	if cmd.Stderr == nil {
		cmd.Stderr = pts
	}
	if cmd.Stdin == nil {
		cmd.Stdin = pts
	}

	if err = cmd.Start(); err != nil {
		_ = ptm.Close()
		return nil, err
	}
	return ptm, err
}

// Interactive mode detection heuristic:
// if there is more than InteractiveModeThreshold value seconds between two readings,
// then interactive mode is detected.
func readInner(ptmx *os.File, cmd *exec.Cmd, output *bytes.Buffer, done chan error) {
	buf := make([]byte, ReadBufferSize)
	timeoutDuration := InteractiveModeThreshold * time.Millisecond

	for {
		timer := time.NewTimer(timeoutDuration)
		readChan := make(chan struct {
			n   int
			err error
		}, 1)

		go func() {
			n, err := ptmx.Read(buf)
			readChan <- struct {
				n   int
				err error
			}{n, err}
		}()

		select {
		case res := <-readChan:
			if !timer.Stop() {
				<-timer.C
			}

			if res.n > 0 {
				output.Write(buf[:res.n])
			}

			if res.err != nil {
				done <- cmd.Wait()
				return
			}
		case <-timer.C:
			_ = cmd.Process.Kill()
			done <- errors.New("interactive mode is not allowed")
			return
		}
	}
}

// read output from pty
func Read(ptmx *os.File, cmd *exec.Cmd) ([]byte, error) {
	output := &bytes.Buffer{}
	done := make(chan error, 1)

	go readInner(ptmx, cmd, output, done)
	if err := <-done; err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
