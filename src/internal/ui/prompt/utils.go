package prompt

import (
	"strings"
)

// This is to generate error objects that can pe nicely printed to UI
type InvalidCmdError struct {
	uiMsg        string
	wrappedError error
}

func (e InvalidCmdError) Error() string {
	if e.wrappedError == nil {
		return e.uiMsg
	}
	return e.wrappedError.Error()
}

func (e InvalidCmdError) Unwrap() error {
	return e.wrappedError
}

func (e InvalidCmdError) UIMessage() string {
	return e.uiMsg
}

func (p *Model) IsOpen() bool {
	return p.open
}

func (p *Model) Validate() bool {
	// Prompt was closed, but textInput was not cleared
	if !p.open && p.textInput.Value() != "" {
		return false
	}
	return true
}

func modeString(shellMode bool) string {
	if shellMode {
		return shellModeString
	}
	return spfModeString
}

func shellPrompt(shellMode bool) string {
	if shellMode {
		return shellPromptChar
	}
	return spfPromptChar
}

// Only allocates memory proportional to first token's size
func getFirstToken(command string) string {
	spaceIndex := strings.IndexByte(command, ' ')
	if spaceIndex == -1 {
		return command
	}
	return command[:spaceIndex]
}

func defaultCommandSlice() []promptCommand {
	return []promptCommand{
		{
			command:     openCommand,
			usage:       openCommand + " <PATH>",
			description: "Open a new panel at a specified path",
		},
		{
			command:     splitCommand,
			usage:       splitCommand,
			description: "Open a new panel at a current file panel's path",
		},
		{
			command:     cdCommand,
			usage:       cdCommand + " <PATH>",
			description: "Change directory of current panel",
		},
	}
}
