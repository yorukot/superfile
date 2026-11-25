package internal

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/pkg/utils"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
)

// Back to parent directory
func (m *model) parentDirectory() {
	err := m.getFocusedFilePanel().ParentDirectory()
	if err != nil {
		slog.Error("Error while changing to parent directory", "error", err)
	}
}

// enterPanel enters a directory, opens a file with default application, or writes
// selected/focused paths to the chooser file and exits if --chooser-file mode is enabled.
func (m *model) enterPanel() {
	panel := m.getFocusedFilePanel()

	if panel.Empty() {
		return
	}
	selectedItem := panel.GetFocusedItem()
	if selectedItem.Directory {
		targetPath := selectedItem.Location

		if selectedItem.Info.Mode()&os.ModeSymlink != 0 {
			var symlinkErr error
			targetPath, symlinkErr = filepath.EvalSymlinks(targetPath)
			if symlinkErr != nil {
				return
			}

			// targetPath shouldn't be a link now, so Stat and Lstat should be same
			if targetInfo, lstatErr := os.Lstat(targetPath); lstatErr != nil || !targetInfo.IsDir() {
				return
			}
		}
		// TODO : Propagate error out from this this function. Return here, instead of logging
		err := m.updateCurrentFilePanelDir(targetPath)
		if err != nil {
			slog.Error("Error while changing to directory", "error", err, "target", targetPath)
		}
		return
	}

	if variable.ChooserFile != "" {
		chooserErr := m.chooserFileWriteAndQuit(m.getSelectedOrFocusedPaths(panel))
		if chooserErr == nil {
			return
		}
		// Continue with preview if file is not writable
		slog.Error("Error while writing to chooser file, continuing with file open", "error", chooserErr)
	}
	m.executeOpenCommand()
}

func (m *model) executeOpenCommand() {
	panel := m.getFocusedFilePanel()

	filePath := panel.GetFocusedItem().Location

	openCommand := "xdg-open"
	switch runtime.GOOS {
	case utils.OsDarwin:
		openCommand = "open"
	case utils.OsWindows:
		dllpath := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
		dllfile := "url.dll,FileProtocolHandler"

		//nolint:gosec // Uses Windows system handler to open the selected file.
		cmd := exec.Command(dllpath, dllfile, filePath)
		err := cmd.Start()
		if err != nil {
			slog.Error("Error while open file with", "error", err)
		}

		return
	}

	// For now open_with works only for mac and linux
	// TODO: Make it in parity with windows.
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
	if extEditor, ok := common.Config.OpenWith[ext]; ok {
		openCommand = extEditor
	}

	// The configured command may carry its own arguments (e.g. "mpv --force-window=yes"),
	// the same way editor commands do. The file path is appended after them.
	program, args, err := buildOpenArgs(openCommand, filePath)
	if err != nil {
		slog.Error("Invalid open command", "error", err, "command", openCommand, "ext", ext)
		return
	}

	cmd := exec.Command(program, args...) //nolint:gosec // Open command is intentionally user-configurable.
	utils.DetachFromTerminal(cmd)
	err = cmd.Start()
	if err != nil {
		// TODO: This kind of errors should go to user facing pop ups
		slog.Error("Error while open file with", "error", err)
	}
}

// buildOpenArgs splits a configured open command into the program to run and
// its arguments, with filePath appended last. The command may include its own
// arguments, e.g. "mpv --force-window=yes".
//
// A value that names an existing file is kept as a single executable path, so
// that app bundles such as "/Applications/My Player.app/Contents/MacOS/player"
// are not split at their spaces. When such a path also needs arguments, quote
// it: "\"/Applications/My Player.app/Contents/MacOS/player\" --fullscreen".
func buildOpenArgs(openCommand, filePath string) (string, []string, error) {
	command := strings.TrimSpace(openCommand)
	if command == "" {
		return "", nil, errors.New("open command is empty")
	}
	if isExistingPath(command) {
		return command, []string{filePath}, nil
	}
	parts := splitOpenCommand(command)
	if len(parts) == 0 {
		return "", nil, errors.New("open command is empty")
	}
	return parts[0], append(parts[1:], filePath), nil
}

// splitOpenCommand splits a command line into tokens on unquoted whitespace.
// Text inside single or double quotes is kept together, so a program path that
// contains spaces can be passed as the first token.
func splitOpenCommand(command string) []string {
	var (
		tokens   []string
		current  strings.Builder
		quote    rune
		hasToken bool
	)
	for _, r := range command {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
			hasToken = true
		case r == ' ' || r == '\t':
			if hasToken {
				tokens = append(tokens, current.String())
				current.Reset()
				hasToken = false
			}
		default:
			current.WriteRune(r)
			hasToken = true
		}
	}
	if hasToken {
		tokens = append(tokens, current.String())
	}
	return tokens
}

// isExistingPath reports whether command names an existing file that can be
// launched as it is. Values without spaces or quotes are skipped, since
// splitting cannot change them.
func isExistingPath(command string) bool {
	if !strings.ContainsAny(command, " \t\"'") {
		// No spaces or quoting, so splitting cannot change the value.
		return false
	}
	info, err := os.Stat(command)
	return err == nil && !info.IsDir()
}

// Switch to the directory where the sidebar cursor is located
func (m *model) sidebarSelectDirectory() {
	// We can't do this when we have only divider directories
	// m.sidebarModel.directories[m.sidebarModel.cursor].location would point to a divider dir.
	if m.sidebarModel.NoActualDir() {
		return
	}
	// TODO(Refactor): Move this to a function m.ResetFocus()
	m.focusPanel = nonePanelFocus
	panel := m.getFocusedFilePanel()

	err := m.updateCurrentFilePanelDir(m.sidebarModel.GetCurrentDirectoryLocation())
	if err != nil {
		slog.Error("Error switching to sidebar directory", "error", err)
	}
	panel.IsFocused = true
}

// Toggle dotfile display or not
func (m *model) toggleDotFileController() {
	m.fileModel.ToggleDotFile()
	err := utils.WriteBoolFile(variable.ToggleDotFile, m.fileModel.DisplayDotFiles)
	if err != nil {
		slog.Error("Error while updating toggleDotFile data", "error", err)
	}
}

// Toggle dotfile display or not
func (m *model) toggleFooterController() tea.Cmd {
	m.toggleFooter = !m.toggleFooter
	err := utils.WriteBoolFile(variable.ToggleFooter, m.toggleFooter)
	if err != nil {
		slog.Error("Error while updating toggleFooter data", "error", err)
	}
	m.setHeightValues()
	return m.updateComponentDimensions()
}

// Focus on search bar
func (m *model) searchBarFocus() {
	panel := m.getFocusedFilePanel()
	if panel.SearchBar.Focused() {
		panel.SearchBar.Blur()
	} else {
		panel.SearchBar.Focus()
		m.firstTextInput = true
	}

	// config search bar width
	panel.SearchBar.SetWidth(m.fileModel.SinglePanelWidth - common.InnerPadding)
}

func (m *model) sidebarSearchBarFocus() {
	if m.sidebarModel.SearchBarFocused() {
		// Ideally Code should never reach here. Once sidebar is focused, we should
		// not cause sidebarSearchBarFocus() event by pressing search key
		slog.Error("sidebarSearchBarFocus() called on Focused sidebar")
		m.sidebarModel.SearchBarBlur()
		return
	}
	m.sidebarModel.SearchBarFocus()
	m.firstTextInput = true
}
