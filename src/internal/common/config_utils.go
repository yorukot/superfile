package common

import (
	"os"
	"runtime"

	"github.com/yorukot/superfile/src/internal/utils"
)

// ResolveEditor returns the command used to open files in an editor.
// Priority: Config.Editor → $EDITOR → OS fallback (non-empty result guaranteed).
func ResolveEditor() string {
	if Config.Editor != "" {
		return Config.Editor
	}
	if envEditor := os.Getenv("EDITOR"); envEditor != "" {
		return envEditor
	}
	if runtime.GOOS == utils.OsWindows {
		return "notepad"
	}
	return "nano"
}
