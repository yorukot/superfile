//go:build windows

package metadata

import (
	"os"
	"syscall"

	"golang.org/x/sys/windows"
)

func getOwnerAndGroup(_ os.FileInfo) (string, string) {
	return "", ""
}

// returns file attributes
func getFileAttributes(path string) (string, bool) {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return "", false
	}
	attrs, err := windows.GetFileAttributes(ptr)
	if err != nil {
		return "", false
	}
	result := ""
	if attrs&syscall.FILE_ATTRIBUTE_READONLY != 0 {
		result += "R"
	}
	if attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0 {
		result += "H"
	}
	if attrs&syscall.FILE_ATTRIBUTE_SYSTEM != 0 {
		result += "S"
	}
	if attrs&syscall.FILE_ATTRIBUTE_ARCHIVE != 0 {
		result += "A"
	}
	if attrs&syscall.FILE_ATTRIBUTE_DIRECTORY != 0 {
		result += "D"
	}
	return result, true
}
