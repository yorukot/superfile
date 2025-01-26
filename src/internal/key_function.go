package internal

import tea "github.com/charmbracelet/bubbletea"

func containsKey(v string, a []string) string {
	for _, i := range a {
		if i == v {
			return v
		}
	}
	return ""
}

// mainKey handles most of key commands in the regular state of the application. For 
// keys that performs actions in multiple panels, like going up or down,
// check the state of model m and handle properly.
func (m *model) mainKey(msg string, cmd tea.Cmd) ( tea.Cmd) {
	switch msg {

    // If move up Key is pressed, check the current state and executes
	case containsKey(msg, hotkeys.ListUp):
		if m.focusPanel == sidebarFocus {
			m.controlSideBarListUp(false)
		} else if m.focusPanel == processBarFocus {
			m.controlProcessbarListUp(false)
		} else if m.focusPanel == metadataFocus {
			m.controlMetadataListUp(false)
		} else if m.focusPanel == nonePanelFocus {
			m.controlFilePanelListUp(false)
			m.fileMetaData.renderIndex = 0
			go func() {
				m.returnMetaData()
			}()
		}

    // If move down Key is pressed, check the current state and executes
	case containsKey(msg, hotkeys.ListDown):
		if m.focusPanel == sidebarFocus {
			m.controlSideBarListDown(false)
		} else if m.focusPanel == processBarFocus {
			m.controlProcessbarListDown(false)
		} else if m.focusPanel == metadataFocus {
			m.controlMetadataListDown(false)
		} else if m.focusPanel == nonePanelFocus {
			m.controlFilePanelListDown(false)
			m.fileMetaData.renderIndex = 0
			go func() {
				m.returnMetaData()
			}()
		}

  case containsKey(msg, hotkeys.PageUp):
    m.controlFilePanelPgUp()

  case containsKey(msg, hotkeys.PageDown):
    m.controlFilePanelPgDown()

	case containsKey(msg, hotkeys.ChangePanelMode):
		m.changeFilePanelMode()

	case containsKey(msg, hotkeys.NextFilePanel):
		m.nextFilePanel()

	case containsKey(msg, hotkeys.PreviousFilePanel):
		m.previousFilePanel()

	case containsKey(msg, hotkeys.CloseFilePanel):
		m.closeFilePanel()

	case containsKey(msg, hotkeys.CreateNewFilePanel):
		m.createNewFilePanel()

	case containsKey(msg, hotkeys.ToggleFilePreviewPanel):
		m.toggleFilePreviewPanel()

	case containsKey(msg, hotkeys.FocusOnSidebar):
		m.focusOnSideBar()

	case containsKey(msg, hotkeys.FocusOnProcessBar):
		m.focusOnProcessBar()

	case containsKey(msg, hotkeys.FocusOnMetaData):
		m.focusOnMetadata()
		go func() {
			m.returnMetaData()
		}()

	case containsKey(msg, hotkeys.PasteItems):
		go func() {
			m.pasteItem()
		}()

	case containsKey(msg, hotkeys.FilePanelItemCreate):
		m.panelCreateNewFile()
	case containsKey(msg, hotkeys.PinnedDirectory):
		m.pinnedDirectory()

	case containsKey(msg, hotkeys.ToggleDotFile):
		m.toggleDotFileController()

	case containsKey(msg, hotkeys.ToggleFooter):
		m.toggleFooterController()

	case containsKey(msg, hotkeys.ExtractFile):
		go func() {
			m.extractFile()
		}()

	case containsKey(msg, hotkeys.CompressFile):
		go func() {
			m.compressFile()
		}()

	case containsKey(msg, hotkeys.OpenHelpMenu):
		m.openHelpMenu()

	case containsKey(msg, hotkeys.OpenSortOptionsMenu):
		m.openSortOptionsMenu()

	case containsKey(msg, hotkeys.ToggleReverseSort):
		m.toggleReverseSort()

	case containsKey(msg, hotkeys.OpenCommandLine):
		m.openCommandLine()

	case containsKey(msg, hotkeys.OpenFileWithEditor):
		cmd = m.openFileWithEditor()

	case containsKey(msg, hotkeys.OpenCurrentDirectoryWithEditor):
		cmd = m.openDirectoryWithEditor()

	default:
		m.normalAndBrowserModeKey(msg)
	}

	return cmd
}

func (m *model) normalAndBrowserModeKey(msg string) {
	// if not focus on the filepanel return
	if m.fileModel.filePanels[m.filePanelFocusIndex].focusType != focus {
		if m.focusPanel == sidebarFocus && (msg == containsKey(msg, hotkeys.Confirm)) {
			m.sidebarSelectDirectory()
		}
		return
	}
	// Check if in the select mode and focusOn filepanel
	if m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
		switch msg {
		case containsKey(msg, hotkeys.Confirm):
			m.singleItemSelect()
		case containsKey(msg, hotkeys.FilePanelSelectModeItemsSelectUp):
			m.itemSelectUp(false)
		case containsKey(msg, hotkeys.FilePanelSelectModeItemsSelectDown):
			m.itemSelectDown(false)
		case containsKey(msg, hotkeys.DeleteItems):
			go func() {
				m.deleteItemWarn()
			}()
		case containsKey(msg, hotkeys.CopyItems):
			m.copyMultipleItem(false)
		case containsKey(msg, hotkeys.CutItems):
			m.copyMultipleItem(true)
		case containsKey(msg, hotkeys.FilePanelSelectAllItem):
			m.selectAllItem()
		}
		return
	}

	switch msg {
	case containsKey(msg, hotkeys.Confirm):
		m.enterPanel()
	case containsKey(msg, hotkeys.ParentDirectory):
		m.parentDirectory()
	case containsKey(msg, hotkeys.DeleteItems):
		go func() {
			m.deleteItemWarn()
		}()
	case containsKey(msg, hotkeys.CopyItems):
		m.copySingleItem(false)
	case containsKey(msg, hotkeys.CutItems):
		m.copySingleItem(true)
	case containsKey(msg, hotkeys.FilePanelItemRename):
		m.panelItemRename()
	case containsKey(msg, hotkeys.SearchBar):
		m.searchBarFocus()
	case containsKey(msg, hotkeys.CopyPath):
		m.copyPath()
  case containsKey(msg, hotkeys.CopyPWD):
    m.copyPWD()
	}
}

// Check the hotkey to cancel operation or create file 
func (m *model) typingModalOpenKey(msg string) {
	switch msg {
	case containsKey(msg, hotkeys.CancelTyping):
		m.cancelTypingModal()
	case containsKey(msg, hotkeys.ConfirmTyping):
		m.createItem()
	}
}

func (m *model) warnModalOpenKey(msg string) {
	switch msg {
	case containsKey(msg, hotkeys.Quit), containsKey(msg, hotkeys.CancelTyping):
		m.cancelWarnModal()
		if m.warnModal.warnType == confirmRenameItem {
			m.cancelRename()
		}
	case containsKey(msg, hotkeys.Confirm):
		m.warnModal.open = false
		switch m.warnModal.warnType {
		case confirmDeleteItem:
			panel := m.fileModel.filePanels[m.filePanelFocusIndex]
			if m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
				if !hasTrash || isExternalDiskPath(panel.location) {
					go func() {
						m.completelyDeleteMultipleItems()
						m.fileModel.filePanels[m.filePanelFocusIndex].selected = m.fileModel.filePanels[m.filePanelFocusIndex].selected[:0]
					}()
				} else {
					go func() {
						m.deleteMultipleItems()
						m.fileModel.filePanels[m.filePanelFocusIndex].selected = m.fileModel.filePanels[m.filePanelFocusIndex].selected[:0]
					}()
				}
			} else {
				if !hasTrash || isExternalDiskPath(panel.location) {
					go func() {
						m.completelyDeleteSingleItem()
					}()
				} else {
					go func() {
						m.deleteSingleItem()
					}()
				}

			}
		case confirmRenameItem:
            m.confirmRename()
		}
	}
}

// Handle key input to confirm or cancel and close quiting warn in SPF
func (m *model) confirmToQuitSuperfile(msg string) bool {
	switch msg {
	case containsKey(msg, hotkeys.Quit), containsKey(msg, hotkeys.CancelTyping):
		m.cancelWarnModal()
		m.confirmToQuit = false
		return false
	case containsKey(msg, hotkeys.Confirm):
		return true
	default:
		return false
	}
}

// Handles key inputs inside sort options menu
func (m *model) sortOptionsKey(msg string) {
	switch msg {
	case containsKey(msg, hotkeys.OpenSortOptionsMenu):
		m.cancelSortOptions()
	case containsKey(msg, hotkeys.Quit):
		m.cancelSortOptions()
	case containsKey(msg, hotkeys.Confirm):
		m.confirmSortOptions()
	case containsKey(msg, hotkeys.ListUp):
		m.sortOptionsListUp()
	case containsKey(msg, hotkeys.ListDown):
		m.sortOptionsListDown()
	}
}

func (m *model) renamingKey(msg string) {
	switch msg {
	case containsKey(msg, hotkeys.CancelTyping):
		m.cancelRename()
	case containsKey(msg, hotkeys.ConfirmTyping):
		if m.IsRenamingConflicting() {
			m.warnModalForRenaming()
		} else {
			m.confirmRename()
		}
	}
}

// Check the key input and cancel or confirms the search
func (m *model) focusOnSearchbarKey(msg string) {
	switch msg {
	case containsKey(msg, hotkeys.CancelTyping):
		m.cancelSearch()
	case containsKey(msg, hotkeys.ConfirmTyping):
		m.confirmSearch()
	}
}

// Check hotkey input in help menu. Possible actions are moving up, down
// and quiting the menu
func (m *model) helpMenuKey(msg string) {
	switch msg {
	case containsKey(msg, hotkeys.ListUp):
		m.helpMenuListUp()
	case containsKey(msg, hotkeys.ListDown):
		m.helpMenuListDown()
	case containsKey(msg, hotkeys.Quit):
		m.quitHelpMenu()
	}
}

// Handle command line keys closing or entering command line
func (m *model) commandLineKey(msg string) {
	switch msg {
	case containsKey(msg, hotkeys.CancelTyping):
		m.closeCommandLine()
	case containsKey(msg, hotkeys.ConfirmTyping):
		m.enterCommandLine()
	}
}
