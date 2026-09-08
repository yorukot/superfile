package internal

import (
	"os/exec"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/notify"
)

// executeCustomCommand executes a custom command defined in CustomCommands.
func (m *model) executeCustomCommand(cmdName string) tea.Cmd {
	rawArgs, ok := common.Config.CustomCommands[cmdName]
	if !ok || len(rawArgs) < 2 {
		return m.warnModalForCustomCommand("Configuration Error", "Command '"+cmdName+"' is not found or invalid in config.toml.")
	}

	mode := rawArgs[0]
	cmdPayload := rawArgs[1:]

	switch mode {
	case "background":
		c := m.prepareCommand(cmdPayload)
		if c == nil {
			return m.warnModalForCustomCommand("Execution Error", "Failed to resolve selected file path for '"+cmdName+"'.")
		}

		return func() tea.Msg {
			err := c.Start()
			if err != nil {
				reqID := m.nextIoReqCnt()
				notifyModel := notify.New(true,
					"Execution Error",
					"Failed to start command '"+cmdPayload[0]+"': "+err.Error(),
					notify.NoAction)
				return NewNotifyModalMsg(notifyModel, reqID)
			}

			err = c.Wait()
			if err != nil {
				reqID := m.nextIoReqCnt()
				notifyModel := notify.New(true,
					"Execution Error",
					"Command '"+cmdPayload[0]+"' exited with error: "+err.Error(),
					notify.NoAction)
				return NewNotifyModalMsg(notifyModel, reqID)
			}
			return nil
		}

	case "foreground":
		c := m.prepareCommand(cmdPayload)
		if c == nil {
			return m.warnModalForCustomCommand("Execution Error", "Failed to resolve selected file path for '"+cmdName+"'.")
		}

		return tea.ExecProcess(c, func(err error) tea.Msg {
			if err != nil {
				reqID := m.nextIoReqCnt()
				notifyModel := notify.New(true,
					"Execution Error",
					"Failed to run command '"+cmdPayload[0]+"': "+err.Error(),
					notify.NoAction)
				return NewNotifyModalMsg(notifyModel, reqID)
			}
			return nil
		})

	default:
		return m.warnModalForCustomCommand("Configuration Error", "Invalid mode '"+mode+"' for command '"+cmdName+"'. Must be 'foreground' or 'background'.")
	}
}

// prepareCommand builds the exec.Cmd structure and replaces templates with file details.
func (m *model) prepareCommand(rawArgs []string) *exec.Cmd {
	if len(rawArgs) == 0 {
		return nil
	}

	panel := m.fileModel.GetFocusedFilePanel()
	if panel.EmptyOrInvalid() {
		return nil
	}

	selectedItem := panel.GetFocusedItem()
	selectedFile := filepath.ToSlash(selectedItem.Location)
	fileName := filepath.Base(selectedFile)
	selectedDir := filepath.ToSlash(panel.Location)

	binary := strings.ToLower(rawArgs[0])
	isShell := binary == "sh" || binary == "bash" || binary == "zsh" || binary == "dash" || binary == "powershell" || binary == "pwsh"

	// Find the shell command string index right after "-c" or "-command"
	shellCmdIdx := -1
	if isShell {
		for i, arg := range rawArgs {
			lowerArg := strings.ToLower(arg)
			if (lowerArg == "-c" || lowerArg == "-command") && i+1 < len(rawArgs) {
				shellCmdIdx = i + 1
				break
			}
		}
	}

	processedArgs := make([]string, len(rawArgs))
	for i, arg := range rawArgs {
		if isShell && i == shellCmdIdx {
			arg = replaceShellPlaceholder(arg, "{filepath}", selectedFile)
			arg = replaceShellPlaceholder(arg, "{filename}", fileName)
			arg = replaceShellPlaceholder(arg, "{dir}", selectedDir)
		} else {
			arg = strings.ReplaceAll(arg, "{filepath}", selectedFile)
			arg = strings.ReplaceAll(arg, "{filename}", fileName)
			arg = strings.ReplaceAll(arg, "{dir}", selectedDir)
		}
		processedArgs[i] = arg
	}

	return exec.Command(processedArgs[0], processedArgs[1:]...)
}

// replaceShellPlaceholder safely substitutes placeholders in shell arguments.
func replaceShellPlaceholder(arg, placeholder, val string) string {
	if !strings.Contains(arg, placeholder) {
		return arg
	}

	// If placeholder is already enclosed in quotes, substitute directly without adding extra quotes
	if strings.Contains(arg, `"`+placeholder+`"`) || strings.Contains(arg, `'`+placeholder+`'`) {
		return strings.ReplaceAll(arg, placeholder, val)
	}

	// Unquoted placeholder: safely wrap value in single quotes
	escapedVal := "'" + strings.ReplaceAll(val, "'", "'\\''") + "'"
	return strings.ReplaceAll(arg, placeholder, escapedVal)
}

// warnModalForCustomCommand creates a tea.Cmd that displays a standard superfile warning modal.
func (m *model) warnModalForCustomCommand(title string, content string) tea.Cmd {
	reqID := m.nextIoReqCnt()
	return func() tea.Msg {
		notifyModel := notify.New(true,
			title,
			content,
			notify.NoAction)
		return NewNotifyModalMsg(notifyModel, reqID)
	}
}