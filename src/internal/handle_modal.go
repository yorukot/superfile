package internal

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Cancel typing modal e.g. create file or directory
func (m *model) cancelTypingModal() {
	m.typingModal.textInput.Blur()
	m.typingModal.open = false
}

// Close warn modal
func (m *model) cancelWarnModal() {
	m.warnModal.open = false
}

// Confirm to create file or directory
func (m *model) createItem() {
	// Reset the typingModal in all cases
	defer func() {
		m.typingModal.open = false
		m.typingModal.textInput.Blur()
	}()
	path := filepath.Join(m.typingModal.location, m.typingModal.textInput.Value())
	if !strings.HasSuffix(m.typingModal.textInput.Value(), string(filepath.Separator)) {
		path, _ = renameIfDuplicate(path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			outPutLog("createItem error during directory creation : ", err)
			return
		}
		f, err := os.Create(path)
		if err != nil {
			outPutLog("createItem error during file creation : ", err)
			return
		}
		defer f.Close()
	} else {
		// Directory creation
		err := os.MkdirAll(path, 0755)
		if err != nil {
			outPutLog("createItem error during directory creation : ", err)
			return
		}
	}
}

// Cancel rename file or directory
func (m *model) cancelRename() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.rename.Blur()
	panel.renaming = false
	m.fileModel.renaming = false
}

// Connfirm rename file or directory
func (m *model) confirmRename() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	oldPath := panel.element[panel.cursor].location
	newPath := filepath.Join(panel.location, panel.rename.Value())

	// Rename the file
	err := os.Rename(oldPath, newPath)
	if err != nil {
		outPutLog("confirmRename error during rename : ", err)
		// Dont return. We have to also reset the panel and model information
	}
	m.fileModel.renaming = false
	panel.rename.Blur()
	panel.renaming = false
}

func (m *model) openSortOptionsMenu() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.sortOptions.open = true
}

func (m *model) cancelSortOptions() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.sortOptions.cursor = panel.sortOptions.data.selected
	panel.sortOptions.open = false
}

func (m *model) confirmSortOptions() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.sortOptions.data.selected = panel.sortOptions.cursor
	panel.sortOptions.open = false
}

// Move the cursor up in the sort options menu
func (m *model) sortOptionsListUp() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.sortOptions.cursor > 0 {
		panel.sortOptions.cursor--
	} else {
		panel.sortOptions.cursor = len(panel.sortOptions.data.options) - 1
	}
}

// Move the cursor down in the sort options menu
func (m *model) sortOptionsListDown() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.sortOptions.cursor < len(panel.sortOptions.data.options)-1 {
		panel.sortOptions.cursor++
	} else {
		panel.sortOptions.cursor = 0
	}
}

func (m *model) toggleReverseSort() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.sortOptions.data.reversed = !panel.sortOptions.data.reversed
}

// Cancel search, this will clear all searchbar input
func (m *model) cancelSearch() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.searchBar.Blur()
	panel.searchBar.SetValue("")
}

// Confirm search. This will exit the search bar and filter the files
func (m *model) confirmSearch() {
	panel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.searchBar.Blur()
}

// Help menu panel list up
func (m *model) helpMenuListUp() {
	if m.helpMenu.cursor > 1 {
		m.helpMenu.cursor--
		if m.helpMenu.cursor < m.helpMenu.renderIndex {
			m.helpMenu.renderIndex--
			if m.helpMenu.data[m.helpMenu.cursor].subTitle != "" {
				m.helpMenu.renderIndex--
			}
		}
		if m.helpMenu.data[m.helpMenu.cursor].subTitle != "" {
			m.helpMenu.cursor--
		}
	} else {
		m.helpMenu.cursor = len(m.helpMenu.data) - 1
		m.helpMenu.renderIndex = len(m.helpMenu.data) - m.helpMenu.height
	}
}

// Help menu panel list down
func (m *model) helpMenuListDown() {
	if len(m.helpMenu.data) == 0 {
		return
	}

	if m.helpMenu.cursor < len(m.helpMenu.data)-1 {
		m.helpMenu.cursor++
		if m.helpMenu.cursor > m.helpMenu.renderIndex+m.helpMenu.height-1 {
			m.helpMenu.renderIndex++
			if m.helpMenu.data[m.helpMenu.cursor].subTitle != "" {
				m.helpMenu.renderIndex++
			}
		}
		if m.helpMenu.data[m.helpMenu.cursor].subTitle != "" {
			m.helpMenu.cursor++
		}
	} else {
		m.helpMenu.cursor = 1
		m.helpMenu.renderIndex = 0
	}
}

// Toggle help menu
func (m *model) openHelpMenu() {
	if m.helpMenu.open {
		m.helpMenu.open = false
		return
	}

	m.helpMenu.open = true
}

// Quit help menu
func (m *model) quitHelpMenu() {
	m.helpMenu.open = false
}

// Command line
func (m *model) openCommandLine() {
	m.firstTextInput = true
	footerHeight--
	m.commandLine.input = generateCommandLineInputBox()
	m.commandLine.input.Width = m.fullWidth - 3
	m.commandLine.input.Focus()
}

func (m *model) closeCommandLine() {
	footerHeight++
	m.commandLine.input.SetValue("")
	m.commandLine.input.Blur()
}

// Exec a command line input inside the pointing file dir. Like opening the
// focused file in the text editor
func (m *model) enterCommandLine() {
	focusPanelDir := ""
	for _, panel := range m.fileModel.filePanels {
		if panel.focusType == focus {
			focusPanelDir = panel.location
		}
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// On Windows, we use PowerShell with -Command flag for single command execution
		cmd = exec.Command("powershell.exe", "-Command", m.commandLine.input.Value())
	default:
		// On Unix-like systems, use bash/sh
		cmd = exec.Command("/bin/sh", "-c", m.commandLine.input.Value())
	}

	cmd.Dir = focusPanelDir // switch to the focused panel directory

	output, err := cmd.CombinedOutput()

	if err != nil {
		slog.Error("Command execution failed", "error", err, "output", string(output))
		return
	}

	m.commandLine.input.SetValue("")
	m.commandLine.input.Blur()
	footerHeight++
}
