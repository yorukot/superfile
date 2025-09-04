package zoxide

import (
	"github.com/charmbracelet/bubbles/textinput"
	zoxidelib "github.com/lazysegtree/go-zoxide"
)

// No need to name it as ZoxideModel. It will me imported as zoxide.Model
type Model struct {

	// Configuration
	headline string
	zClient  *zoxidelib.Client

	// State
	open        bool
	justOpened  bool // Flag to ignore the opening keystroke
	textInput   textinput.Model
	results     []zoxidelib.Result
	cursor      int // Index of currently selected result for keyboard navigation
	renderIndex int // Index of first visible result in scrollable list

	// Dimensions - Exported, since model will be dynamically adjusting them
	width int
	// Height is dynamically adjusted based on content
	maxHeight int
}
