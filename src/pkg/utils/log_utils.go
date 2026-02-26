package utils

import (
	"fmt"
	"log/slog"
	"os"
)

// Print line to stderr and exit with status 1
// Cannot use log.Fataln() as slog.SetDefault() causes those lines to
// go into log file
func PrintlnAndExit(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

// Print formatted output line to stderr and exit with status 1
// Cannot use log.Fataln() as slog.SetDefault() causes those lines to
// go into log file
func PrintfAndExitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
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

// Used in unit test
func SetRootLoggerToDiscarded() {
	slog.SetDefault(slog.New(slog.DiscardHandler))
}
