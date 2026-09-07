//go:build !unix

package internal

import "os"

// newThemeReloadSignals returns nil on platforms without SIGUSR1 (e.g. Windows).
func newThemeReloadSignals() chan os.Signal {
	return nil
}
