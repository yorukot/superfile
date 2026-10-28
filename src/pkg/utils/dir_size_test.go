package utils

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirSizeWithContext(t *testing.T) {
	root := t.TempDir()
	nested := filepath.Join(root, "a", "b")
	SetupDirectories(t, nested, filepath.Join(root, "empty"))
	SetupFilesWithData(t, []byte("12345"), filepath.Join(root, "f1"))
	SetupFilesWithData(t, []byte("1234567890"), filepath.Join(root, "a", "f2"), filepath.Join(nested, "f3"))

	t.Run("Sums sizes of all nested files", func(t *testing.T) {
		size, partial, err := DirSizeWithContext(context.Background(), root)
		require.NoError(t, err)
		assert.Equal(t, int64(25), size)
		assert.False(t, partial)
	})

	t.Run("Empty directory", func(t *testing.T) {
		size, partial, err := DirSizeWithContext(context.Background(), filepath.Join(root, "empty"))
		require.NoError(t, err)
		assert.Equal(t, int64(0), size)
		assert.False(t, partial)
	})

	t.Run("Non existent directory", func(t *testing.T) {
		_, _, err := DirSizeWithContext(context.Background(), filepath.Join(root, "missing"))
		require.Error(t, err)
	})

	t.Run("Cancelled context stops the walk", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		size, _, err := DirSizeWithContext(ctx, root)
		require.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, int64(0), size)
	})
}

func TestDirSizeWithContextPermissionDenied(t *testing.T) {
	// Skip on Windows as permission handling differs
	if runtime.GOOS == OsWindows {
		t.Skip("Skipping permission test on Windows")
	}

	// Skip when running as root since root can read directories with 000 permissions
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	root := t.TempDir()
	locked := filepath.Join(root, "locked")
	SetupDirectories(t, locked)
	SetupFilesWithData(t, []byte("12345"), filepath.Join(root, "f1"), filepath.Join(locked, "f2"))

	require.NoError(t, os.Chmod(locked, 0))
	defer os.Chmod(locked, UserDirPerm) // Reset permissions for cleanup

	t.Run("Unreadable subdirectory gives a partial result", func(t *testing.T) {
		size, partial, err := DirSizeWithContext(context.Background(), root)
		require.NoError(t, err)
		assert.Equal(t, int64(5), size)
		assert.True(t, partial)
	})

	t.Run("Unreadable root directory is an error", func(t *testing.T) {
		_, _, err := DirSizeWithContext(context.Background(), locked)
		require.Error(t, err)
	})
}
