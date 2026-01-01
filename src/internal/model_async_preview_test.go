package internal

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

// TestAsyncPreviewPanelSync validates preview panel sync after file panel operations
func TestAsyncPreviewPanelSync(t *testing.T) {
	utils.SetRootLoggerToStdout(testing.Verbose())
	testDir := t.TempDir()

	// Create test files
	utils.SetupFiles(t,
		filepath.Join(testDir, "file1.txt"),
		filepath.Join(testDir, "file2.txt"),
		filepath.Join(testDir, "file3.txt"),
	)

	t.Run("Preview syncs after panel creation", func(t *testing.T) {
		// Save and restore config
		originalPreviewWidth := common.Config.FilePreviewWidth
		common.Config.FilePreviewWidth = 10
		t.Cleanup(func() {
			common.Config.FilePreviewWidth = originalPreviewWidth
		})

		m := defaultTestModel(testDir)
		p := NewTestTeaProgWithEventLoop(t, m)

		// Set window size to ensure preview is visible
		p.Send(tea.WindowSizeMsg{Width: 100, Height: 40})

		// Create new panel
		p.SendKey(common.Hotkeys.CreateNewFilePanel[0])

		// Verify preview syncs with new panel's focused item
		assert.Eventually(t, func() bool {
			if len(m.fileModel.FilePanels) != 2 {
				return false
			}
			focusedItem := m.getFocusedFilePanel().GetFocusedItem()
			return m.fileModel.FilePreview.GetLocation() == focusedItem.Location
		}, DefaultTestTimeout, DefaultTestTick)
	})

	t.Run("Preview syncs after panel deletion", func(t *testing.T) {
		// Save and restore config
		originalPreviewWidth := common.Config.FilePreviewWidth
		common.Config.FilePreviewWidth = 10
		t.Cleanup(func() {
			common.Config.FilePreviewWidth = originalPreviewWidth
		})

		m := defaultTestModel(testDir)
		// Create 3 panels
		_, _ = m.fileModel.CreateNewFilePanel(testDir)
		_, _ = m.fileModel.CreateNewFilePanel(testDir)
		p := NewTestTeaProgWithEventLoop(t, m)

		p.Send(tea.WindowSizeMsg{Width: 150, Height: 40})

		// Focus middle panel and delete it
		m.fileModel.FocusedPanelIndex = 1
		p.SendKey(common.Hotkeys.CloseFilePanel[0])

		// Verify preview syncs with remaining focused panel
		assert.Eventually(t, func() bool {
			if len(m.fileModel.FilePanels) != 2 {
				return false
			}
			focusedItem := m.getFocusedFilePanel().GetFocusedItem()
			return m.fileModel.FilePreview.GetLocation() == focusedItem.Location
		}, DefaultTestTimeout, DefaultTestTick)
	})
}

// TestAsyncPreviewContent validates preview panel content changes
func TestAsyncPreviewContent(t *testing.T) {
	utils.SetRootLoggerToStdout(testing.Verbose())
	testDir := t.TempDir()

	// Create test files with different content
	file1 := filepath.Join(testDir, "file1.txt")
	file2 := filepath.Join(testDir, "file2.txt")
	utils.SetupFilesWithData(t, []byte("Content of file 1"), file1)
	utils.SetupFilesWithData(t, []byte("Content of file 2"), file2)

	originalPreviewWidth := common.Config.FilePreviewWidth
	common.Config.FilePreviewWidth = 10
	t.Cleanup(func() {
		common.Config.FilePreviewWidth = originalPreviewWidth
	})

	m := defaultTestModel(testDir)
	p := NewTestTeaProgWithEventLoop(t, m)

	p.Send(tea.WindowSizeMsg{Width: 100, Height: 40})

	// Navigate to file2
	p.SendKey(common.Hotkeys.ListDown[0])

	// Verify preview content changes to show file2
	assert.Eventually(t, func() bool {
		focusedItem := m.getFocusedFilePanel().GetFocusedItem()
		if focusedItem.Name != "file2.txt" {
			return false
		}
		return m.fileModel.FilePreview.GetLocation() == focusedItem.Location
	}, DefaultTestTimeout, DefaultTestTick)
}

// TestAsyncPreviewWidthChange validates preview panel width adjustments
func TestAsyncPreviewWidthChange(t *testing.T) {
	utils.SetRootLoggerToStdout(testing.Verbose())
	testDir := t.TempDir()

	utils.SetupFiles(t, filepath.Join(testDir, "file1.txt"))

	originalPreviewWidth := common.Config.FilePreviewWidth
	common.Config.FilePreviewWidth = 10
	t.Cleanup(func() {
		common.Config.FilePreviewWidth = originalPreviewWidth
	})

	m := defaultTestModel(testDir)
	p := NewTestTeaProgWithEventLoop(t, m)

	// Initial window size
	initialWidth := 100
	p.Send(tea.WindowSizeMsg{Width: initialWidth, Height: 40})

	// Wait for initial setup
	assert.Eventually(t, func() bool {
		return m.fileModel.FilePreview.GetContentWidth() > 0
	}, DefaultTestTimeout, DefaultTestTick)

	initialPreviewWidth := m.fileModel.FilePreview.GetContentWidth()

	// Resize window
	newWidth := 150
	p.Send(tea.WindowSizeMsg{Width: newWidth, Height: 40})

	// Verify preview width adjusts
	assert.Eventually(t, func() bool {
		currentWidth := m.fileModel.FilePreview.GetContentWidth()
		return currentWidth != initialPreviewWidth && currentWidth > 0
	}, DefaultTestTimeout, DefaultTestTick)
}
