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
	reqCnt      int

	width     int
	maxHeight int
}

type UpdateMsg struct {
	query   string
	results []Directory
	reqID   int
}

func NewUpdateMsg(query string, results []Directory, reqID int) UpdateMsg {
	return UpdateMsg{
		query:   query,
		results: results,
		reqID:   reqID,
	}
}
