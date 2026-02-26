package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	zoxidelib "github.com/lazysegtree/go-zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/ui/filepanel"

	"github.com/yorukot/superfile/src/internal/common"
)

const DefaultTestTick = 10 * time.Millisecond
const DefaultTestTimeout = time.Second
const DefaultTestModelWidth = 2 * common.MinimumWidth
const DefaultTestModelHeight = 2 * common.MinimumHeight

// -------------------- Model setup utils

func defaultTestModel(dirs ...string) *model {
	m := defaultModelConfig(false, false, false, dirs, nil)
	return setModelParamsForTest(m, true)
}

// TODO: Change this to a better API. passing opts
// WithZClient(), WithFooter()
func defaultTestModelWithZClient(zClient *zoxidelib.Client, dirs ...string) *model {
	m := defaultModelConfig(false, false, false, dirs, zClient)
	return setModelParamsForTest(m, true)
}

func defaultTestModelWithFooterAndFilePreview(dirs ...string) *model {
	m := defaultModelConfig(false, true, false, dirs, nil)
	return setModelParamsForTest(m, false)
}

func defaultTestModelWithFilePreview(dirs ...string) *model {
	m := defaultModelConfig(false, false, false, dirs, nil)
	return setModelParamsForTest(m, false)
}

func setModelParamsForTest(m *model, disablePreview bool) *model {
	m.disableMetadata = true
	if disablePreview {
		m.fileModel.FilePreview.Close()
	}
	// async size updates like preview panel content update
	// will not be done
	TeaUpdate(m, tea.WindowSizeMsg{Width: DefaultTestModelWidth, Height: DefaultTestModelHeight})
	return m
}

// Helper function to setup panel mode and selection
func setupPanelModeAndSelection(t *testing.T, m *model, useSelectMode bool, itemName string, selectedItems []string) {
	t.Helper()
	panel := m.getFocusedFilePanel()

	if useSelectMode {
		// Switch to select mode and set selected items
		m.getFocusedFilePanel().ChangeFilePanelMode()
		require.Equal(t, filepanel.SelectMode, panel.PanelMode)
		for _, item := range selectedItems {
			panel.SetSelected(item)
		}
	} else {
		// Find the item in browser mode
		setFilePanelSelectedItemByName(t, panel, itemName)
	}
}

// --------------------  Bubletea utilities

// TODO : Should we validate that returned value is of type *model ?
// and equal to m ? We are assuming that to be true as of now
func TeaUpdate(m *model, msg tea.Msg) tea.Cmd {
	_, cmd := m.Update(msg)
	return cmd
}

// Is the command tea.quit, or a batch that contains tea.quit
func IsTeaQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	// Ignore commands with longer IO Operations, which waits on a channel
	msg := ExecuteTeaCmdWithTimeout(cmd, time.Millisecond)
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

func ExecuteTeaCmdWithTimeout(cmd tea.Cmd, timeout time.Duration) tea.Msg {
	result := make(chan tea.Msg, 1)
	go func() {
		result <- cmd()
	}()
	select {
	case msg := <-result:
		return msg
	case <-time.After(timeout):
		return nil
	}
}

// Helper function to perform copy or cut operation
func performCopyOrCutOperation(t *testing.T, m *model, isCut bool) {
	t.Helper()
	if isCut {
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.CutItems[0]))
	} else {
		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.CopyItems[0]))
	}
}

// -------------- Validation Utilities

// Helper function to verify clipboard state after copy/cut
func verifyClipboardState(t *testing.T, m *model, isCut bool, useSelectMode bool, selectedItemsCount int) {
	t.Helper()
	assert.Equal(t, isCut, m.clipboard.IsCut(), "Clipboard cut state should match operation")
	if useSelectMode {
		assert.Len(t, m.clipboard.GetItems(), selectedItemsCount, "Clipboard should contain all selected items")
	} else {
		assert.Len(t, m.clipboard.GetItems(), 1, "Clipboard should contain one item")
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
	}, DefaultTestTimeout, DefaultTestTick, message)
}

// Helper function to verify expected destination files exist
func verifyDestinationFiles(t *testing.T, targetDir string, expectedDestFiles []string) {
	t.Helper()
	for _, expectedFile := range expectedDestFiles {
		destPath := filepath.Join(targetDir, expectedFile)
		assert.Eventually(t, func() bool {
			_, err := os.Stat(destPath)
			return err == nil
		}, DefaultTestTimeout, DefaultTestTick, "%s should exist in destination", expectedFile)
	}
}

// Helper function to verify prevented paste results
func verifyPreventedPasteResults(t *testing.T, m *model, originalPath string) {
	t.Helper()
	if originalPath != "" {
		verifyPathExists(t, originalPath, "Original file should still exist when paste is prevented")
	}
	// Clipboard should not be cleared when paste is prevented
	assert.NotEmpty(t, m.clipboard.GetItems(), "Clipboard should not be cleared when paste is prevented")
}

// Helper function to verify successful paste results
func verifySuccessfulPasteResults(t *testing.T, targetDir string, expectedDestFiles []string,
	originalPath string, shouldOriginalExist bool) {
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
}

// -------------- Other utilities
// Helper function to navigate to target directory if different from start
func navigateToTargetDir(t *testing.T, m *model, startDir, targetDir string) {
	t.Helper()
	if targetDir != startDir {
		err := m.updateCurrentFilePanelDir(targetDir)
		require.NoError(t, err)
		TeaUpdate(m, nil)
	}
}

// Helper function to get original path for existence check
func getOriginalPath(useSelectMode bool, itemName, startDir string) string {
	if !useSelectMode && itemName != "" {
		return filepath.Join(startDir, itemName)
	}
	return ""
}

func setFilePanelSelectedItemByLocation(t *testing.T, panel *filepanel.Model, filePath string) {
	t.Helper()
	idx := panel.FindElementIndexByLocation(filePath)
	require.NotEqual(t, -1, idx, "%s should be found in panel", filePath)
	panel.SetCursorPosition(idx)
}

func setFilePanelSelectedItemByName(t *testing.T, panel *filepanel.Model, fileName string) {
	t.Helper()
	idx := panel.FindElementIndexByName(fileName)
	require.NotEqual(t, -1, idx, "%s should be found in panel", fileName)
	panel.SetCursorPosition(idx)
}

func splitPanelAsync(p *TeaProg) {
	p.SendKey(common.Hotkeys.OpenSPFPrompt[0])
	p.SendKey("split")
	p.Send(tea.KeyMsg{Type: tea.KeyEnter})
	p.Send(tea.KeyMsg{Type: tea.KeyEsc})
}
