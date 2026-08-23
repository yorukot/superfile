//go:build !windows

package internal

import (
	"os"
	"os/signal"
	"syscall"
)

// newThemeReloadSignals returns a channel that receives SIGUSR1.
func newThemeReloadSignals() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGUSR1)
	return signals
}
