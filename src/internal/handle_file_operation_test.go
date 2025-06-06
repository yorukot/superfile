package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

func TestCompressSelectedFiles(t *testing.T) {
	curTestDir := t.TempDir()
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(curTestDir, "file1.txt")
	file2 := filepath.Join(dir1, "file2.txt")

	setupDirectories(t, curTestDir, dir1, dir2)
	setupFiles(t, file1, file2)

	// Note this is to validate the end to end user interface, not to extensively validate
	// that compress works as expected. For that, we have TestZipSources

	// Need to test that
	// 1 - Compress single file (Browser Mode)
	// 2 - Compress Single directory with files (Browser Mode)
	// 3 - Compress single file where cursor is pointed when nothing is selected (Select Mode)
	// 4 - Compress single selected file in Select Mode where cursor points to different file
	// 5 - Compress multiple selected files and directories
	// 6 - Pressing compress hotkey on empty panel doesn't do anything or crashes on both browser/select mode

	// Copied from CopyTest. Todo - work on it.

	testdata := []struct {
		name              string
		startDir          string
		cursor            int
		selectMode        bool
		selectedElem      []string
		expectedZipName   string
		cursorIndexForZip int
		extractedDirName  string
		// Relative to extractedDir
		expectedFilesAfterExtract []string
	}{
		{
			name:                      "Single File Compress",
			startDir:                  curTestDir,
			cursor:                    2,
			selectMode:                false,
			selectedElem:              nil,
			expectedZipName:           "file1.zip",
			cursorIndexForZip:         3,
			extractedDirName:          "file1",
			expectedFilesAfterExtract: []string{"file1.txt"},
		},
		{
			name:                      "Single Directory Compress",
			startDir:                  curTestDir,
			cursor:                    0,
			selectMode:                false,
			selectedElem:              nil,
			expectedZipName:           "dir1.zip",
			cursorIndexForZip:         2,
			extractedDirName:          "dir1(1)",
			expectedFilesAfterExtract: []string{filepath.Join("dir1", "file2.txt")},
		},
		{
			name:                      "Single File Compress with select mode without selection",
			startDir:                  curTestDir,
			cursor:                    2,
			selectMode:                true,
			selectedElem:              []string{},
			expectedZipName:           "file1.zip",
			cursorIndexForZip:         3,
			extractedDirName:          "file1",
			expectedFilesAfterExtract: []string{"file1.txt"},
		},
		{
			name:                      "Single File Compress with select mode with different cursor and selection",
			startDir:                  curTestDir,
			cursor:                    0, // points to dir1
			selectMode:                true,
			selectedElem:              []string{file1},
			expectedZipName:           "file1.zip",
			cursorIndexForZip:         3,
			extractedDirName:          "file1",
			expectedFilesAfterExtract: []string{"file1.txt"},
		},
		{
			name:                      "Multi file compression",
			startDir:                  curTestDir,
			cursor:                    0, // points to dir1
			selectMode:                true,
			selectedElem:              []string{dir2, dir1, file1},
			expectedZipName:           "dir2.zip",
			cursorIndexForZip:         2,
			extractedDirName:          "dir2(1)",
			expectedFilesAfterExtract: []string{"dir2", filepath.Join("dir1", "file2.txt"), "file1.txt"},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			m := defaultTestModel(tt.startDir)
			require.Greater(t, len(m.getFocusedFilePanel().element), tt.cursor)
			// Update cursor
			m.getFocusedFilePanel().cursor = tt.cursor

			require.Equal(t, browserMode, m.getFocusedFilePanel().panelMode)
			if tt.selectMode {
				m.changeFilePanelMode()
				m.getFocusedFilePanel().selected = tt.selectedElem
			}

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.CompressFile[0]))
			zipFile := filepath.Join(tt.startDir, tt.expectedZipName)
			// Actual compress may take time, since its an os operations
			assert.Eventually(t, func() bool {
				_, err := os.Lstat(zipFile)
				return err == nil
			}, time.Second, 10*time.Millisecond)

			// No-op update to get the filepanel updated
			// Todo - This should not be needed. Only operation finish SPF should refresh
			// on its own
			TeaUpdateWithErrCheck(t, &m, nil)

			require.Greater(t, len(m.getFocusedFilePanel().element), tt.cursorIndexForZip)
			assert.Equal(t, zipFile, m.getFocusedFilePanel().element[tt.cursorIndexForZip].location,
				"%s does not exists at index %d among %v", zipFile, tt.cursorIndexForZip,
				m.getFocusedFilePanel().element)

			m.getFocusedFilePanel().cursor = tt.cursorIndexForZip

			TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.ExtractFile[0]))
			// File extraction is supposedly async. So function's return doesn't means its done.
			extractedDir := filepath.Join(tt.startDir, tt.extractedDirName)
			assert.Eventually(t, func() bool {
				for _, f := range tt.expectedFilesAfterExtract {
					_, err := os.Stat(filepath.Join(extractedDir, f))
					if err != nil {
						return false
					}
				}
				return true
			}, time.Second, 10*time.Millisecond, "Extraction of files failed Required - [%s]+%v",
				extractedDir, tt.expectedFilesAfterExtract)

			require.NoError(t, os.RemoveAll(extractedDir))
			require.NoError(t, os.RemoveAll(zipFile))
		})
	}

	t.Run("Compress on Empty panel", func(t *testing.T) {
		m := defaultTestModel(dir2)
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.CompressFile[0]))
		// Should not crash. Nothing should happen. If there is a crash, it will be caught
		entries, err := os.ReadDir(dir2)
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}
