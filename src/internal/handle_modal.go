package internal

import (
	"os"
	"path/filepath"
)

// Cancel typing modal e.g. create file or directory
func cancelTypingModal(m model) model {
	m.typingModal.textInput.Blur()
	m.typingModal.open = false
	return m
}

// Close warn modal
func cancelWarnModal(m model) model {
	m.warnModal.open = false
	return m
}

// Confirm to create file or directory
func createItem(m model) model {
	if m.typingModal.itemType == newFile {
		path := filepath.Join(m.typingModal.location, m.typingModal.textInput.Value())
		path, _ = renameIfDuplicate(path)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			outPutLog("Create item function error", err)
		}
		f, err := os.Create(path)
		if err != nil {
			outPutLog("Create item function create file error", err)
		}
		defer f.Close()
	} else {
		path := m.typingModal.location + "/" + m.typingModal.textInput.Value()
		err := os.MkdirAll(path, 0755)
		if err != nil {
			outPutLog("Create item function create folder error", err)
		}
	}
	m.typingModal.open = false
	m.typingModal.textInput.Blur()
	return m
}

// Cancel rename file or directory
func cancelReanem(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.rename.Blur()
	panel.renaming = false
	m.fileModel.renaming = false
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

// Connfirm rename file or directory
func confirmRename(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	oldPath := panel.element[panel.cursor].location
	newPath := panel.location + "/" + panel.rename.Value()

	// Rename the file
	err := os.Rename(oldPath, newPath)
	if err != nil {
		outPutLog("Confirm function rename error", err)
	}

	m.fileModel.renaming = false
	panel.rename.Blur()
	panel.renaming = false
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

// Cancel search, this will clear all searchbar input
func cancelSearch(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.searchBar.Blur()
	panel.searchBar.SetValue("")
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

// Confirm search
func confirmSearch(m model) model {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	panel.searchBar.Blur()
	m.fileModel.filePanels[m.filePanelFocusIndex] = panel
	return m
}

// Help menu panel list up
func helpMenuListUp(m model) model {
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

	return m
}

// Help menu panel list down
func helpMenuListDown(m model) model {
	if len(m.helpMenu.data) == 0 {
		return m
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

	return m
}

// Quit help menu
func quitHelpMenu(m model) model {
	m.helpMenu.open = false
	return m
}
