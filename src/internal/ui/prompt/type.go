package prompt

import (
	"github.com/charmbracelet/bubbles/textinput"
)

type PromptModal struct {
	headline string
	open     bool
	// whether its shellMode or spfMode
	shellMode bool
	textInput textinput.Model
	errorMsg  string
}
