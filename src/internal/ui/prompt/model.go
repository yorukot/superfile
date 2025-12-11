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

func (m *Model) handleNormalKeyInput(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	switch {
	case m.textInput.Value() == "" && msg.String() == m.spfPromptHotkey:
		m.setShellMode(false)
	case m.textInput.Value() == "" && msg.String() == m.shellPromptHotkey:
		m.setShellMode(true)
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

func (m *Model) Render() string {
	r := ui.PromptRenderer(m.maxHeight, m.width)
	r.SetBorderTitle(m.headline + " " + modeString(m.shellMode))
	r.AddLines(" " + m.textInput.View())

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
	m.textInput.Width = width - promptInputPadding
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
