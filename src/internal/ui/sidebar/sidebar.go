package sidebar

import (
	"log/slog"
	"slices"

	tea "github.com/charmbracelet/bubbletea"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
)

// Rename file where the cursor is located
func (s *Model) PinnedItemRename() {
	pinnedBegin, pinnedEnd := s.pinnedIndexRange()
	// We have not selected a pinned directory, rename is not allowed
	if s.cursor < pinnedBegin || s.cursor > pinnedEnd {
		return
	}

	nameLen := len(s.directories[s.cursor].Name)
	cursorPos := nameLen

	s.renaming = true
	s.rename = common.GeneratePinnedRenameTextInput(cursorPos, s.directories[s.cursor].Name)
}

// Cancel rename pinned directory
func (s *Model) CancelSidebarRename() {
	s.rename.Blur()
	s.renaming = false
}

// Confirm rename pinned directory
func (s *Model) ConfirmSidebarRename() {
	itemLocation := s.directories[s.cursor].Location
	newItemName := s.rename.Value()
	// This is needed to update the current pinned directory data loaded into memory
	s.directories[s.cursor].Name = newItemName

	// recover the state of rename
	s.CancelSidebarRename()

	pinnedDirs := s.pinnedMgr.Load()
	// Call getPinnedDirectories, instead of using what is stored in sidebar.directories
	// sidebar.directories could have less directories in case a search filter is used
	for i := range pinnedDirs {
		// Considering the situation when many
		if pinnedDirs[i].Location == itemLocation {
			pinnedDirs[i].Name = newItemName
		}
	}

	if err := s.pinnedMgr.Save(pinnedDirs); err != nil {
		slog.Error("error saving pinned directories", "error", err)
	}
}

// UpdateState handles the sidebar's state updates
func (s *Model) UpdateState(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	if s.renaming {
		s.rename, cmd = s.rename.Update(msg)
	} else if s.searchBar.Focused() {
		s.searchBar, cmd = s.searchBar.Update(msg)
	}

	if s.cursor < 0 {
		s.cursor = 0
	}
	return cmd
}

// HandleSearchBarKey handles key events for the sidebar search bar
func (s *Model) HandleSearchBarKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		s.SearchBarBlur()
		s.searchBar.SetValue("")
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		s.SearchBarBlur()
		s.resetCursor()
	}
}

// UpdateDirectories updates the directories list based on search value
// This is a bit inefficient, as we already had the directories when we
// initialized the sidebar. We call the directory fetching logic many times
// which is a disk heavy operation.
func (s *Model) UpdateDirectories() {
	if s.Disabled() {
		return
	}
	if s.searchBar.Value() != "" {
		s.directories = getFilteredDirectories(s.searchBar.Value(), s.pinnedMgr)
	} else {
		s.directories = getDirectories(s.pinnedMgr)
	}
	// This is needed, as due to filtering, the cursor might be invalid
	if s.isCursorInvalid() {
		s.resetCursor()
	}
}

func (s *Model) TogglePinnedDirectory(dir string) error {
	return s.pinnedMgr.Toggle(dir)
}

// New creates a new sidebar model with the given parameters
func New() Model {
	if common.Config.SidebarWidth == 0 {
		return Model{
			disabled: true,
		}
	}
	// pinnedMgr is created here, can be done higher up in the call chain
	pinnedMgr := NewPinnedFileManager(variable.PinnedFile)
	res := Model{
		renderIndex: 0,
		directories: getDirectories(&pinnedMgr),
		searchBar:   common.GenerateSearchBar(),
		pinnedMgr:   &pinnedMgr,
		width:       common.Config.SidebarWidth + common.BorderPadding,
		height:      minHeight,
		disabled:    false,
	}

	res.searchBar.Width = res.width - common.BorderPadding - searchBarPadding
	res.searchBar.Placeholder = "(" + common.Hotkeys.SearchBar[0] + ")" + " Search"
	return res
}
