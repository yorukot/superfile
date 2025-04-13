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
