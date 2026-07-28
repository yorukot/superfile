//go:build windows

package internal

// sameDeviceID is a no-op on Windows. Drive letter comparison is used instead
// in isSamePartition(). This function should never be called on Windows.
func sameDeviceID(_, _ string) (bool, error) {
	// On Windows, isSamePartition uses getDriveLetter() and returns before
	// reaching this function. This stub exists only to satisfy the compiler.
	return false, nil
}
