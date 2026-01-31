//go:build linux
// +build linux

package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	"github.com/yorukot/superfile/src/internal/utils/pty"

	"golang.org/x/sys/unix"
)

const DefaultCols = 80
const DefaultRows = 16

// Choose correct shell as per OS
func ExecuteCommandInShell(winSize *Winsize,
	timeLimit time.Duration,
	cmdDir string,
	shellCommand string) (int, string, error) {
	baseCmd := "/bin/sh"
	args := []string{"-c", shellCommand}
	var ws *unix.Winsize
	if winSize == nil {
		ws = &unix.Winsize{Row: DefaultRows, Col: DefaultCols}
	} else {
		ws = &unix.Winsize{Row: winSize.Row, Col: winSize.Col, Xpixel: 0, Ypixel: 0}
	}
	return ExecuteCommand(ws, timeLimit, cmdDir, baseCmd, args...)
}

func ExecuteCommand(winSize *unix.Winsize,
	timeLimit time.Duration,
	cmdDir string,
	baseCmd string,
	args ...string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	cmd := exec.CommandContext(ctx, baseCmd, args...)
	cmd.Dir = cmdDir
	retCode := -1

	ptmx, err := pty.Start(cmd, winSize)
	if err != nil {
		return retCode, "", fmt.Errorf("unexpected Error in pty start command : %w", err)
	}
	outputBytes, err := pty.Read(ptmx, cmd)
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		slog.Error("User's command timed out", "outputBytes", outputBytes,
			"cmd error", err, "ctx error", ctx.Err())
		return retCode, string(outputBytes), ctx.Err()
	}
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok { //nolint: errorlint // We dont expect error to be Wrapped here
			retCode = exitErr.ExitCode()
		} else {
			err = fmt.Errorf("unexpected Error in command execution : %w", err)
		}
		return retCode, string(outputBytes), err
	}
	retCode = 0
	return retCode, string(outputBytes), err
}
