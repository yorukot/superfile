package internal

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/internal/utils"

	variable "github.com/yorukot/superfile/src/config"
)

// Back to parent directory
func (m *model) parentDirectory() {
	err := m.getFocusedFilePanel().parentDirectory()
	if err != nil {
		slog.Error("Error while changing to parent directory", "error", err)
	}
}

// Enter directory or open file with default application
// TODO: Unit test this
func (m *model) enterPanel() {
	panel := m.getFocusedFilePanel()

	if len(panel.element) == 0 {
		return
	}
	selectedItem := panel.getSelectedItem()
	if selectedItem.directory {
		// TODO : Propagate error out from this this function. Return here, instead of logging
		err := m.updateCurrentFilePanelDir(selectedItem.location)
		if err != nil {
			slog.Error("Error while changing to directory", "error", err, "target", selectedItem.location)
		}
		return
	}
	fileInfo, err := os.Lstat(selectedItem.location)
	if err != nil {
		slog.Error("Error while getting file info", "error", err)
		return
	}

	if fileInfo.Mode()&os.ModeSymlink != 0 {
		targetPath, symlinkErr := filepath.EvalSymlinks(selectedItem.location)
		if symlinkErr != nil {
			return
		}

		targetInfo, lstatErr := os.Lstat(targetPath)

		if lstatErr != nil {
			return
		}

		if targetInfo.IsDir() {
			err = m.updateCurrentFilePanelDir(targetPath)
			if err != nil {
				slog.Error("Error while changing to directory", "error", err, "target", targetPath)
			}
			return
		}
	}

	if variable.ChooserFile != "" {
		chooserErr := m.chooserFileWriteAndQuit(panel.element[panel.cursor].location)
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
	openCommand := "xdg-open"
	switch runtime.GOOS {
	case utils.OsDarwin:
		openCommand = "open"
	case utils.OsWindows:
		dllpath := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "rundll32.exe")
		dllfile := "url.dll,FileProtocolHandler"

		cmd := exec.Command(dllpath, dllfile, panel.element[panel.cursor].location)
		err := cmd.Start()
		if err != nil {
			slog.Error("Error while open file with", "error", err)
		}

		return
	}

	cmd := exec.Command(openCommand, panel.element[panel.cursor].location)
	utils.DetachFromTerminal(cmd)
	err := cmd.Start()
	if err != nil {
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
	panel.isFocused = true
}

// Select all item in the file panel (only work on select mode)
func (m *model) selectAllItem() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	for _, item := range panel.element {
		panel.selected = append(panel.selected, item.location)
	}
}

// Select the item where cursor located (only work on select mode)
func (panel *filePanel) singleItemSelect() {
	if len(panel.element) > 0 && panel.cursor >= 0 && panel.cursor < len(panel.element) {
		elementLocation := panel.element[panel.cursor].location

		if arrayContains(panel.selected, elementLocation) {
			// This is inefficient. Once you select 1000 items,
			// each select / deselect operation can take 1000 operations
			// It can be easily made constant time.
			// TODO : (performance)convert panel.selected to a set (map[string]struct{})
			panel.selected = removeElementByValue(panel.selected, elementLocation)
		} else {
			panel.selected = append(panel.selected, elementLocation)
		}
	}
}

// Toggle dotfile display or not
func (m *model) toggleDotFileController() {
	m.toggleDotFile = !m.toggleDotFile
	m.updatedToggleDotFile = true
	err := utils.WriteBoolFile(variable.ToggleDotFile, m.toggleDotFile)
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
	// TODO : Revisit this. Is this really need here, is this correct ?
	m.setHeightValues(m.fullHeight)
	// File preview panel requires explicit height update, unlike sidebar/file panels
	// which receive height as render parameters and update automatically on each frame
	if m.fileModel.filePreview.IsOpen() {
		m.setFilePreviewPanelSize()
		// Force re-render of preview content with new dimensions
		return m.getFilePreviewCmd(true)
	}
	return nil
}

// Focus on search bar
func (m *model) searchBarFocus() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.searchBar.Focused() {
		panel.searchBar.Blur()
	} else {
		panel.searchBar.Focus()
		m.firstTextInput = true
	}

	// config search bar width
	panel.searchBar.Width = m.fileModel.width - 4
}

func (m *model) sidebarSearchBarFocus() {
	if m.sidebarModel.SearchBarFocused() {
		// Ideally Code should never reach here. Once sidebar is focussed, we should
		// not cause sidebarSearchBarFocus() event by pressing search key
		slog.Error("sidebarSearchBarFocus() called on Focussed sidebar")
		m.sidebarModel.SearchBarBlur()
		return
	}
	m.sidebarModel.SearchBarFocus()
	m.firstTextInput = true
}
