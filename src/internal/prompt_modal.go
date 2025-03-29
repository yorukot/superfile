package internal

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
)

type PromptCommand struct {
	renderPrefix  string
	renderHint    string
	handleCommand func(input string, p *promptModalType, m *model) bool
}
type PromptCommandPrefix = string

const (
	PROMPT_COMMAND_COMMAND      PromptCommandPrefix = ""
	PROMPT_COMMAND_SHELL        PromptCommandPrefix = "$"
	PROMPT_COMMAND_CD           PromptCommandPrefix = "cd"
	PROMPT_COMMAND_NEWFILEPANEL PromptCommandPrefix = "open"
	PROMPT_COMMAND_SPLIT        PromptCommandPrefix = "split"
)

var promptCommands map[PromptCommandPrefix]PromptCommand

func init() {

	promptCommands = make(map[PromptCommandPrefix]PromptCommand)

	promptCommands[PROMPT_COMMAND_COMMAND] = PromptCommand{
		renderPrefix: "> ",
		handleCommand: func(input string, p *promptModalType, m *model) bool {

			fields := strings.Fields(input)
			inputCmd := fields[0]

			if inputCmd != string(PROMPT_COMMAND_COMMAND) {
				if cmd, ok := promptCommands[inputCmd]; ok {
					return cmd.handleCommand(strings.Join(fields[1:], " "), p, m)
				}
			}

			p.errormsg = "not a supershell - command"
			return false

		},
	}

	promptCommands[PROMPT_COMMAND_SHELL] = PromptCommand{
		renderPrefix: "$ ",
		renderHint:   "Bash/Powershell - Command",
		handleCommand: func(input string, _ *promptModalType, m *model) bool {

			// Make base work inside the prompt modal

			basePath := m.fileModel.filePanels[m.filePanelFocusIndex].location
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "windows":
				// On Windows, we use PowerShell with -Command flag for single command execution
				cmd = exec.Command("powershell.exe", "-Command", input)
			default:
				// On Unix-like systems, use bash/sh
				cmd = exec.Command("/bin/sh", "-c", input)
			}
			cmd.Dir = basePath // switch to the focused panel directory

			output, err := cmd.CombinedOutput()

			if err != nil {
				slog.Error("Command execution failed", "error", err, "output", string(output))
				return true
			}

			return true

		},
	}

	promptCommands[PROMPT_COMMAND_CD] = PromptCommand{
		renderPrefix: "CD > ",
		renderHint:   "CD current FilePanel",
		handleCommand: func(input string, p *promptModalType, m *model) bool {
			basePath := m.fileModel.filePanels[m.filePanelFocusIndex].location
			path := strings.TrimSpace(input)
			if !filepath.IsAbs(path) {
				path = basePath + string(os.PathSeparator) + path
			}

			if dir, err := os.Stat(path); err == nil {
				if dir.IsDir() {
					m.fileModel.filePanels[m.filePanelFocusIndex].location = path
					return true

				} else {
					p.errormsg = "not a directory"

				}

			} else {
				p.errormsg = "given path does not exist"

			}

			return false
		},
	}

	promptCommands[PROMPT_COMMAND_NEWFILEPANEL] = PromptCommand{
		renderPrefix: "OPEN > ",
		renderHint:   "new Filepanel at given path",
		handleCommand: func(input string, p *promptModalType, m *model) bool {

			basePath := m.fileModel.filePanels[m.filePanelFocusIndex].location
			path := strings.TrimSpace(input)
			if !filepath.IsAbs(path) {
				path = basePath + string(os.PathSeparator) + path
			}

			if dir, err := os.Stat(path); err == nil {
				if dir.IsDir() {
					m.createNewFilePanel(path)
					return true

				} else {
					p.errormsg = "not a directory"

				}
			} else {
				p.errormsg = "given path does not exist"

			}

			return false

		},
	}

	promptCommands[PROMPT_COMMAND_SPLIT] = PromptCommand{
		renderHint: "new Filepanel at current location",
		handleCommand: func(_ string, _ *promptModalType, m *model) bool {

			location := m.fileModel.filePanels[m.filePanelFocusIndex].location
			m.createNewFilePanel(location)

			return true
		},
	}

}

func (p *promptModalType) Open(m *model, cmdPrefix PromptCommandPrefix) {

	prompt, ok := promptCommands[cmdPrefix]
	if !ok {
		log.Fatalf("this should not happen during Runtime. Please fix your code: promptModel.Open called with invalid cmdPrefix")
	}

	if len(prompt.renderPrefix) == 0 {
		log.Fatalf("this should not happen during Runtime. Command '%s' is not meant to have text input.", cmdPrefix)
	}

	p.cmd = prompt

	p.textInput = textinput.New()
	p.textInput.Prompt = ""
	p.textInput.CharLimit = 156
	p.textInput.SetValue("")

	p.textInput.Cursor.Style = modalCursorStyle
	p.textInput.Cursor.TextStyle = modalStyle
	p.textInput.TextStyle = modalStyle
	p.textInput.PlaceholderStyle = modalStyle

	suggestions := make([]string, 0, len(promptCommands)-1)
	for cmd := range promptCommands {
		if PROMPT_COMMAND_COMMAND == cmd {
			continue
		}

		suggestions = append(suggestions, string(cmd))
	}

	p.textInput.SetSuggestions(suggestions)
	p.textInput.ShowSuggestions = true

	p.open = true
}

func (p *promptModalType) Close() {
	p.open = false
	p.errormsg = ""
	p.textInput.SetValue("")
}

func (p *promptModalType) Confirm(m *model) bool {
	return p.cmd.handleCommand(p.textInput.Value(), p, m)
}

func (p *promptModalType) Render(width int) string {

	var content, promptLine string
	text := p.textInput.Value()
	suggestions := p.textInput.CurrentSuggestion()

	if len(text) == 0 {

		suggestions = ""

		for _, s := range p.textInput.AvailableSuggestions() {
			suggestion := fmt.Sprintf("%s%*s%s", s, 10-len(s), "", promptCommands[s].renderHint)
			suggestions += fmt.Sprintf("%s%*s\n", suggestion, width-len(suggestion), "")
		}

	}

	content += fmt.Sprintf(
		"%s%*s\n\n",
		p.headline,
		width-len(p.headline), "",
	)

	promptLine += p.cmd.renderPrefix + text
	content += fmt.Sprintf("%s%*s\n%s\n", promptLine, width-len(promptLine), "", strings.Repeat(Config.BorderTop, width))
	content += fmt.Sprintf("%s%*s\n", suggestions, width-len(suggestions), "")
	if len(p.errormsg) > 0 {
		content += fmt.Sprintf("%s\n", strings.Repeat(Config.BorderTop, width))
		content += fmt.Sprintf("%s%*s", p.errormsg, width-len(p.errormsg), "")
	}

	return modalBorderStyle(1, width+2).Render(content)

}
