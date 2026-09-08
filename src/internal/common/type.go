package common

import "fmt"

// Placeholder inteface for now, might later move 'model' type to commons and have
// and add an execute(model) function to this
type ModelAction interface {
	String() string
}

type NoAction struct {
}

func (n NoAction) String() string {
	return "NoAction"
}

type ShellCommandAction struct {
	Command     string
	Interactive bool // if true, use tea.ExecProcess for full terminal access
}

func (s ShellCommandAction) String() string {
	mode := "quick"
	if s.Interactive {
		mode = "interactive"
	}
	return fmt.Sprintf("ShellCommandAction[%s] for command %s", mode, s.Command)
}

// We could later move 'model' type to commons and have
// these actions implement an execute(model) interface
type SplitPanelAction struct{}

func (s SplitPanelAction) String() string {
	return "SplitPanelAction"
}

type CDCurrentPanelAction struct {
	Location string
}

func (c CDCurrentPanelAction) String() string {
	return "CDCurrentPanelAction to " + c.Location
}

type OpenPanelAction struct {
	Location string
}

func (o OpenPanelAction) String() string {
	return "OpenPanelAction at " + o.Location
}
