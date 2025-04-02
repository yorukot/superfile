package prompt

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/common/utils"
	"log/slog"
	"os"
	"testing"
)

var setupDone = false

// Initialize the globals we need for testing
func initGlobals(initLogging bool) {
	// Maybe we could use doOnce go constructs
	if setupDone {
		return
	}
	if initLogging {
		slog.SetDefault(slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}
	common.Hotkeys.ConfirmTyping = []string{"enter"}
	common.Hotkeys.CancelTyping = []string{"ctrl+c", "esc"}
	setupDone = true
}

func TestModel_HandleUpdate(t *testing.T) {
	const defaultCwd = "/"
	// Could take this as an test arguement
	initGlobals(true)
	// We want to test
	// 1. Handle update called on closed Model
	// 2. Pressing confirm on empty input
	// 3. Three conditions of getPromptAction
	// 4. Cancel typing
	// 5. Switching between shell and SPF mode
	// 6. Updating text input with text - single and multi char string
	// 7. No keyMesage like cursor.BlinkMsg, cursor.blinkCanceled
	// Validate blink m.textInput.Cursor.Blink
	// Dont test getPromptAction here

	t.Run("Handle update called on closed Model", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, true)
		action, _ := m.HandleUpdate(utils.TeaRuneKeyMsg("x"), "/")
		assert.False(t, m.IsOpen())
		assert.Equal(t, common.NoAction{}, action)
	})

	t.Run("Pressing confirm on empty input", func(t *testing.T) {

		actualTest := func(closeOnSuccess bool, openAfterEnter bool) {
			m := GenerateModel(spfPromptChar, shellPromptChar, closeOnSuccess)
			m.Open(true)
			assert.True(t, m.IsOpen())

			action, _ := m.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter}, defaultCwd)
			assert.Equal(t, openAfterEnter, m.IsOpen())
			assert.Equal(t, common.NoAction{}, action)
			assert.Equal(t, "", m.resultMsg)
			assert.Equal(t, true, m.actionSuccess)
		}

		actualTest(true, false)
		actualTest(false, true)

	})

	t.Run("Validate Prompt Actions", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, true)
		m.Open(false)

		action, _ := m.HandleUpdate(utils.TeaRuneKeyMsg(splitCommand), defaultCwd)
		assert.Equal(t, common.NoAction{}, action)

		action, _ = m.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter}, defaultCwd)
		assert.Equal(t, common.SplitPanelAction{}, action)

		action, _ = m.HandleUpdate(utils.TeaRuneKeyMsg("bad_command"), defaultCwd)
		action, _ = m.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter}, defaultCwd)
		assert.Equal(t, common.NoAction{}, action)
		assert.False(t, m.actionSuccess)
		assert.NotEmpty(t, m.resultMsg)

		m.shellMode = true
		command := "abc def /xyz"
		action, _ = m.HandleUpdate(utils.TeaRuneKeyMsg(command), defaultCwd)
		action, _ = m.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter}, defaultCwd)
		assert.Equal(t, common.ShellCommandAction{Command: command}, action)

		// Todo : test third error condition .

	})
}
