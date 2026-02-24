//go:build windows

package utils

import (
	"context"
	"os/exec"
)

func formCommand(cmdDir string, ctx context.Context, baseCmd string, args []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, baseCmd, args...)
	cmd.Dir = cmdDir
	return cmd
}
