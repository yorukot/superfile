package find

import (
	"charm.land/bubbles/v2/textinput"
)

type FindResult struct {
	Path  string
	IsDir bool
}

// No need to name it as FindModel. It will me imported as find.Model
type Model struct {

	// Configuration
	headline  string
	searchDir string

	// State
	open        bool
	justOpened  bool // Flag to ignore the opening keystroke
	fdFound     bool
	textInput   textinput.Model
	results     []FindResult
	errMsg      string
	cursor      int // Index of currently selected result for keyboard navigation
	renderIndex int // Index of first visible result in scrollable list

	// Dimensions - Exported, since model will be dynamically adjusting them
	width int
	// Height is dynamically adjusted based on content
	maxHeight int

	// Request tracking for async queries
	reqCnt int
}

// UpdateMsg represents an async query result
type UpdateMsg struct {
	query   string
	results []FindResult
	errMsg  string
	reqID   int
}

func NewUpdateMsg(query string, results []FindResult, errMsg string, reqID int) UpdateMsg {
	return UpdateMsg{
		query:   query,
		results: results,
		errMsg:  errMsg,
		reqID:   reqID,
	}
}

func (msg UpdateMsg) GetReqID() int {
	return msg.reqID
}
