package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFileAtomicallyRemovesReplacementAfterPublishFailure(t *testing.T) {
	parentDir := t.TempDir()
	destinationDir := filepath.Join(parentDir, "config.toml")
	require.NoError(t, os.Mkdir(destinationDir, 0o700))

	err := writeFileAtomically(destinationDir, []byte("debug = true\n"), ConfigFilePerm)
	require.Error(t, err)

	info, statErr := os.Stat(destinationDir)
	require.NoError(t, statErr)
	assert.True(t, info.IsDir(), "failed publication must leave the destination unchanged")

	replacementFiles, globErr := filepath.Glob(filepath.Join(parentDir, ".config.toml.tmp-*"))
	require.NoError(t, globErr)
	assert.Empty(t, replacementFiles, "failed publication must remove its temporary replacement")
}

func TestWriteFileAtomicallyPreservesSymlinkAndReplacesTarget(t *testing.T) {
	parentDir := t.TempDir()
	targetDir := filepath.Join(parentDir, "config")
	require.NoError(t, os.Mkdir(targetDir, 0o700))

	targetPath := filepath.Join(targetDir, "config.toml")
	require.NoError(t, os.WriteFile(targetPath, []byte("debug = false\n"), 0o640))

	symlinkPath := filepath.Join(parentDir, "config.toml")
	if err := os.Symlink(targetPath, symlinkPath); err != nil {
		t.Skipf("creating symlinks is not supported in this environment: %v", err)
	}

	repairedData := []byte("debug = true\n")
	require.NoError(t, writeFileAtomically(symlinkPath, repairedData, 0o640))

	symlinkInfo, err := os.Lstat(symlinkPath)
	require.NoError(t, err)
	assert.NotZero(t, symlinkInfo.Mode()&os.ModeSymlink, "atomic replacement must preserve the symlink")

	targetData, err := os.ReadFile(targetPath)
	require.NoError(t, err)
	assert.Equal(t, repairedData, targetData)
}
