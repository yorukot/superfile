//go:build windows

package utils

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

func replaceFile(replacementPath, replacedPath string) error {
	replacedPathPtr, err := windows.UTF16PtrFromString(replacedPath)
	if err != nil {
		return err
	}
	replacementPathPtr, err := windows.UTF16PtrFromString(replacementPath)
	if err != nil {
		return err
	}

	result, _, callErr := windows.NewLazySystemDLL("kernel32.dll").NewProc("ReplaceFileW").Call(
		uintptr(unsafe.Pointer(replacedPathPtr)),
		uintptr(unsafe.Pointer(replacementPathPtr)),
		0,
		0,
		0,
		0,
	)
	if result == 0 {
		if callErr != syscall.Errno(0) {
			return callErr
		}
		return syscall.EINVAL
	}
	return nil
}
