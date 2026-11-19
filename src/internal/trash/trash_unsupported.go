//go:build !linux && !darwin && !windows

package trash

func Init() error {
	return ErrUnsupported
}

func Available(_ string) bool {
	return false
}

func Move(path string) (Result, error) {
	return Result{OriginalPath: path}, ErrUnsupported
}
