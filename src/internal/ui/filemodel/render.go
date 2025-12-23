package filemodel

import "github.com/charmbracelet/lipgloss"

func (m *Model) Render() string {
	f := make([]string, m.PanelCount()+1)
	for i, filePanel := range m.FilePanels {
		f[i] = filePanel.Render(filePanel.IsFocused)
	}
	if m.FilePreview.IsOpen() {
		f[m.PanelCount()] = m.FilePreview.GetContent()
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, f...)
}
