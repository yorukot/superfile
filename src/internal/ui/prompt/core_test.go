package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/common/utils"
	"testing"
)

var setupDone = false

// Initialize the globals we need for testing
func initGlobals() {
	// Maybe we could use doOnce go constructs
	if setupDone {
		return
	}
	common.Hotkeys.ConfirmTyping = []string{"enter"}
	common.Hotkeys.ConfirmTyping = []string{"ctrl+c", "esc"}
	setupDone = true
}

func TestModel_HandleUpdate(t *testing.T) {
	// We want to test
	// 1. Handle update called on closed Model
	// 2. Pressing confirm on empty input
	// 3. Three conditions of getPromptAction -
	// 4. Cancel typing
	// 5. Switching between shell and SPF mode
	// 6. Updating text input with text - single and multi char string
	// 7. No keyMesage like cursor.BlinkMsg, cursor.blinkCanceled
	// Validate blink m.textInput.Cursor.Blink
	// Dont test getPromptAction here

	testdata := []struct {
		m              Model
		input          []tea.Msg
		initialization func(Model)
		validator      func(Model) bool
		expectedAction []common.ModelAction
		name           string
	}{
		{
			// Default Model is closed
			m:     GenerateModel(spfPromptChar, shellPromptChar, true),
			input: []tea.Msg{utils.TeaRuneKeyMsg("x")},
			initialization: func(m Model) {
				m.Open(true)
			},
			validator: func(m Model) bool {
				return true
			},
			expectedAction: []common.ModelAction{common.NoAction{}},
			name:           "Handle update called on closed Model",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialization(tt.m)
			for i, msg := range tt.input {
				// Todo : Replace cwd with this file's directory
				action, _ := tt.m.HandleUpdate(msg, "/")
				assert.Equal(t, tt.expectedAction[i], action)
			}
			assert.True(t, tt.validator(tt.m))
		})
	}
}
