package internal

import (
	"encoding/json"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/yorukot/superfile/src/config/icon"

	variable "github.com/yorukot/superfile/src/config"
)

// Rename file where the cusror is located
func (m *model) pinnedItemRename() {
	sidebar := m.sidebarModel

	pinnedBegin, pinnedEnd := pinnedIndexRange()
	if sidebar.cursor < pinnedBegin || sidebar.cursor >= pinnedEnd {
		return
	}

	nameLen := len(sidebar.directories[sidebar.cursor].name)
	cursorPos := nameLen

	ti := textinput.New()
	ti.Cursor.Style = filePanelCursorStyle
	ti.Cursor.TextStyle = filePanelStyle
	ti.Prompt = filePanelCursorStyle.Render(icon.Cursor + " ")
	ti.TextStyle = modalStyle
	ti.Cursor.Blink = true
	ti.Placeholder = "New name"
	ti.PlaceholderStyle = modalStyle
	ti.SetValue(sidebar.directories[sidebar.cursor].name)
	ti.SetCursor(cursorPos)
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = Config.SidebarWidth - 4

	m.sidebarModel.renaming = true
	m.sidebarModel.rename = ti
}

// Cancel rename pinned directory
func (m *model) cancelSidebarRename() {
	sidebar := &m.sidebarModel
	sidebar.rename.Blur()
	sidebar.renaming = false
}

// Confirm rename pinned directory
func (m *model) confirmSidebarRename() {
	sidebar := m.sidebarModel

	sidebar.directories[sidebar.cursor].name = sidebar.rename.Value()
	// recover the state of rename
	m.cancelSidebarRename()

	type pinnedDir struct {
		Location string `json:"location"`
		Name     string `json:"name"`
	}
	var pinnedDirs []pinnedDir

	pinnedBegin, pinnedEnd := pinnedIndexRange()
	dirs := sidebar.directories[pinnedBegin:pinnedEnd]
	for _, dir := range dirs {
		pinnedDirs = append(pinnedDirs, pinnedDir{Location: dir.location, Name: dir.name})
	}

	jsonData, err := json.Marshal(pinnedDirs)
	if err != nil {
		outPutLog("Error marshaling pinned directories data")
	}

	err = os.WriteFile(variable.PinnedFile, jsonData, 0644)
	if err != nil {
		outPutLog("Error updating pinned directories data", err)
	}
}

func pinnedIndexRange() (int, int) {
	// pinned directories start after well-known directories and the divider
	pinnedIndex := len(getWellKnownDirectories()) + 1
	pinnedLen := len(getPinnedDirectories())
	return pinnedIndex, pinnedIndex + pinnedLen
}
