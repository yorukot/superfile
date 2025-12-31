package filepanel

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yorukot/superfile/src/internal/utils"
)

func (m *Model) ChangeFilePanelMode() {
	switch m.PanelMode {
	case SelectMode:
		m.ResetSelected()
		m.PanelMode = BrowserMode
	case BrowserMode:
		m.PanelMode = SelectMode
	default:
		slog.Error("Unexpected panelMode", "panelMode", m.PanelMode)
	}
}

func (m *Model) SetSelected(location string, state bool) {
	if state {
		m.selectOrderCounter++
		m.selected[location] = m.selectOrderCounter
	} else {
		delete(m.selected, location)
	}
}

func (m *Model) SetSelectedAll(locations []string, state bool) {
	for _, location := range locations {
		m.SetSelected(location, state)
	}
}

// This should be the function that is always called whenever we are updating a directory.
func (m *Model) UpdateCurrentFilePanelDir(path string) error {
	slog.Debug("updateCurrentFilePanelDir", "panel.location", m.Location, "path", path)
	// In case non Absolute path is passed, make sure to resolve it.
	path = utils.ResolveAbsPath(m.Location, path)

	// Ignore if its the same directory. It prevents resetting of searchBar
	if path == m.Location {
		return nil
	}

	// NOTE: This could be a configurable feature
	// Update the cursor and render status in case we switch back to this.
	m.DirectoryRecords[m.Location] = directoryRecord{
		directoryCursor: m.Cursor,
		directoryRender: m.RenderIndex,
	}

	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s : no such file or directory, stats err : %w", path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	// Switch to "path"
	m.Location = path

	// TODO(BUG) : We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), we must validate
	// the cursor and render values.
	curDirectoryRecord, hasRecord := m.DirectoryRecords[m.Location]
	if hasRecord {
		m.Cursor = curDirectoryRecord.directoryCursor
		m.RenderIndex = curDirectoryRecord.directoryRender
	} else {
		m.Cursor = 0
		m.RenderIndex = 0
	}

	slog.Debug("updateCurrentFilePanelDir : After update", "cursor", m.Cursor, "render", m.RenderIndex)

	// Reset the searchbar Value
	// TODO(Refactoring) : Have a common searchBar type for sidebar and this search bar.
	m.SearchBar.SetValue("")

	return nil
}

func (m *Model) ParentDirectory() error {
	return m.UpdateCurrentFilePanelDir("..")
}
