package sidebar

func (d directory) IsDivider() bool {
	return d == pinnedDividerDir || d == diskDividerDir
}
func (d directory) RequiredHeight() int {
	if d.IsDivider() {
		return defaultRenderHeight
	}
	return 1
}

// True if only dividers are in directories slice,
// but no actual directories
// This will be pretty quick. But we can replace it with
// len(s.directories) <= 2 - More hacky and hardcoded-like, but faster
func (s *Model) NoActualDir() bool {
	for _, d := range s.directories {
		if !d.IsDivider() {
			return false
		}
	}
	return true
}

func (s *Model) isCursorInvalid() bool {
	return s.cursor < 0 || s.cursor >= len(s.directories) || s.directories[s.cursor].IsDivider()
}

func (s *Model) resetCursor() {
	s.cursor = 0
	// Move to first non Divider dir
	for i, d := range s.directories {
		if !d.IsDivider() {
			s.cursor = i
			return
		}
	}
	// If all directories are divider, code will reach here. and s.cursor will stay 0
	// Or s.directories is empty
}

// SearchBarFocused returns whether the search bar is focused
func (s *Model) SearchBarFocused() bool {
	return s.searchBar.Focused()
}

// SearchBarBlur removes focus from the search bar
func (s *Model) SearchBarBlur() {
	s.searchBar.Blur()
}

// SearchBarFocus sets focus on the search bar
func (s *Model) SearchBarFocus() {
	s.searchBar.Focus()
}

// IsRenaming returns whether the sidebar is currently in renaming mode
func (s *Model) IsRenaming() bool {
	return s.renaming
}

// GetCurrentDirectoryLocation returns the location of the currently selected directory
func (s *Model) GetCurrentDirectoryLocation() string {
	if s.isCursorInvalid() || s.NoActualDir() {
		return ""
	}
	return s.directories[s.cursor].Location
}

func (s *Model) pinnedIndexRange() (int, int) {
	// pinned directories start after well-known directories and the divider
	// Can't use getPinnedDirectories() here, as if we are in search mode, we would be showing
	// and having less directories in sideBar.directories slice

	// TODO : This is inefficient to iterate each time for this.
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
