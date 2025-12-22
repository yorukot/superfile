package filemodel

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

func (m *FileModel) CreateNewFilePanel(location string) error {
	if m.PanelCount() >= m.MaxFilePanel {
		return ErrMaximumPanelCount
	}

	if _, err := os.Stat(location); err != nil {
		return fmt.Errorf("cannot access location : %s", location)
	}

	m.FilePanels = append(m.FilePanels, filepanel.New(
		location, m.GetFocusedFilePanel().SortOptions, false, ""))

	newPanelIndex := m.PanelCount() - 1

	m.FilePanels[m.FocusedPanelIndex].IsFocused = false
	m.FilePanels[newPanelIndex].IsFocused = true
	m.FocusedPanelIndex = newPanelIndex

	m.UpdateChildComponentWidth()

	return nil
}

func (m *FileModel) CloseFilePanel() {
	if m.PanelCount() <= 1 {
		slog.Error("CloseFilePanel called on with panelCount <= 1")
		return
	}

	m.FilePanels = append(m.FilePanels[:m.FocusedPanelIndex],
		m.FilePanels[m.FocusedPanelIndex+1:]...)

	if m.FocusedPanelIndex != 0 {
		m.FocusedPanelIndex--
	}
	m.FilePanels[m.FocusedPanelIndex].IsFocused = true

	m.UpdateChildComponentWidth()
}
