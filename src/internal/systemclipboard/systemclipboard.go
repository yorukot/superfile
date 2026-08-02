// Package systemclipboard provides copy/paste of real files and directories
// through the operating system clipboard, so that items copied inside superfile
// can be pasted into the native file manager (Finder, Explorer, Nautilus, ...)
// and vice-versa.
//
// The heavy lifting is platform specific and lives in the build-tagged files in
// this package, mirroring the layout of the sibling `trash` package:
//
//   - systemclipboard_windows.go       CF_HDROP via Win32 syscalls (no cgo)
//   - systemclipboard_darwin.go / .m   NSPasteboard file URLs via cgo
//   - systemclipboard_darwin_nocgo.go  graceful fallback when built without cgo
//   - systemclipboard_linux.go         shells out to wl-clipboard / xclip
//   - systemclipboard_unsupported.go   every other platform
//
// Each platform file implements the same three functions:
//
//	func Available() bool
//	func CopyFiles(paths []string, cut bool) error
//	func PasteFiles() (paths []string, cut bool, err error)
//
// Notes on the `cut` flag: not every platform can represent a "cut" (move) on
// the system clipboard. macOS Finder in particular has no notion of cutting
// files, so on macOS a cut degrades to a copy. Callers should treat the `cut`
// value returned by PasteFiles as advisory.
package systemclipboard

import "errors"

// ErrUnsupported is returned when the current platform (or environment) cannot
// transfer files through the system clipboard.
var ErrUnsupported = errors.New("system clipboard file transfer is not supported on this platform")

// ErrNoFiles is returned by PasteFiles when the system clipboard does not hold
// any file references.
var ErrNoFiles = errors.New("no files found on the system clipboard")
