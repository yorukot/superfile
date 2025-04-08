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

	pinnedBegin, pinnedEnd := m.sidebarModel.pinnedIndexRange()
	// We have not selected a pinned directory, rename is not allowed
	if sidebar.cursor < pinnedBegin || sidebar.cursor > pinnedEnd {
		return
	}

	nameLen := len(sidebar.directories[sidebar.cursor].name)
	cursorPos := nameLen

	m.sidebarModel.renaming = true
	m.sidebarModel.rename = common.GeneratePinnedRenameTextInput(cursorPos, sidebar.directories[sidebar.cursor].name)
}

// Cancel rename pinned directory
func (m *model) cancelSidebarRename() {
	sidebar := &m.sidebarModel
	sidebar.rename.Blur()
	sidebar.renaming = false
}

// Confirm rename pinned directory
func (m *model) confirmSidebarRename() {
	sidebar := &m.sidebarModel

	itemLocation := sidebar.directories[sidebar.cursor].location
	newItemName := sidebar.rename.Value()
	// This is needed to update the current pinned directory data loaded into memory
	sidebar.directories[sidebar.cursor].name = newItemName

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
		if dir.location == itemLocation {
			dir.name = newItemName
		}
		pinnedDirs = append(pinnedDirs, pinnedDir{Location: dir.location, Name: dir.name})
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

func (s *sidebarModel) pinnedIndexRange() (int, int) {
	// pinned directories start after well-known directories and the divider
	// Can't use getPinnedDirectories() here, as if we are in search mode, we would be showing
	// and having less directories in sideBar.directories slice

	// Todo : This is inefficient to iterate each time for this.
	// This information can be kept precomputed
	pinnedDividerIdx := -1
	diskDividerIdx := -1
	for i, d := range s.directories {
		if d == pinnedDividerDir {
			pinnedDividerIdx = i
		}
		if d == diskDividerDir {
			diskDividerIdx = i
			break
		}
	}
	return pinnedDividerIdx + 1, diskDividerIdx - 1
}
