package filemodel

import (
	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
)

func (m *Model) PreviewViewportHeight() int {
	height := m.Height
	if common.Config.EnableFilePreviewBorder {
		height -= common.BorderPadding
	}
	return height
}

// RefreshPreviewScroll re-renders the current preview in-process. Unlike
// GetFilePreviewCmd, this avoids the loading placeholder and stale async updates
// that cause flicker during repeated scroll keypresses.
func (m *Model) RefreshPreviewScroll() tea.Cmd {
	if !m.FilePreview.IsOpen() {
		return nil
	}
	panel := m.GetFocusedFilePanel()
	if panel.EmptyOrInvalid() {
		return nil
	}
	selectedItem := panel.GetFocusedItem()

	fullModalWidth := m.Width + common.Config.SidebarWidth
	if common.Config.SidebarWidth != 0 {
		fullModalWidth += common.BorderPadding
	}
	width := m.ExpectedPreviewWidth
	height := m.Height

	content, rawTransmit := m.FilePreview.RenderWithPath(
		selectedItem.Location, width, height, fullModalWidth)
	m.FilePreview.UpdateRenderedContent(content, width, height, selectedItem.Location)

	if rawTransmit != "" {
		return tea.Raw(rawTransmit)
	}
	return nil
}
