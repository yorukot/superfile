package internal

import (
	"os"
	"path/filepath"
	"strings"
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

// Helper function to find item index in panel by name
func findItemIndexInPanel(panel *filePanel, itemName string) int {
	for i, elem := range panel.element {
		if elem.name == itemName {
			return i
		}
	}
	return -1
}

// Helper function to setup panel mode and selection
func setupPanelModeAndSelection(t *testing.T, m *model, useSelectMode bool, itemName string, selectedItems []string) {
	t.Helper()
	panel := m.getFocusedFilePanel()

	if useSelectMode {
		// Switch to select mode and set selected items
		m.changeFilePanelMode()
		require.Equal(t, selectMode, panel.panelMode)
		panel.selected = selectedItems
	} else {
		// Find the item in browser mode
		itemIndex := findItemIndexInPanel(panel, itemName)
		require.NotEqual(t, -1, itemIndex, "%s should be found in panel", itemName)
		panel.cursor = itemIndex
	}
}

// Helper function to perform copy or cut operation
func performCopyOrCutOperation(t *testing.T, m *model, isCut bool) {
	t.Helper()
	if isCut {
		TeaUpdateWithErrCheck(t, m, utils.TeaRuneKeyMsg(common.Hotkeys.CutItems[0]))
	} else {
		TeaUpdateWithErrCheck(t, m, utils.TeaRuneKeyMsg(common.Hotkeys.CopyItems[0]))
	}
}

// Helper function to verify clipboard state after copy/cut
func verifyClipboardState(t *testing.T, m *model, isCut bool, useSelectMode bool, selectedItemsCount int) {
	t.Helper()
	assert.Equal(t, isCut, m.copyItems.cut, "Clipboard cut state should match operation")
	if useSelectMode {
		assert.Len(t, m.copyItems.items, selectedItemsCount, "Clipboard should contain all selected items")
	} else {
		assert.Len(t, m.copyItems.items, 1, "Clipboard should contain one item")
	}
}

// Helper function to navigate to target directory if different from start
func navigateToTargetDir(t *testing.T, m *model, startDir, targetDir string) {
	t.Helper()
	if targetDir != startDir {
		m.updateCurrentFilePanelDir(targetDir)
		TeaUpdateWithErrCheck(t, m, nil)
	}
}

// Helper function to get original path for existence check
func getOriginalPath(useSelectMode bool, itemName, startDir string) string {
	if !useSelectMode && itemName != "" {
		return filepath.Join(startDir, itemName)
	}
	return ""
}

// Helper function to verify file or directory exists
func verifyPathExists(t *testing.T, path, message string) {
	t.Helper()
	if strings.Contains(path, ".txt") {
		assert.FileExists(t, path, message)
	} else {
		assert.DirExists(t, path, message)
	}
}

// Helper function to verify file or directory doesn't exist after cut
func verifyPathNotExistsEventually(t *testing.T, path, message string) {
	t.Helper()
	assert.Eventually(t, func() bool {
		_, err := os.Stat(path)
		return os.IsNotExist(err)
	}, time.Second, 10*time.Millisecond, message)
}

// Helper function to verify expected destination files exist
func verifyDestinationFiles(t *testing.T, targetDir string, expectedDestFiles []string) {
	t.Helper()
	for _, expectedFile := range expectedDestFiles {
		destPath := filepath.Join(targetDir, expectedFile)
		assert.Eventually(t, func() bool {
			_, err := os.Stat(destPath)
			return err == nil
		}, time.Second, 10*time.Millisecond, "%s should exist in destination", expectedFile)
	}
}

// Helper function to verify prevented paste results
func verifyPreventedPasteResults(t *testing.T, m *model, originalPath string) {
	t.Helper()
	if originalPath != "" {
		verifyPathExists(t, originalPath, "Original file should still exist when paste is prevented")
	}
	// Clipboard should not be cleared when paste is prevented
	assert.NotEqual(t, 0, len(m.copyItems.items), "Clipboard should not be cleared when paste is prevented")
}

// Helper function to verify successful paste results
func verifySuccessfulPasteResults(t *testing.T, targetDir string, expectedDestFiles []string, originalPath string, shouldOriginalExist bool, m *model) {
	t.Helper()
	// Verify expected files were created in destination
	verifyDestinationFiles(t, targetDir, expectedDestFiles)

	// Verify original file existence based on operation type
	if originalPath != "" {
		if shouldOriginalExist {
			verifyPathExists(t, originalPath, "Original file should exist after copy operation")
		} else {
			verifyPathNotExistsEventually(t, originalPath, "Original file should not exist after cut operation")
		}
	}

	// TODO: Need to add a test to verify clipboard state.
}

// Helper function to setup model and perform copy/cut operation
func setupModelAndPerformOperation(t *testing.T, startDir string, useSelectMode bool, itemName string, selectedItems []string, isCut bool) *model {
	t.Helper()
	m := defaultTestModel(startDir)
	TeaUpdateWithErrCheck(t, &m, nil)

	setupPanelModeAndSelection(t, &m, useSelectMode, itemName, selectedItems)
	performCopyOrCutOperation(t, &m, isCut)

	selectedItemsCount := len(selectedItems)
	if !useSelectMode {
		selectedItemsCount = 1
	}
	verifyClipboardState(t, &m, isCut, useSelectMode, selectedItemsCount)

	return &m
}

func TestPasteItem(t *testing.T) {
	curTestDir := t.TempDir()
	sourceDir := filepath.Join(curTestDir, "source")
	destDir := filepath.Join(curTestDir, "dest")
	subDir := filepath.Join(sourceDir, "subdir")
	file1 := filepath.Join(sourceDir, "file1.txt")
	file2 := filepath.Join(sourceDir, "file2.txt")
	dirFile1 := filepath.Join(subDir, "dirfile1.txt")

	setupDirectories(t, curTestDir, sourceDir, destDir, subDir)
	setupFiles(t, file1, file2, dirFile1)

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

			// Navigate to target directory
			navigateToTargetDir(t, m, tt.startDir, tt.targetDir)

			// Get original file path for existence check
			originalPath := getOriginalPath(tt.selectMode, tt.itemName, tt.startDir)

			// Perform paste operation
			TeaUpdateWithErrCheck(t, m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

			// Verify results based on whether paste should be prevented
			if tt.shouldPreventPaste {
				verifyPreventedPasteResults(t, m, originalPath)
			} else {
				verifySuccessfulPasteResults(t, tt.targetDir, tt.expectedDestFiles, originalPath, tt.shouldOriginalExist, m)
			}
		})
	}

	// Special test cases that don't fit the table-driven pattern
	t.Run("Paste with Empty Clipboard", func(t *testing.T) {
		emptyTestDir := filepath.Join(curTestDir, "empty_test")
		setupDirectories(t, emptyTestDir)

		m := defaultTestModel(emptyTestDir)
		TeaUpdateWithErrCheck(t, &m, nil)

		// Ensure clipboard is empty
		m.copyItems.items = []string{}

		// Get initial count
		entriesBefore, err := os.ReadDir(emptyTestDir)
		require.NoError(t, err)

		// Attempt to paste (should do nothing)
		TeaUpdateWithErrCheck(t, &m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Should not crash and no new files should be created
		entriesAfter, err := os.ReadDir(emptyTestDir)
		require.NoError(t, err)

		assert.Equal(t, len(entriesBefore), len(entriesAfter), "No new files should be created when pasting with empty clipboard")
	})

	t.Run("Multiple Items Copy and Paste", func(t *testing.T) {
		// Create fresh files for this test
		multiFile1 := filepath.Join(sourceDir, "multi1.txt")
		multiFile2 := filepath.Join(sourceDir, "multi2.txt")
		setupFiles(t, multiFile1, multiFile2)

		selectedItems := []string{multiFile1, multiFile2}
		m := setupModelAndPerformOperation(t, sourceDir, true, "", selectedItems, false)

		// Navigate to destination
		navigateToTargetDir(t, m, sourceDir, destDir)

		// Paste items
		TeaUpdateWithErrCheck(t, m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Verify both files were copied
		expectedDestFiles := []string{"multi1.txt", "multi2.txt"}
		verifyDestinationFiles(t, destDir, expectedDestFiles)
	})

	t.Run("Cut into Subdirectory Prevention", func(t *testing.T) {
		// Create a separate subdirectory for this test to avoid conflicts with table-driven tests
		testSubDir := filepath.Join(sourceDir, "testsubdir")
		testDirFile := filepath.Join(testSubDir, "testdirfile.txt")
		setupDirectories(t, testSubDir)
		setupFiles(t, testDirFile)

		// Test the logic that prevents cutting a directory into its subdirectory
		m := setupModelAndPerformOperation(t, sourceDir, false, "testsubdir", nil, true)

		// Navigate into the subdirectory and try to paste there (should be prevented)
		navigateToTargetDir(t, m, sourceDir, testSubDir)
		TeaUpdateWithErrCheck(t, m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Directory should still exist in original location after prevention
		assert.DirExists(t, testSubDir, "Directory should still exist after failed paste into subdirectory")
	})

	t.Run("Duplicate File Handling", func(t *testing.T) {
		// Create a file to copy
		dupFile := filepath.Join(sourceDir, "duplicate.txt")
		setupFiles(t, dupFile)

		m := setupModelAndPerformOperation(t, sourceDir, false, "duplicate.txt", nil, false)

		// Navigate to destination and paste
		navigateToTargetDir(t, m, sourceDir, destDir)
		TeaUpdateWithErrCheck(t, m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Verify first copy
		verifyDestinationFiles(t, destDir, []string{"duplicate.txt"})

		// Paste again to test duplicate handling
		TeaUpdateWithErrCheck(t, m, utils.TeaRuneKeyMsg(common.Hotkeys.PasteItems[0]))

		// Verify duplicate file with different name
		verifyDestinationFiles(t, destDir, []string{"duplicate(1).txt"})
	})
}
