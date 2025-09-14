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
	err := sidebar.TogglePinnedDirectory(panel.Location)
	if err != nil {
		slog.Error("Error while toggling pinned directory", "error", err)
	}
}

// Create new file panel
func (m *model) createNewFilePanel(location string) error {
	// In case we have model width and height zero, maxFilePanel would be 0
	// But we would have len() here as 1. Hence there would be discrepency here.
	// Although this is not possible in actual usage, and can be only reproduced in tests.
	if len(m.fileModel.filePanels) == m.fileModel.maxFilePanel {
		// TODO : Define as a predefined error in errors.go
		return errors.New("maximum panel count reached")
	}

	if location == "" {
		location = variable.HomeDir
	}

	if _, err := os.Stat(location); err != nil {
		return fmt.Errorf("cannot access location : %s", location)
	}

	m.fileModel.filePanels = append(m.fileModel.filePanels, FilePanel{
		Location:         location,
		SortOptions:      m.fileModel.filePanels[m.filePanelFocusIndex].SortOptions,
		PanelMode:        BrowserMode,
		isFocused:        false,
		DirectoryRecords: make(map[string]DirectoryRecord),
		SearchBar:        common.GenerateSearchBar(),
	})

	if m.fileModel.filePreview.IsOpen() {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth -
				(4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1))
		} else {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = false
	m.fileModel.filePanels[m.filePanelFocusIndex+1].isFocused = returnFocusType(m.focusPanel)
	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth() -
		(4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.filePanelFocusIndex++

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth()) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].SetSearchBarWidth(m.fileModel.width)
	}
	return nil
}

// Close current focus file panel
func (m *model) closeFilePanel() {
	if len(m.fileModel.filePanels) == 1 {
		return
	}

	m.fileModel.filePanels = append(m.fileModel.filePanels[:m.filePanelFocusIndex],
		m.fileModel.filePanels[m.filePanelFocusIndex+1:]...)

	if m.fileModel.filePreview.IsOpen() {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth -
				(4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1))
		} else {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	if m.filePanelFocusIndex != 0 {
		m.filePanelFocusIndex--
	}

	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth() -
		(4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = returnFocusType(m.focusPanel)

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth()) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].SetSearchBarWidth(m.fileModel.width)
	}
}

func (m *model) toggleFilePreviewPanel() {
	m.fileModel.filePreview.ToggleOpen()
	m.fileModel.filePreview.SetWidth(0)
	m.fileModel.filePreview.SetHeight(m.mainPanelHeight + 2)
	if m.fileModel.filePreview.IsOpen() {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth -
				(4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1))
		} else {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth() -
		(4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth()) / 20

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].SetSearchBarWidth(m.fileModel.width)
	}
}

// Focus on next file panel
func (m *model) nextFilePanel() {
	m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = false
	if m.filePanelFocusIndex == (len(m.fileModel.filePanels) - 1) {
		m.filePanelFocusIndex = 0
	} else {
		m.filePanelFocusIndex++
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = returnFocusType(m.focusPanel)
}

// Focus on previous file panel
func (m *model) previousFilePanel() {
	m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = false
	if m.filePanelFocusIndex == 0 {
		m.filePanelFocusIndex = (len(m.fileModel.filePanels) - 1)
	} else {
		m.filePanelFocusIndex--
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = returnFocusType(m.focusPanel)
}

// Focus on sidebar
func (m *model) focusOnSideBar() {
	if common.Config.SidebarWidth == 0 {
		return
	}
	if m.focusPanel == sidebarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = true
	} else {
		m.focusPanel = sidebarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = false
	}
}

// Focus on processbar
func (m *model) focusOnProcessBar() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == processBarFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = true
	} else {
		m.focusPanel = processBarFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = false
	}
}

// focus on metadata
func (m *model) focusOnMetadata() {
	if !m.toggleFooter {
		return
	}

	if m.focusPanel == metadataFocus {
		m.focusPanel = nonePanelFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = true
	} else {
		m.focusPanel = metadataFocus
		m.fileModel.filePanels[m.filePanelFocusIndex].isFocused = false
	}
}
