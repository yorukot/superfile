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
}

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
