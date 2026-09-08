//go:build linux

package internal

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

func TestPasteDirCutAcrossFilesystems(t *testing.T) {
	//nolint:usetesting // The regression needs a source directory on the /dev/shm mount.
	sourceDir, err := os.MkdirTemp("/dev/shm", "superfile-cross-device-")
	if err != nil {
		t.Skipf("cannot create a temporary directory on /dev/shm: %v", err)
	}
	t.Cleanup(func() { assert.NoError(t, os.RemoveAll(sourceDir)) })

	destinationDir := t.TempDir()
	sourceInfo, err := os.Stat(sourceDir)
	require.NoError(t, err)
	destinationInfo, err := os.Stat(destinationDir)
	require.NoError(t, err)
	sourceStat, ok := sourceInfo.Sys().(*syscall.Stat_t)
	require.True(t, ok)
	destinationStat, ok := destinationInfo.Sys().(*syscall.Stat_t)
	require.True(t, ok)
	if sourceStat.Dev == destinationStat.Dev {
		t.Skip("/dev/shm is not a separate filesystem")
	}

	sourceFile := filepath.Join(sourceDir, "hiyo.txt")
	destinationFile := filepath.Join(destinationDir, "hiyo.txt")
	require.NoError(t, os.WriteFile(sourceFile, []byte("cross-device paste"), 0o600))

	process := processbar.NewProcess("test", "hiyo.txt", processbar.OpCut, 1)
	processBarModel := processbar.New()
	require.NoError(t, pasteDir(sourceFile, destinationFile, &process, true, &processBarModel))

	_, err = os.Stat(sourceFile)
	assert.True(t, os.IsNotExist(err))
	content, err := os.ReadFile(destinationFile)
	require.NoError(t, err)
	assert.Equal(t, "cross-device paste", string(content))
}
