package prompt

import (
	"fmt"
	"strings"

	"github.com/yorukot/superfile/src/internal/common"
)

func getPromptAction(shellMode bool, value string, cwdLocation string) (common.ModelAction, error) {
	noAction := common.NoAction{}
	if value == "" {
		return noAction, nil
	}
	if shellMode {
		return common.ShellCommandAction{
			Command: value,
		}, nil
	}

	promptArgs, err := tokenizePromptCommand(value, cwdLocation)
	if err != nil {
		return noAction, invalidCmdError{
			uiMsg:        tokenizationError + " : " + err.Error(),
			wrappedError: fmt.Errorf("error during tokenization : %w", err),
		}
	}

	switch promptArgs[0] {
	case "split":
		if len(promptArgs) != 1 {
			return noAction, invalidCmdError{
				uiMsg: splitCommandArgError,
			}
		}
		return common.SplitPanelAction{}, nil
	case "cd":
		if len(promptArgs) != ExpectedArgCount {
			return noAction, invalidCmdError{
				uiMsg: fmt.Sprintf("cd command needs exactly one argument, received %d",
					len(promptArgs)-1),
			}
		}
		return common.CDCurrentPanelAction{
			Location: promptArgs[1],
		}, nil
	case "open":
		if len(promptArgs) != ExpectedArgCount {
			return noAction, invalidCmdError{
				uiMsg: fmt.Sprintf("open command needs exactly one argument, received %d",
					len(promptArgs)-1),
			}
		}
		return common.OpenPanelAction{
			Location: promptArgs[1],
		}, nil

	default:
		return noAction, invalidCmdError{
			uiMsg: "Invalid spf command : " + promptArgs[0],
		}
	}
}

// Only allocates memory proportional to first token's size
// Only works for space right now. Does not splits command based on
// \n or \t , etc
func getFirstToken(command string) string {
	command = strings.TrimSpace(command)
	spaceIndex := strings.IndexByte(command, ' ')
	if spaceIndex == -1 {
		return command
	}
	return command[:spaceIndex]
}
