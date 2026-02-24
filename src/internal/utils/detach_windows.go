//go:build windows

package utils

import "os/exec"

func DetachFromTerminal(cmd *exec.Cmd, keepStdoutAndStderr bool) {
	// No-op: current Windows path uses rundll32 and returns immediately.
	// If needed later, set CreationFlags/HideWindow via syscall.SysProcAttr.
}
