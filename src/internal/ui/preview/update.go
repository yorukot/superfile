package preview

// UpdateMsg represents an async query result

type UpdateMsg struct {
	// location can contain either the path of current content's file
	// or path of file whose preview request is in flight.
	// It should not have past data
	location string

	// preview panel's content needs to be in sync with its width/height
	// you cannot update width/height without updating the content
	content       string
	contentWidth  int
	contentHeight int
	reqID         int

	// rawTransmit is non-empty for Kitty protocol images. It contains
	// APC escape sequences that must be sent directly to the terminal
	// via tea.Raw(), bypassing the cell-based renderer.
	rawTransmit string
}

func NewUpdateMsg(location string, content string, rawTransmit string, width int, height int, reqID int) UpdateMsg {
	return UpdateMsg{
		location:      location,
		content:       content,
		rawTransmit:   rawTransmit,
		contentWidth:  width,
		contentHeight: height,
		reqID:         reqID,
	}
}

func (msg UpdateMsg) GetReqID() int {
	return msg.reqID
}

func (m *Model) Apply(msg UpdateMsg) {
	m.setContent(msg.content, msg.contentWidth, msg.contentHeight, msg.location)
}

func (msg UpdateMsg) GetLocation() string {
	return msg.location
}

func (msg UpdateMsg) GetContentWidth() int {
	return msg.contentWidth
}

func (msg UpdateMsg) GetContentHeight() int {
	return msg.contentHeight
}

func (msg UpdateMsg) GetRawTransmit() string {
	return msg.rawTransmit
}
