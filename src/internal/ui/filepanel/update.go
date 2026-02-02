package filepanel

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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
		directoryCursor: m.cursor,
		directoryRender: m.renderIndex,
	}

	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s : no such file or directory, stats err : %w", path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	// In case of switching to parent, explicitly set focus.
	// This is to handle when there isn't a DirectoryRecord, yet.
	if filepath.Dir(m.Location) == path {
		m.TargetFile = filepath.Base(m.Location)
	}
	// Switch to "path"
	m.Location = path

	// NOTE: We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), on dire change
	// we do a element fetch and validate the cursor and render values. But the filepane could
	// stay in invalid state till that and operations done before the update may fail
	curDirectoryRecord, hasRecord := m.DirectoryRecords[m.Location]
	if hasRecord {
		m.cursor = curDirectoryRecord.directoryCursor
		m.renderIndex = curDirectoryRecord.directoryRender
	} else {
		m.cursor = 0
		m.renderIndex = 0
	}

	slog.Debug("updateCurrentFilePanelDir : After update", "cursor", m.cursor, "render", m.renderIndex)

	// Reset the searchbar Value
	// TODO(Refactoring) : Have a common searchBar type for sidebar and this search bar.
	m.SearchBar.SetValue("")

	return nil
}

func (m *Model) ParentDirectory() error {
	return m.UpdateCurrentFilePanelDir("..")
}

// Select all item in the file panel (only work on select mode)
func (m *Model) SelectAllItem() {
	for _, item := range m.element {
		m.SetSelected(item.Location)
	}
}
