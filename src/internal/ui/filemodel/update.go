package filemodel

import (
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/preview"
)

func (m *Model) CreateNewFilePanel(location string) (tea.Cmd, error) {
	if m.PanelCount() >= m.MaxFilePanel {
		return nil, ErrMaximumPanelCount
	}

	if _, err := os.Stat(location); err != nil {
		return nil, fmt.Errorf("cannot access location : %s", location)
	}

	m.FilePanels = append(m.FilePanels, filepanel.New(
		location, m.GetFocusedFilePanel().SortOptions, false, ""))

	newPanelIndex := m.PanelCount() - 1

	m.FilePanels[m.FocusedPanelIndex].IsFocused = false
	m.FilePanels[newPanelIndex].IsFocused = true
	m.FilePanels[newPanelIndex].SetHeight(m.Height)
	m.FocusedPanelIndex = newPanelIndex

	m.updateChildComponentWidth()
	return m.ensurePreviewDimensionsSync(), nil
}

func (m *Model) CloseFilePanel() (tea.Cmd, error) {
	if m.PanelCount() <= 1 {
		return nil, ErrMinimumPanelCount
	}

	m.FilePanels = append(m.FilePanels[:m.FocusedPanelIndex],
		m.FilePanels[m.FocusedPanelIndex+1:]...)

	if m.FocusedPanelIndex != 0 {
		m.FocusedPanelIndex--
	}
	m.FilePanels[m.FocusedPanelIndex].IsFocused = true
	m.updateChildComponentWidth()

	return m.ensurePreviewDimensionsSync(), nil
}

func (m *Model) ToggleFilePreviewPanel() tea.Cmd {
	m.FilePreview.ToggleOpen()
	m.updateChildComponentWidth()
	return m.ensurePreviewDimensionsSync()
}

func (m *Model) UpdatePreviewPanel(msg preview.UpdateMsg) {
	selectedItem := m.GetFocusedFilePanel().GetFocusedItemPtr()
	if selectedItem == nil {
		slog.Debug("Panel empty or cursor invalid. Ignoring FilePreviewUpdateMsg")
		return
	}
	if selectedItem.Location != msg.GetLocation() {
		slog.Debug("FilePreviewUpdateMsg for older files. Ignoring",
			"curLocation", selectedItem.Location, "msgLocation", msg.GetLocation())
		return
	}

	if m.ExpectedPreviewWidth != msg.GetContentWidth() ||
		m.Height != msg.GetContentHeight() {
		slog.Debug("FilePreviewUpdateMsg for older dimensions. Ignoring",
			"curW", m.ExpectedPreviewWidth, "curH", m.Height,
			"msgW", msg.GetContentWidth(), "msgH", msg.GetContentHeight())
		return
	}
	m.FilePreview.Apply(msg)
}

func (m *Model) GetFilePreviewCmd(forcePreviewRender bool) tea.Cmd {
	if !m.FilePreview.IsOpen() {
		return nil
	}
	panel := m.GetFocusedFilePanel()
	if panel.EmptyOrInvalid() {
		// Sync call because this will be fast
		m.FilePreview.SetEmptyWithDimensions(m.ExpectedPreviewWidth, m.Height)
		return nil
	}
	selectedItem := panel.GetFocusedItem()
	if m.FilePreview.GetLocation() == selectedItem.Location && !forcePreviewRender {
		return nil
	}

	m.FilePreview.SetLocation(selectedItem.Location)
	m.FilePreview.SetLoading()

	// HACK!!!. fileModel must not be aware of other dimensions. but...
	// Unfortunately, previewPanel isn't completely 'under' fileModel
	// Note: Must save the dimensions for the closure of the Cmd to avoid
	// problems
	fullModalWidth := m.Width + common.Config.SidebarWidth
	if common.Config.SidebarWidth != 0 {
		fullModalWidth += common.BorderPadding
	}
	width := m.ExpectedPreviewWidth
	height := m.Height

	reqCnt := m.ioReqCnt
	m.ioReqCnt++
	slog.Debug("Submitting file preview render request", "id", reqCnt,
		"path", selectedItem.Location, "w", width, "h", height)

	return func() tea.Msg {
		content := m.FilePreview.RenderWithPath(selectedItem.Location, width, height, fullModalWidth)
		return preview.NewUpdateMsg(selectedItem.Location, content,
			width, height, reqCnt)
	}
}

func (m *Model) ToggleDotFile() {
	m.DisplayDotFiles = !m.DisplayDotFiles
	m.UpdateFilePanelsIfNeeded(true)
}

func (m *Model) UpdateFilePanelsIfNeeded(force bool) {
	for i := range m.FilePanels {
		m.FilePanels[i].UpdateElementsIfNeeded(force, m.DisplayDotFiles)
	}
}
