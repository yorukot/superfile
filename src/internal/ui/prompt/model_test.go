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
	// Could take this as a test argument
	initGlobals(true)
	// We want to test
	// 1. Handle update called on closed Model
	// 2. Pressing confirm on empty input
	// 3. Three conditions of getPromptAction
	// 4. Cancel typing
	// 5. Switching between shell and SPF mode
	// 6. Updating text input with text - Tested in above.
	// 7. keyMsg like cursor.BlinkMsg, cursor.blinkCanceled
	// Validate blink m.textInput.Cursor.Blink
	// Dont test getPromptAction here. It will be a separate test

	t.Run("Handle update called on closed Model", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, true)
		action, _ := m.HandleUpdate(utils.TeaRuneKeyMsg("x"), defaultCwd)
		assert.Empty(t, m.textInput.Value())
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

		// Todo : test third error condition
		// Right now, I dont know what could cause that unexpected error

	})

	t.Run("Validate Cancel typing", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, true)

		actualTest := func(closeKey tea.KeyMsg, shouldBeOpen bool) {
			m.Open(true)
			action, _ := m.HandleUpdate(utils.TeaRuneKeyMsg("xyz"), defaultCwd)
			action, _ = m.HandleUpdate(closeKey, defaultCwd)
			assert.Equal(t, common.NoAction{}, action)
			assert.Equal(t, shouldBeOpen, m.IsOpen())
		}

		actualTest(tea.KeyMsg{Type: tea.KeyCtrlC}, false)
		actualTest(tea.KeyMsg{Type: tea.KeyEscape}, false)
		actualTest(tea.KeyMsg{Type: tea.KeyCtrlD}, true)

	})

	t.Run("Switching between shell and SPF mode", func(t *testing.T) {

		actualTest := func(promptChar string, shellChar string) {
			m := GenerateModel(promptChar, shellChar, true)
			m.Open(true)
			assert.True(t, m.shellMode)

			// Shell to prompt
			action, _ := m.HandleUpdate(utils.TeaRuneKeyMsg(promptChar), defaultCwd)
			assert.False(t, m.shellMode)
			assert.True(t, m.actionSuccess)
			assert.Equal(t, common.NoAction{}, action)

			// Prompt to shell
			action, _ = m.HandleUpdate(utils.TeaRuneKeyMsg(shellChar), defaultCwd)
			assert.True(t, m.shellMode)
			assert.True(t, m.actionSuccess)
			assert.Equal(t, common.NoAction{}, action)

			// Pressing shellChar when you are already on shell shouldn't to anything
			action, _ = m.HandleUpdate(utils.TeaRuneKeyMsg(shellChar), defaultCwd)
			assert.True(t, m.shellMode)
			assert.True(t, m.actionSuccess)
			assert.Equal(t, common.NoAction{}, action)
		}
		actualTest(">", ":")
		actualTest("$", "#")
	})

	t.Run("Validate Cursor Blink update", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, true)
		m.Open(true)
		assert.False(t, m.textInput.Cursor.Blink)

		blinkMsg := m.textInput.Cursor.BlinkCmd()()
		action, _ := m.HandleUpdate(blinkMsg, defaultCwd)
		assert.Equal(t, common.NoAction{}, action)
		assert.True(t, m.textInput.Cursor.Blink)

		blinkMsg = m.textInput.Cursor.BlinkCmd()()
		action, _ = m.HandleUpdate(blinkMsg, defaultCwd)
		assert.Equal(t, common.NoAction{}, action)
		assert.False(t, m.textInput.Cursor.Blink)

		blinkMsg = m.textInput.Cursor.BlinkCmd()()
		action, _ = m.HandleUpdate(blinkMsg, defaultCwd)
		assert.Equal(t, common.NoAction{}, action)
		assert.True(t, m.textInput.Cursor.Blink)

		// We could test BlinkCancelled and initialBlink as well, but that's too much for now
	})
}

func TestMode_HandleResults(t *testing.T) {
	t.Run("Verify Shell results update", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, true)
		m.Open(true)
		m.HandleShellCommandResults(0, "")

		// Validate close happens when closeOnSuccess is true
		assert.True(t, m.actionSuccess)
		assert.Equal(t, m.resultMsg, "Command exited with status 0")
		assert.False(t, m.IsOpen())

		m.Open(true)
		m.HandleShellCommandResults(1, "")
		assert.False(t, m.actionSuccess)
		assert.Equal(t, m.resultMsg, "Command exited with status 1")
		assert.True(t, m.IsOpen())

		m.closeOnSuccess = false
		m.HandleShellCommandResults(0, "")
		// Validate that close does not happen when closeOnSuccess is true
		assert.True(t, m.actionSuccess)
		assert.Equal(t, m.resultMsg, "Command exited with status 0")
		assert.True(t, m.IsOpen())

	})

	t.Run("Verify SPF results update", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, false)
		m.Open(true)
		msg := "Test message"
		m.HandleSPFActionResults(true, msg)

		assert.True(t, m.actionSuccess)
		assert.Equal(t, msg, m.resultMsg)
		assert.True(t, m.IsOpen())

		m.closeOnSuccess = true
		// Validate close happens when closeOnSuccess is true
		m.HandleSPFActionResults(true, "")
		assert.False(t, m.IsOpen())

	})
}
