package sidebar

import (
	"github.com/yorukot/superfile/src/pkg/utils"
)

// These are effectively consts
// Had to use `var` as go doesn't allows const structs
var homeDividerDir = directory{ //nolint: gochecknoglobals // This is more like a const.
	Name:     "",
	Location: "Home+-*/=?",
}

var pinnedDividerDir = directory{ //nolint: gochecknoglobals // This is more like a const.
	Name:     "",
	Location: "Pinned+-*/=?",
}

var diskDividerDir = directory{ //nolint: gochecknoglobals // This is more like a const.
	Name:     "",
	Location: "Disks+-*/=?",
}

var defaultSectionSlice = []string{ //nolint: gochecknoglobals // This is more like a const.
	utils.SidebarSectionHome, utils.SidebarSectionPinned, utils.SidebarSectionDisks,
}

// superfile logo + blank line + search bar
const sideBarInitialHeight = 3

// Lines of the sidebar that sit above the directory list, counted within the
// sidebar's content area and therefore excluding its border.
const (
	// directoryHeaderLines covers the superfile logo line and the blank line
	// that follows it.
	directoryHeaderLines = 2
	// searchBarHeight is the single row taken by the search bar, but only while
	// the search bar is actually being drawn.
	searchBarHeight = 1
)

// UI dimension constants for sidebar
const (
	// searchBarPadding is the total padding for search bar (borders + prompt + extra char)
	searchBarPadding = 5 // 2 (borders) + 2 (prompt) + 1 (extra char)

	directoryCapacityForDividers = 2

	// dividerDirHeight is the default height when no height is available
	dividerDirHeight = 3

	minHeight = 5
	minWidth  = 7
)
