package common

/* To work with PromptModal */
type ActionType = int

// Constants for actions
const (
	SplitPanel ActionType = iota
	RunShellCommand
)

type ModelUpdateAction struct {
	action ActionType
	args   []string
}
