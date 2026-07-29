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

		// Return an asynchronous command to run the process and monitor exit status
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

			// Wait for completion in background to catch any exit errors safely
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
		// Explicitly handle unsupported execution modes instead of falling through silently
		return m.warnModalForCustomCommand("Configuration Error", "Invalid mode '"+mode+"' for command '"+cmdName+"'. Must be 'foreground' or 'background'.")
	}
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

	// Shell injection prevention: escape arguments if explicitly invoking a POSIX shell
	binary := rawArgs[0]
	isShell := binary == "sh" || binary == "bash" || binary == "zsh" || binary == "dash"

	processedArgs := make([]string, len(rawArgs))
	for i, arg := range rawArgs {
		fileVal := selectedFile
		nameVal := fileName
		dirVal := selectedDir

		// If running via shell, wrap placeholders in single quotes and escape existing single quotes
		if isShell {
			fileVal = "'" + strings.ReplaceAll(selectedFile, "'", "'\\''") + "'"
			nameVal = "'" + strings.ReplaceAll(fileName, "'", "'\\''") + "'"
			dirVal = "'" + strings.ReplaceAll(selectedDir, "'", "'\\''") + "'"
		}

		arg = strings.ReplaceAll(arg, "{filepath}", fileVal)
		arg = strings.ReplaceAll(arg, "{filename}", nameVal)
		arg = strings.ReplaceAll(arg, "{dir}", dirVal)
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
