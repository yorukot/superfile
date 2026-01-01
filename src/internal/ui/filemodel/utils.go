package filemodel

import "github.com/yorukot/superfile/src/internal/ui/filepanel"

func (m *Model) GetFocusedFilePanel() *filepanel.Model {
	return &m.FilePanels[m.FocusedPanelIndex]
}
