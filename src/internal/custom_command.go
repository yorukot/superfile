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
	if len(rawArgs) == 0 {
		return nil
	}

	panel := m.fileModel.GetFocusedFilePanel()
	if panel.EmptyOrInvalid() {
		return nil
	}

	selectedItem := panel.GetFocusedItem()

	// use filepath.ToSlash to convert '\' into '/' for Windows.
	selectedFile := filepath.ToSlash(selectedItem.Location)
	fileName := filepath.Base(selectedFile)
	selectedDir := filepath.ToSlash(panel.Location)

	// Shell injection prevention: identify if invoking a POSIX shell
	binary := rawArgs[0]
	isShell := binary == "sh" || binary == "bash" || binary == "zsh" || binary == "dash"

	// Find the shell command string index right after "-c" flag if present
	shellCmdIdx := -1
	if isShell {
		for i, arg := range rawArgs {
			if arg == "-c" && i+1 < len(rawArgs) {
				shellCmdIdx = i + 1
				break
			}
		}
	}

	processedArgs := make([]string, len(rawArgs))
	for i, arg := range rawArgs {
		if isShell && i == shellCmdIdx {
			// Apply contextual shell quoting ONLY to the command string argument after -c
			arg = replaceShellPlaceholder(arg, "{filepath}", selectedFile)
			arg = replaceShellPlaceholder(arg, "{filename}", fileName)
			arg = replaceShellPlaceholder(arg, "{dir}", selectedDir)
		} else {
			// Direct substitution for binary name, flags, script paths, and positional arguments
			arg = strings.ReplaceAll(arg, "{filepath}", selectedFile)
			arg = strings.ReplaceAll(arg, "{filename}", fileName)
			arg = strings.ReplaceAll(arg, "{dir}", selectedDir)
		}
		processedArgs[i] = arg
	}

	return exec.Command(processedArgs[0], processedArgs[1:]...)
}

// replaceShellPlaceholder replaces a template placeholder in a shell command argument,
// taking into account surrounding single or double quotes for proper escaping.
func replaceShellPlaceholder(arg, placeholder, val string) string {
	singleQuoted := "'" + placeholder + "'"
	doubleQuoted := `"` + placeholder + `"`

	if strings.Contains(arg, singleQuoted) {
		// Surrounded by single quotes in arg: escape single quotes inside val, no extra outer quotes
		escapedVal := strings.ReplaceAll(val, "'", "'\\''")
		return strings.ReplaceAll(arg, singleQuoted, "'"+escapedVal+"'")
	}

	if strings.Contains(arg, doubleQuoted) {
		// Surrounded by double quotes in arg: escape special double-quote characters (\, ", $, `)
		escapedVal := val
		escapedVal = strings.ReplaceAll(escapedVal, `\`, `\\`)
		escapedVal = strings.ReplaceAll(escapedVal, `"`, `\"`)
		escapedVal = strings.ReplaceAll(escapedVal, `$`, `\$`)
		escapedVal = strings.ReplaceAll(escapedVal, "`", "\\`")
		return strings.ReplaceAll(arg, doubleQuoted, `"`+escapedVal+`"`)
	}

	// Unquoted placeholder: wrap in single quotes and escape internal single quotes
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