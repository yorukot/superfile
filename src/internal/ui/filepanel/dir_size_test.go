package filepanel

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Returns a panel with cursor on directory "dir" (with 2 files of 5 bytes) and a file "file"
func setupDirSizePanel(t *testing.T) (Model, string) {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, "dir")
	file := filepath.Join(root, "file")
	utils.SetupDirectories(t, filepath.Join(dir, "sub"))
	utils.SetupFilesWithData(t, []byte("12345"), filepath.Join(dir, "f1"), filepath.Join(dir, "sub", "f2"), file)

	dirInfo, err := os.Stat(dir)
	require.NoError(t, err)
	fileInfo, err := os.Stat(file)
	require.NoError(t, err)

	panel := testModel(0, 0, 12, BrowserMode, []Element{
		{Name: "dir", Location: dir, Directory: true, Info: dirInfo},
		{Name: "file", Location: file, Directory: false, Info: fileInfo},
	})
	panel.Location = root
	panel.DirectoryRecords = make(map[string]directoryRecord)
	return panel, dir
}

func TestDirSizeCalculation(t *testing.T) {
	t.Run("Shows placeholder while calculating, then the size", func(t *testing.T) {
		panel, dir := setupDirSizePanel(t)
		calc := panel.StartDirSizeCalc(1)
		require.NotNil(t, calc)
		assert.Equal(t, dirSizeCalculatingText, panel.dirSizeText(panel.element[0]))

		res := calc()
		require.NoError(t, res.Err)
		assert.Equal(t, dir, res.Location)
		assert.Equal(t, int64(10), res.Size)
		assert.True(t, panel.ApplyDirSizeResult(res))
		assert.Equal(t, common.FormatFileSize(10), panel.dirSizeText(panel.element[0]))
	})

	t.Run("Not started for files", func(t *testing.T) {
		panel, _ := setupDirSizePanel(t)
		panel.cursor = 1
		assert.Nil(t, panel.StartDirSizeCalc(1))
	})

	t.Run("Not started twice while calculating", func(t *testing.T) {
		panel, _ := setupDirSizePanel(t)
		require.NotNil(t, panel.StartDirSizeCalc(1))
		assert.Nil(t, panel.StartDirSizeCalc(2))
	})

	t.Run("Partial result is prefixed", func(t *testing.T) {
		panel, dir := setupDirSizePanel(t)
		require.NotNil(t, panel.StartDirSizeCalc(1))
		assert.True(t, panel.ApplyDirSizeResult(DirSizeResult{Location: dir, ReqID: 1, Size: 10, Partial: true}))
		assert.Equal(t, dirSizePartialPrefix+common.FormatFileSize(10), panel.dirSizeText(panel.element[0]))
	})

	t.Run("Error clears the entry", func(t *testing.T) {
		panel, dir := setupDirSizePanel(t)
		require.NotNil(t, panel.StartDirSizeCalc(1))
		assert.True(t, panel.ApplyDirSizeResult(DirSizeResult{Location: dir, ReqID: 1, Err: errors.New("denied")}))
		assert.Empty(t, panel.dirSizeText(panel.element[0]))
	})
}

func TestDirSizeCaching(t *testing.T) {
	t.Run("Result is cached and recalculated on demand", func(t *testing.T) {
		panel, dir := setupDirSizePanel(t)
		require.True(t, panel.ApplyDirSizeResult(panel.StartDirSizeCalc(1)()))

		utils.SetupFilesWithData(t, []byte("1234567890"), filepath.Join(dir, "sub", "f3"))
		// Cached value is still shown, as the focused directory's own mtime did not change
		assert.Equal(t, common.FormatFileSize(10), panel.dirSizeText(panel.element[0]))

		calc := panel.StartDirSizeCalc(2)
		require.NotNil(t, calc)
		require.True(t, panel.ApplyDirSizeResult(calc()))
		assert.Equal(t, common.FormatFileSize(20), panel.dirSizeText(panel.element[0]))
	})

	t.Run("Stale result is hidden when directory mtime changes", func(t *testing.T) {
		panel, dir := setupDirSizePanel(t)
		require.True(t, panel.ApplyDirSizeResult(panel.StartDirSizeCalc(1)()))

		newTime := panel.element[0].Info.ModTime().Add(time.Hour)
		require.NoError(t, os.Chtimes(dir, newTime, newTime))
		info, err := os.Stat(dir)
		require.NoError(t, err)
		panel.element[0].Info = info
		assert.Empty(t, panel.dirSizeText(panel.element[0]))
	})

	t.Run("Outdated result is ignored", func(t *testing.T) {
		panel, dir := setupDirSizePanel(t)
		require.NotNil(t, panel.StartDirSizeCalc(1))
		assert.False(t, panel.ApplyDirSizeResult(DirSizeResult{Location: dir, ReqID: 0, Size: 99}))
		assert.Equal(t, dirSizeCalculatingText, panel.dirSizeText(panel.element[0]))
	})
}

func TestDirSizeCancellation(t *testing.T) {
	t.Run("Changing directory cancels calculation", func(t *testing.T) {
		panel, dir := setupDirSizePanel(t)
		calc := panel.StartDirSizeCalc(1)
		require.NotNil(t, calc)

		require.NoError(t, panel.UpdateCurrentFilePanelDir(dir))
		res := calc()
		require.ErrorIs(t, res.Err, context.Canceled)
		assert.False(t, panel.ApplyDirSizeResult(res))
		assert.Empty(t, panel.dirSizes)
	})

	t.Run("Finished results are kept after changing directory", func(t *testing.T) {
		panel, dir := setupDirSizePanel(t)
		require.True(t, panel.ApplyDirSizeResult(panel.StartDirSizeCalc(1)()))

		require.NoError(t, panel.UpdateCurrentFilePanelDir(dir))
		require.NoError(t, panel.ParentDirectory())
		assert.Equal(t, common.FormatFileSize(10), panel.dirSizeText(panel.element[0]))
	})
}
