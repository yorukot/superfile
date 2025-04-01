package common

// ActionType : To work with PromptModal
type ActionType = int

// Constants for actions
const (
	NoAction ActionType = iota
	ShellCommandAction
	SplitPanelAction
)

type PromptAction struct {
	Action ActionType
	Args   []string
}
