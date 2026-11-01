package prompt

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/yorukot/superfile/src/internal/ui"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func DefaultModel(maxHeight int, width int) Model {
	return GenerateModel(common.Hotkeys.OpenSPFPrompt[0],
		common.Hotkeys.OpenCommandLine[0], common.Config.ShellCloseOnSuccess, maxHeight, width)
}

func GenerateModel(spfPromptHotkey string, shellPromptHotkey string, closeOnSuccess bool,
	maxHeight int, width int) Model {
	m := Model{
		headline:          icon.Terminal + icon.Space + promptHeadlineText,
		open:              false,
		shellMode:         true,
		textInput:         common.GeneratePromptTextInput(),
		commands:          defaultCommandSlice(),
		spfPromptHotkey:   spfPromptHotkey,
		shellPromptHotkey: shellPromptHotkey,
		actionSuccess:     true,
		closeOnSuccess:    closeOnSuccess,
	}
	m.SetMaxHeight(maxHeight)
	m.SetWidth(width)
	return m
}

func (m *Model) HandleUpdate(msg tea.Msg, cwdLocation string) (common.ModelAction, tea.Cmd) {
	var action common.ModelAction
	action = common.NoAction{}
	var cmd tea.Cmd
	if !m.IsOpen() {
		slog.Error("HandleUpdate called on closed prompt")
		return action, cmd
	}

	// Keep cwd updated for tab completion
	m.cwd = cwdLocation

	switch msg := msg.(type) {
	case completionsReadyMsg:
		m.handleCompletionsReady(msg)
	case tea.KeyPressMsg:
		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()):
			m.completions = nil
			m.completionDone = false
			action = m.handleConfirm(cwdLocation)
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
			m.completions = nil
			m.completionDone = false
			m.Close()
		default:
			cmd = m.handleNormalKeyInput(msg)
		}
	default:
		// Non keypress updates like Cursor Blink
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return action, cmd
}

func (m *Model) handleConfirm(cwdLocation string) common.ModelAction {
	// Pressing confirm on empty prompt will trigger close
	if m.textInput.Value() == "" {
		m.CloseOnSuccessIfNeeded()
	}

	// Create Action based on input
	var err error
	action, err := getPromptAction(m.shellMode, m.textInput.Value(), cwdLocation)
	if err == nil {
		m.resultMsg = ""
		m.actionSuccess = true
	} else if cmdErr, ok := err.(invalidCmdError); ok { //nolint: errorlint // We don't expect a wrapped error here
		slog.Error("Error from getPromptAction", "error", cmdErr, "uiMsg", cmdErr.uiMsg)
		m.resultMsg = cmdErr.uiMessage()
		m.actionSuccess = false
	} else {
		slog.Error("Unexpected error from getPromptAction", "error", err)
		m.resultMsg = err.Error()
		m.actionSuccess = false
	}
	m.textInput.SetValue("")
	return action
}

func (m *Model) handleNormalKeyInput(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case m.textInput.Value() == "" && msg.String() == m.spfPromptHotkey:
		m.setShellMode(false)
	case m.textInput.Value() == "" && msg.String() == m.shellPromptHotkey:
		m.setShellMode(true)
	case msg.String() == "tab":
		cmd = m.handleTabCompletion()
	case msg.String() == "shift+tab":
		m.handleReverseTabCompletion()
	default:
		// Any non-Tab key clears completions
		m.completions = nil
		m.completionDone = false
		m.textInput, cmd = m.textInput.Update(msg)
	}
	m.resultMsg = ""
	m.actionSuccess = true
	return cmd
}

// handleTabCompletion fires when the user presses Tab. If we already have
// completions, cycle forward. Otherwise fetch new completions from the shell.
func (m *Model) handleTabCompletion() tea.Cmd {
	if m.completionDone {
		// User pressed Tab after previously completing — fetch fresh completions
		m.completions = nil
		m.completionDone = false
	}
	if len(m.completions) > 0 {
		// Cycle to next completion
		m.completionIdx = (m.completionIdx + 1) % len(m.completions)
		completion := m.completions[m.completionIdx]
		m.textInput.SetValue(applyCompletion(m.textInput.Value(), completion))
		m.textInput.SetCursor(len(m.textInput.Value()))
		return nil
	}
	// First Tab — fetch completions (async)
	m.completionIdx = 0
	return fetchCompletions(m.textInput.Value(), m.cwd)
}

// handleReverseTabCompletion cycles backward through completions on Shift+Tab.
func (m *Model) handleReverseTabCompletion() {
	if len(m.completions) == 0 {
		return
	}
	m.completionIdx--
	if m.completionIdx < 0 {
		m.completionIdx = len(m.completions) - 1
	}
	completion := m.completions[m.completionIdx]
	m.textInput.SetValue(applyCompletion(m.textInput.Value(), completion))
		m.textInput.SetCursor(len(m.textInput.Value()))
}

// handleCompletionsReady applies the first completion result and populates
// the dropdown list. If there is only one match, it auto-applies immediately.
func (m *Model) handleCompletionsReady(msg completionsReadyMsg) {
	if len(msg.completions) == 0 {
		m.completions = nil
		return
	}
	m.completions = msg.completions
	m.completionIdx = 0

	// Check if the current input already changed while we were fetching
	currentPrefix, _ := getLastToken(m.textInput.Value())
	if currentPrefix != msg.prefix {
		// Input changed — completions are stale
		m.completions = nil
		return
	}

	if len(msg.completions) == 1 {
		// Single match — apply immediately, mark done
		m.textInput.SetValue(applyCompletion(m.textInput.Value(), msg.completions[0]))
		m.textInput.SetCursor(len(m.textInput.Value()))
		m.completions = nil
		m.completionDone = true
		return
	}

	// Multiple matches — find common prefix and extend input
	prefix := commonPrefix(msg.completions)
	if prefix != msg.prefix {
		m.textInput.SetValue(applyCompletion(m.textInput.Value(), prefix))
		m.textInput.SetCursor(len(m.textInput.Value()))
		// Re-fetch with the extended prefix for a more precise list
		m.completions = nil
	}
}

// After action is performed, model will update the Model with results
func (m *Model) HandleShellCommandResults(retCode int, output string) {
	m.actionSuccess = retCode == 0
	m.resultMsg = fmt.Sprintf("Command exited with status %d", retCode)

	output = strings.TrimSpace(common.MakePrintableWithEscCheck(output, false))
	if output != "" {
		m.resultMsg += ", Output:\n" + output
	} else {
		m.resultMsg += " (No output)"
	}
	m.CloseOnSuccessIfNeeded()
}

// After action is performed, model will update the prompt.Model with results
// In case of NoAction, this method should not be called.
func (m *Model) HandleSPFActionResults(success bool, msg string) {
	m.actionSuccess = success
	m.resultMsg = msg
	m.CloseOnSuccessIfNeeded()
}

func (m *Model) Render() string {
	r := ui.PromptRenderer(m.maxHeight, m.width)
	r.SetBorderTitle(m.headline + " " + modeString(m.shellMode))
	r.AddLines(" " + m.textInput.View())

	// Show completion dropdown if we have candidates
	if len(m.completions) > 0 {
		r.AddSection()
		for i, c := range m.completions {
			line := " " + c
			if i == m.completionIdx {
				line = common.PromptSuccessStyle.Render(" >" + c)
			}
			r.AddLines(line)
		}
	}

	if !m.shellMode {
		// To make sure its added one time only per render call
		hintSectionAdded := false
		if m.textInput.Value() == "" {
			if !hintSectionAdded {
				r.AddSection()
				hintSectionAdded = true
			}
			r.AddLines(" '" + m.shellPromptHotkey + "' - Get into Shell mode")
		}
		command := getFirstToken(m.textInput.Value())
		for _, cmd := range m.commands {
			if strings.HasPrefix(cmd.command, command) {
				if !hintSectionAdded {
					r.AddSection()
					hintSectionAdded = true
				}
				r.AddLines(" '" + cmd.usage + "' - " + cmd.description)
			}
		}
	} else if m.textInput.Value() == "" {
		r.AddSection()
		r.AddLines(" '" + m.spfPromptHotkey + "' - Get into SPF mode")
		r.AddLines(" '!' prefix - Full terminal mode (for ssh, vim, etc.)")
	}

	if m.resultMsg != "" {
		msgPrefix := successMessagePrefix
		resultStyle := common.PromptSuccessStyle
		if !m.actionSuccess {
			resultStyle = common.PromptFailureStyle
			msgPrefix = failureMessagePrefix
		}
		r.AddSection()
		r.AddLines(resultStyle.Render(" " + msgPrefix + " : " + m.resultMsg))
	}
	return r.Render()
}

func (m *Model) Open(shellMode bool) {
	m.open = true
	m.setShellMode(shellMode)
	_ = m.textInput.Focus()
}

func (m *Model) setShellMode(shellMode bool) {
	m.shellMode = shellMode
	m.textInput.Prompt = shellPrompt(m.shellMode) + " "
}

func (m *Model) Close() {
	m.open = false
	m.setShellMode(true)
	m.textInput.SetValue("")
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) IsShellMode() bool {
	return m.shellMode
}

func (m *Model) LastActionSucceeded() bool {
	return m.actionSuccess
}

func (m *Model) GetWidth() int {
	return m.width
}

func (m *Model) GetMaxHeight() int {
	return m.maxHeight
}

func (m *Model) SetWidth(width int) {
	if width < PromptMinWidth {
		slog.Warn("Prompt initialized with too less width", "width", width)
		width = PromptMinWidth
	}
	m.width = width
	// Excluding borders(2), SpacePadding(1), Prompt(2), and one extra character that is appended
	// by textInput.View()
	m.textInput.SetWidth(width - promptInputPadding)
}

func (m *Model) SetMaxHeight(maxHeight int) {
	if maxHeight < PromptMinHeight {
		slog.Warn("Prompt initialized with too less maxHeight", "maxHeight", maxHeight)
		maxHeight = PromptMinHeight
	}
	m.maxHeight = maxHeight
}

func (m *Model) validate() bool {
	// Prompt was closed, but textInput was not cleared
	if !m.open && m.textInput.Value() != "" {
		return false
	}
	return true
}

func (m *Model) CloseOnSuccessIfNeeded() {
	if m.closeOnSuccess && m.actionSuccess {
		m.Close()
	}
}
