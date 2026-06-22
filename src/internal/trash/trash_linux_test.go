//go:build linux

package trash

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEscapeTrashPath(t *testing.T) {
	assert.Equal(t, "/tmp/a%20b/%25file%0A", escapeTrashPath("/tmp/a b/%file\n"))
	assert.Equal(t, "relative/path", escapeTrashPath("relative/path"))
	assert.Equal(t, "%C3%A9", escapeTrashPath("é"))
}

func TestMoveCreatesFreeDesktopTrashInfo(t *testing.T) {
	dataHome := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dataHome)

	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "file with %.txt")
	require.NoError(t, os.WriteFile(src, []byte("content"), 0o644))

	result, err := Move(src)
	require.NoError(t, err)
	assert.Equal(t, BackendFreeDesktop, result.Backend)
	assert.True(t, result.StrictlyRecycled)
	assert.NoFileExists(t, src)
	assert.FileExists(t, result.TrashedPath)

	infoPath := filepath.Join(dataHome, "Trash", "info", filepath.Base(result.TrashedPath)+".trashinfo")
	info, err := os.ReadFile(infoPath)
	require.NoError(t, err)
	infoString := string(info)
	assert.True(t, strings.HasPrefix(infoString, "[Trash Info]\n"))
	assert.Contains(t, infoString, "Path="+escapeTrashPath(src)+"\n")

	dateLine := ""
	for _, line := range strings.Split(infoString, "\n") {
		if strings.HasPrefix(line, "DeletionDate=") {
			dateLine = strings.TrimPrefix(line, "DeletionDate=")
		}
	}
	require.NotEmpty(t, dateLine)
	_, err = time.ParseInLocation(trashInfoDateLayout, dateLine, time.Local)
	assert.NoError(t, err)
}

func TestMoveDoesNotOverwriteDuplicateBasename(t *testing.T) {
	dataHome := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dataHome)

	firstDir := t.TempDir()
	secondDir := t.TempDir()
	first := filepath.Join(firstDir, "duplicate.txt")
	second := filepath.Join(secondDir, "duplicate.txt")
	require.NoError(t, os.WriteFile(first, []byte("first"), 0o644))
	require.NoError(t, os.WriteFile(second, []byte("second"), 0o644))

	firstResult, err := Move(first)
	require.NoError(t, err)
	secondResult, err := Move(second)
	require.NoError(t, err)

	assert.NotEqual(t, firstResult.TrashedPath, secondResult.TrashedPath)
	assert.FileExists(t, firstResult.TrashedPath)
	assert.FileExists(t, secondResult.TrashedPath)
}

func TestMoveBoundsLongBasename(t *testing.T) {
	dataHome := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dataHome)

	srcDir := t.TempDir()
	src := filepath.Join(srcDir, strings.Repeat("a", linuxMaxFilenameBytes))
	require.NoError(t, os.WriteFile(src, []byte("content"), 0o644))

	result, err := Move(src)
	require.NoError(t, err)

	trashName := filepath.Base(result.TrashedPath)
	assert.LessOrEqual(t, len(trashName), maxTrashEntryNameBytes)
	assert.LessOrEqual(t, len(trashName+trashInfoSuffix), linuxMaxFilenameBytes)
	assert.FileExists(t, result.TrashedPath)
	assert.FileExists(t, filepath.Join(dataHome, "Trash", "info", trashName+trashInfoSuffix))
}

func TestMoveBoundsDuplicateBasename(t *testing.T) {
	dataHome := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dataHome)

	firstDir := t.TempDir()
	secondDir := t.TempDir()
	base := strings.Repeat("b", maxTrashEntryNameBytes)
	first := filepath.Join(firstDir, base)
	second := filepath.Join(secondDir, base)
	require.NoError(t, os.WriteFile(first, []byte("first"), 0o644))
	require.NoError(t, os.WriteFile(second, []byte("second"), 0o644))

	firstResult, err := Move(first)
	require.NoError(t, err)
	secondResult, err := Move(second)
	require.NoError(t, err)

	firstName := filepath.Base(firstResult.TrashedPath)
	secondName := filepath.Base(secondResult.TrashedPath)
	assert.NotEqual(t, firstName, secondName)
	assert.LessOrEqual(t, len(firstName), maxTrashEntryNameBytes)
	assert.LessOrEqual(t, len(secondName), maxTrashEntryNameBytes)
	assert.FileExists(t, firstResult.TrashedPath)
	assert.FileExists(t, secondResult.TrashedPath)
}

func TestAvailableReturnsFalseInsideTrash(t *testing.T) {
	dataHome := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dataHome)
	require.NoError(t, Init())

	assert.True(t, Available(t.TempDir()))
	assert.False(t, Available(filepath.Join(dataHome, "Trash")))
	assert.False(t, Available(filepath.Join(dataHome, "Trash", "files")))
	assert.False(t, Available(filepath.Join(dataHome, "Trash", "info")))
}

func TestMoveRejectsTrashSymlink(t *testing.T) {
	dataHome := t.TempDir()
	t.Setenv("XDG_DATA_HOME", dataHome)

	trashRoot := filepath.Join(dataHome, "Trash")
	require.NoError(t, os.Symlink(t.TempDir(), trashRoot))

	src := filepath.Join(t.TempDir(), "file.txt")
	require.NoError(t, os.WriteFile(src, []byte("content"), 0o644))

	_, err := Move(src)
	require.Error(t, err)
	assert.FileExists(t, src)
}
