package sidebar

// These are effectively consts
// Had to use `var` as go doesn't allows const structs
var pinnedDividerDir = directory{ //nolint: gochecknoglobals // This is more like a const.
	Name:     "",
	Location: "Pinned+-*/=?",
}

var diskDividerDir = directory{ //nolint: gochecknoglobals // This is more like a const.
	Name:     "",
	Location: "Disks+-*/=?",
}

// superfile logo + blank line + search bar
const sideBarInitialHeight = 3

// UI dimension constants for sidebar
const (
	// searchBarPadding is the total padding for search bar (borders + prompt + extra char)
	searchBarPadding = 5 // 2 (borders) + 2 (prompt) + 1 (extra char)

	// directoryCapacityExtra is extra capacity for separator lines in directory list
	directoryCapacityExtra = 2

	// defaultRenderHeight is the default height when no height is available
	defaultRenderHeight = 3
)
