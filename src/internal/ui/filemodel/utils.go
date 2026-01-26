package filemodel

import (
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/preview"
)

func (m *Model) GetFocusedFilePanel() *filepanel.Model {
	return &m.FilePanels[m.FocusedPanelIndex]
}

func New(firstPanelPaths []string, toggleDotFile bool) Model {
	return Model{
		FilePanels:       filepanel.FilePanelSlice(firstPanelPaths),
		FilePreview:      preview.New(),
		SinglePanelWidth: common.DefaultFilePanelWidth,
		DisplayDotFiles:  toggleDotFile,
	}
}
