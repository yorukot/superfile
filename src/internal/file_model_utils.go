package internal

import "github.com/yorukot/superfile/src/internal/ui/filepanel"

func (m *FileModel) GetFocusedFilePanel() *filepanel.Model {
	return &m.FilePanels[m.FocusedPanelIndex]
}
