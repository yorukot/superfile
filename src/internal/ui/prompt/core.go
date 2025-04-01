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
		headline:  icon.Terminal + " superfile - Prompt",
		open:      false,
		shellMode: true,
		textInput: common.GeneratePromptTextInput(),
	}
}

func (p *PromptModal) HandleMessage(msg tea.Msg) (common.PromptAction, tea.Cmd) {
	slog.Debug("promptModal HandleMessage()", "msg", msg,
		"textInput", p.textInput.Value())
	action := common.NoPromptAction()
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
			}
			p.textInput.SetValue("")
		} else if slices.Contains(common.Hotkeys.CancelTyping, msg.String()) {
			p.Close()
		} else {
			if msg.String() == ">" && p.textInput.Value() == "" {
				p.shellMode = false
			} else if msg.String() == ":" && p.textInput.Value() == "" {
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

func (p *PromptModal) Open() {
	p.open = true
	_ = p.textInput.Focus()
}

func (p *PromptModal) Close() {
	p.open = false
	p.shellMode = true
	p.textInput.SetValue("")
}

func (p *PromptModal) Render(width int) string {

	content := " " + p.headline + modeString(p.shellMode)
	content += "\n" + strings.Repeat(common.Config.BorderTop, width)
	content += "\n" + " " + shellPrompt(p.shellMode) + " " + p.textInput.View()
	return common.ModalBorderStyleLeft(1, width+2).Render(content)
}

func getPromptAction(shellMode bool, value string) (common.PromptAction, error) {
	if value == "" {
		return common.NoPromptAction(), nil
	}
	if shellMode {
		return common.PromptAction{
			Action: common.ShellCommandAction,
			Args:   []string{value},
		}, nil
	}

	promptArgs := strings.Fields(value)

	switch promptArgs[0] {
	case "split":
		return common.PromptAction{
			Action: common.SplitPanelAction,
		}, nil
	default:
		return common.NoPromptAction(), fmt.Errorf("invalid spf prompt command")
	}

}
