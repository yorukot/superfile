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
// The first element of the array must be either "foreground" or "background".
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

		err := c.Start()
		if err != nil {
			return m.warnModalForCustomCommand("Execution Error", "Failed to start command '"+cmdPayload[0]+"': "+err.Error())
		}

		go func() {
			_ = c.Wait()
		}()
		return nil

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
	}

	return nil
}

// prepareCommand builds the exec.Cmd structure and replaces templates with file details.
func (m *model) prepareCommand(rawArgs []string) *exec.Cmd {
	panel := m.fileModel.GetFocusedFilePanel()
	if panel.EmptyOrInvalid() {
		return nil
	}

	selectedItem := panel.GetFocusedItem()
	selectedFile := selectedItem.Location
	fileName := filepath.Base(selectedFile)
	selectedDir := panel.Location

	processedArgs := make([]string, len(rawArgs))
	for i, arg := range rawArgs {
		arg = strings.ReplaceAll(arg, "{filepath}", selectedFile)
		arg = strings.ReplaceAll(arg, "{filename}", fileName)
		arg = strings.ReplaceAll(arg, "{dir}", selectedDir)
		processedArgs[i] = arg
	}

	return exec.Command(processedArgs[0], processedArgs[1:]...)
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
