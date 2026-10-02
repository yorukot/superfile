//go:build !linux && !darwin && !windows

package systemclipboard

// Available reports false on platforms without a file clipboard backend.
func Available() bool { return false }

// CopyFiles returns ErrUnsupported without changing the system clipboard.
func CopyFiles(_ []string, _ bool) error { return ErrUnsupported }

// PasteFiles returns no paths and ErrUnsupported on this platform.
func PasteFiles() ([]string, bool, error) { return nil, false, ErrUnsupported }
