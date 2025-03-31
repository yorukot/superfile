package prompt

import (
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

func (p *PromptModal) HandleMessage(msg string) {
	slog.Debug("promptModal HandleMessage()", "msg", msg,
		"textInput", p.textInput.Value(),
		"inputView", p.textInput.View())
	if slices.Contains(common.Hotkeys.ConfirmTyping, msg) {
		p.textInput.SetValue("")
	} else if slices.Contains(common.Hotkeys.CancelTyping, msg) {
		p.Close()
	}
}

func (p *PromptModal) HandleUpdate(msg tea.Msg) tea.Cmd {
	slog.Debug("promptModal HandleUpdate()", "msg", msg,
		"textInput", p.textInput.Value(),
		"inputView", p.textInput.View())
	var cmd tea.Cmd
	ignoreInput := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == ">" && p.textInput.Value() == "" {
			p.shellMode = false
			ignoreInput = true
		}
		if msg.String() == ":" && p.textInput.Value() == "" {
			p.shellMode = true
			ignoreInput = true
		}
	}

	if !ignoreInput {
		p.textInput, cmd = p.textInput.Update(msg)
	}
	return cmd
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
