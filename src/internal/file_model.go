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
			(common.InnerPadding + (len(m.fileModel.FilePanels))*common.BorderPadding)) / (len(m.fileModel.FilePanels) + 1)
	}
	return (m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth
}

// Proper set panels size. Assure that panels do not overlap
func (m *model) setFileModelDimensions() {

	width := m.fullWidth
	if m.fileModel.FilePreview.IsOpen() {
		m.fileModel.FilePreview.SetWidth(m.getFilePreviewWidth())
		m.fileModel.FilePreview.SetHeight(m.mainPanelHeight + common.BorderPadding)
	}
	// set each file panel size and max file panel amount
	m.fileModel.Width = (width - common.Config.SidebarWidth - m.fileModel.FilePreview.GetWidth() -
		(common.InnerPadding + (len(m.fileModel.FilePanels)-1)*common.BorderPadding)) / len(m.fileModel.FilePanels)
	m.fileModel.MaxFilePanel = (width - common.Config.SidebarWidth - m.fileModel.FilePreview.GetWidth()) / common.FilePanelWidthUnit
	for i := range m.fileModel.FilePanels {
		m.fileModel.FilePanels[i].UpdateDimensions(
			m.fileModel.Width+common.BorderPadding,
			m.mainPanelHeight+common.BorderPadding,
		)
	}
	if m.fileModel.MaxFilePanel >= common.FilePanelMax {
		m.fileModel.MaxFilePanel = common.FilePanelMax
	}
}

// Create new file panel
func (m *model) createNewFilePanel(location string) error {
	// In case we have model width and height zero, maxFilePanel would be 0
	// But we would have len() here as 1. Hence there would be discrepency here.
	// Although this is not possible in actual usage, and can be only reproduced in tests.
	if len(m.fileModel.FilePanels) == m.fileModel.MaxFilePanel {
		// TODO : Define as a predefined error in errors.go
		return errors.New("maximum panel count reached")
	}

	if location == "" {
		location = variable.HomeDir
	}

	if _, err := os.Stat(location); err != nil {
		return fmt.Errorf("cannot access location : %s", location)
	}

	m.fileModel.FilePanels = append(m.fileModel.FilePanels, filepanel.New(
		location, m.getFocusedFilePanel().SortOptions, false, ""))

	if m.fileModel.FilePreview.IsOpen() {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.FilePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth -
				(common.InnerPadding + (len(m.fileModel.FilePanels))*common.BorderPadding)) / (len(m.fileModel.FilePanels) + 1))
		} else {
			m.fileModel.FilePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = false
	m.fileModel.FilePanels[m.filePanelFocusIndex+1].IsFocused = returnFocusType(m.focusPanel)
	m.fileModel.Width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.FilePreview.GetWidth() -
		(common.InnerPadding + (len(m.fileModel.FilePanels)-1)*common.BorderPadding)) / len(m.fileModel.FilePanels)
	m.filePanelFocusIndex++

	m.fileModel.MaxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.FilePreview.GetWidth()) / common.FilePanelWidthUnit

	for i := range m.fileModel.FilePanels {
		m.fileModel.FilePanels[i].SearchBar.Width = m.fileModel.Width - common.InnerPadding
	}
	return nil
}

// Close current focus file panel
func (m *model) closeFilePanel() {
	if len(m.fileModel.FilePanels) == 1 {
		return
	}

	m.fileModel.FilePanels = append(m.fileModel.FilePanels[:m.filePanelFocusIndex],
		m.fileModel.FilePanels[m.filePanelFocusIndex+1:]...)

	if m.fileModel.FilePreview.IsOpen() {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.FilePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth -
				(common.InnerPadding + (len(m.fileModel.FilePanels))*common.BorderPadding)) / (len(m.fileModel.FilePanels) + 1))
		} else {
			m.fileModel.FilePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	if m.filePanelFocusIndex != 0 {
		m.filePanelFocusIndex--
	}

	m.fileModel.Width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.FilePreview.GetWidth() -
		(common.InnerPadding + (len(m.fileModel.FilePanels)-1)*common.BorderPadding)) / len(m.fileModel.FilePanels)
	m.fileModel.FilePanels[m.filePanelFocusIndex].IsFocused = returnFocusType(m.focusPanel)

	m.fileModel.MaxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.FilePreview.GetWidth()) / common.FilePanelWidthUnit

	for i := range m.fileModel.FilePanels {
		m.fileModel.FilePanels[i].SearchBar.Width = m.fileModel.Width - common.InnerPadding
	}
}

func (m *model) toggleFilePreviewPanel() {
	m.fileModel.FilePreview.ToggleOpen()
	m.fileModel.FilePreview.SetWidth(0)
	m.fileModel.FilePreview.SetHeight(m.mainPanelHeight + common.BorderPadding)
	if m.fileModel.FilePreview.IsOpen() {
		// File preview panel width same as file panel
		if common.Config.FilePreviewWidth == 0 {
			m.fileModel.FilePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth -
				(common.InnerPadding + (len(m.fileModel.FilePanels))*common.BorderPadding)) / (len(m.fileModel.FilePanels) + 1))
		} else {
			m.fileModel.FilePreview.SetWidth((m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth)
		}
	}

	m.fileModel.Width = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.FilePreview.GetWidth() -
		(common.InnerPadding + (len(m.fileModel.FilePanels)-1)*common.BorderPadding)) / len(m.fileModel.FilePanels)

	m.fileModel.MaxFilePanel = (m.fullWidth - common.Config.SidebarWidth - m.fileModel.FilePreview.GetWidth()) / common.FilePanelWidthUnit

	for i := range m.fileModel.FilePanels {
		m.fileModel.FilePanels[i].SearchBar.Width = m.fileModel.Width - common.InnerPadding
	}
}

// TODO : Replace all usage of "m.fileModel.filePanels[m.filePanelFocusIndex]" with this
// There are many usage
func (m *model) getFocusedFilePanel() *filepanel.Model {
	return &m.fileModel.FilePanels[m.filePanelFocusIndex]
}
