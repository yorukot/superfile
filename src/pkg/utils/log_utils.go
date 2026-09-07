package utils

import (
	"fmt"
	"io"
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
	setRootLoggerTo(os.Stdout, debug)
}

// Used before the log file is open, so that diagnostics never mix into the
// stdout of commands whose output is consumed by scripts
func SetRootLoggerToStderr(debug bool) {
	setRootLoggerTo(os.Stderr, debug)
}

func setRootLoggerTo(w io.Writer, debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(
		w, &slog.HandlerOptions{Level: level})))
}

// Used in unit test
func SetRootLoggerToDiscarded() {
	slog.SetDefault(slog.New(slog.DiscardHandler))
}
