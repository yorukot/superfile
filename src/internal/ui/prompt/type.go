package prompt

import "github.com/charmbracelet/bubbles/textinput"

// No need to name it as PromptModel. It will me imported as prompt.Model
type Model struct {
	headline          string
	commands          []promptCommand
	spfPromptHotkey   string
	shellPromptHotkey string

	open bool
	// whether its shellMode or spfMode
	shellMode bool
	textInput textinput.Model
	resultMsg string

	// Whether the user intended action was successful
	actionSuccess bool
}

// This is only used to render suggestions
// Should not be exported
type promptCommand struct {
	command     string
	usage       string
	description string
}
