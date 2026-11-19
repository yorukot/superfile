//go:build darwin && !cgo

package trash

import "fmt"

func Init() error {
	return fmt.Errorf("%w: macOS trash requires cgo for Foundation FileManager", ErrUnsupported)
}

func Available(_ string) bool {
	return false
}

func Move(path string) (Result, error) {
	return Result{OriginalPath: path, Backend: BackendMacOS}, fmt.Errorf("%w: macOS trash requires cgo for Foundation FileManager", ErrUnsupported)
}
