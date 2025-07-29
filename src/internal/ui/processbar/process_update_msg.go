package processbar

type updateMsg interface {
	Apply(m *Model) error
	GetReqID() int
}

// TODO: Can we remove this duplication with model_msg ?
type BaseMsg struct {
	reqID int
}

func (msg BaseMsg) GetReqID() int {
	return msg.reqID
}

type newProcessMsg struct {
	BaseMsg
	NewProcess Process
}

func (msg newProcessMsg) Apply(m *Model) error {
	return m.AddProcess(msg.NewProcess)
}

type updateProcessMsg struct {
	BaseMsg
	NewProcess Process
}

func (msg updateProcessMsg) Apply(m *Model) error {
	return m.UpdateExistingProcess(msg.NewProcess)
}

// Construction will be options UpdateName(), UpdateDone(), etc..

type stopListeningMsg struct {
	BaseMsg
}

func (msg stopListeningMsg) Apply(_ *Model) error {
	return nil
}
