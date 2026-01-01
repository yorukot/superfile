package filemodel

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

// Use SetDimensions if you want to update both
// it will prevent duplicate file preview commands and hence, is efficient
func (m *Model) SetDimensions(width int, height int) tea.Cmd {
	m.Height = max(height, FileModelMinHeight)
	m.Width = max(width, FileModelMinWidth)
	m.updateChildComponentWidth()
	m.updateChildComponentHeight()
	return m.ensurePreviewDimensionsSync()
}
func (m *Model) SetHeight(height int) tea.Cmd {
	m.Height = max(height, FileModelMinHeight)
	m.updateChildComponentHeight()
	return m.ensurePreviewDimensionsSync()
}

func (m *Model) SetWidth(width int) tea.Cmd {
	m.Width = max(width, FileModelMinWidth)
	m.updateChildComponentWidth()
	return m.ensurePreviewDimensionsSync()
}

func (m *Model) PanelCount() int {
	return len(m.FilePanels)
}

func (m *Model) updateChildComponentHeight() {
	for i := range m.FilePanels {
		m.FilePanels[i].SetHeight(m.Height)
	}
}

func (m *Model) updateChildComponentWidth() {
	// TODO: programatically ensure that this becomes impossible
	if m.PanelCount() == 0 {
		slog.Error("Unexpected error: fileModel with 0 panels")
		return
	}
	panelCount := len(m.FilePanels)
	widthForPanels := m.Width

	if m.FilePreview.IsOpen() {
		// Need to give some width to preview
		if common.Config.FilePreviewWidth == 0 {
			// FileModel will be split among `panelCount+1`
			m.ExpectedPreviewWidth = m.Width / (panelCount + 1)
		} else {
			m.ExpectedPreviewWidth = m.Width / common.Config.FilePreviewWidth
		}
		widthForPanels -= m.ExpectedPreviewWidth
	}

	panelWidth := widthForPanels / panelCount
	lastPanelWidth := widthForPanels - (panelCount-1)*panelWidth

	for i := range panelCount {
		if i == panelCount-1 {
			m.FilePanels[i].SetWidth(lastPanelWidth)
		} else {
			m.FilePanels[i].SetWidth(panelWidth)
		}
	}

	m.SinglePanelWidth = panelWidth
	m.MaxFilePanel = widthForPanels / filepanel.MinWidth
	// Cap at the system maximum
	if m.MaxFilePanel > common.FilePanelMax {
		m.MaxFilePanel = common.FilePanelMax
	}
}

func (m *Model) ensurePreviewDimensionsSync() tea.Cmd {
	if m.FilePreview.GetContentWidth() != m.ExpectedPreviewWidth ||
		m.FilePreview.GetContentHeight() != m.Height {
		return m.GetFilePreviewCmd(true)
	}
	return nil
}
