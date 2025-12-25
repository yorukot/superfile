package filemodel

import "github.com/charmbracelet/lipgloss"

func (m *Model) Render() string {
	f := make([]string, m.PanelCount()+1)
	for i, filePanel := range m.FilePanels {
		f[i] = filePanel.Render(filePanel.IsFocused)
	}
	f[m.PanelCount()] = m.GetFilePreviewRender()
	return lipgloss.JoinHorizontal(lipgloss.Top, f...)
}

func (m *Model) GetFilePreviewRender() string {
	if !m.FilePreview.IsOpen() {
		return ""
	}
	// Check if width and height have been synced yet
	if m.FilePreview.GetHeight() == m.Height &&
		m.FilePreview.GetWidth() == m.ExpectedPreviewWidth {
		if m.FilePreview.IsLoading() {
			return m.FilePreview.RenderText(FilePreviewLoadingText)
		}
		if m.FilePreview.IsEmpty() {
			return m.FilePreview.RenderText("")
		}
		return m.FilePreview.GetContent()
	}

	// Placeholder resizing text till they get synced
	return m.FilePreview.RenderTextWithDimension(
		FilePreviewResizingText, m.Height, m.ExpectedPreviewWidth)
}
