package prompt

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type PromptModal struct {
	headline string
	open     bool
	// whether its shellMode or spfMode
	shellMode         bool
	textInput         textinput.Model
	errorMsg          string
	commands          []promptCommand
	spfPromptHotkey   string
	shellPromptHotkey string
}

// This is only used to render suggestions
// Should not be exported
type promptCommand struct {
	command     string
	usage       string
	description string
}
