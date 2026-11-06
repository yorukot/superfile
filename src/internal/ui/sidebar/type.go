package sidebar

import (
	"charm.land/bubbles/v2/textinput"

	"github.com/yorukot/superfile/src/internal/disks"
)

type directory struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Section  string `json:"-"`
	Icon     string `json:"-"`
}

type Model struct {
	diskManager disks.Manager
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
