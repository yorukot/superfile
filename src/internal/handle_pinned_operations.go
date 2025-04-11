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

	pinnedDirs := common.GetPinnedDirectories()
	// Call getPinnedDirectories, instead of using what is stored in sidebar.directories
	// sidebar.directories could have less directories in case a search filter is used
	for i := range pinnedDirs {
		// Considering the situation when many
		if pinnedDirs[i].Location == itemLocation {
			pinnedDirs[i].Name = newItemName
		}
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
