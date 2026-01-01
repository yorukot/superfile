package internal

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

func TestAsyncPreviewPanelSync(t *testing.T) {
	origPreviewWidth := common.Config.FilePreviewWidth
	defer func() {
		common.Config.FilePreviewWidth = origPreviewWidth
	}()
	common.Config.FilePreviewWidth = 10

	testDir := t.TempDir()
	utils.SetupFiles(t,
		filepath.Join(testDir, "file1.txt"),
		filepath.Join(testDir, "file2.txt"),
		filepath.Join(testDir, "file3.txt"),
	)

	t.Run("Preview syncs after panel creation", func(t *testing.T) {
		m := defaultTestModel(testDir)
		p := NewTestTeaProgWithEventLoop(t, m)

		p.SendKey(common.Hotkeys.CreateNewFilePanel[0])

		assert.Eventually(t, func() bool {
			if len(m.fileModel.FilePanels) != 2 {
				return false
			}
			focusedItem := m.getFocusedFilePanel().GetFocusedItem()
			return m.fileModel.FilePreview.GetLocation() == focusedItem.Location
		}, DefaultTestTimeout, DefaultTestTick, "Preview should sync with new panel's focused item")
	})

	t.Run("Preview syncs after panel deletion", func(t *testing.T) {
		m := defaultTestModel(testDir)
		_, _ = m.fileModel.CreateNewFilePanel(testDir)
		_, _ = m.fileModel.CreateNewFilePanel(testDir)
		p := NewTestTeaProgWithEventLoop(t, m)

		m.fileModel.FocusedPanelIndex = 1
		p.SendKey(common.Hotkeys.CloseFilePanel[0])

		assert.Eventually(t, func() bool {
			if len(m.fileModel.FilePanels) != 2 {
				return false
			}
			focusedItem := m.getFocusedFilePanel().GetFocusedItem()
			return m.fileModel.FilePreview.GetLocation() == focusedItem.Location
		}, DefaultTestTimeout, DefaultTestTick, "Preview should sync after panel deletion")
	})
}

func TestAsyncPreviewContent(t *testing.T) {
	origPreviewWidth := common.Config.FilePreviewWidth
	defer func() {
		common.Config.FilePreviewWidth = origPreviewWidth
	}()
	common.Config.FilePreviewWidth = 10

	testDir := t.TempDir()
	file1 := filepath.Join(testDir, "file1.txt")
	file2 := filepath.Join(testDir, "file2.txt")
	utils.SetupFilesWithData(t, []byte("Content of file 1"), file1)
	utils.SetupFilesWithData(t, []byte("Content of file 2"), file2)

	m := defaultTestModel(testDir)
	p := NewTestTeaProgWithEventLoop(t, m)

	p.SendKey(common.Hotkeys.ListDown[0])

	assert.Eventually(t, func() bool {
		focusedItem := m.getFocusedFilePanel().GetFocusedItem()
		if focusedItem.Name != "file2.txt" {
			return false
		}
		return m.fileModel.FilePreview.GetLocation() == focusedItem.Location
	}, DefaultTestTimeout, DefaultTestTick, "Preview should change to file2 content")
}

func TestAsyncPreviewWidthChange(t *testing.T) {
	origPreviewWidth := common.Config.FilePreviewWidth
	defer func() {
		common.Config.FilePreviewWidth = origPreviewWidth
	}()
	common.Config.FilePreviewWidth = 10

	testDir := t.TempDir()
	utils.SetupFiles(t, filepath.Join(testDir, "file1.txt"))

	m := defaultTestModel(testDir)
	p := NewTestTeaProgWithEventLoop(t, m)

	assert.Eventually(t, func() bool {
		return m.fileModel.FilePreview.GetContentWidth() > 0
	}, DefaultTestTimeout, DefaultTestTick, "Initial preview width should be set")

	initialPreviewWidth := m.fileModel.FilePreview.GetContentWidth()

	p.Send(tea.WindowSizeMsg{Width: DefaultTestModelWidth + 50, Height: DefaultTestModelHeight})

	assert.Eventually(t, func() bool {
		currentWidth := m.fileModel.FilePreview.GetContentWidth()
		return currentWidth != initialPreviewWidth && currentWidth > 0
	}, DefaultTestTimeout, DefaultTestTick, "Preview width should adjust after window resize")
}
