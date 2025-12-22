package internal

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/common"
)

// Pinned directory
func (m *model) pinnedDirectory() {
	panel := &m.fileModel.FilePanels[m.filePanelFocusIndex]
	err := m.sidebarModel.TogglePinnedDirectory(panel.Location)
	if err != nil {
		slog.Error("Error while toggling pinned directory", "error", err)
	}
}

// Focus on next file panel
func (m *model) nextFilePanel() {
	m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = false
	if m.filePanelFocusIndex == (len(m.fileModel.FilePanels) - 1) {
		m.filePanelFocusIndex = 0
	} else {
		m.filePanelFocusIndex++
	}

	m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = returnFocusType(m.focusPanel)
}

// Focus on previous file panel
func (m *model) previousFilePanel() {
	m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = false
	if m.filePanelFocusIndex == 0 {
		m.filePanelFocusIndex = (len(m.fileModel.FilePanels) - 1)
	} else {
		m.filePanelFocusIndex--
	}

	m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = returnFocusType(m.focusPanel)
}

// Focus on sidebar
func (m *model) focusOnSideBar() {
	if common.Config.SidebarWidth == 0 {
		return
	}
	if m.focusPanel == sidebarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = true
	} else {
		m.focusPanel = sidebarFocus
		m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = false
	}
}

// Focus on processbar
func (m *model) focusOnProcessBar() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == processBarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = true
	} else {
		m.focusPanel = processBarFocus
		m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = false
	}
}

// focus on metadata
func (m *model) focusOnMetadata() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == metadataFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = true
	} else {
		m.focusPanel = metadataFocus
		m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = false
	}
}
