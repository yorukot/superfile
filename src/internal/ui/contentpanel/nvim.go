package contentpanel

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/creack/pty"
)

// NvimState holds the state for an embedded nvim session.
type NvimState struct {
	ptyFile *os.File       // PTY master file descriptor
	cmd     *exec.Cmd      // nvim process
	running bool           // whether nvim is currently running
	mu      sync.Mutex     // protects the output buffer
	output  strings.Builder // accumulated raw output (ANSI-stripped for now)
	done    chan struct{}   // closed when nvim exits
}

// StartNvim launches nvim in a PTY sized to the given dimensions.
// Returns the NvimState or an error.
func StartNvim(filePath string, width, height int) (*NvimState, error) {
	nvimCmd := "nvim"
	// Fall back to vim if nvim is not installed
	if _, err := exec.LookPath("nvim"); err != nil {
		nvimCmd = "vim"
	}

	cmd := exec.Command(nvimCmd, filePath)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	// Create PTY with the given dimensions
	ptyFile, err := pty.StartWithSize(cmd, &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start nvim in PTY: %w", err)
	}

	ns := &NvimState{
		ptyFile: ptyFile,
		cmd:     cmd,
		running: true,
		done:    make(chan struct{}),
	}

	// Start output reader goroutine
	go ns.readOutput()

	// Monitor nvim exit
	go func() {
		err := cmd.Wait()
		ns.mu.Lock()
		ns.running = false
		ns.mu.Unlock()
		close(ns.done)
		if err != nil {
			slog.Debug("nvim exited", "error", err)
		}
	}()

	slog.Debug("nvim started in PTY", "path", filePath, "width", width, "height", height)
	return ns, nil
}

// readOutput runs in a goroutine, reading PTY output and stripping ANSI codes.
func (ns *NvimState) readOutput() {
	buf := make([]byte, 4096)
	for {
		n, err := ns.ptyFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				slog.Debug("nvim PTY read error", "error", err)
			}
			return
		}
		if n > 0 {
			ns.mu.Lock()
			ns.output.Write(buf[:n])
			// Limit buffer size to prevent memory bloat
			if ns.output.Len() > 1_000_000 {
				// Truncate oldest half
				old := ns.output.String()
				ns.output.Reset()
				ns.output.WriteString(old[len(old)/2:])
			}
			ns.mu.Unlock()
		}
	}
}

// GetOutput returns the current PTY output with ANSI codes stripped.
func (ns *NvimState) GetOutput() string {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	raw := ns.output.String()
	return stripANSI(raw)
}

// Resize changes the PTY dimensions to match the panel.
func (ns *NvimState) Resize(width, height int) error {
	if ns.ptyFile == nil {
		return nil
	}
	return pty.Setsize(ns.ptyFile, &pty.Winsize{
		Rows: uint16(height),
		Cols: uint16(width),
	})
}

// WriteKey sends a key sequence to nvim's PTY.
func (ns *NvimState) WriteKey(data []byte) error {
	if ns.ptyFile == nil || !ns.running {
		return nil
	}
	_, err := ns.ptyFile.Write(data)
	return err
}

// IsRunning returns true while nvim is still active.
func (ns *NvimState) IsRunning() bool {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	return ns.running
}

// WaitDone blocks until nvim exits with a timeout.
func (ns *NvimState) WaitDone(timeout time.Duration) {
	select {
	case <-ns.done:
	case <-time.After(timeout):
	}
}

// Stop forcefully terminates nvim and closes the PTY.
func (ns *NvimState) Stop() {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	if !ns.running {
		return
	}
	ns.running = false
	if ns.ptyFile != nil {
		ns.ptyFile.Close()
	}
	if ns.cmd != nil && ns.cmd.Process != nil {
		ns.cmd.Process.Kill()
	}
}

// stripANSI removes ANSI escape sequences from the output for plain-text display.
func stripANSI(input string) string {
	var result strings.Builder
	result.Grow(len(input))
	inEscape := false
	i := 0
	runes := []rune(input)
	for i < len(runes) {
		r := runes[i]
		if !inEscape {
			if r == '\x1b' && i+1 < len(runes) && runes[i+1] == '[' {
				inEscape = true
				i += 2
				continue
			}
			result.WriteRune(r)
			i++
		} else {
			// Skip until we find a letter (end of CSI sequence)
			for i < len(runes) && inEscape {
				r = runes[i]
				if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
					inEscape = false
				}
				i++
			}
		}
	}
	// Also strip carriage returns
	out := strings.ReplaceAll(result.String(), "\r", "")
	return out
}
