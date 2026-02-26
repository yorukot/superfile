package utils

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"runtime"
	"time"
)

// Choose correct shell as per OS
func ExecuteCommandInShell(timeLimit time.Duration, cmdDir string, shellCommand string) (int, string, error) {
	// Linux and Darwin
	baseCmd := "/bin/sh"
	args := []string{"-c", shellCommand}

	if runtime.GOOS == OsWindows {
		baseCmd = "powershell.exe"
		args[0] = "-Command"
	}

	return ExecuteCommand(timeLimit, cmdDir, baseCmd, args...)
}

func ExecuteCommand(timeLimit time.Duration, cmdDir string, baseCmd string, args ...string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeLimit)
	defer cancel()

	cmd := exec.CommandContext(ctx, baseCmd, args...)
	cmd.Dir = cmdDir
	DetachFromTerminal(cmd)
	outputBytes, err := cmd.CombinedOutput()
	retCode := -1

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		slog.Error("User's command timed out", "outputBytes", outputBytes,
			"cmd error", err, "ctx error", ctx.Err())
		return retCode, string(outputBytes), ctx.Err()
	}

	if err == nil {
		retCode = 0
	} else if exitErr, ok := err.(*exec.ExitError); ok { //nolint: errorlint // We dont expect error to be Wrapped here
		retCode = exitErr.ExitCode()
	} else {
		err = fmt.Errorf("unexpected Error in command execution : %w", err)
	}

	return retCode, string(outputBytes), err
}
