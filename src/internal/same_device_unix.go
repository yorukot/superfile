//go:build !windows

package internal

import (
	"fmt"
	"syscall"
)

// sameDeviceID checks whether two paths reside on the same filesystem device
// by comparing their device IDs from stat(2). This is the standard Unix
// approach for detecting cross-device boundaries.
func sameDeviceID(path1, path2 string) (bool, error) {
	var stat1, stat2 syscall.Stat_t

	if err := syscall.Stat(path1, &stat1); err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path1, err)
	}
	if err := syscall.Stat(path2, &stat2); err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path2, err)
	}
	return stat1.Dev == stat2.Dev, nil
}
