package internal

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/yorukot/superfile/src/internal/common"

	variable "github.com/yorukot/superfile/src/config"
)

// Rename file where the cusror is located
func (m *model) pinnedItemRename() {
	sidebar := m.sidebarModel

	pinnedBegin, pinnedEnd := m.sidebarModel.PinnedIndexRange()
	// We have not selected a pinned directory, rename is not allowed
	if sidebar.Cursor < pinnedBegin || sidebar.Cursor > pinnedEnd {
		return
	}

	nameLen := len(sidebar.Directories[sidebar.Cursor].Name)
	cursorPos := nameLen

	m.sidebarModel.Renaming = true
	m.sidebarModel.Rename = common.GeneratePinnedRenameTextInput(cursorPos, sidebar.Directories[sidebar.Cursor].Name)
}

// Cancel rename pinned directory
func (m *model) cancelSidebarRename() {
	sidebar := &m.sidebarModel
	sidebar.Rename.Blur()
	sidebar.Renaming = false
}

// Confirm rename pinned directory
func (m *model) confirmSidebarRename() {
	sidebar := &m.sidebarModel

	itemLocation := sidebar.Directories[sidebar.Cursor].Location
	newItemName := sidebar.Rename.Value()
	// This is needed to update the current pinned directory data loaded into memory
	sidebar.Directories[sidebar.Cursor].Name = newItemName

	// recover the state of rename
	m.cancelSidebarRename()

	type pinnedDir struct {
		Location string `json:"location"`
		Name     string `json:"name"`
	}
	var pinnedDirs []pinnedDir

	// Call getPinnedDirectories, instead of using what is stored in sidebar.directories
	// sidebar.directories could have less directories in case a search filter is used
	for _, dir := range getPinnedDirectories() {
		// Considering the situation when many
		if dir.Location == itemLocation {
			dir.Name = newItemName
		}
		pinnedDirs = append(pinnedDirs, pinnedDir{Location: dir.Location, Name: dir.Name})
	}

	jsonData, err := json.Marshal(pinnedDirs)
	if err != nil {
		slog.Error("Error marshaling pinned directories data", "error", err)
	}

	err = os.WriteFile(variable.PinnedFile, jsonData, 0644)
	if err != nil {
		slog.Error("Error updating pinned directories data", "error", err)
	}
}

func (s *SidebarModel) PinnedIndexRange() (int, int) {
	// pinned directories start after well-known directories and the divider
	// Can't use getPinnedDirectories() here, as if we are in search mode, we would be showing
	// and having less directories in sideBar.directories slice

	// Todo : This is inefficient to iterate each time for this.
	// This information can be kept precomputed
	pinnedDividerIdx := -1
	diskDividerIdx := -1
	for i, d := range s.Directories {
		if d == PinnedDividerDir {
			pinnedDividerIdx = i
		}
		if d == DiskDividerDir {
			diskDividerIdx = i
			break
		}
	}
	return pinnedDividerIdx + 1, diskDividerIdx - 1
}
