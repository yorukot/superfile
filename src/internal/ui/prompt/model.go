package prompt

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yorukot/superfile/src/internal/common"
	"log/slog"
	"slices"
	"strings"

	"github.com/yorukot/superfile/src/config/icon"
)

func DefaultModel() Model {
	return GenerateModel(common.Hotkeys.OpenSPFPrompt[0],
		common.Hotkeys.OpenCommandLine[0], common.Config.ShellCloseOnSucsess)
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
			// Pressing confirm on empty prompt will trigger close
			if m.textInput.Value() == "" {
				m.CloseOnSuccessIfNeeded()
			}

			// Create Action based on input
			var err error
			action, err = getPromptAction(m.shellMode, m.textInput.Value(), cwdLocation)
			if err == nil {
				m.resultMsg = ""
				m.actionSuccess = true
			} else if cmdErr, ok := err.(invalidCmdError); ok {
				// We dont expect a wrapped error here, so using type assertion
				m.resultMsg = cmdErr.uiMessage()
				m.actionSuccess = false
			} else {
				slog.Error("Unexpected error from getPromptAction", "error", err)
				m.resultMsg = err.Error()
				m.actionSuccess = false
			}
			m.textInput.SetValue("")
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
			m.Close()
		default:
			if m.textInput.Value() == "" && msg.String() == m.spfPromptHotkey {
				m.shellMode = false
			} else if m.textInput.Value() == "" && msg.String() == m.shellPromptHotkey {
				m.shellMode = true
			} else {
				m.textInput, cmd = m.textInput.Update(msg)
			}
			m.resultMsg = ""
			m.actionSuccess = true
		}
	default:
		m.textInput, cmd = m.textInput.Update(msg)

	}
	return action, cmd
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

// Todo : Unit test this
func (m *Model) Render(width int) string {

	divider := strings.Repeat(common.Config.BorderTop, width)
	content := " " + m.headline + modeString(m.shellMode)
	content += "\n" + divider
	content += "\n" + " " + shellPrompt(m.shellMode) + " " + m.textInput.View()
	suggestionText := ""

	if !m.shellMode {
		if m.textInput.Value() == "" {
			suggestionText += "\n '" + m.shellPromptHotkey + "' - Get into Shell mode"
		}
		command := getFirstToken(m.textInput.Value())
		for _, cmd := range m.commands {
			if strings.HasPrefix(cmd.command, command) {
				suggestionText += "\n '" + cmd.usage + "' - " + cmd.description
			}
		}
	} else if m.textInput.Value() == "" {
		suggestionText += "\n '" + m.spfPromptHotkey + "' - Get into SPF Prompt mode"
	}

	if suggestionText != "" {
		content += "\n" + divider
		content += suggestionText
	}
	// Rendering error Message is a but fuzzy right now. Todo Fix this.
	// Todo : Handle error being multiline or being too long
	if m.resultMsg != "" {
		msgPrefix := successMessagePrefix
		resultStyle := common.PromptSuccessStyle
		if !m.actionSuccess {
			resultStyle = common.PromptFailureStyle
			msgPrefix = failureMessagePrefix
		}
		content += "\n" + divider
		content += "\n " + resultStyle.Render(msgPrefix+" : "+m.resultMsg)
	}
	return common.ModalBorderStyleLeft(1, width+1).Render(content)
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
