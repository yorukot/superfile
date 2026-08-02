//go:build !linux && !darwin && !windows

package systemclipboard

func Available() bool { return false }

func CopyFiles(_ []string, _ bool) error { return ErrUnsupported }

func PasteFiles() ([]string, bool, error) { return nil, false, ErrUnsupported }
