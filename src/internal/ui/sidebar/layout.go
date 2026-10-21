package sidebar

import (
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/mouse"
)

// directoryRow records where one sidebar entry ended up on screen.
type directoryRow struct {
	// index is the position of the entry within the directories slice.
	index int
	// top is the terminal row the entry starts on, counted from the top of the
	// sidebar and including its border.
	top int
	// height is the number of terminal rows the entry occupies.
	height int
}

// searchBarRendered reports whether the search bar takes a line in the rendered
// output. The sidebar draws it when focused, whenever it holds a query, and
// whenever the sidebar itself has focus.
//
// Render and mouse hit-testing read this predicate, so a click always addresses
// the row the user actually sees rather than a shifted guess.
func (s *Model) searchBarRendered(sidebarFocused bool) bool {
	return s.searchBar.Focused() || s.searchBar.Value() != "" || sidebarFocused
}

// directoriesTop returns the terminal row, counted from the top of the sidebar
// and including its border, on which the directory list starts.
func (s *Model) directoriesTop(sidebarFocused bool) int {
	top := common.BorderTopHeight + directoryHeaderLines
	if s.searchBarRendered(sidebarFocused) {
		top += searchBarHeight
	}
	return top
}

// layoutDirectoryRows walks the directory list the same way rendering does and
// reports the on-screen position of every entry that fits. Render and mouse
// hit-testing both consume this, so they cannot disagree about the scroll offset
// or about which entries were dropped for lack of space.
//
// Dividers occupy rows but are reported like any other entry; callers that only
// care about selectable entries skip them with isDivider.
func (s *Model) layoutDirectoryRows(sidebarFocused bool) []directoryRow {
	if s.NoActualDir() {
		return nil
	}

	mainPanelHeight := s.height - common.BorderPadding
	// The budget always reserves the search bar row, which keeps the number of
	// entries drawn unchanged from before. When the search bar is absent the
	// row stays unused at the bottom of the sidebar, as it always has.
	budget := sideBarInitialHeight
	top := s.directoriesTop(sidebarFocused)
	rows := make([]directoryRow, 0, len(s.directories))

	for i := s.renderIndex; i < len(s.directories); i++ {
		height := s.directories[i].requiredHeight()
		if budget+height > mainPanelHeight {
			break
		}
		budget += height
		rows = append(rows, directoryRow{index: i, top: top, height: height})
		top += height
	}
	return rows
}

// DirectoryIndexAtRow returns the index within directories of the entry drawn on
// the given terminal row, or NoIndex when the row holds no entry. Section
// dividers resolve to NoIndex as well since they are not selectable.
func (s *Model) DirectoryIndexAtRow(sidebarFocused bool, y int) int {
	for _, row := range s.layoutDirectoryRows(sidebarFocused) {
		if y < row.top || y >= row.top+row.height {
			continue
		}
		if s.directories[row.index].isDivider() {
			return mouse.NoIndex
		}
		return row.index
	}
	return mouse.NoIndex
}

// MouseRegions registers one region per selectable entry the sidebar currently
// draws. originX and originY are the terminal coordinates of the sidebar's top
// left corner, its border included.
func (s *Model) MouseRegions(reg *mouse.Registry, sidebarFocused bool, originX, originY int) {
	if s.Disabled() {
		return
	}

	for _, row := range s.layoutDirectoryRows(sidebarFocused) {
		if s.directories[row.index].isDivider() {
			continue
		}
		top := originY + row.top
		reg.Add(mouse.SidebarDirectoryTarget(row.index),
			originX, top, originX+s.GetWidth(), top+row.height)
	}
}
