//go:build !windows
// +build !windows

package filepreview

import (
	"syscall"
	"unsafe"
)

// getTerminalCellSizeViaIoctl uses ioctl system call to get terminal size
func getTerminalCellSizeViaIoctl() (TerminalCellSize, bool) {
	// Try different file descriptors in order of preference
	fds := []uintptr{
		1, // stdout
		0, // stdin
		2, // stderr
	}

	for _, fd := range fds {
		if cellSize, ok := getTerminalSizeFromFd(fd); ok {
			return cellSize, true
		}
	}

	return TerminalCellSize{}, false
}

// getTerminalSizeFromFd gets terminal size from a specific file descriptor
func getTerminalSizeFromFd(fd uintptr) (TerminalCellSize, bool) {
	var ws winsize

	// TIOCGWINSZ ioctl call to get window size
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)),
	)

	if errno != 0 {
		return TerminalCellSize{}, false
	}

	// Check if we got valid pixel dimensions
	if ws.Xpixel > 0 && ws.Ypixel > 0 && ws.Col > 0 && ws.Row > 0 {
		pixelsPerColumn := int(ws.Xpixel) / int(ws.Col)
		pixelsPerRow := int(ws.Ypixel) / int(ws.Row)

		// Sanity check the values
		if pixelsPerColumn > 0 && pixelsPerRow > 0 &&
			pixelsPerColumn < 100 && pixelsPerRow < 100 {
			return TerminalCellSize{
				PixelsPerColumn: pixelsPerColumn,
				PixelsPerRow:    pixelsPerRow,
			}, true
		}
	}

	return TerminalCellSize{}, false
}
