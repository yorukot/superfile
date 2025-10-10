package processbar

import (
	"log/slog"
)

// Only used in tests, to have processbar used in a standalone way without model
func (m *Model) ListenForChannelUpdates() {
	// A goroutine running forever
	go func() {
		for {
			msg, ok := <-m.msgChan
			if !ok {
				slog.Debug("Channel closed, stopping listener")
				return
			}
			slog.Debug("Received message", "id", msg.GetReqID())
			if _, ok := msg.(stopListeningMsg); ok {
				return
			}
			_, err := msg.Apply(m)
			if err != nil {
				slog.Error("Could not apply update to processbar", "error", err)
			}
		}
	}()
}

// An IO Operation, that will wait forever on msgChannel
func (m *Model) GetListenCmd() Cmd {
	return func() UpdateMsg {
		return <-m.msgChan
	}
}

// Might add options to drain the channel in case msg is high priority.
func (m *Model) trySendMsgToChannel(msg UpdateMsg) error {
	select {
	case m.msgChan <- msg:
		return nil
	default:
		// Process queue full with messages. Cannot add new process
		return &ProcessChannelFullError{}
	}
}

// Block till message is sent
func (m *Model) sendMsgToChannelBlocking(msg UpdateMsg) {
	m.msgChan <- msg
}

func (m *Model) sendMsgToChannel(msg UpdateMsg, blocking bool) error {
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
		// Return zero-value process to indicate failure
		return Process{}, err
	}
	return p, nil
}

func (m *Model) SendUpdateProcessMsg(p Process, blockingSend bool) error {
	msg := updateProcessMsg{NewProcess: p, BaseMsg: BaseMsg{reqID: m.newReqCnt()}}
	return m.sendMsgToChannel(msg, blockingSend)
}

// Non Blocking and can fail
func (m *Model) TrySendingUpdateProcessMsg(p Process) {
	msg := updateProcessMsg{NewProcess: p, BaseMsg: BaseMsg{reqID: m.newReqCnt()}}
	err := m.sendMsgToChannel(msg, false)
	if err != nil {
		slog.Error("Failed to send message to channel", "reqID", msg.GetReqID(), "error", err)
	}
}

func (m *Model) SendStopListeningMsgBlocking() {
	m.sendMsgToChannelBlocking(stopListeningMsg{BaseMsg: BaseMsg{reqID: m.newReqCnt()}})
}
