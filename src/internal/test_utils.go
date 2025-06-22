package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

var SampleDataBytes = []byte("This is sample") //nolint: gochecknoglobals // Effectively const

func setupDirectories(t *testing.T, dirs ...string) {
	t.Helper()
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
	}
}

func setupFilesWithData(t *testing.T, data []byte, files ...string) {
	t.Helper()
	for _, file := range files {
		err := os.WriteFile(file, data, 0644)
		require.NoError(t, err)
	}
}

func setupFiles(t *testing.T, files ...string) {
	setupFilesWithData(t, SampleDataBytes, files...)
}

// -------------------- Model setup utils

func defaultTestModel(dirs ...string) model {
	m := defaultModelConfig(false, false, false, dirs)
	_, _ = TeaUpdate(&m, tea.WindowSizeMsg{Width: 2 * common.MinimumWidth, Height: 2 * common.MinimumHeight})
	return m
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

// --------------------  Bubletea utilities

// TeaUpdate : Utility to send update to model , majorly used in tests
// Not using pointer receiver as this is more like a utility, than
// a member function of model
// Todo : Consider wrapping TeaUpdate with a helper that both forwards the return
// values and does a require.NoError(t, err)
func TeaUpdate(m *model, msg tea.Msg) (tea.Cmd, error) {
	resModel, cmd := m.Update(msg)

	mObj, ok := resModel.(model)
	if !ok {
		return cmd, fmt.Errorf("unexpected model type: %T", resModel)
	}
	*m = mObj
	return cmd, nil
}

func TeaUpdateWithErrCheck(t *testing.T, m *model, msg tea.Msg) tea.Cmd {
	cmd, err := TeaUpdate(m, msg)
	require.NoError(t, err)
	return cmd
}

// Is the command tea.quit, or a batch that contains tea.quit
func IsTeaQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	msg := cmd()
	switch msg := msg.(type) {
	case tea.QuitMsg:
		return true
	case tea.BatchMsg:
		for _, curCmd := range msg {
			if IsTeaQuit(curCmd) {
				return true
			}
		}
		return false
	default:
		return false
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

// -------------- Validation Utilities

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

// Helper function to verify file or directory exists
func verifyPathExists(t *testing.T, path, message string) {
	t.Helper()
	info, err := os.Stat(path)
	require.NoError(t, err, message)
	if info.IsDir() {
		assert.DirExists(t, path, message)
	} else {
		assert.FileExists(t, path, message)
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
	assert.NotEmpty(t, m.copyItems.items, "Clipboard should not be cleared when paste is prevented")
}

// Helper function to verify successful paste results
func verifySuccessfulPasteResults(t *testing.T, targetDir string, expectedDestFiles []string, originalPath string, shouldOriginalExist bool) {
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

// -------------- Other utilities

// Helper function to find item index in panel by name
func findItemIndexInPanel(panel *filePanel, itemName string) int {
	for i, elem := range panel.element {
		if elem.name == itemName {
			return i
		}
	}
	return -1
}

// Helper function to navigate to target directory if different from start
func navigateToTargetDir(t *testing.T, m *model, startDir, targetDir string) {
	t.Helper()
	if targetDir != startDir {
		err := m.updateCurrentFilePanelDir(targetDir)
		require.NoError(t, err)
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
