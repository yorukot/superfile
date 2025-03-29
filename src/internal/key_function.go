package internal

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	variable "github.com/yorukot/superfile/src/config"
)

// Todo : Replace usage of this with direct slices.Contains call
// This fucntion is not required
func containsKey(v string, a []string) string {
	if slices.Contains(a, v) {
		return v
	}
	return ""
}

// mainKey handles most of key commands in the regular state of the application. For
// keys that performs actions in multiple panels, like going up or down,
// check the state of model m and handle properly.
func (m *model) mainKey(msg string, cmd tea.Cmd) tea.Cmd {

	switch msg {

	// If move up Key is pressed, check the current state and executes
	case containsKey(msg, Hotkeys.ListUp):
		if m.focusPanel == sidebarFocus {
			m.sidebarModel.listUp(m.mainPanelHeight)
		} else if m.focusPanel == processBarFocus {
			m.processBarModel.listUp(m.footerHeight)
		} else if m.focusPanel == metadataFocus {
			m.fileMetaData.listUp()
		} else if m.focusPanel == nonePanelFocus {
			m.fileModel.filePanels[m.filePanelFocusIndex].listUp(m.mainPanelHeight)
			m.fileMetaData.renderIndex = 0
			go func() {
				m.returnMetaData()
			}()
		}

		// If move down Key is pressed, check the current state and executes
	case containsKey(msg, Hotkeys.ListDown):
		if m.focusPanel == sidebarFocus {
			m.sidebarModel.listDown(m.mainPanelHeight)
		} else if m.focusPanel == processBarFocus {
			m.processBarModel.listDown(m.footerHeight)
		} else if m.focusPanel == metadataFocus {
			m.fileMetaData.listDown()
		} else if m.focusPanel == nonePanelFocus {
			m.fileModel.filePanels[m.filePanelFocusIndex].listDown(m.mainPanelHeight)
			m.fileMetaData.renderIndex = 0
			go func() {
				m.returnMetaData()
			}()
		}

	case containsKey(msg, Hotkeys.PageUp):
		m.fileModel.filePanels[m.filePanelFocusIndex].pgUp(m.mainPanelHeight)

	case containsKey(msg, Hotkeys.PageDown):
		m.fileModel.filePanels[m.filePanelFocusIndex].pgDown(m.mainPanelHeight)

	case containsKey(msg, Hotkeys.ChangePanelMode):
		m.changeFilePanelMode()

	case containsKey(msg, Hotkeys.NextFilePanel):
		m.nextFilePanel()

	case containsKey(msg, Hotkeys.PreviousFilePanel):
		m.previousFilePanel()

	case containsKey(msg, Hotkeys.CloseFilePanel):
		m.closeFilePanel()

	case containsKey(msg, Hotkeys.CreateNewFilePanel):
		m.createNewFilePanel(variable.HomeDir)

	case containsKey(msg, Hotkeys.ToggleFilePreviewPanel):
		m.toggleFilePreviewPanel()

	case containsKey(msg, Hotkeys.FocusOnSidebar):
		m.focusOnSideBar()

	case containsKey(msg, Hotkeys.FocusOnProcessBar):
		m.focusOnProcessBar()

	case containsKey(msg, Hotkeys.FocusOnMetaData):
		m.focusOnMetadata()
		go func() {
			m.returnMetaData()
		}()

	case containsKey(msg, Hotkeys.PasteItems):
		go func() {
			m.pasteItem()
		}()

	case containsKey(msg, Hotkeys.FilePanelItemCreate):
		m.panelCreateNewFile()
	case containsKey(msg, Hotkeys.PinnedDirectory):
		m.pinnedDirectory()

	case containsKey(msg, Hotkeys.ToggleDotFile):
		m.toggleDotFileController()

	case containsKey(msg, Hotkeys.ToggleFooter):
		m.toggleFooterController()

	case containsKey(msg, Hotkeys.ExtractFile):
		go func() {
			m.extractFile()
		}()

	case containsKey(msg, Hotkeys.CompressFile):
		go func() {
			m.compressFile()
		}()

	case containsKey(msg, Hotkeys.OpenSPFPrompt):
		m.promptModal.Open(m, PROMPT_COMMAND_COMMAND)

	case containsKey(msg, Hotkeys.OpenHelpMenu):
		m.openHelpMenu()

	case containsKey(msg, Hotkeys.OpenSortOptionsMenu):
		m.openSortOptionsMenu()

	case containsKey(msg, Hotkeys.ToggleReverseSort):
		m.toggleReverseSort()

	case containsKey(msg, Hotkeys.OpenFileWithEditor):
		cmd = m.openFileWithEditor()

	case containsKey(msg, Hotkeys.OpenCurrentDirectoryWithEditor):
		cmd = m.openDirectoryWithEditor()

	default:
		m.normalAndBrowserModeKey(msg)
	}

	return cmd
}

func (m *model) normalAndBrowserModeKey(msg string) {
	// if not focus on the filepanel return
	if m.fileModel.filePanels[m.filePanelFocusIndex].focusType != focus {
		if m.focusPanel == sidebarFocus && (msg == containsKey(msg, Hotkeys.Confirm)) {
			m.sidebarSelectDirectory()
		}
		if m.focusPanel == sidebarFocus && (msg == containsKey(msg, Hotkeys.FilePanelItemRename)) {
			m.pinnedItemRename()
		}
		if m.focusPanel == sidebarFocus && (msg == containsKey(msg, Hotkeys.SearchBar)) {
			m.sidebarSearchBarFocus()
		}
		return
	}
	// Check if in the select mode and focusOn filepanel
	if m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
		switch msg {
		case containsKey(msg, Hotkeys.Confirm):
			m.fileModel.filePanels[m.filePanelFocusIndex].singleItemSelect()
		case containsKey(msg, Hotkeys.FilePanelSelectModeItemsSelectUp):
			m.fileModel.filePanels[m.filePanelFocusIndex].itemSelectUp(m.mainPanelHeight)
		case containsKey(msg, Hotkeys.FilePanelSelectModeItemsSelectDown):
			m.fileModel.filePanels[m.filePanelFocusIndex].itemSelectDown(m.mainPanelHeight)
		case containsKey(msg, Hotkeys.DeleteItems):
			go func() {
				m.deleteItemWarn()
			}()
		case containsKey(msg, Hotkeys.CopyItems):
			m.copyMultipleItem(false)
		case containsKey(msg, Hotkeys.CutItems):
			m.copyMultipleItem(true)
		case containsKey(msg, Hotkeys.FilePanelSelectAllItem):
			m.selectAllItem()
		}
		return
	}

	switch msg {
	case containsKey(msg, Hotkeys.Confirm):
		m.enterPanel()
	case containsKey(msg, Hotkeys.ParentDirectory):
		m.parentDirectory()
	case containsKey(msg, Hotkeys.DeleteItems):
		go func() {
			m.deleteItemWarn()
		}()
	case containsKey(msg, Hotkeys.CopyItems):
		m.copySingleItem(false)
	case containsKey(msg, Hotkeys.CutItems):
		m.copySingleItem(true)
	case containsKey(msg, Hotkeys.FilePanelItemRename):
		m.panelItemRename()
	case containsKey(msg, Hotkeys.SearchBar):
		m.searchBarFocus()
	case containsKey(msg, Hotkeys.CopyPath):
		m.copyPath()
	case containsKey(msg, Hotkeys.CopyPWD):
		m.copyPWD()
	}
}

// Check the hotkey to cancel operation or create file
func (m *model) typingModalOpenKey(msg string) {
	switch msg {
	case containsKey(msg, Hotkeys.CancelTyping):
		m.cancelTypingModal()
	case containsKey(msg, Hotkeys.ConfirmTyping):
		m.createItem()
	}

}

// Todo : There is a lot of duplication for these models, each one of them has to handle
// ConfirmTyping and CancleTyping in a similar way. There is a scope of some good refactoring here.
func (m *model) promptModalOpenKey(msg string) {

}

func (m *model) warnModalOpenKey(msg string) {
	switch msg {
	case containsKey(msg, Hotkeys.Quit), containsKey(msg, Hotkeys.CancelTyping):
		m.cancelWarnModal()
		if m.warnModal.warnType == confirmRenameItem {
			m.cancelRename()
		}
	case containsKey(msg, Hotkeys.Confirm):
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
	case containsKey(msg, Hotkeys.Quit), containsKey(msg, Hotkeys.CancelTyping):
		m.cancelWarnModal()
		m.confirmToQuit = false
		return false
	case containsKey(msg, Hotkeys.Confirm):
		return true
	default:
		return false
	}
}

// Handles key inputs inside sort options menu
func (m *model) sortOptionsKey(msg string) {
	switch msg {
	case containsKey(msg, Hotkeys.OpenSortOptionsMenu):
		m.cancelSortOptions()
	case containsKey(msg, Hotkeys.Quit):
		m.cancelSortOptions()
	case containsKey(msg, Hotkeys.Confirm):
		m.confirmSortOptions()
	case containsKey(msg, Hotkeys.ListUp):
		m.sortOptionsListUp()
	case containsKey(msg, Hotkeys.ListDown):
		m.sortOptionsListDown()
	}
}

func (m *model) renamingKey(msg string) {
	switch msg {
	case containsKey(msg, Hotkeys.CancelTyping):
		m.cancelRename()
	case containsKey(msg, Hotkeys.ConfirmTyping):
		if m.IsRenamingConflicting() {
			m.warnModalForRenaming()
		} else {
			m.confirmRename()
		}
	}
}

func (m *model) sidebarRenamingKey(msg string) {
	switch msg {
	case containsKey(msg, Hotkeys.CancelTyping):
		m.cancelSidebarRename()
	case containsKey(msg, Hotkeys.ConfirmTyping):
		m.confirmSidebarRename()
	}
}

func (m *model) sidebarSearchBarKey(msg string) {
	switch msg {
	case containsKey(msg, Hotkeys.CancelTyping):
		m.sidebarModel.searchBar.Blur()
		m.sidebarModel.searchBar.SetValue("")
	case containsKey(msg, Hotkeys.ConfirmTyping):
		m.sidebarModel.searchBar.Blur()
		m.sidebarModel.resetCursor()
	}
}

// Check the key input and cancel or confirms the search
func (m *model) focusOnSearchbarKey(msg string) {
	switch msg {
	case containsKey(msg, Hotkeys.CancelTyping):
		m.cancelSearch()
	case containsKey(msg, Hotkeys.ConfirmTyping):
		m.confirmSearch()
	}
}

// Check hotkey input in help menu. Possible actions are moving up, down
// and quiting the menu
func (m *model) helpMenuKey(msg string) {
	switch msg {
	case containsKey(msg, Hotkeys.ListUp):
		m.helpMenuListUp()
	case containsKey(msg, Hotkeys.ListDown):
		m.helpMenuListDown()
	case containsKey(msg, Hotkeys.Quit):
		m.quitHelpMenu()
	}
}
