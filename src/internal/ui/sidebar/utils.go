package sidebar

import (
	"log/slog"

	"github.com/yorukot/superfile/src/pkg/utils"
)

// isDivider returns true if the directory is one of the section dividers.
func (d directory) isDivider() bool {
	return d == homeDividerDir || d == pinnedDividerDir || d == diskDividerDir
}

// requiredHeight returns the number of terminal lines required to render this item.
func (d directory) requiredHeight() int {
	if d.isDivider() {
		return dividerDirHeight
	}
	return 1
}

// NoActualDir returns true if the sidebar contains only dividers and no actual directories.
func (s *Model) NoActualDir() bool {
	for _, d := range s.directories {
		if !d.isDivider() {
			return false
		}
	}
	return true
}

// isCursorInvalid returns true if the current cursor position is out of bounds or points to a divider.
func (s *Model) isCursorInvalid() bool {
	return s.cursor < 0 || s.cursor >= len(s.directories) || s.directories[s.cursor].isDivider()
}

// resetCursor moves the cursor to the first selectable directory in the sidebar.
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

// pinnedIndexRange calculates the start and end indices of the pinned directories section.
// Returns (-1, -1) if the section is missing or empty.
func (s *Model) pinnedIndexRange() (int, int) {
	begin, end := -1, -1
	for i, d := range s.directories {
		if d.Section == utils.SidebarSectionPinned {
			if begin == -1 {
				begin = i
			}
			end = i
		}
	}

	return begin, end
}

// GetWidth returns the current width of the sidebar.
func (m *Model) GetWidth() int {
	return m.width
}

// GetHeight returns the current height of the sidebar.
func (m *Model) GetHeight() int {
	return m.height
}

// SetHeight updates the height of the sidebar, ensuring it meets the minimum requirement.
func (m *Model) SetHeight(height int) {
	if height < minHeight {
		slog.Error("Attempted to set too low height to sidebar", "height", height)
		return
	}
	m.height = height
}

// Disabled returns true if the sidebar is currently disabled.
func (m *Model) Disabled() bool {
	return m.disabled
}
