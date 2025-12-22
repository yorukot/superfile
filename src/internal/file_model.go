package internal

import (
	"errors"
	"fmt"
	"os"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

// Set file preview panel Widht to width. Assure that
func (m *model) getFilePreviewWidth() int {
	if common.Config.FilePreviewWidth == 0 {
		return (m.fullWidth - common.Config.SidebarWidth -
			(common.InnerPadding + (len(m.fileModel.filePanels))*common.BorderPadding)) / (len(m.fileModel.filePanels) + 1)
	}
	return (m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth
}

// Proper set panels size. Assure that panels do not overlap
func (m *model) setFileModelDimensions() {
	width := m.fullWidth
	if m.fileModel.filePreview.IsOpen() {
		m.fileModel.filePreview.SetWidth(m.getFilePreviewWidth())
		m.fileModel.filePreview.SetHeight(m.mainPanelHeight + common.BorderPadding)
	}
	// set each file panel size and max file panel amount
	m.fileModel.width = (width - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth() -
		(common.InnerPadding + (len(m.fileModel.filePanels)-1)*common.BorderPadding)) / len(m.fileModel.filePanels)
	m.fileModel.maxFilePanel = (width - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth()) / common.FilePanelWidthUnit
	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].UpdateDimensions(
			m.fileModel.width+common.BorderPadding,
			m.mainPanelHeight+common.BorderPadding,
		)
	}
	if m.fileModel.maxFilePanel >= common.FilePanelMax {
		m.fileModel.maxFilePanel = common.FilePanelMax
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

	m.fileModel.filePanels = append(m.fileModel.filePanels, filepanel.New(
		location, m.getFocusedFilePanel().SortOptions, false, ""))

	if m.fileModel.filePreview.IsOpen() {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth -
				(common.InnerPadding + (len(m.fileModel.filePanels))*common.BorderPadding)) / (len(m.fileModel.filePanels) + 1))
		} else {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	m.fileModel.filePanels[m.filePanelFocusIndex].IsFocused = false
	m.fileModel.filePanels[m.filePanelFocusIndex+1].IsFocused = returnFocusType(m.focusPanel)
	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth() -
		(common.InnerPadding + (len(m.fileModel.filePanels)-1)*common.BorderPadding)) / len(m.fileModel.filePanels)
	m.filePanelFocusIndex++

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth()) / common.FilePanelWidthUnit

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].SearchBar.Width = m.fileModel.width - common.InnerPadding
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
				(common.InnerPadding + (len(m.fileModel.filePanels))*common.BorderPadding)) / (len(m.fileModel.filePanels) + 1))
		} else {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	if m.filePanelFocusIndex != 0 {
		m.filePanelFocusIndex--
	}

	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth() -
		(common.InnerPadding + (len(m.fileModel.filePanels)-1)*common.BorderPadding)) / len(m.fileModel.filePanels)
	m.fileModel.filePanels[m.filePanelFocusIndex].IsFocused = returnFocusType(m.focusPanel)

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth()) / common.FilePanelWidthUnit

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].SearchBar.Width = m.fileModel.width - common.InnerPadding
	}
}

func (m *model) toggleFilePreviewPanel() {
	m.fileModel.filePreview.ToggleOpen()
	m.fileModel.filePreview.SetWidth(0)
	m.fileModel.filePreview.SetHeight(m.mainPanelHeight + common.BorderPadding)
	if m.fileModel.filePreview.IsOpen() {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth -
				(common.InnerPadding + (len(m.fileModel.filePanels))*common.BorderPadding)) / (len(m.fileModel.filePanels) + 1))
		} else {
			m.fileModel.filePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	m.fileModel.width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth() -
		(common.InnerPadding + (len(m.fileModel.filePanels)-1)*common.BorderPadding)) / len(m.fileModel.filePanels)

	m.fileModel.maxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth()) / common.FilePanelWidthUnit

	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].SearchBar.Width = m.fileModel.width - common.InnerPadding
	}
}

// TODO : Replace all usage of "m.fileModel.filePanels[m.filePanelFocusIndex]" with this
// There are many usage
func (m *model) getFocusedFilePanel() *filepanel.Model {
	return &m.fileModel.filePanels[m.filePanelFocusIndex]
}
