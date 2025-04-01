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

func DefaultPrompt() Model {
	return Model{
		headline:          icon.Terminal + " " + promptHeadlineText,
		open:              false,
		shellMode:         true,
		textInput:         common.GeneratePromptTextInput(),
		commands:          defaultCommandSlice(),
		spfPromptHotkey:   common.Hotkeys.OpenSPFPrompt[0],
		shellPromptHotkey: common.Hotkeys.OpenCommandLine[0],
	}
}

func (p *Model) HandleUpdate(msg tea.Msg) (common.ModelAction, tea.Cmd) {
	slog.Debug("promptModal HandleUpdate()", "msg", msg,
		"textInput", p.textInput.Value())
	var action common.ModelAction
	action = common.NoAction{}
	var cmd tea.Cmd
	if !p.IsOpen() {
		slog.Error("HandleUpdate called on closed prompt")
		return action, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()):
			var err error
			action, err = getPromptAction(p.shellMode, p.textInput.Value())
			if err == nil {
				p.resultMsg = ""
				p.actionSuccess = true
			} else if cmdErr, ok := err.(InvalidCmdError); ok {
				// We dont expect a wrapped error here, so using type assertion
				p.resultMsg = cmdErr.UIMessage()
				p.actionSuccess = false
			} else {
				p.resultMsg = err.Error()
				p.actionSuccess = false
			}
			p.textInput.SetValue("")
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
			p.Close()
		default:
			if p.textInput.Value() == "" && msg.String() == p.spfPromptHotkey {
				p.shellMode = false
			} else if p.textInput.Value() == "" && msg.String() == p.shellPromptHotkey {
				p.shellMode = true
			} else {
				p.textInput, cmd = p.textInput.Update(msg)
			}
			p.resultMsg = ""
			p.actionSuccess = true
		}
	default:
		p.textInput, cmd = p.textInput.Update(msg)
	}
	return action, cmd
}

// After action is performed, model will update the Model with results
func (p *Model) HandleShellCommandResults(retCode int, _ string) {
	p.actionSuccess = retCode == 0
	// Not allowing user to see output yet. This needs to be sanitized and
	// need to be made sure that it doesn't breaks layout
	// Hence we are ignoring output for now
	p.resultMsg = fmt.Sprintf("Command exited withs status %d", retCode)
}

// After action is performed, model will update the Model with results
// In case of NoAction, this method should not be called.
func (p *Model) HandleSPFActionResults(success bool, msg string) {
	p.actionSuccess = success
	p.resultMsg = msg
}

func getPromptAction(shellMode bool, value string) (common.ModelAction, error) {
	noAction := common.NoAction{}
	if value == "" {
		return noAction, nil
	}
	if shellMode {
		return common.ShellCommandAction{
			Command: value,
		}, nil
	}

	// Todo - Add tokenization for $() and ${} args
	promptArgs := strings.Fields(value)

	switch promptArgs[0] {
	case "split":
		if len(promptArgs) != 1 {
			return noAction, InvalidCmdError{
				uiMsg: "split command should not be given arguments",
			}
		}
		return common.SplitPanelAction{}, nil
	case "cd":
		if len(promptArgs) != 2 {
			return noAction, InvalidCmdError{
				uiMsg: fmt.Sprintf("cd command needs exactly one argument, received %d",
					len(promptArgs)-1),
			}
		}
		return common.CDCurrentPanelAction{
			Location: promptArgs[1],
		}, nil
	case "open":
		if len(promptArgs) != 2 {
			return noAction, InvalidCmdError{
				uiMsg: fmt.Sprintf("open command needs exactly one argument, received %d",
					len(promptArgs)-1),
			}
		}
		return common.OpenPanelAction{
			Location: promptArgs[1],
		}, nil

	default:
		return noAction, InvalidCmdError{
			uiMsg: "Invalid spf prompt command : " + promptArgs[0],
		}
	}

}

func (p *Model) Open(shellMode bool) {
	p.open = true
	p.shellMode = shellMode
	_ = p.textInput.Focus()
}

func (p *Model) Close() {
	p.open = false
	p.shellMode = true
	p.textInput.SetValue("")
}

func (p *Model) Render(width int) string {

	// Todo fix divider being a bit smaller
	divider := strings.Repeat(common.Config.BorderTop, width)
	content := " " + p.headline + modeString(p.shellMode)
	content += "\n" + divider
	content += "\n" + " " + shellPrompt(p.shellMode) + " " + p.textInput.View()
	suggestionText := ""

	if !p.shellMode {
		if p.textInput.Value() == "" {
			suggestionText += "\n '" + p.shellPromptHotkey + "' - Get into Shell mode"
		}
		command := getFirstToken(p.textInput.Value())
		for _, cmd := range p.commands {
			if strings.HasPrefix(cmd.command, command) {
				suggestionText += "\n '" + cmd.usage + "' - " + cmd.description
			}
		}
	} else if p.textInput.Value() == "" {
		suggestionText += "\n '" + p.spfPromptHotkey + "' - Get into SPF Prompt mode"
	}

	if suggestionText != "" {
		content += "\n" + divider
		content += suggestionText
	}
	// Rendering error Message is a but fuzzy right now. Todo Fix this.
	// Todo : Handle error being multiline or being too long
	if p.resultMsg != "" {
		msgPrefix := successMessagePrefix
		resultStyle := common.PromptSuccessStyle
		if !p.actionSuccess {
			resultStyle = common.PromptFailureStyle
			msgPrefix = failureMessagePrefix
		}
		content += "\n" + divider
		content += "\n " + resultStyle.Render(msgPrefix+" : "+p.resultMsg)
	}
	return common.ModalBorderStyleLeft(1, width+2).Render(content)
}
