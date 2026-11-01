package prompt

import "charm.land/bubbles/v2/textinput"

// No need to name it as PromptModel. It will me imported as prompt.Model
type Model struct {

	// Configuration
	headline          string
	commands          []promptCommand
	spfPromptHotkey   string
	shellPromptHotkey string
	closeOnSuccess    bool

	// State
	open bool
	// whether its shellMode or spfMode
	// Always use setShellMode to adjust
	shellMode bool
	textInput textinput.Model
	resultMsg string

	// Tab completion state
	completions     []string // current completion candidates
	completionIdx   int      // index into completions for cycling
	completionDone  bool     // set when user has confirmed a completion
	cwd             string   // current working directory for completions

	// Whether the user intended action was successful
	actionSuccess bool

	// Dimensions - Exported, since model will be dynamically adjusting them
	width int
	// Height is dynamically adjusted based on content
	maxHeight int
}

// This is only used to render suggestions
// Should not be exported
type promptCommand struct {
	command     string
	usage       string
	description string
}
