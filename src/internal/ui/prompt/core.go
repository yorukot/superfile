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

func DefaultPrompt() PromptModal {
	return PromptModal{
		headline:          icon.Terminal + " superfile - Prompt",
		open:              false,
		shellMode:         true,
		textInput:         common.GeneratePromptTextInput(),
		commands:          defaultCommandSlice(),
		spfPromptHotkey:   common.Hotkeys.OpenSPFPrompt[0],
		shellPromptHotkey: common.Hotkeys.OpenCommandLine[0],
	}
}

func (p *PromptModal) HandleMessage(msg tea.Msg) (common.ModelAction, tea.Cmd) {
	slog.Debug("promptModal HandleMessage()", "msg", msg,
		"textInput", p.textInput.Value())
	var action common.ModelAction
	action = common.NoAction{}
	var cmd tea.Cmd
	if !p.IsOpen() {
		slog.Error("HandleMessage called on closed prompt")
		return action, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()) {
			var err error
			action, err = getPromptAction(p.shellMode, p.textInput.Value())
			if err != nil {
				// Todo create and error that wraps user facing message
				p.errorMsg = err.Error()
			} else {
				p.errorMsg = ""
			}
			p.textInput.SetValue("")
		} else if slices.Contains(common.Hotkeys.CancelTyping, msg.String()) {
			p.Close()
		} else {
			if p.textInput.Value() == "" && msg.String() == p.spfPromptHotkey {
				p.shellMode = false
			} else if p.textInput.Value() == "" && msg.String() == p.shellPromptHotkey {
				p.shellMode = true
			} else {
				p.textInput, cmd = p.textInput.Update(msg)
			}
		}
	default:
		p.textInput, cmd = p.textInput.Update(msg)
	}

	return action, cmd

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
		return common.SplitPanelAction{}, nil
	case "cd":
		if len(promptArgs) != 2 {
			return noAction, fmt.Errorf("cd prompts needs exactly one arguement, received %d",
				len(promptArgs)-1)
		}
		return common.CDCurrentPanelAction{
			Location: promptArgs[1],
		}, nil
	case "open":
		// Todo : Duplication. Fix this
		if len(promptArgs) != 2 {
			return noAction, fmt.Errorf("open prompts needs exactly one arguement, received %d",
				len(promptArgs)-1)
		}
		return common.OpenPanelAction{
			Location: promptArgs[1],
		}, nil

	default:
		return noAction, fmt.Errorf("invalid spf prompt command")
	}

}

func (p *PromptModal) Open(shellMode bool) {
	p.open = true
	p.shellMode = shellMode
	_ = p.textInput.Focus()
}

func (p *PromptModal) Close() {
	p.open = false
	p.shellMode = true
	p.textInput.SetValue("")
}

func (p *PromptModal) Render(width int) string {

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
		for i := 0; i < len(p.commands); i++ {
			if strings.HasPrefix(p.commands[i].command, command) {
				suggestionText += "\n '" + p.commands[i].usage + "' - " + p.commands[i].description
			}
		}
	} else {
		if p.textInput.Value() == "" {
			suggestionText += " '" + p.spfPromptHotkey + "' - Get into SPF Prompt mode"
		}
	}

	if suggestionText != "" {
		content += "\n" + divider
		content += suggestionText
	}
	// Rendering error Message is a but fuzzy right now. Todo Fix this.
	if p.errorMsg != "" {
		content += "\n" + divider
		content += "\n " + p.errorMsg
	}
	return common.ModalBorderStyleLeft(1, width+2).Render(content)
}
