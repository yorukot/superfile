package filemodel

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

func (m *Model) SetHeight(height int) tea.Cmd {
	if height < FileModelMinHeight {
		height = FileModelMinHeight
	}
	m.Height = height
	return m.UpdateChildComponentHeight()
}

func (m *Model) SetWidth(width int) tea.Cmd {
	if width < FileModelMinWidth {
		width = FileModelMinWidth
	}
	m.Width = width
	return m.UpdateChildComponentWidth()
}

func (m *Model) PanelCount() int {
	return len(m.FilePanels)
}

func (m *Model) UpdateChildComponentHeight() tea.Cmd {
	for i := range m.FilePanels {
		m.FilePanels[i].SetHeight(m.Height)
	}

	if m.FilePreview.GetHeight() == m.Height {
		return nil
	}

	return m.GetFilePreviewCmd(true)
}

func (m *Model) UpdateChildComponentWidth() tea.Cmd {
	// TODO: programatically ensure that this becomes impossible
	if m.PanelCount() == 0 {
		slog.Error("Unexpected error: fileModel with 0 panels")
		return nil
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

	// Whether needs to re-render preview?
	if m.FilePreview.GetWidth() != m.ExpectedPreviewWidth {
		return m.GetFilePreviewCmd(true)
	}
	return nil
}
