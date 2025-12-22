package internal

import "log/slog"

func (m *FileModel) NextFilePanel() {
	m.MoveFocusedPanelBy(1)
}

func (m *FileModel) PreviousFilePanel() {
	m.MoveFocusedPanelBy(-1)
}

func (m *FileModel) MoveFocusedPanelBy(delta int) {
	if m.PanelCount() == 0 {
		slog.Error("Unexpected error: fileModel with 0 panels")
		return
	}
	m.GetFocusedFilePanel().IsFocused = false
	m.FocusedPanelIndex = (m.FocusedPanelIndex + delta + m.PanelCount()) % m.PanelCount()
	m.FilePanels[m.FocusedPanelIndex].IsFocused = true
}
