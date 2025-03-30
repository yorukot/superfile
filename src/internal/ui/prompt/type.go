package prompt

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/yorukot/superfile/src/internal/common"
)


type PromptCommandPrefix = string


type PromptCommand struct {
	renderPrefix  string
	renderHint    string
	handleCommand func(input string) (common.ModelUpdateAction, error)
}


type PromptModal struct {
	headline  string
	open      bool
	textInput textinput.Model
	errormsg  string
	
	// Defined here for decoupling between internal and prompt package
	// We will later take hotkeys to the common pacakge 
	confirmHotkeys []string 
	cancelHotkeys []string

	commandList map[PromptCommandPrefix]PromptCommand
}