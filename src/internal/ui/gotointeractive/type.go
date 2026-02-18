package gotointeractive

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Result struct {
	Name  string
	IsDir bool
}

type Model struct {
	headline    string
	open        bool
	textInput   textinput.Model
	results     []Result
	cursor      int
	renderIndex int
	currentPath string

	width      int
	maxHeight  int
	justOpened bool
	reqCnt     int
}

type UpdateMsg struct {
	query   string
	results []Result
	reqID   int
	path    string
}

func NewUpdateMsg(query string, results []Result, reqID int, path string) UpdateMsg {
	return UpdateMsg{
		query:   query,
		results: results,
		reqID:   reqID,
		path:    path,
	}
}

func (msg UpdateMsg) GetReqID() int {
	return msg.reqID
}
