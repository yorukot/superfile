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
