package internal

// Todo , merge this and predefined variables file
// These are effectively consts
// Had to use `var` as go doesn't allows const structs
var pinnedDividerDir = directory{ //nolint: gochecknoglobals // This is more like a const.
	name:     "",
	location: "Pinned+-*/=?",
}

var diskDividerDir = directory{ //nolint: gochecknoglobals // This is more like a const.
	name:     "",
	location: "Disks+-*/=?",
}

// superfile logo + blank line + search bar
const sideBarInitialHeight = 3

const invalidTypeString = "InvalidType"
