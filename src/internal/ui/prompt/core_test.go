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

type modelValidator func(*testing.T, Model)

type update struct {
	input          tea.Msg
	validator      modelValidator
	expectedAction common.ModelAction
}

// Validator
func alwaysOk(_ *testing.T, m Model) {
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
		name           string
		m              Model
		initialization func(*Model)

		updates []update
	}{
		{
			// Default Model is closed
			name: "Handle update called on closed Model",
			m:    GenerateModel(spfPromptChar, shellPromptChar, true),

			initialization: func(m *Model) {},
			updates: []update{
				{
					input: utils.TeaRuneKeyMsg("x"),
					validator: func(t *testing.T, m Model) {
						assert.False(t, m.IsOpen())
					},
					expectedAction: common.NoAction{},
				},
			},
		},
		{
			// Default Model is closed
			name: "Pressing confirm on empty input",
			m:    GenerateModel(spfPromptChar, shellPromptChar, true),

			initialization: func(m *Model) {
				m.Open(true)
			},
			updates: []update{
				{
					input: tea.KeyMsg{Type: tea.KeyEnter},
					validator: func(t *testing.T, m Model) {
						assert.False(t, m.IsOpen())
					},
					expectedAction: common.NoAction{},
				},
			},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialization(&tt.m)
			for _, u := range tt.updates {
				// Todo : Replace cwd with this file's directory
				action, _ := tt.m.HandleUpdate(u.input, "/")
				assert.Equal(t, u.expectedAction, action)
				u.validator(t, tt.m)
			}

		})
	}
}
