package sidebar

import "log/slog"

func (d directory) isDivider() bool {
	return d == homeDividerDir || d == pinnedDividerDir || d == diskDividerDir
}
func (d directory) requiredHeight() int {
	if d.isDivider() {
		return dividerDirHeight
	}
	return 1
}

// True if only dividers are in directories slice,
// but no actual directories
// This will be pretty quick. But we can replace it with
// len(s.directories) <= 2 - More hacky and hardcoded-like, but faster
func (s *Model) NoActualDir() bool {
	for _, d := range s.directories {
		if !d.isDivider() {
			return false
		}
	}
	return true
}

func (s *Model) isCursorInvalid() bool {
	return s.cursor < 0 || s.cursor >= len(s.directories) || s.directories[s.cursor].isDivider()
}

func (s *Model) resetCursor() {
	s.cursor = 0
	// Move to first non Divider dir
	for i, d := range s.directories {
		if !d.isDivider() {
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
	pinnedDividerIdx := -1
	for i, d := range s.directories {
		if d == pinnedDividerDir {
			pinnedDividerIdx = i
			break
		}
	}

	if pinnedDividerIdx == -1 {
		return -1, -1
	}

	pinnedEndIdx := len(s.directories) - 1
	for i := pinnedDividerIdx + 1; i < len(s.directories); i++ {
		if s.directories[i].isDivider() {
			pinnedEndIdx = i - 1
			break
		}
	}

	return pinnedDividerIdx + 1, pinnedEndIdx
}

// TODO: There are some utils like this that are common in all models
// Come up with a way to prevent all this code duplication
func (m *Model) GetWidth() int {
	return m.width
}

func (m *Model) GetHeight() int {
	return m.height
}

func (m *Model) SetHeight(height int) {
	if height < minHeight {
		slog.Error("Attempted to set too low height to sidebar", "height", height)
		return
	}
	m.height = height
}

func (m *Model) Disabled() bool {
	return m.disabled
}
