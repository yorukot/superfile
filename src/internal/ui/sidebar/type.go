package sidebar

import "github.com/charmbracelet/bubbles/textinput"

type Directory struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

type Model struct {
	directories []Directory
	renderIndex int
	cursor      int
	rename      textinput.Model
	renaming    bool
	searchBar   textinput.Model
}
