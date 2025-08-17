package internal

import (
	"github.com/yorukot/superfile/src/internal/ui"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

// TODO : Move this to a separate package. Keeping all relevant stuff here for now.

type FilePreviewPanel struct {
	open           bool
	width          int
	height         int
	location       string
	content        string
	imagePreviewer *filepreview.ImagePreviewer
}

func (m *FilePreviewPanel) SetWidth(width int) {
	m.width = width
}

func (m *FilePreviewPanel) SetHeight(height int) {
	m.height = height
}

// Simple rendered string with given text
func (m *FilePreviewPanel) RenderText(text string) string {
	// TODO : This width adjustment must not be done inside render function. It should
	// only be triggered via Update()
	// TODO : Remove this code duplication with above method
	// by moving file preview to a separate package

	return ui.FilePreviewPanelRenderer(m.height, m.width).
		AddLines(text).
		Render() + m.imagePreviewer.ClearKittyImages()
}

func (m *FilePreviewPanel) SetContextWithRenderText(text string) {
	m.content = m.RenderText(text)
}
