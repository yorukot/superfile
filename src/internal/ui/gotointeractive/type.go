package gotointeractive

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	headline    string
	open        bool
	textInput   textinput.Model
	results     []string
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
	results []string
	reqID   int
}

func NewUpdateMsg(query string, results []string, reqID int) UpdateMsg {
	return UpdateMsg{
		query:   query,
		results: results,
		reqID:   reqID,
	}
}

func (msg UpdateMsg) GetReqID() int {
	return msg.reqID
}
