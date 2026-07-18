//go:build linux

package filesystem

import "golang.org/x/sys/unix"

func renameNoReplace(source, destination string) error {
	err := unix.Renameat2(unix.AT_FDCWD, source, unix.AT_FDCWD, destination, unix.RENAME_NOREPLACE)
	if err == unix.ENOSYS || err == unix.EINVAL || err == unix.EOPNOTSUPP {
		return errNoReplaceUnsupported
	}
	return err
}
