package filemodel

import (
	"errors"
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

	return m.UpdateChildComponentWidth(), nil
}

func (m *Model) CloseFilePanel() (tea.Cmd, error) {
	if m.PanelCount() <= 1 {
		return nil, errors.New("CloseFilePanel called on with panelCount <= 1")
	}

	m.FilePanels = append(m.FilePanels[:m.FocusedPanelIndex],
		m.FilePanels[m.FocusedPanelIndex+1:]...)

	if m.FocusedPanelIndex != 0 {
		m.FocusedPanelIndex--
	}
	m.FilePanels[m.FocusedPanelIndex].IsFocused = true

	return m.UpdateChildComponentWidth(), nil
}

func (m *Model) UpdatePreviewPanel(msg preview.UpdateMsg) {
	selectedItem := m.GetFocusedFilePanel().GetSelectedItemPtr()
	if selectedItem == nil {
		slog.Debug("Panel empty or cursor invalid. Ignoring FilePreviewUpdateMsg")
		return
	}
	if selectedItem.Location != msg.GetLocation() {
		slog.Debug("FilePreviewUpdateMsg for older files. Ignoring")
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
		m.FilePreview.SetEmpty()
		return nil
	}
	selectedItem := panel.GetSelectedItem()
	if m.FilePreview.GetLocation() == selectedItem.Location && !forcePreviewRender {
		return nil
	}

	m.FilePreview.SetLocation(selectedItem.Location)
	m.FilePreview.SetLoading()
	reqCnt := m.ioReqCnt
	m.ioReqCnt++
	slog.Debug("Submitting file preview render request", "id", reqCnt, "path", selectedItem.Location)

	// HACK!!!. fileModel must not be aware of other dimensions. but...
	// Unfortunately, previewPanel isn't completely 'under' fileModel
	fullModalWidth := m.Width + common.Config.SidebarWidth
	if common.Config.SidebarWidth != 0 {
		fullModalWidth += common.BorderPadding
	}

	return func() tea.Msg {
		content := m.FilePreview.RenderWithPath(selectedItem.Location, m.ExpectedPreviewWidth, m.Height, fullModalWidth)
		return preview.NewUpdateMsg(selectedItem.Location, content,
			m.ExpectedPreviewWidth, m.Height, reqCnt)
	}
}
