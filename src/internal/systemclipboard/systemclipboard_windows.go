//go:build windows

package systemclipboard

import (
	"fmt"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"
)

const (
	cfHDROP = 15

	gmemMoveable = 0x0002

	dropEffectCopy = 0x00000001
	dropEffectMove = 0x00000002
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	shell32  = syscall.NewLazyDLL("shell32.dll")

	procOpenClipboard            = user32.NewProc("OpenClipboard")
	procCloseClipboard           = user32.NewProc("CloseClipboard")
	procEmptyClipboard           = user32.NewProc("EmptyClipboard")
	procGetClipboardData         = user32.NewProc("GetClipboardData")
	procSetClipboardData         = user32.NewProc("SetClipboardData")
	procRegisterClipboardFormatW = user32.NewProc("RegisterClipboardFormatW")

	procGlobalAlloc      = kernel32.NewProc("GlobalAlloc")
	procGlobalFree       = kernel32.NewProc("GlobalFree")
	procGlobalLock       = kernel32.NewProc("GlobalLock")
	procGlobalUnlock     = kernel32.NewProc("GlobalUnlock")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")

	procDragQueryFileW = shell32.NewProc("DragQueryFileW")
)

// dropfiles mirrors the Win32 DROPFILES structure that prefixes a CF_HDROP
// payload. Total size is 20 bytes (DWORD + POINT + 2*BOOL).
type dropfiles struct {
	pFiles uint32 // offset (in bytes) to the file list
	ptX    int32
	ptY    int32
	fNC    int32
	fWide  int32 // 1 -> the file list is UTF-16
}

// Available reports that this build includes the native Windows clipboard backend.
func Available() bool { return true }

// CopyFiles places the given paths on the clipboard as a CF_HDROP payload and
// sets the "Preferred DropEffect" so that a paste into Explorer performs a copy
// or a move accordingly.
func CopyFiles(paths []string, cut bool) error {
	abs := make([]string, 0, len(paths))
	for _, p := range paths {
		a, err := filepath.Abs(p)
		if err != nil {
			return err
		}
		abs = append(abs, a)
	}
	if len(abs) == 0 {
		return ErrNoFiles
	}

	// The clipboard is thread-affine; keep everything on one OS thread.
	errCh := make(chan error, 1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		errCh <- copyFiles(abs, cut)
	}()
	return <-errCh
}

// copyFiles publishes absolute paths and their preferred drop effect through
// Win32. The caller must remain on one OS thread for the entire operation.
func copyFiles(abs []string, cut bool) error {
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd == 0 {
		return fmt.Errorf("GetConsoleWindow failed: no console window")
	}
	if r, _, err := procOpenClipboard.Call(hwnd); r == 0 {
		return fmt.Errorf("OpenClipboard failed: %w", err)
	}
	defer procCloseClipboard.Call()

	if r, _, err := procEmptyClipboard.Call(); r == 0 {
		return fmt.Errorf("EmptyClipboard failed: %w", err)
	}

	hDrop, err := globalFromBytes(buildHDropBytes(abs))
	if err != nil {
		return err
	}
	if r, _, e := procSetClipboardData.Call(cfHDROP, hDrop); r == 0 {
		procGlobalFree.Call(hDrop)
		return fmt.Errorf("SetClipboardData(CF_HDROP) failed: %w", e)
	}

	hEffect, err := globalFromBytes(buildDropEffectBytes(cut))
	if err != nil {
		return err
	}
	cfEffect, err := registerFormat("Preferred DropEffect")
	if err != nil {
		procGlobalFree.Call(hEffect)
		return err
	}
	if r, _, e := procSetClipboardData.Call(uintptr(cfEffect), hEffect); r == 0 {
		procGlobalFree.Call(hEffect)
		procEmptyClipboard.Call()
		return fmt.Errorf("SetClipboardData(Preferred DropEffect) failed: %w", e)
	}
	return nil
}

// buildHDropBytes builds the raw CF_HDROP payload: a DROPFILES header followed
// by a double-null-terminated list of UTF-16 paths. Everything is assembled in
// ordinary Go memory so no unsafe pointer arithmetic is required.
func buildHDropBytes(paths []string) []byte {
	var chars []uint16
	for _, p := range paths {
		// UTF16FromString errors only on embedded NULs; skip the current path
		// if it cannot be encoded.
		u16, err := syscall.UTF16FromString(p)
		if err != nil {
			continue
		}
		chars = append(chars, u16...) // includes the trailing NUL for this path
	}
	chars = append(chars, 0) // extra NUL to terminate the list

	headerSize := int(unsafe.Sizeof(dropfiles{}))
	buf := make([]byte, headerSize+len(chars)*2)

	header := dropfiles{pFiles: uint32(headerSize), fWide: 1}
	*(*dropfiles)(unsafe.Pointer(&buf[0])) = header

	list := unsafe.Slice((*uint16)(unsafe.Pointer(&buf[headerSize])), len(chars))
	copy(list, chars)
	return buf
}

// buildDropEffectBytes encodes the copy or move effect as a Win32 DWORD payload.
func buildDropEffectBytes(cut bool) []byte {
	effect := uint32(dropEffectCopy)
	if cut {
		effect = dropEffectMove
	}
	buf := make([]byte, 4)
	*(*uint32)(unsafe.Pointer(&buf[0])) = effect
	return buf
}

// globalFromBytes copies buf into a freshly allocated moveable global memory
// block suitable for SetClipboardData, returning its handle.
func globalFromBytes(buf []byte) (uintptr, error) {
	hMem, _, err := procGlobalAlloc.Call(gmemMoveable, uintptr(len(buf)))
	if hMem == 0 {
		return 0, fmt.Errorf("GlobalAlloc failed: %w", err)
	}
	ptr, _, _ := procGlobalLock.Call(hMem)
	if ptr == 0 {
		procGlobalFree.Call(hMem)
		return 0, fmt.Errorf("GlobalLock failed")
	}
	// GlobalLock returns a stable pointer into non-movable locked memory.
	dst := unsafe.Slice((*byte)(unsafe.Pointer(ptr)), len(buf)) //nolint:govet // uintptr from GlobalLock points to locked, non-movable memory
	copy(dst, buf)
	procGlobalUnlock.Call(hMem)
	return hMem, nil
}

// PasteFiles reads a CF_HDROP payload from the clipboard and reports whether the
// source marked it as a move (cut).
func PasteFiles() ([]string, bool, error) {
	type result struct {
		paths []string
		cut   bool
		err   error
	}
	ch := make(chan result, 1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		paths, cut, err := pasteFiles()
		ch <- result{paths, cut, err}
	}()
	r := <-ch
	return r.paths, r.cut, r.err
}

// pasteFiles reads dropped paths and their preferred move flag through Win32,
// returning ErrNoFiles when none are available. The caller must pin its OS thread.
func pasteFiles() ([]string, bool, error) {
	if r, _, err := procOpenClipboard.Call(0); r == 0 {
		return nil, false, fmt.Errorf("OpenClipboard failed: %w", err)
	}
	defer procCloseClipboard.Call()

	hDrop, _, _ := procGetClipboardData.Call(cfHDROP)
	if hDrop == 0 {
		return nil, false, ErrNoFiles
	}

	count, _, _ := procDragQueryFileW.Call(hDrop, 0xFFFFFFFF, 0, 0)
	if count == 0 {
		return nil, false, ErrNoFiles
	}

	var paths []string
	for i := uintptr(0); i < count; i++ {
		n, _, _ := procDragQueryFileW.Call(hDrop, i, 0, 0)
		if n == 0 {
			continue
		}
		buf := make([]uint16, n+1) // include the terminating NUL
		n, _, _ = procDragQueryFileW.Call(
			hDrop, i,
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(len(buf)),
		)
		if n == 0 {
			continue
		}
		paths = append(paths, syscall.UTF16ToString(buf[:n]))
	}
	if len(paths) == 0 {
		return nil, false, ErrNoFiles
	}

	return paths, readDropEffectCut(), nil
}

// readDropEffectCut reports whether the open clipboard requests a move, defaulting
// to copy semantics if the preferred drop effect cannot be read.
func readDropEffectCut() bool {
	cfEffect, err := registerFormat("Preferred DropEffect")
	if err != nil {
		return false
	}
	hEffect, _, _ := procGetClipboardData.Call(uintptr(cfEffect))
	if hEffect == 0 {
		return false
	}
	ptr, _, _ := procGlobalLock.Call(hEffect)
	if ptr == 0 {
		return false
	}
	defer procGlobalUnlock.Call(hEffect)
	// ptr comes from GlobalLock and points to locked, non-movable memory.
	effect := *(*uint32)(unsafe.Pointer(ptr)) //nolint:govet // stable locked-memory pointer
	return effect&dropEffectMove != 0
}

// registerFormat obtains the Win32 clipboard format identifier for name,
// registering the format if it does not already exist.
func registerFormat(name string) (uint32, error) {
	u16, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}
	r, _, e := procRegisterClipboardFormatW.Call(uintptr(unsafe.Pointer(u16)))
	if r == 0 {
		return 0, fmt.Errorf("RegisterClipboardFormat(%q) failed: %w", name, e)
	}
	return uint32(r), nil
}
