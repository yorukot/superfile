package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDirSize covers the directory-size walker. It pins down two things:
//   - the sum across nested files (including an empty subdir and a regular
//     file at the root) is correct;
//   - an inaccessible root (non-existent path) returns 0 without panicking.
//     WalkDir delivers err != nil with a nil entry there; the walker must
//     return before dereferencing entry (regression for a nil-pointer panic).
func TestDirSize(t *testing.T) {
	root := t.TempDir()

	// Two regular files of known sizes in different subdirectories.
	subDir := filepath.Join(root, "sub")
	require.NoError(t, os.Mkdir(subDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "top.txt"),
		make([]byte, 512), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(subDir, "leaf.txt"),
		make([]byte, 128), 0o644))

	// An empty directory contributes nothing; it must not break the walk.
	require.NoError(t, os.Mkdir(filepath.Join(root, "empty"), 0o755))

	got := DirSize(root)
	assert.Equal(t, int64(640), got)

	// An empty directory on its own sums to zero.
	assert.Equal(t, int64(0), DirSize(filepath.Join(root, "empty")))
}

// TestDirSizeNonExistentRoot is the regression case for the nil-entry panic:
// WalkDir invokes the callback with err != nil and a nil os.DirEntry when the
// root cannot be walked. Previously DirSize dereferenced entry unconditionally
// and crashed here; it must now return 0 cleanly.
func TestDirSizeNonExistentRoot(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "does-not-exist")
	// The parent temp dir exists, the child does not, so the first WalkDir
	// callback is exactly the (err != nil, entry == nil) case.
	if runtime.GOOS == "windows" {
		missing = `C:\superfile-does-not-exist-dirsize`
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("DirSize panicked on non-existent root: %v", r)
		}
	}()
	assert.Equal(t, int64(0), DirSize(missing))
}
