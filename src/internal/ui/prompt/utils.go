package prompt

import (
	"strings"
)

// This is to generate error objects that can pe nicely printed to UI
type invalidCmdError struct {
	uiMsg        string
	wrappedError error
}

func (e invalidCmdError) Error() string {
	if e.wrappedError == nil {
		return e.uiMsg
	}
	return e.wrappedError.Error()
}

func (e invalidCmdError) Unwrap() error {
	return e.wrappedError
}

func (e invalidCmdError) uiMessage() string {
	return e.uiMsg
}

func (m *Model) Open(shellMode bool) {
	m.open = true
	m.shellMode = shellMode
	_ = m.textInput.Focus()
}

func (m *Model) Close() {
	m.open = false
	m.shellMode = true
	m.textInput.SetValue("")
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) validate() bool {
	// Prompt was closed, but textInput was not cleared
	if !m.open && m.textInput.Value() != "" {
		return false
	}
	return true
}

func (m *Model) CloseOnSuccessIfNeeded() {
	if m.closeOnSuccess && m.actionSuccess {
		m.Close()
	}
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
