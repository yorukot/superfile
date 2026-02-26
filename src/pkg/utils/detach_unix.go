//go:build !windows

package utils

import (
	"os/exec"
	"syscall"
)

func DetachFromTerminal(cmd *exec.Cmd) {
	// Start new session so child isn't tied to the TTY (prevents SIGHUP on terminal close).
	// This also prevents programs like sudo to read/write to tty
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	// Optionally, redirect stdio to avoid terminal hangups
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
}
