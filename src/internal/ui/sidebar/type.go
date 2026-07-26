package sidebar

import "charm.land/bubbles/v2/textinput"

// deprecated. Replace it by public struct SideBarItem
type directory struct {
	Location string  `json:"location"`
	Name     string  `json:"name"`
	Section  string  `json:"-"`
	FreeSize *uint64 `json:"-"` // nil if size not calculated
}

type Model struct {
	directories []directory
	renderIndex int
	cursor      int
	rename      textinput.Model
	renaming    bool
	searchBar   textinput.Model
	pinnedMgr   *PinnedManager
	width       int
	height      int
	disabled    bool
	sections    []string
}

type SideBarItem struct {
	Location string
	Name     string
	Type     string
}
