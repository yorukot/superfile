package prompt

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// Initialize the globals we need for testing
func initGlobals() {
	// Updating globals for test is not a good idea and can lead to all sorts of issue
	// When multiple tests depend on same global variable and want different values
	// Since this is config that would likely stay same, maybe this is okay.
	// Also, this is done in main model's test too.
	// We need to find a better way to do this
	err := common.PopulateGlobalConfigs()
	if err != nil {
		fmt.Printf("error while populating config, err : %v", err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	for env, val := range testEnvValues {
		err := os.Setenv(env, val)
		if err != nil {
			fmt.Printf("Could not set env variables, error : %v", err)
			os.Exit(1)
		}
	}
	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}

	initGlobals()
	m.Run()
}

func defaultTestModel() Model {
	return GenerateModel(spfPromptChar, shellPromptChar, true, defaultTestMaxHeight, defaultTestWidth)
}

func TestModel_HandleUpdate(t *testing.T) {
	// We don't test getPromptAction here. It is a separate test
	t.Run("Handle update called on closed Model", func(t *testing.T) {
		m := defaultTestModel()
		action, _ := m.HandleUpdate(utils.TeaRuneKeyMsg("x"), defaultTestCwd)
		assert.Empty(t, m.textInput.Value())
		assert.True(t, m.validate())
		assert.False(t, m.IsOpen())
		assert.Equal(t, common.NoAction{}, action)
	})

	t.Run("Pressing confirm on empty input", func(t *testing.T) {
		actualTest := func(closeOnSuccess bool, openAfterEnter bool) {
			m := GenerateModel(spfPromptChar, shellPromptChar, closeOnSuccess, defaultTestMaxHeight, defaultTestWidth)
			m.Open(true)
			assert.True(t, m.IsOpen())

			action, _ := m.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter}, defaultTestCwd)
			assert.Equal(t, openAfterEnter, m.IsOpen())
			assert.Equal(t, common.NoAction{}, action)
			assert.Empty(t, m.resultMsg)
			assert.True(t, m.LastActionSucceeded())
			assert.True(t, m.validate())
		}

		actualTest(true, false)
		actualTest(false, true)
	})

	t.Run("Validate Prompt Actions", func(t *testing.T) {
		m := defaultTestModel()
		m.Open(false)

		action, _ := m.HandleUpdate(utils.TeaRuneKeyMsg(SplitCommand), defaultTestCwd)
		assert.Equal(t, common.NoAction{}, action)

		action, _ = m.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter}, defaultTestCwd)
		assert.Equal(t, common.SplitPanelAction{}, action)

		_, _ = m.HandleUpdate(utils.TeaRuneKeyMsg("bad_command"), defaultTestCwd)
		action, _ = m.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter}, defaultTestCwd)
		assert.Equal(t, common.NoAction{}, action)
		assert.False(t, m.LastActionSucceeded())
		assert.NotEmpty(t, m.resultMsg)

		m.setShellMode(true)
		command := "abc def /xyz"
		_, _ = m.HandleUpdate(utils.TeaRuneKeyMsg(command), defaultTestCwd)
		action, _ = m.HandleUpdate(tea.KeyMsg{Type: tea.KeyEnter}, defaultTestCwd)
		assert.Equal(t, common.ShellCommandAction{Command: command}, action)
	})

	t.Run("Validate Cancel typing", func(t *testing.T) {
		m := defaultTestModel()

		actualTest := func(closeKey tea.KeyMsg, shouldBeOpen bool) {
			m.Open(true)
			_, _ = m.HandleUpdate(utils.TeaRuneKeyMsg("xyz"), defaultTestCwd)
			action, _ := m.HandleUpdate(closeKey, defaultTestCwd)
			assert.Equal(t, common.NoAction{}, action)
			assert.Equal(t, shouldBeOpen, m.IsOpen())
		}

		actualTest(tea.KeyMsg{Type: tea.KeyCtrlC}, false)
		actualTest(tea.KeyMsg{Type: tea.KeyEscape}, false)
		actualTest(tea.KeyMsg{Type: tea.KeyCtrlD}, true)
	})

	t.Run("Switching between shell and SPF mode", func(t *testing.T) {
		actualTest := func(promptChar string, shellChar string) {
			m := GenerateModel(promptChar, shellChar, true, defaultTestMaxHeight, defaultTestWidth)
			m.Open(true)
			assert.True(t, m.IsShellMode())

			// Shell to prompt
			action, _ := m.HandleUpdate(utils.TeaRuneKeyMsg(promptChar), defaultTestCwd)
			assert.False(t, m.IsShellMode())
			assert.True(t, m.LastActionSucceeded())
			assert.Equal(t, common.NoAction{}, action)

			// Prompt to shell
			action, _ = m.HandleUpdate(utils.TeaRuneKeyMsg(shellChar), defaultTestCwd)
			assert.True(t, m.IsShellMode())
			assert.True(t, m.LastActionSucceeded())
			assert.Equal(t, common.NoAction{}, action)

			// Pressing shellChar when you are already on shell shouldn't to anything
			action, _ = m.HandleUpdate(utils.TeaRuneKeyMsg(shellChar), defaultTestCwd)
			assert.True(t, m.IsShellMode())
			assert.True(t, m.LastActionSucceeded())
			assert.Equal(t, common.NoAction{}, action)
		}
		actualTest(">", ":")
		actualTest("$", "#")
	})

	t.Run("Validate Cursor Blink update", func(t *testing.T) {
		m := defaultTestModel()
		m.Open(true)
		assert.False(t, m.textInput.Cursor.Blink)

		blinkMsg := m.textInput.Cursor.BlinkCmd()()
		action, _ := m.HandleUpdate(blinkMsg, defaultTestCwd)
		assert.Equal(t, common.NoAction{}, action)
		assert.True(t, m.textInput.Cursor.Blink)

		blinkMsg = m.textInput.Cursor.BlinkCmd()()
		action, _ = m.HandleUpdate(blinkMsg, defaultTestCwd)
		assert.Equal(t, common.NoAction{}, action)
		assert.False(t, m.textInput.Cursor.Blink)

		blinkMsg = m.textInput.Cursor.BlinkCmd()()
		action, _ = m.HandleUpdate(blinkMsg, defaultTestCwd)
		assert.Equal(t, common.NoAction{}, action)
		assert.True(t, m.textInput.Cursor.Blink)

		// We could test BlinkCancelled and initialBlink as well, but that's too much for now
	})
}

func TestModel_HandleResults(t *testing.T) {
	t.Run("Verify Shell results update", func(t *testing.T) {
		m := defaultTestModel()
		m.Open(true)
		m.HandleShellCommandResults(0, "")

		// Validate close happens when closeOnSuccess is true
		assert.True(t, m.LastActionSucceeded())
		assert.Equal(t, "Command exited with status 0 (No output)", m.resultMsg)
		assert.False(t, m.IsOpen())

		m.Open(true)
		m.HandleShellCommandResults(1, "")
		assert.False(t, m.LastActionSucceeded())
		assert.Equal(t, "Command exited with status 1 (No output)", m.resultMsg)
		assert.True(t, m.IsOpen())

		m.closeOnSuccess = false
		m.HandleShellCommandResults(0, "")
		// Validate that close does not happen when closeOnSuccess is true
		assert.True(t, m.LastActionSucceeded())
		assert.Equal(t, "Command exited with status 0 (No output)", m.resultMsg)
		assert.True(t, m.IsOpen())
	})

	t.Run("Verify Shell command output is displayed", func(t *testing.T) {
		m := defaultTestModel()
		m.closeOnSuccess = false
		m.Open(true)

		// Test with single line output
		m.HandleShellCommandResults(0, "hello world")
		assert.Equal(t, "Command exited with status 0, Output:\nhello world", m.resultMsg)

		// Test with multi-line output
		m.HandleShellCommandResults(0, "line1\nline2\nline3")
		assert.Equal(t, "Command exited with status 0, Output:\nline1\nline2\nline3", m.resultMsg)

		// Test output is trimmed
		m.HandleShellCommandResults(0, "  trimmed output  \n")
		assert.Equal(t, "Command exited with status 0, Output:\ntrimmed output", m.resultMsg)

		m.HandleShellCommandResults(0, "ESC SEQ\x1b[2;6H")
		assert.Equal(t, "Command exited with status 0, Output:\nESC SEQ[2;6H", m.resultMsg)

		// Test with failed command and output
		m.HandleShellCommandResults(1, "error message")
		assert.False(t, m.LastActionSucceeded())
		assert.Equal(t, "Command exited with status 1, Output:\nerror message", m.resultMsg)
	})

	t.Run("Verify SPF results update", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, false, defaultTestMaxHeight, defaultTestWidth)
		m.Open(true)
		msg := "Test message"
		m.HandleSPFActionResults(true, msg)

		assert.True(t, m.LastActionSucceeded())
		assert.Equal(t, msg, m.resultMsg)
		assert.True(t, m.IsOpen())

		m.closeOnSuccess = true
		// Validate close happens when closeOnSuccess is true
		m.HandleSPFActionResults(true, "")
		assert.False(t, m.IsOpen())
	})
}

func TestModel_Render(t *testing.T) {
	// Test
	// 1 - Default view with shell mode and spf prompt mode
	// 2 - User input
	// 3 - User input, that is truncated due to being too large
	// 4 - User input with special characters, emojies, etc.
	// 5 - Prompt mode suggestion with these prefixes
	//   - "cd"
	//   - "c"
	//   - "open <PATH>"
	//   - "open <PATH> <Extra arg>"
	//   - "non_existent_command"

	// 6 - Model with result message (Without is tested above)
	// 7 - Color of result message green on success, red on failure
	// This one is hard, and we will likely not do it soon.
	// Needs global style variables

	// Challenges - needs border config strings for render test
	t.Run("Basic Render Checks", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, true, 10, 40)
		m.setShellMode(true)
		res := ansi.Strip(m.Render())
		exp := "" +
			"╭─┤ " + icon.Terminal + " superfile Prompt (Shell Mode) ├──╮\n" +
			// 23--------4------------56789012345678901234567890123456789
			"│ :                                    │\n" +
			// 23456789012345678901234567890123456789
			"├──────────────────────────────────────┤\n" +
			"│ '>' - Get into SPF mode              │\n" +
			"╰──────────────────────────────────────╯"
		assert.Equal(t, exp, res)
		m.setShellMode(false)
		res = ansi.Strip(m.Render())
		exp = "" +
			"╭─┤ " + icon.Terminal + " superfile Prompt (SPF Mode) ├────╮\n" +
			// 23--------4------------56789012345678901234567890123456789
			"│ >                                    │\n" +
			"├──────────────────────────────────────┤\n" +
			"│ ':' - Get into Shell mode            │\n" +
			"│ 'open <PATH>' - Open a new panel at a│\n" +
			"│ 'split' - Open a new panel at a curre│\n" +
			"│ 'cd <PATH>' - Change directory of cur│\n" +
			"╰──────────────────────────────────────╯"
		assert.Equal(t, exp, res)
	})

	t.Run("Test User Input", func(t *testing.T) {
		execute := func(input string, expected string) {
			// Changing this will need test adjustments
			width := 10
			m := GenerateModel(spfPromptChar, shellPromptChar, true, 10, width)
			m.Open(true)
			m.textInput.SetValue(input)
			m.textInput.Cursor.Blink = false
			res := ansi.Strip(m.Render())
			inputLine := strings.Split(res, "\n")[1]
			require.Equal(t, width, ansi.StringWidth(inputLine))
			// | : xxxx |
			// 0123456789
			content := strings.TrimPrefix(inputLine, "│ : ")
			content = strings.TrimSuffix(content, " │")

			assert.Equal(t, expected, content)
		}
		execute("abc", "abc ")
		execute("0123456789", "6789")
		execute("✅1✅2", "1✅2")
		execute("✅1✅2✅", "2✅ ")
	})

	t.Run("Result Message", func(t *testing.T) {
		m := GenerateModel(spfPromptChar, shellPromptChar, true, 10, 50)
		m.setShellMode(true)
		m.HandleShellCommandResults(0, "")
		res := ansi.Strip(m.Render())
		exp := "" +
			"╭─┤ " + icon.Terminal + " superfile Prompt (Shell Mode) ├────────────╮\n" +
			// 23--------4------------567890123456789012345678901234567890123456789
			"│ :                                              │\n" +
			// 234567890123456789012345678901234567890123456789
			"├────────────────────────────────────────────────┤\n" +
			"│ '>' - Get into SPF mode                        │\n" +
			"├────────────────────────────────────────────────┤\n" +
			"│ Success : Command exited with status 0 (No outp│\n" +
			"╰────────────────────────────────────────────────╯"
		assert.Equal(t, exp, res)
		m.HandleShellCommandResults(1, "")
		res = ansi.Strip(m.Render())
		exp = "" +
			"╭─┤ " + icon.Terminal + " superfile Prompt (Shell Mode) ├────────────╮\n" +
			// 23--------4------------567890123456789012345678901234567890123456789
			"│ :                                              │\n" +
			// 234567890123456789012345678901234567890123456789
			"├────────────────────────────────────────────────┤\n" +
			"│ '>' - Get into SPF mode                        │\n" +
			"├────────────────────────────────────────────────┤\n" +
			"│ Error : Command exited with status 1 (No output│\n" +
			"╰────────────────────────────────────────────────╯"
		assert.Equal(t, exp, res)
	})
	shellModeSuggestion := "':' - Get into Shell mode"
	var openCmdSuggestion string
	var splitCmdSuggestion string
	var cdCmdSuggestion string
	for _, cmd := range defaultCommandSlice() {
		curSuggestion := "'" + cmd.usage + "' - " + cmd.description

		switch cmd.command {
		case OpenCommand:
			openCmdSuggestion = curSuggestion
		case SplitCommand:
			splitCmdSuggestion = curSuggestion
		case CdCommand:
			cdCmdSuggestion = curSuggestion
		default:
			assert.Fail(t, "Unknow command")
		}
	}

	testdataSuggestions := []struct {
		name                string
		textInput           string
		expectedSuggestions []string
	}{
		{
			name:      "No Input",
			textInput: "",
			expectedSuggestions: []string{
				shellModeSuggestion,
				openCmdSuggestion,
				splitCmdSuggestion,
				cdCmdSuggestion,
			},
		},
		{
			name:      "Command without args",
			textInput: "cd",
			expectedSuggestions: []string{
				cdCmdSuggestion,
			},
		},
		{
			name:      "Incomplete Command",
			textInput: "c",
			expectedSuggestions: []string{
				cdCmdSuggestion,
			},
		},
		{
			name:      "Command with args",
			textInput: "open /abc",
			expectedSuggestions: []string{
				openCmdSuggestion,
			},
		},
		{
			name:      "Command with extra args",
			textInput: "open /abc /abc",
			expectedSuggestions: []string{
				openCmdSuggestion,
			},
		},
		{
			name:                "Invalid command",
			textInput:           "non_existent_command",
			expectedSuggestions: []string{},
		},
	}

	for _, tt := range testdataSuggestions {
		t.Run(tt.name, func(t *testing.T) {
			m := DefaultModel(defaultTestMaxHeight, defaultTestWidth)
			m.Open(false)
			m.textInput.SetValue(tt.textInput)
			res := ansi.Strip(m.Render())
			resLines := strings.Split(res, "\n")
			if len(tt.expectedSuggestions) == 0 {
				require.Len(t, resLines, 3)
				return
			}

			require.Len(t, resLines, 4+len(tt.expectedSuggestions))
			suggestionLines := resLines[3 : len(resLines)-1]
			require.Len(t, suggestionLines, len(tt.expectedSuggestions))

			for i := range tt.expectedSuggestions {
				exp := tt.expectedSuggestions[i]
				actualLine := suggestionLines[i]
				actualLine = strings.TrimPrefix(actualLine, "│ ")
				actualLine = strings.TrimSuffix(actualLine, "│")
				actualLine = strings.TrimSpace(actualLine)
				assert.Equal(t, exp, actualLine)
			}
		})
	}
}
