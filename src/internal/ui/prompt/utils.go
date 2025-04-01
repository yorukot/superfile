package prompt

import "github.com/yorukot/superfile/src/internal/common"

func (p *PromptModal) IsOpen() bool {
	return p.open
}

func (p *PromptModal) Validate() bool {
	// Prompt was closed, but textInput was not cleared
	if !p.open && p.textInput.Value() != "" {
		return false
	}
	return true
}

func modeString(shellMode bool) string {
	if shellMode {
		return "(Shell Mode)"
	}
	return "(Prompt Mode)"
}

func shellPrompt(shellMode bool) string {
	if shellMode {
		return ":"
	}
	return ">"
}

func getPromptAction(shellMode bool, value string) common.PromptAction {
	if shellMode {
		return common.PromptAction{
			Action: common.ShellCommandAction,
			Args:   []string{value},
		}
	}
	return common.NoPromptAction()
}
