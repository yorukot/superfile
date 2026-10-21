package filemodel

import (
	"log/slog"

	tea "charm.land/bubbletea/v2"

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

	m.panelOriginsX = make([]int, 0, panelCount)
	originX := 0
	for i := range panelCount {
		m.panelOriginsX = append(m.panelOriginsX, originX)

		panelWidthForPanel := panelWidth
		if i == panelCount-1 {
			panelWidthForPanel = lastPanelWidth
		}
		m.FilePanels[i].SetWidth(panelWidthForPanel)
		originX += panelWidthForPanel
	}

	m.SinglePanelWidth = panelWidth
	m.MaxFilePanel = m.maxPanelCount()
}

// maxPanelCount returns how many file panels fit in the current width while
// keeping every panel at least filepanel.MinWidth wide.
//
// The preview width depends on the panel count when Config.FilePreviewWidth is
// 0, so the space left over for the current panel count can't be reused for the
// limit - that reserves too much for the preview and blocks panel creation that
// would actually render fine. Each candidate count is checked against the
// preview width it would really get.
func (m *Model) maxPanelCount() int {
	maxPanels := 1
	for n := 1; n <= common.FilePanelMax; n++ {
		widthForPanels := m.Width
		if m.FilePreview.IsOpen() {
			widthForPanels -= m.previewWidthForPanelCount(n)
		}
		if widthForPanels/n < filepanel.MinWidth {
			break
		}
		maxPanels = n
	}
	return maxPanels
}

func (m *Model) previewWidthForPanelCount(panelCount int) int {
	if common.Config.FilePreviewWidth == 0 {
		return m.Width / (panelCount + 1)
	}
	return m.Width / common.Config.FilePreviewWidth
}

func (m *Model) ensurePreviewDimensionsSync() tea.Cmd {
	if m.FilePreview.GetContentWidth() != m.ExpectedPreviewWidth ||
		m.FilePreview.GetContentHeight() != m.Height {
		return m.GetFilePreviewCmd(true)
	}
	return nil
}
