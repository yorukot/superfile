package internal

import (
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

func (m *FileModel) GetFocusedFilePanel() *filepanel.Model {
	return &m.FilePanels[m.FocusedPanelIndex]
}

func (m *model) createNewFilePanel(location string) error {
	return m.fileModel.createNewFilePanel(location)
}

// Close current focus file panel
func (m *model) closeFilePanel() {
	m.fileModel.CloseFilePanel()
}

func (m *model) toggleFilePreviewPanel() {
	m.fileModel.FilePreview.ToggleOpen()
	m.fileModel.UpdateChildComponentWidth()
}

// TODO : Replace all usage of "m.fileModel.filePanels[m.filePanelFocusIndex]" with this
// There are many usage
func (m *model) getFocusedFilePanel() *filepanel.Model {
	return m.fileModel.GetFocusedFilePanel()
}
