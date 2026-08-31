//go:build windows

package internal

import (
	"strings"

	"golang.org/x/sys/windows"
)

func extendedLengthPath(path string) (string, error) {
	path = strings.ReplaceAll(path, "/", `\`)

	if path == "" || strings.HasPrefix(path, `\\?\`) || strings.HasPrefix(path, `\??\`) ||
		strings.HasPrefix(path, `\\.\`) {
		return path, nil
	}

	fullPath, err := windows.FullPath(path)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(fullPath, `\\`) {
		return `\\?\UNC\` + strings.TrimPrefix(fullPath, `\\`), nil
	}

	return `\\?\` + fullPath, nil
}

func renameNoReplace(src, dst string) error {
	src, err := extendedLengthPath(src)
	if err != nil {
		return err
	}
	dst, err = extendedLengthPath(dst)
	if err != nil {
		return err
	}

	srcPtr, err := windows.UTF16PtrFromString(src)
	if err != nil {
		return err
	}
	dstPtr, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(srcPtr, dstPtr, 0)
}
