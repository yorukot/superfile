package internal

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

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

// Enter directory or open file with default application
// TODO: Unit test this
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
		chooserErr := m.chooserFileWriteAndQuit(panel.GetFocusedItem().Location)
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

	cmd := exec.Command(openCommand, filePath)
	utils.DetachFromTerminal(cmd)
	err := cmd.Start()
	if err != nil {
		// TODO: This kind of errors should go to user facing pop ups
		slog.Error("Error while open file with", "error", err)
	}
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
	panel.SearchBar.Width = m.fileModel.SinglePanelWidth - common.InnerPadding
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
