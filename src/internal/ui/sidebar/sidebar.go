package sidebar

import (
	"log/slog"
	"slices"
	"time"

	tea "charm.land/bubbletea/v2"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
)

// PinnedItemRename initiates the rename process for the currently selected pinned directory.
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

// CancelSidebarRename aborts the rename process for a pinned directory.
func (s *Model) CancelSidebarRename() {
	s.rename.Blur()
	s.renaming = false
}

// ConfirmSidebarRename finalizes the rename process and saves changes to the pinned directories file.
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

// UpdateState handles the sidebar's state updates in response to Bubble Tea messages.
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

// HandleSearchBarKey processes key events specifically for the sidebar's search bar.
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

// UpdateDirectoriesIfNeeded refreshes the sidebar, at most once every
// sidebarRefreshInterval. A changed search query, or an explicit force after
// pinning or renaming, is reflected at once - the rest is only polling for pins
// and mounts changed outside superfile, and doing that on every message made
// every keystroke wait on the filesystem.
func (s *Model) UpdateDirectoriesIfNeeded(force bool) {
	if s.Disabled() {
		return
	}
	query := s.searchBar.Value()
	if !force && query == s.loadedQuery && time.Since(s.lastRefresh) < sidebarRefreshInterval {
		return
	}
	s.loadedQuery = query
	s.UpdateDirectories()
	// Stamped after the read, not before: a read that takes as long as the
	// interval would otherwise leave the stamp already expired on return, and the
	// next message would start another one straight away.
	s.lastRefresh = time.Now()
}

// UpdateDirectories refreshes the list of directories based on the search query or section configuration.
// ponytail: still synchronous. Move it to a tea.Cmd if a stat of a pinned
// directory on a hung mount starts blocking the UI - that needs Load() to stop
// saving the cleaned list back, otherwise it races Toggle() over the same file.
func (s *Model) UpdateDirectories() {
	if s.Disabled() {
		return
	}
	if s.searchBar.Value() != "" {
		s.directories = getFilteredDirectories(s.searchBar.Value(), s.pinnedMgr, s.sections)
	} else {
		s.directories = getDirectories(s.pinnedMgr, s.sections)
	}
	// This is needed, as due to filtering, the cursor might be invalid
	if s.isCursorInvalid() {
		s.resetCursor()
	}
}

// TogglePinnedDirectory adds or removes a directory from the pinned list.
func (s *Model) TogglePinnedDirectory(dir string) error {
	return s.pinnedMgr.Toggle(dir)
}

// New initializes and returns a new Model for the sidebar correctly set up with configuration.
func New() Model {
	if common.Config.SidebarWidth == 0 {
		return Model{
			disabled: true,
		}
	}
	// pinnedMgr is created here, can be done higher up in the call chain
	pinnedMgr := NewPinnedFileManager(variable.PinnedFile)
	s := Model{
		renderIndex: 0,
		searchBar:   common.GenerateSearchBar(),
		pinnedMgr:   &pinnedMgr,
		width:       common.Config.SidebarWidth + common.BorderPadding,
		height:      minHeight,
		disabled:    false,
		sections:    common.Config.SidebarSections,
	}

	s.directories = getDirectories(&pinnedMgr, s.sections)
	s.searchBar.SetWidth(s.width - common.BorderPadding - searchBarPadding)
	s.searchBar.Placeholder = "(" + common.Hotkeys.SearchBar[0] + ")" + " Search"
	return s
}
