package utils

import "log/slog"
import "os"

// Todo : Eventually we want to remove all such usage that can result in app exiting abruptly
func LogAndExit(msg string, values ...any) {
	slog.Error(msg, values...)
	os.Exit(1)
}

// Used in unit test
func SetRootLoggerToStdout(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stdout, &slog.HandlerOptions{Level: level})))
}
