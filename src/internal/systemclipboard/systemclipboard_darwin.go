//go:build darwin && cgo

package systemclipboard

/*
#cgo LDFLAGS: -framework Foundation -framework AppKit
#include <stdlib.h>

typedef struct {
	char *data;
	char *errorMessage;
} SPFClipResult;

char *spf_clipboard_copy_files(const char **paths, int count);
SPFClipResult spf_clipboard_paste_files(void);
void spf_clip_free_string(char *value);
*/
import "C"

import (
	"errors"
	"path/filepath"
	"strings"
	"unsafe"
)

// Available reports that this build includes the native macOS pasteboard backend.
func Available() bool { return true }

// CopyFiles writes the given paths to the macOS general pasteboard as file URLs.
//
// The cut flag is intentionally ignored: macOS Finder has no concept of cutting
// files (paste is copy, "Move Item Here" is a separate action), and the
// pasteboard cannot represent a move. A cut therefore degrades to a copy.
func CopyFiles(paths []string, _ bool) error {
	if len(paths) == 0 {
		return ErrNoFiles
	}
	abs := make([]string, 0, len(paths))
	for _, p := range paths {
		a, err := filepath.Abs(p)
		if err != nil {
			return err
		}
		abs = append(abs, a)
	}

	cStrings := make([]*C.char, len(abs))
	for i, p := range abs {
		cStrings[i] = C.CString(p)
	}
	defer func() {
		for _, cs := range cStrings {
			C.free(unsafe.Pointer(cs))
		}
	}()

	cErr := C.spf_clipboard_copy_files(
		(**C.char)(unsafe.Pointer(&cStrings[0])),
		C.int(len(cStrings)),
	)
	if cErr != nil {
		defer C.spf_clip_free_string(cErr)
		return errors.New(C.GoString(cErr))
	}
	return nil
}

// PasteFiles reads file paths from the macOS general pasteboard. The returned
// cut flag is always false (see CopyFiles).
func PasteFiles() ([]string, bool, error) {
	result := C.spf_clipboard_paste_files()
	defer C.spf_clip_free_string(result.data)
	defer C.spf_clip_free_string(result.errorMessage)

	if result.errorMessage != nil {
		return nil, false, errors.New(C.GoString(result.errorMessage))
	}

	var joined string
	if result.data != nil {
		joined = C.GoString(result.data)
	}
	joined = strings.TrimSpace(joined)
	if joined == "" {
		return nil, false, ErrNoFiles
	}

	var paths []string
	for _, line := range strings.Split(joined, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			paths = append(paths, line)
		}
	}
	if len(paths) == 0 {
		return nil, false, ErrNoFiles
	}
	return paths, false, nil
}
