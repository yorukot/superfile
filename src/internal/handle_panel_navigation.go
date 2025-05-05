package internal

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/yorukot/superfile/src/internal/ui/sidebar"

	"github.com/yorukot/superfile/src/internal/common"

	variable "github.com/yorukot/superfile/src/config"
)

// Pinned directory
func (m *model) pinnedDirectory() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	err := sidebar.TogglePinnedDirectory(panel.location)
	if err != nil {
		slog.Error("Error while toggling pinned directory", "error", err)
	}
}

// Create new file panel
func (m *model) createNewFilePanel(location string) error {
	if len(m.fileModel.filePanels) == m.fileModel.maxFilePanel {
		return errors.New("maximum panel count reached")
	}

	if location == "" {
		location = variable.HomeDir
	}

	if _, err := os.Stat(location); err != nil {
		return fmt.Errorf("cannot access location : %s", location)
	}

	m.fileModel.filePanels = append(m.fileModel.filePanels, filePanel{

		location:         location,
		sortOptions:      m.fileModel.filePanels[m.filePanelFocusIndex].sortOptions,
		panelMode:        browserMode,
		focusType:        secondFocus,
		directoryRecords: make(map[string]directoryRecord),
		searchBar:        common.GenerateSearchBar(),
	})

	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - common.Config.SidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {
			m.fileModel.filePreview.width = (m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	m.fileModel.filePanels[m.filePanelFocusIndex+1].focusType = returnFocusType(m.focusPanel)
	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.filePanelFocusIndex++

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
	return nil
}

// Close current focus file panel
func (m *model) closeFilePanel() {
	if len(m.fileModel.filePanels) == 1 {
		return
	}

	m.fileModel.filePanels = append(m.fileModel.filePanels[:m.filePanelFocusIndex], m.fileModel.filePanels[m.filePanelFocusIndex+1:]...)

	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - common.Config.SidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {
			m.fileModel.filePreview.width = (m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth
		}
	}

	if m.filePanelFocusIndex != 0 {
		m.filePanelFocusIndex--
	}

	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
}

func (m *model) toggleFilePreviewPanel() {
	m.fileModel.filePreview.open = !m.fileModel.filePreview.open
	m.fileModel.filePreview.width = 0
	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.width = (m.fullWidth - common.Config.SidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
		} else {
			m.fileModel.filePreview.width = (m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth
		}
	}

	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.width) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
}

// Focus on next file panel
func (m *model) nextFilePanel() {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == (len(m.fileModel.filePanels) - 1) {
		m.filePanelFocusIndex = 0
	} else {
		m.filePanelFocusIndex++
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
}

// Focus on previous file panel
func (m *model) previousFilePanel() {
	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
	if m.filePanelFocusIndex == 0 {
		m.filePanelFocusIndex = (len(m.fileModel.filePanels) - 1)
	} else {
		m.filePanelFocusIndex--
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.focusPanel)
}

// Focus on sidebar
func (m *model) focusOnSideBar() {
	if common.Config.SidebarWidth == 0 {
		return
	}
	if m.focusPanel == sidebarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = sidebarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
}

// Focus on processbar
func (m *model) focusOnProcessBar() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == processBarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = processBarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
}

// focus on metadata
func (m *model) focusOnMetadata() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == metadataFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
	} else {
		m.focusPanel = metadataFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
	}
}
