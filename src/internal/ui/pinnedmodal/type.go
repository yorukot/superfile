package pinnedmodal

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Directory struct {
	Location string
	Name     string
}

type Model struct {
	headline string

	open        bool
	justOpened  bool
	textInput   textinput.Model
	results     []Directory
	allDirs     []Directory
	cursor      int
	renderIndex int

	width     int
	maxHeight int
}

type UpdateMsg struct {
	query   string
	results []Directory
}

func NewUpdateMsg(query string, results []Directory) UpdateMsg {
	return UpdateMsg{
		query:   query,
		results: results,
	}
}
