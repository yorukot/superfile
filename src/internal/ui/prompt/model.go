package prompt

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/yorukot/superfile/src/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func DefaultModel() Model {
	return GenerateModel(common.Hotkeys.OpenSPFPrompt[0],
		common.Hotkeys.OpenCommandLine[0], common.Config.ShellCloseOnSuccess)
}

func GenerateModel(spfPromptHotkey string, shellPromptHotkey string, closeOnSuccess bool) Model {
	return Model{
		headline:          icon.Terminal + " " + promptHeadlineText,
		open:              false,
		shellMode:         true,
		textInput:         common.GeneratePromptTextInput(),
		commands:          defaultCommandSlice(),
		spfPromptHotkey:   spfPromptHotkey,
		shellPromptHotkey: shellPromptHotkey,
		actionSuccess:     true,
		closeOnSuccess:    closeOnSuccess,
	}
}

func (m *Model) HandleUpdate(msg tea.Msg, cwdLocation string) (common.ModelAction, tea.Cmd) {
	slog.Debug("prompt.Model HandleUpdate()", "msg", msg,
		"textInput", m.textInput.Value(),
		"cursorBlink", m.textInput.Cursor.Blink)
	var action common.ModelAction
	action = common.NoAction{}
	var cmd tea.Cmd
	if !m.IsOpen() {
		slog.Error("HandleUpdate called on closed prompt")
		return action, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()):
			action = m.handleConfirm(cwdLocation)
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
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
	} else if cmdErr, ok := err.(invalidCmdError); ok { //nolint: errorlint // We don't expect a wrapped error here, so using type assertion
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

func (m *Model) handleNormalKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case m.textInput.Value() == "" && msg.String() == m.spfPromptHotkey:
		m.shellMode = false
	case m.textInput.Value() == "" && msg.String() == m.shellPromptHotkey:
		m.shellMode = true
	default:
		m.textInput, cmd = m.textInput.Update(msg)
	}
	m.resultMsg = ""
	m.actionSuccess = true
	return cmd
}

// After action is performed, model will update the Model with results
func (m *Model) HandleShellCommandResults(retCode int, _ string) {
	m.actionSuccess = retCode == 0
	// Not allowing user to see output yet. This needs to be sanitized and
	// need to be made sure that it doesn't breaks layout
	// Hence we are ignoring output for now
	m.resultMsg = fmt.Sprintf("Command exited with status %d", retCode)
	m.CloseOnSuccessIfNeeded()
}

// After action is performed, model will update the prompt.Model with results
// In case of NoAction, this method should not be called.
func (m *Model) HandleSPFActionResults(success bool, msg string) {
	m.actionSuccess = success
	m.resultMsg = msg
	m.CloseOnSuccessIfNeeded()
}

func (m *Model) Render(maxHeight int, width int) string {
	r := ui.PromptRenderer(maxHeight, width)
	r.SetBorderTitle(m.headline + modeString(m.shellMode))
	r.AddLines(" " + shellPrompt(m.shellMode) + " " + m.textInput.View())

	if !m.shellMode {
		r.AddSection()
		if m.textInput.Value() == "" {
			r.AddLines(" '" + m.shellPromptHotkey + "' - Get into Shell mode")
		}
		command := getFirstToken(m.textInput.Value())
		for _, cmd := range m.commands {
			if strings.HasPrefix(cmd.command, command) {
				r.AddLines(" '" + cmd.usage + "' - " + cmd.description)
			}
		}
	} else if m.textInput.Value() == "" {
		r.AddSection()
		r.AddLines(" '" + m.spfPromptHotkey + "' - Get into SPF Prompt mode")
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
	m.shellMode = shellMode
	_ = m.textInput.Focus()
}

func (m *Model) Close() {
	m.open = false
	m.shellMode = true
	m.textInput.SetValue("")
}

func (m *Model) IsOpen() bool {
	return m.open
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
