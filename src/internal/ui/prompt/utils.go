package prompt

import "strings"

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
