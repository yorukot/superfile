package sidebar

import (
	"time"

	"charm.land/bubbles/v2/textinput"
)

type directory struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Section  string `json:"-"`
	Icon     string `json:"-"`
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

	// loadedQuery is the search query `directories` was built for
	loadedQuery string
	lastRefresh time.Time
}
