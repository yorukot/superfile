package filepanel

import (
	"context"
	"fmt"
	"log/slog"
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
	currentLocation := m.CurrentLocation()
	targetPath := resolveLocationPath(currentLocation, path)

	// Ignore if its the same directory. It prevents resetting of searchBar
	if targetPath.String() == currentLocation.Path.String() {
		return nil
	}

	session, err := m.paneSession()
	if err != nil {
		return err
	}
	info, err := session.Stat(context.Background(), targetPath)
	if err != nil {
		return err
	}
	if !info.IsDir {
		return fmt.Errorf("%s is not a directory", targetPath.String())
	}
	return m.ApplyCurrentFilePanelDir(targetPath.String())
}

func (m *Model) ApplyCurrentFilePanelDir(path string) error {
	currentLocation := m.CurrentLocation()
	targetPath := resolveLocationPath(currentLocation, path)
	if targetPath.String() == currentLocation.Path.String() {
		return nil
	}
	// Update the cursor and render status in case we switch back to this.
	m.DirectoryRecords[locationKey(currentLocation)] = directoryRecord{
		directoryCursor: m.cursor,
		directoryRender: m.renderIndex,
	}

	// In case of switching to parent, explicitly set focus.
	// This is to handle when there isn't a DirectoryRecord, yet.
	if parentPath(currentLocation.Path).String() == targetPath.String() {
		m.TargetFile = baseName(currentLocation.Path)
	}
	// Switch to "path"
	location := currentLocation
	location.Path = targetPath
	m.SetPaneLocation(location)

	// NOTE: We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), on dire change
	// we do a element fetch and validate the cursor and render values. But the filepane could
	// stay in invalid state till that and operations done before the update may fail
	curDirectoryRecord, hasRecord := m.DirectoryRecords[locationKey(location)]
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
