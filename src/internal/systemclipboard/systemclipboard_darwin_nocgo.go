//go:build darwin && !cgo

package systemclipboard

import "fmt"

// Available reports false because the macOS pasteboard backend requires cgo.
func Available() bool { return false }

// CopyFiles returns ErrUnsupported with a hint to build with cgo enabled.
func CopyFiles(_ []string, _ bool) error {
	return fmt.Errorf("%w: macOS clipboard file transfer requires cgo (NSPasteboard)", ErrUnsupported)
}

// PasteFiles returns no paths and ErrUnsupported because cgo is disabled.
func PasteFiles() ([]string, bool, error) {
	return nil, false, fmt.Errorf("%w: macOS clipboard file transfer requires cgo (NSPasteboard)", ErrUnsupported)
}
