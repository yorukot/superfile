package preview

// UpdateMsg represents an async query result

type UpdateMsg struct {
	location string

	// preview panel's content needs to be in sync with its width/height
	// you cannot update width/height without updating the content
	content string
	width   int
	height  int
	reqID   int
}

func NewUpdateMsg(location string, content string, width int, height int, reqID int) UpdateMsg {
	return UpdateMsg{
		location: location,
		content:  content,
		width:    width,
		height:   height,
		reqID:    reqID,
	}
}

func (msg UpdateMsg) GetReqID() int {
	return msg.reqID
}

func (m *Model) Apply(msg UpdateMsg) {
	m.width = msg.width
	m.height = msg.height
	m.content = msg.content
	m.location = msg.location
	m.loading = false
}

func (msg UpdateMsg) GetLocation() string {
	return msg.location
}
