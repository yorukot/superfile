package filemodel

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

func (m *FileModel) SetHeight(height int) {
	if height < FileModelMinHeight {
		height = FileModelMinHeight
	}
	m.Height = height
	m.UpdateChildComponentHeight()
}

func (m *FileModel) SetWidth(width int) {
	if width < FileModelMinWidth {
		width = FileModelMinWidth
	}
	m.Width = width
	m.UpdateChildComponentWidth()
}

func (m *FileModel) PanelCount() int {
	return len(m.FilePanels)
}

func (m *FileModel) UpdateChildComponentHeight() {
	for i := range m.FilePanels {
		m.FilePanels[i].SetHeight(m.Height)
	}
	m.FilePreview.SetHeight(m.Height)
}

func (m *FileModel) UpdateChildComponentWidth() {
	// TODO: programatically ensure that this becomes impossible
	if m.PanelCount() == 0 {
		slog.Error("Unexpected error: fileModel with 0 panels")
		return
	}
	panelCount := len(m.FilePanels)
	widthForPanels := m.Width

	if m.FilePreview.IsOpen() {
		// Need to give some width to preview
		var previewWidth int
		if common.Config.FilePreviewWidth == 0 {
			// FileModel will be split among `panelCount+1`
			previewWidth = m.Width / (panelCount + 1)
		} else {
			previewWidth = m.Width / common.Config.FilePreviewWidth
		}
		m.FilePreview.SetWidth(previewWidth)
		widthForPanels -= previewWidth
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
}
