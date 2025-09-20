package sidebar

import "github.com/charmbracelet/bubbles/textinput"

type directory struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

type Model struct {
	directories []directory
	renderIndex int
	cursor      int
	rename      textinput.Model
	renaming    bool
	searchBar   textinput.Model
	pinnedMgr   *PinnedManager
}
