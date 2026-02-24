//go:build !windows

package utils

import (
	"os/exec"
	"syscall"
)

func DetachFromTerminal(cmd *exec.Cmd, keepStdoutAndStderr bool) {
	// Start new session so child isn't tied to the TTY (prevents SIGHUP on terminal close).
	// This also prevents programs like sudo to read/write to tty
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}

	// Optionally, redirect stdio to avoid terminal hangups
	// Stdin set to nil is also needed to prevent interactive commands in shell mode
	cmd.Stdin = nil
	if !keepStdoutAndStderr {
		cmd.Stdout = nil
		cmd.Stderr = nil
	}
}
