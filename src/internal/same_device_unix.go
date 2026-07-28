//go:build !windows

package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// sameDeviceID checks whether two paths reside on the same filesystem device
// by comparing their device IDs from stat(2). This is the standard Unix
// approach for detecting cross-device boundaries.
//
// If path2 does not exist yet (e.g. a move destination that hasn't been
// created), the function stats its parent directory instead, since a new
// file inherits the device of its parent.
func sameDeviceID(path1, path2 string) (bool, error) {
	var stat1, stat2 syscall.Stat_t

	if err := syscall.Stat(path1, &stat1); err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path1, err)
	}
	if err := syscall.Stat(path2, &stat2); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("failed to stat %s: %w", path2, err)
		}
		// Destination doesn't exist yet — stat parent directory instead.
		if err := syscall.Stat(filepath.Dir(path2), &stat2); err != nil {
			return false, fmt.Errorf("failed to stat destination parent: %w", err)
		}
	}
	return stat1.Dev == stat2.Dev, nil
}
