package prompt

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/yorukot/superfile/src/internal/common"
	"log/slog"
	"path/filepath"
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
			} else {
				p.errorMsg = ""
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

	// Rendering error Message is a but fuzzy right now. Todo Fix this.
	if p.errorMsg != "" {
		content += "\n" + strings.Repeat(common.Config.BorderTop, width)
		content += "\n" + " " + p.errorMsg
	}
	return common.ModalBorderStyleLeft(1, width+2).Render(content)
}

func getPromptAction(shellMode bool, value string) (common.PromptAction, error) {
	noAction := common.NoPromptAction()
	if value == "" {
		return noAction, nil
	}
	if shellMode {
		return common.PromptAction{
			Action: common.ShellCommandAction,
			Args:   []string{value},
		}, nil
	}

	// Todo - Add tokenization for $() and ${} args
	promptArgs := strings.Fields(value)

	switch promptArgs[0] {
	case "split":
		return common.PromptAction{
			Action: common.SplitPanelAction,
		}, nil
	case "cd":
		if len(promptArgs) != 2 {
			return noAction, fmt.Errorf("cd prompts needs exactly one arguement, received %d",
				len(promptArgs)-1)
		}
		cdPath := ""
		if path, err := filepath.Abs(promptArgs[1]); err == nil {
			cdPath = path
		} else {
			return noAction, fmt.Errorf("invalid cd path : %s", path)
		}

		// Todo : Instead, we can have a function that creates this object
		return common.PromptAction{
			Action: common.CDCurrentPanelAction,
			Args:   []string{cdPath},
		}, nil
	case "open":
		// Todo : Duplication. Fix this
		if len(promptArgs) != 2 {
			return noAction, fmt.Errorf("open prompts needs exactly one arguement, received %d",
				len(promptArgs)-1)
		}
		newPanelPath := ""
		if path, err := filepath.Abs(promptArgs[1]); err == nil {
			newPanelPath = path
		} else {
			return noAction, fmt.Errorf("invalid open path : %s", path)
		}
		// Todo : Instead, we can have a function that creates this object
		return common.PromptAction{
			Action: common.OpenPanelAction,
			Args:   []string{newPanelPath},
		}, nil

	default:
		return noAction, fmt.Errorf("invalid spf prompt command")
	}

}
