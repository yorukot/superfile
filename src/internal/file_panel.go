package internal

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yorukot/superfile/src/internal/utils"
)

func (panel *filePanel) getSelectedItem() element {
	if panel.cursor < 0 || len(panel.element) <= panel.cursor {
		return element{}
	}
	return panel.element[panel.cursor]
}

func (panel *filePanel) resetSelected() {
	panel.selected = panel.selected[:0]
}

// For modification. Make sure to do a nil check
func (panel *filePanel) getSelectedItemPtr() *element {
	if panel.cursor < 0 || len(panel.element) <= panel.cursor {
		return nil
	}
	return &panel.element[panel.cursor]
}

// Note : This will soon be moved to its own package.
func (panel *filePanel) changeFilePanelMode() {
	switch panel.panelMode {
	case selectMode:
		panel.selected = panel.selected[:0]
		panel.panelMode = browserMode
	case browserMode:
		panel.panelMode = selectMode
	default:
		slog.Error("Unexpected panelMode", "panelMode", panel.panelMode)
	}
}

// This should be the function that is always called whenever we are updating a directory.
func (panel *filePanel) updateCurrentFilePanelDir(path string) error {
	slog.Debug("updateCurrentFilePanelDir", "panel.location", panel.location, "path", path)
	// In case non Absolute path is passed, make sure to resolve it.
	path = utils.ResolveAbsPath(panel.location, path)

	// Ignore if its the same directory. It prevents resetting of searchBar
	if path == panel.location {
		return nil
	}

	// NOTE: This could be a configurable feature
	// Update the cursor and render status in case we switch back to this.
	panel.directoryRecords[panel.location] = directoryRecord{
		directoryCursor: panel.cursor,
		directoryRender: panel.render,
	}

	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s : no such file or directory, stats err : %w", path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	// Switch to "path"
	panel.location = path

	// TODO(BUG) : We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), we must validate
	// the cursor and render values.
	curDirectoryRecord, hasRecord := panel.directoryRecords[panel.location]
	if hasRecord {
		panel.cursor = curDirectoryRecord.directoryCursor
		panel.render = curDirectoryRecord.directoryRender
	} else {
		panel.cursor = 0
		panel.render = 0
	}

	slog.Debug("updateCurrentFilePanelDir : After update", "cursor", panel.cursor, "render", panel.render)

	// Reset the searchbar Value
	// TODO(Refactoring) : Have a common searchBar type for sidebar and this search bar.
	panel.searchBar.SetValue("")

	return nil
}

func (panel *filePanel) parentDirectory() error {
	return panel.updateCurrentFilePanelDir("..")
}
