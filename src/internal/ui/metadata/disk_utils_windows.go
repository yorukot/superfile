//go:build windows

package metadata

import (
	"golang.org/x/sys/windows"
)

func getFreeSpace(path string) (*DriveSize, error) {
	var freeBytes uint64
	var totalBytes uint64
	err := windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr(path), &freeBytes, &totalBytes, nil)
	if err != nil {
		return nil, err
	}
	return &DriveSize{Total: totalBytes, Available: freeBytes}, nil
}
