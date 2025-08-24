package internal

import (
	"os"
	"path/filepath"
	"testing"

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

	utils.SetupDirectories(t, curTestDir, dir1, dir2)
	utils.SetupFiles(t, file1, file2)

	// Note this is to validate the end to end user interface, not to extensively validate
	// that compress works as expected. For that, we have TestZipSources

	// Need to test that
	// 1 - Compress single file (Browser Mode)
	// 2 - Compress Single directory with files (Browser Mode)
	// 3 - Compress single file where cursor is pointed when nothing is selected (Select Mode)
	// 4 - Compress single selected file in Select Mode where cursor points to different file
	// 5 - Compress multiple selected files and directories
	// 6 - Pressing compress hotkey on empty panel doesn't do anything or crashes on both browser/select mode

	// Copied from CopyTest. TODO - work on it.

	testdata := []struct {
		name             string
		startDir         string
		cursor           int
		selectMode       bool
		selectedElem     []string
		expectedZipName  string
		extractedDirName string
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
			extractedDirName:          "dir2(1)",
			expectedFilesAfterExtract: []string{"dir2", filepath.Join("dir1", "file2.txt"), "file1.txt"},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			m := defaultTestModel(tt.startDir)
			p := NewTestTeaProgWithEventLoop(t, m)
			require.Greater(t, len(m.getFocusedFilePanel().element), tt.cursor)
			// Update cursor
			m.getFocusedFilePanel().cursor = tt.cursor

			require.Equal(t, browserMode, m.getFocusedFilePanel().panelMode)
			if tt.selectMode {
				m.getFocusedFilePanel().changeFilePanelMode()
				m.getFocusedFilePanel().selected = tt.selectedElem
			}

			p.SendKey(common.Hotkeys.CompressFile[0])
			zipFile := filepath.Join(tt.startDir, tt.expectedZipName)
			// Actual compress may take time, since its an os operations
			assert.Eventually(t, func() bool {
				_, err := os.Lstat(zipFile)
				return err == nil
			}, DefaultTestTimeout, DefaultTestTick)

			// Assert zip file exists right after compression
			require.FileExists(t, zipFile, "Expected zip file does not exist after compression")

			// No-op update to get the filepanel updated
			// TODO - This should not be needed. Only operation finish SPF should refresh
			// on its own
			p.SendDirectly(nil)

			setFilePanelSelectedItemByLocation(t, m.getFocusedFilePanel(), zipFile)
			t.Logf("Panel elements : %v, cursor : %v",
				m.getFocusedFilePanel().element, m.getFocusedFilePanel().cursor)
			selectedItemLocation := m.getFocusedFilePanel().getSelectedItem().location
			assert.Equal(t, zipFile, selectedItemLocation)
			// Ensure we are extracting the zip file, not a directory
			fileInfo, err := os.Stat(selectedItemLocation)
			require.NoError(t, err, "Failed to stat panel location before extraction")
			require.False(t, fileInfo.IsDir(),
				"Panel location for extraction is a directory, expected a zip file: %s", selectedItemLocation)

			p.SendKey(common.Hotkeys.ExtractFile[0])
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
			}, DefaultTestTimeout, DefaultTestTick, "Extraction of files failed Required - [%s]+%v",
				extractedDir, tt.expectedFilesAfterExtract)

			require.NoError(t, os.RemoveAll(extractedDir))
			require.NoError(t, os.RemoveAll(zipFile))
		})
	}

	t.Run("Compress on Empty panel", func(t *testing.T) {
		NewTestTeaProgWithEventLoop(t, defaultTestModel(dir2)).
			SendKey(common.Hotkeys.CompressFile[0])
		// Should not crash. Nothing should happen. If there is a crash, it will be caught
		entries, err := os.ReadDir(dir2)
		require.NoError(t, err)
		assert.Empty(t, entries)
	})
}

func TestPasteItem(t *testing.T) {
	curTestDir := t.TempDir()
	sourceDir := filepath.Join(curTestDir, "source")
	destDir := filepath.Join(curTestDir, "dest")
	subDir := filepath.Join(sourceDir, "subdir")
	file1 := filepath.Join(sourceDir, "file1.txt")
	file2 := filepath.Join(sourceDir, "file2.txt")
	dirFile1 := filepath.Join(subDir, "dirfile1.txt")

	utils.SetupDirectories(t, curTestDir, sourceDir, destDir, subDir)
	utils.SetupFiles(t, file1, file2, dirFile1)

	testdata := []struct {
		name                 string
		startDir             string
		targetDir            string
		itemName             string
		isCut                bool
		selectMode           bool
		selectedItems        []string
		shouldClipboardClear bool
		shouldOriginalExist  bool
		expectedDestFiles    []string
		shouldPreventPaste   bool
		description          string
	}{
		{
			name:                 "Copy Single File",
			startDir:             sourceDir,
			targetDir:            destDir,
			itemName:             "file1.txt",
			isCut:                false,
			selectMode:           false,
			selectedItems:        nil,
			shouldClipboardClear: false,
			shouldOriginalExist:  true,
			expectedDestFiles:    []string{"file1.txt"},
			shouldPreventPaste:   false,
			description:          "Copy a single file from source to destination",
		},
		{
			name:                 "Cut Single File",
			startDir:             sourceDir,
			targetDir:            destDir,
			itemName:             "file2.txt",
			isCut:                true,
			selectMode:           false,
			selectedItems:        nil,
			shouldClipboardClear: true,
			shouldOriginalExist:  false,
			expectedDestFiles:    []string{"file2.txt"},
			shouldPreventPaste:   false,
			description:          "Cut a single file from source to destination",
		},
		{
			name:                 "Cut Directory into Same Location",
			startDir:             sourceDir,
			targetDir:            sourceDir, // Same directory
			itemName:             "subdir",
			isCut:                true,
			selectMode:           false,
			selectedItems:        nil,
			shouldClipboardClear: false,      // Should not clear because paste should be prevented
			shouldOriginalExist:  true,       // Should still exist because paste prevented
			expectedDestFiles:    []string{}, // No files should be created in dest
			shouldPreventPaste:   true,
			description:          "Cutting directory into same location should be prevented",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			m := setupModelAndPerformOperation(t, tt.startDir, tt.selectMode, tt.itemName, tt.selectedItems, tt.isCut)
			p := NewTestTeaProgWithEventLoop(t, m)
			// Navigate to target directory
			navigateToTargetDir(t, m, tt.startDir, tt.targetDir)

			// Get original file path for existence check
			originalPath := getOriginalPath(tt.selectMode, tt.itemName, tt.startDir)

			// Perform paste operation
			p.SendKey(common.Hotkeys.PasteItems[0])

			// Verify results based on whether paste should be prevented
			if tt.shouldPreventPaste {
				verifyPreventedPasteResults(t, m, originalPath)
			} else {
				verifySuccessfulPasteResults(t, tt.targetDir, tt.expectedDestFiles, originalPath, tt.shouldOriginalExist)
			}
			// Checking separately, as this is something independent of tt.shouldPreventPaste
			if tt.shouldClipboardClear {
				assert.Empty(t, p.m.copyItems.items, "Clipboard should be cleared after successful cut-paste")
			} else {
				assert.NotEmpty(t, p.m.copyItems.items, "Clipboard should remain after copy-paste")
			}
		})
	}

	// Special test cases that don't fit the table-driven pattern
	t.Run("Paste with Empty Clipboard", func(t *testing.T) {
		emptyTestDir := filepath.Join(curTestDir, "empty_test")
		utils.SetupDirectories(t, emptyTestDir)
		m := defaultTestModel(emptyTestDir)
		p := NewTestTeaProgWithEventLoop(t, m)

		// Ensure clipboard is empty
		m.copyItems.items = []string{}

		// Get initial count
		entriesBefore, err := os.ReadDir(emptyTestDir)
		require.NoError(t, err)

		// Attempt to paste (should do nothing)
		p.SendKey(common.Hotkeys.PasteItems[0])

		// Should not crash and no new files should be created
		entriesAfter, err := os.ReadDir(emptyTestDir)
		require.NoError(t, err)

		assert.Len(t, entriesAfter, len(entriesBefore),
			"No new files should be created when pasting with empty clipboard")
	})

	t.Run("Multiple Items Copy and Paste", func(t *testing.T) {
		// Create fresh files for this test
		multiFile1 := filepath.Join(sourceDir, "multi1.txt")
		multiFile2 := filepath.Join(sourceDir, "multi2.txt")
		utils.SetupFiles(t, multiFile1, multiFile2)

		selectedItems := []string{multiFile1, multiFile2}
		m := setupModelAndPerformOperation(t, sourceDir, true, "", selectedItems, false)
		p := NewTestTeaProgWithEventLoop(t, m)

		// Navigate to destination
		navigateToTargetDir(t, m, sourceDir, destDir)

		// Paste items
		p.SendKey(common.Hotkeys.PasteItems[0])

		// Verify both files were copied
		expectedDestFiles := []string{"multi1.txt", "multi2.txt"}
		verifyDestinationFiles(t, destDir, expectedDestFiles)
	})

	t.Run("Cut into Subdirectory Prevention", func(t *testing.T) {
		// Create a separate subdirectory for this test to avoid conflicts with table-driven tests
		testSubDir := filepath.Join(sourceDir, "testsubdir")
		testDirFile := filepath.Join(testSubDir, "testdirfile.txt")
		utils.SetupDirectories(t, testSubDir)
		utils.SetupFiles(t, testDirFile)

		// Test the logic that prevents cutting a directory into its subdirectory
		m := setupModelAndPerformOperation(t, sourceDir, false, "testsubdir", nil, true)
		p := NewTestTeaProgWithEventLoop(t, m)

		// Navigate into the subdirectory and try to paste there (should be prevented)
		navigateToTargetDir(t, m, sourceDir, testSubDir)
		p.SendKey(common.Hotkeys.PasteItems[0])

		// Directory should still exist in original location after prevention
		assert.DirExists(t, testSubDir, "Directory should still exist after failed paste into subdirectory")
	})

	t.Run("Duplicate File Handling", func(t *testing.T) {
		// Create a file to copy
		dupFile := filepath.Join(sourceDir, "duplicate.txt")
		utils.SetupFiles(t, dupFile)

		m := setupModelAndPerformOperation(t, sourceDir, false, "duplicate.txt", nil, false)
		p := NewTestTeaProgWithEventLoop(t, m)
		// Navigate to destination and paste
		navigateToTargetDir(t, m, sourceDir, destDir)
		p.SendKey(common.Hotkeys.PasteItems[0])

		// Verify first copy
		verifyDestinationFiles(t, destDir, []string{"duplicate.txt"})

		// Paste again to test duplicate handling
		p.SendKey(common.Hotkeys.PasteItems[0])

		// Verify duplicate file with different name
		verifyDestinationFiles(t, destDir, []string{"duplicate(1).txt"})
	})
}

// ------  Very specific utilities that are required for this test case file only

// Helper function to setup model and perform copy/cut operation
func setupModelAndPerformOperation(t *testing.T, startDir string, useSelectMode bool,
	itemName string, selectedItems []string, isCut bool) *model {
	t.Helper()
	m := defaultTestModel(startDir)
	TeaUpdateWithErrCheck(m, nil)

	setupPanelModeAndSelection(t, m, useSelectMode, itemName, selectedItems)
	performCopyOrCutOperation(t, m, isCut)

	selectedItemsCount := len(selectedItems)
	if !useSelectMode {
		selectedItemsCount = 1
	}
	verifyClipboardState(t, m, isCut, useSelectMode, selectedItemsCount)

	return m
}
