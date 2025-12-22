package filemodel

import "github.com/charmbracelet/lipgloss"

func (m *FileModel) Render() string {
	f := make([]string, m.PanelCount())
	for i, filePanel := range m.FilePanels {
		f[i] = filePanel.Render(filePanel.IsFocused)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, f...)
}
