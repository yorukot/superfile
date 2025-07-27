package internal

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yorukot/superfile/src/internal/utils"
)

type IgnorerWriter struct{}

func (w IgnorerWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

type TeaProg struct {
	m    *model
	prog *tea.Program
}

// If you use this, make sure to handle cleanup
func NewTeaProg(m *model, eventLoop bool) *TeaProg {
	p := &TeaProg{m: m, prog: tea.NewProgram(m, tea.WithInput(nil), tea.WithOutput(IgnorerWriter{}))}
	if eventLoop {
		p.StartEventLoop()
	}
	return p
}

func (p *TeaProg) GetModel() *model {
	return p.m
}

func NewTestTeaProgWithEventLoop(t *testing.T, m *model) *TeaProg {
	p := NewTeaProg(m, true)
	t.Cleanup(func() {
		p.Close()
	})
	return p
}

func (p *TeaProg) StartEventLoop() {
	go p.prog.Run()
	// Send nil to block for start of event loop
	p.prog.Send(nil)
}

func (p *TeaProg) Send(msgs ...tea.Msg) {
	for _, msg := range msgs {
		p.prog.Send(msg)
	}
}

func (p *TeaProg) SendKey(key string) {
	p.Send(utils.TeaRuneKeyMsg(key))
}

// Dont use eventloop and dont care about the tea.Cmd returned by Update()
func (p *TeaProg) SendDirectly(msgs ...tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(msgs))
	for i, msg := range msgs {
		var retModel tea.Model
		retModel, cmds[i] = p.m.Update(msg)
		if m, ok := retModel.(*model); ok {
			p.m = m
		} else {
			// This should never happen as we return *model on Update()
			panic("model is not of type *model")
		}
	}

	return tea.Batch(cmds...)
}

func (p *TeaProg) SendKeyDirectly(key string) tea.Cmd {
	return p.SendDirectly(utils.TeaRuneKeyMsg(key))
}

func (p *TeaProg) Close() {
	p.prog.Quit()
}
