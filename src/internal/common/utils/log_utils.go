package utils 

import "log/slog"
import "os"

// Todo : Eventually we want to remove all such usage that can result in app exiting abruptly
func LogAndExit(msg string, values ...any) {
	slog.Error(msg, values...)
	os.Exit(1)
}