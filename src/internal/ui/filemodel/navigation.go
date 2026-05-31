package filemodel

import "log/slog"
import "os"

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

	// Updates working directory to properly pass cwd to newly spawned processes (terminals,editors, programs)
	_ = os.Chdir(m.GetFocusedFilePanel().Location)

}
