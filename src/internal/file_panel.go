package internal

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/yorukot/superfile/src/internal/utils"
)

func (panel *FilePanel) GetSelectedItem() Element {
	if panel.Cursor < 0 || len(panel.Element) <= panel.Cursor {
		return Element{}
	}
	return panel.Element[panel.Cursor]
}

func (panel *FilePanel) ResetSelected() {
	panel.Selected = panel.Selected[:0]
}

// For modification. Make sure to do a nil check
func (panel *FilePanel) GetSelectedItemPtr() *Element {
	if panel.Cursor < 0 || len(panel.Element) <= panel.Cursor {
		return nil
	}
	return &panel.Element[panel.Cursor]
}

// Note : This will soon be moved to its own package.
func (panel *FilePanel) ChangeFilePanelMode() {
	switch panel.PanelMode {
	case SelectMode:
		panel.Selected = panel.Selected[:0]
		panel.PanelMode = BrowserMode
	case BrowserMode:
		panel.PanelMode = SelectMode
	default:
		slog.Error("Unexpected panelMode", "panelMode", panel.PanelMode)
	}
}

// This should be the function that is always called whenever we are updating a directory.
func (panel *FilePanel) UpdateCurrentFilePanelDir(path string) error {
	slog.Debug("updateCurrentFilePanelDir", "panel.location", panel.Location, "path", path)
	// In case non Absolute path is passed, make sure to resolve it.
	path = utils.ResolveAbsPath(panel.Location, path)

	// Ignore if its the same directory. It prevents resetting of searchBar
	if path == panel.Location {
		return nil
	}

	// NOTE: This could be a configurable feature
	// Update the cursor and render status in case we switch back to this.
	panel.DirectoryRecords[panel.Location] = DirectoryRecord{
		DirectoryCursor: panel.Cursor,
		DirectoryRender: panel.RenderIndex,
	}

	if info, err := os.Stat(path); err != nil {
		return fmt.Errorf("%s : no such file or directory, stats err : %w", path, err)
	} else if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	// Switch to "path"
	panel.Location = path

	// TODO(BUG) : We are fetching the cursor and render from cache, but this could become invalid
	// in case user deletes some items in the directory via another file manager and then switch back
	// Basically this directoryRecords cache can be invalid. On each Update(), we must validate
	// the cursor and render values.
	curDirectoryRecord, hasRecord := panel.DirectoryRecords[panel.Location]
	if hasRecord {
		panel.Cursor = curDirectoryRecord.DirectoryCursor
		panel.RenderIndex = curDirectoryRecord.DirectoryRender
	} else {
		panel.Cursor = 0
		panel.RenderIndex = 0
	}

	slog.Debug("updateCurrentFilePanelDir : After update", "cursor", panel.Cursor, "render", panel.RenderIndex)

	// Reset the searchbar Value
	// TODO(Refactoring) : Have a common searchBar type for sidebar and this search bar.
	panel.SearchBar.SetValue("")

	return nil
}

func (panel *FilePanel) ParentDirectory() error {
	return panel.UpdateCurrentFilePanelDir("..")
}

func (panel *FilePanel) HandleResize(height int) {
	// Min render cursor that keeps the cursor in view
	minVisibleRenderCursor := panel.Cursor - panelElementHeight(height) + 1
	// Max render cursor. This ensures all elements are rendered if there is space
	maxRenderCursor := max(len(panel.Element)-panelElementHeight(height), 0)

	if panel.RenderIndex > maxRenderCursor {
		panel.RenderIndex = maxRenderCursor
	}
	if panel.RenderIndex < minVisibleRenderCursor {
		panel.RenderIndex = minVisibleRenderCursor
	}
}

// Select the item where cursor located (only work on select mode)
func (panel *FilePanel) SingleItemSelect() {
	if len(panel.Element) > 0 && panel.Cursor >= 0 && panel.Cursor < len(panel.Element) {
		elementLocation := panel.Element[panel.Cursor].Location

		if arrayContains(panel.Selected, elementLocation) {
			// This is inefficient. Once you select 1000 items,
			// each select / deselect operation can take 1000 operations
			// It can be easily made constant time.
			// TODO : (performance)convert panel.selected to a set (map[string]struct{})
			panel.Selected = removeElementByValue(panel.Selected, elementLocation)
		} else {
			panel.Selected = append(panel.Selected, elementLocation)
		}
	}
}
