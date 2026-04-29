package helpmenu

import "charm.land/bubbles/v2/textinput"

type hotkeyType int

const (
	globalType hotkeyType = iota
	normalType
	selectType
)

// Modal
type Model struct {
	height       int
	width        int
	opened       bool
	renderIndex  int
	cursor       int
	data         []hotkeydata
	filteredData []hotkeydata
	searchBar    textinput.Model
}

type hotkeydata struct {
	hotkey         []string
	description    string
	hotkeyWorkType hotkeyType
	subTitle       string
}
