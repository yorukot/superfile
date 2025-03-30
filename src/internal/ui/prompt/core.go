package prompt

import (
	"fmt"
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
	}
}

func (p *PromptModal) HandleMessage(msg string) {
	slog.Debug("promptModal HandleMessage()", "msg", msg)
	if slices.Contains(common.Hotkeys.ConfirmTyping, msg) {
		p.textInput.SetValue("")
	} else if slices.Contains(common.Hotkeys.CancelTyping, msg) {
		p.Close()
	} else {
		p.textInput.Focus()
	}
}

func (p *PromptModal) Open() {
	p.open = true
}

func (p *PromptModal) Close() {
	p.open = false
	p.textInput.SetValue("")
}

func (p *PromptModal) Render(width int) string {

	var content, promptLine string
	text := p.textInput.Value()

	content += fmt.Sprintf(
		"%s%*s\n\n",
		p.headline,
		width-len(p.headline), "",
	)

	promptLine += text
	content += fmt.Sprintf("%s%*s\n%s\n", promptLine, width-len(promptLine), "", strings.Repeat(common.Config.BorderTop, width))

	return common.ModalBorderStyle(1, width+2).Render(content)

}
