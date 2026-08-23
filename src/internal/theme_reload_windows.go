//go:build windows

package internal

import "os"

// newThemeReloadSignals returns nil on Windows, which has no SIGUSR1.
func newThemeReloadSignals() chan os.Signal {
	return nil
}
