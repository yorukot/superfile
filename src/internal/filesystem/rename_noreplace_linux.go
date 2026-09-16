//go:build linux

package filesystem

import (
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

func renameNoReplace(source, destination string) error {
	err := unix.Renameat2(unix.AT_FDCWD, source, unix.AT_FDCWD, destination, unix.RENAME_NOREPLACE)
	if err == unix.EINVAL && destinationWithinSource(source, destination) {
		return err
	}
	if err == unix.ENOSYS || err == unix.EINVAL || err == unix.EOPNOTSUPP {
		return errNoReplaceUnsupported
	}
	return err
}

func destinationWithinSource(source, destination string) bool {
	sourcePath, err := filepath.Abs(source)
	if err != nil {
		return false
	}
	destinationPath, err := filepath.Abs(destination)
	if err != nil {
		return false
	}
	relativePath, err := filepath.Rel(sourcePath, destinationPath)
	return err == nil && relativePath != "." && relativePath != ".." &&
		!strings.HasPrefix(relativePath, ".."+string(filepath.Separator))
}
