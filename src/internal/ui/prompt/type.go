package prompt

import "github.com/charmbracelet/bubbles/textinput"

type PromptCommand struct {
	renderPrefix  string
	renderHint    string
	handleCommand func(input string) bool
}
type PromptCommandPrefix = string


type PromptModal struct {
	headline  string
	open      bool
	cmd       PromptCommand
	textInput textinput.Model
	errormsg  string

	commandList map[PromptCommandPrefix]PromptCommand
}