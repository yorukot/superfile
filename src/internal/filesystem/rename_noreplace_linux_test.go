//go:build linux

package filesystem

import (
	"errors"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
)

func TestRenameNoReplacePreservesInvalidDescendantError(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source")
	if err := unix.Mkdir(source, 0o700); err != nil {
		t.Fatalf("create source directory: %v", err)
	}
	destination := filepath.Join(source, "nested")

	err := renameNoReplace(source, destination)
	if !errors.Is(err, unix.EINVAL) {
		t.Fatalf("renameNoReplace() error = %v, want EINVAL", err)
	}
}
