package internal

import (
	"os"
	"os/exec"
	"path/filepath"
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
	if !strings.HasSuffix(m.typingModal.textInput.Value(), "/") {
		path := filepath.Join(m.typingModal.location, m.typingModal.textInput.Value())
		path, _ = renameIfDuplicate(path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			outPutLog("Create item func (m *model)tion error", err)
		}
		f, err := os.Create(path)
		if err != nil {
			outPutLog("Create item func (m *model)tion create file error", err)
		}
		defer f.Close()
	} else {
		path := m.typingModal.location + "/" + m.typingModal.textInput.Value()
		err := os.MkdirAll(path, 0755)
		if err != nil {
			outPutLog("Create item func (m *model)tion create folder error", err)
		}
	}
	m.typingModal.open = false
	m.typingModal.textInput.Blur()
}

// Cancel rename file or directory
func (m *model) cancelReanem() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.rename.Blur()
	panel.renaming = false
	m.fileModel.renaming = false
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Connfirm rename file or directory
func (m *model) confirmRename() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	oldPath := panel.element[panel.cursor].location
	newPath := panel.location + "/" + panel.rename.Value()

	// Rename the file
	err := os.Rename(oldPath, newPath)
	if err != nil {
		outPutLog("Confirm func (m *model)tion rename error", err)
	}

	m.fileModel.renaming = false
	panel.rename.Blur()
	panel.renaming = false
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Cancel search, this will clear all searchbar input
func (m *model) cancelSearch() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.searchBar.Blur()
	panel.searchBar.SetValue("")
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
}

// Confirm search
func (m *model) confirmSearch() {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.searchBar.Blur()
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
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

func (m *model) enterCommandLine() {
	focusPanelDir := ""
	for _, panel := range m.fileModel.filePanels {
		if panel.focusType == focus {
			focusPanelDir = panel.location
		}
	}
	cd := "cd "+focusPanelDir+" && "
	cmd := exec.Command("/bin/sh", "-c", cd + m.commandLine.input.Value())
	_, err := cmd.CombinedOutput()
	m.commandLine.input.SetValue("")
	if err != nil {
		return
	}
	m.commandLine.input.Blur()
	footerHeight++
}

func (m *model) openSftpMenu() {
	if m.sftpPanelModal.open {
		m.sftpPanelModal.open = false
		return
	}

	m.sftpPanelModal.open = true
}