package sidebar

// Todo , merge this and predefined variables file
// These are effectively consts
// Had to use `var` as go doesn't allows const structs
var PinnedDividerDir = Directory{ //nolint: gochecknoglobals // This is more like a const.
	Name:     "",
	Location: "Pinned+-*/=?",
}

var DiskDividerDir = Directory{ //nolint: gochecknoglobals // This is more like a const.
	Name:     "",
	Location: "Disks+-*/=?",
}

// superfile logo + blank line + search bar
const SideBarInitialHeight = 3
