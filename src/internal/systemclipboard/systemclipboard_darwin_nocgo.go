//go:build darwin && !cgo

package systemclipboard

import "fmt"

func Available() bool { return false }

func CopyFiles(_ []string, _ bool) error {
	return fmt.Errorf("%w: macOS clipboard file transfer requires cgo (NSPasteboard)", ErrUnsupported)
}

func PasteFiles() ([]string, bool, error) {
	return nil, false, fmt.Errorf("%w: macOS clipboard file transfer requires cgo (NSPasteboard)", ErrUnsupported)
}
