package processbar

import (
	"log/slog"
)

// Very hard to do this via tea.Cmd, where we want to send progress updates in the middle of
// io operatins. Like 100 update messages in during a copy.
func (m *Model) ListenForUpdates() {
	m.isListening = true
	// A goroutine running forever
	go func() {
		for {
			msg := <-m.msgChan
			slog.Debug("Received message", "id", msg.GetReqID())
			if _, ok := msg.(stopListeningMsg); ok {
				m.isListening = false
				return
			}
			err := msg.Apply(m)
			if err != nil {
				slog.Error("Could not apply update to processbar", "error", err)
			}
			// TODO: We could consider adding a way to gracefully stop
			// the goroutine (context cancellation or stop channel)
		}
	}()
}

func (m *Model) IsListeningForUpdates() bool {
	return m.isListening
}

// Might add options to drain the channel in case msg is high priority.
func (m *Model) trySendMsgToChannel(msg updateMsg) error {
	select {
	case m.msgChan <- msg:
		return nil
	default:
		// Process queue full with messages. Cannot add new process
		return &ProcessChannelFullError{}
	}
}

// Block till message is sent
func (m *Model) sendMsgToChannelBlocking(msg updateMsg) {
	m.msgChan <- msg
}

func (m *Model) sendMsgToChannel(msg updateMsg, blocking bool) error {
	if blocking {
		m.sendMsgToChannelBlocking(msg)
		return nil
	}
	return m.trySendMsgToChannel(msg)
}

func (m *Model) SendAddProcessMsg(name string, total int, blockingSend bool) (Process, error) {
	id := m.newUUIDForProcess()
	p := NewProcess(id, name, total)
	msg := newProcessMsg{
		NewProcess: p,
		BaseMsg:    BaseMsg{reqID: m.newReqCnt()},
	}
	err := m.sendMsgToChannel(msg, blockingSend)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (m *Model) SendUpdateProcessNameMsg(p Process, blockingSend bool) error {
	msg := updateProcessMsg{NewProcess: p, BaseMsg: BaseMsg{reqID: m.newReqCnt()}}
	return m.sendMsgToChannel(msg, blockingSend)
}

// Non Blocking and can fail
func (m *Model) TrySendingUpdateProcessNameMsg(p Process) {
	msg := updateProcessMsg{NewProcess: p, BaseMsg: BaseMsg{reqID: m.newReqCnt()}}
	err := m.sendMsgToChannel(msg, false)
	if err != nil {
		slog.Error("Failed to send message to channel", "reqID", msg.GetReqID(), "error", err)
	}
}

func (m *Model) SendStopListeningMsgBlocking() {
	m.sendMsgToChannelBlocking(stopListeningMsg{BaseMsg: BaseMsg{reqID: m.newReqCnt()}})
}
