package common

// ActionType : To work with PromptModal
type ActionType = int

// Constants for actions
// Todo : Shouldn't we be using inheritance here ?
const (
	NoAction ActionType = iota
	ShellCommandAction
	SplitPanelAction
	CDCurrentPanelAction
	OpenPanelAction
)

type PromptAction struct {
	Action ActionType
	Args   []string
}
