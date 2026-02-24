//go:build !windows

package utils

import (
	"context"
	"os/exec"
	"syscall"
)

func formCommand(cmdDir string, ctx context.Context, baseCmd string, args []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, baseCmd, args...)
	cmd.Dir = cmdDir

	// Set input to /dev/null to return EOF on 'read'
	cmd.Stdin = nil
	// Detach the process from the tty. So that programs like sudo cannot open and read/write to tty
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	return cmd
}
