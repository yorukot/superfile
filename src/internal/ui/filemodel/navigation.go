package filemodel

import "log/slog"

// SetFocusedPanelIndex moves focus to the panel at the given index, keeping the
// per panel focus flags in step with FocusedPanelIndex.
//
// Unlike MoveFocusedPanelBy it is idempotent: an index that is already focused
// changes nothing, which is what a click on the focused panel should do.
func (m *Model) SetFocusedPanelIndex(index int) {
	if m.PanelCount() == 0 {
		slog.Error("Unexpected error: fileModel with 0 panels")
		return
	}
	if index < 0 || index >= m.PanelCount() {
		slog.Error("Unexpected panel index for focus change",
			"index", index, "panel count", m.PanelCount())
		return
	}
	if m.FocusedPanelIndex == index {
		return
	}

	m.FilePanels[m.FocusedPanelIndex].IsFocused = false
	m.FocusedPanelIndex = index
	m.FilePanels[index].IsFocused = true
}

func (m *Model) NextFilePanel() {
	m.MoveFocusedPanelBy(1)
}

func (m *Model) PreviousFilePanel() {
	m.MoveFocusedPanelBy(-1)
}

func (m *Model) MoveFocusedPanelBy(delta int) {
	if m.PanelCount() == 0 {
		slog.Error("Unexpected error: fileModel with 0 panels")
		return
	}
	m.GetFocusedFilePanel().IsFocused = false
	m.FocusedPanelIndex = (m.FocusedPanelIndex + delta + m.PanelCount()) % m.PanelCount()
	m.FilePanels[m.FocusedPanelIndex].IsFocused = true
}
