package internal

import (
	"log/slog"
	"slices"

	"github.com/yorukot/superfile/src/internal/common"

	tea "github.com/charmbracelet/bubbletea"
	variable "github.com/yorukot/superfile/src/config"
)

// mainKey handles most of key commands in the regular state of the application. For
// keys that performs actions in multiple panels, like going up or down,
// check the state of model m and handle properly.
func (m *model) mainKey(msg string) tea.Cmd {
	switch {
	// If move up Key is pressed, check the current state and executes
	case slices.Contains(common.Hotkeys.ListUp, msg):
		switch m.focusPanel {
		case sidebarFocus:
			m.sidebarModel.ListUp(m.mainPanelHeight)
		case processBarFocus:
			m.processBarModel.listUp(m.footerHeight)
		case metadataFocus:
			m.fileMetaData.ListUp()
		case nonePanelFocus:
			m.fileModel.filePanels[m.filePanelFocusIndex].listUp(m.mainPanelHeight)
		}

		// If move down Key is pressed, check the current state and executes
	case slices.Contains(common.Hotkeys.ListDown, msg):
		switch m.focusPanel {
		case sidebarFocus:
			m.sidebarModel.ListDown(m.mainPanelHeight)
		case processBarFocus:
			m.processBarModel.listDown(m.footerHeight)
		case metadataFocus:
			m.fileMetaData.ListDown()
		case nonePanelFocus:
			m.fileModel.filePanels[m.filePanelFocusIndex].listDown(m.mainPanelHeight)
		}

	case slices.Contains(common.Hotkeys.PageUp, msg):
		m.fileModel.filePanels[m.filePanelFocusIndex].pgUp(m.mainPanelHeight)

	case slices.Contains(common.Hotkeys.PageDown, msg):
		m.fileModel.filePanels[m.filePanelFocusIndex].pgDown(m.mainPanelHeight)

	case slices.Contains(common.Hotkeys.ChangePanelMode, msg):
		m.getFocusedFilePanel().changeFilePanelMode()

	case slices.Contains(common.Hotkeys.NextFilePanel, msg):
		m.nextFilePanel()

	case slices.Contains(common.Hotkeys.PreviousFilePanel, msg):
		m.previousFilePanel()

	case slices.Contains(common.Hotkeys.CloseFilePanel, msg):
		m.closeFilePanel()

	case slices.Contains(common.Hotkeys.CreateNewFilePanel, msg):
		err := m.createNewFilePanel(variable.HomeDir)
		if err != nil {
			slog.Error("error while creating new panel", "error", err)
		}
	case slices.Contains(common.Hotkeys.ToggleFilePreviewPanel, msg):
		m.toggleFilePreviewPanel()

	case slices.Contains(common.Hotkeys.FocusOnSidebar, msg):
		m.focusOnSideBar()

	case slices.Contains(common.Hotkeys.FocusOnProcessBar, msg):
		m.focusOnProcessBar()

	case slices.Contains(common.Hotkeys.FocusOnMetaData, msg):
		m.focusOnMetadata()

	case slices.Contains(common.Hotkeys.PasteItems, msg):
		return m.getPasteItemCmd()

	case slices.Contains(common.Hotkeys.FilePanelItemCreate, msg):
		m.panelCreateNewFile()
	case slices.Contains(common.Hotkeys.PinnedDirectory, msg):
		m.pinnedDirectory()

	case slices.Contains(common.Hotkeys.ToggleDotFile, msg):
		m.toggleDotFileController()

	case slices.Contains(common.Hotkeys.ToggleFooter, msg):
		m.toggleFooterController()

	case slices.Contains(common.Hotkeys.ExtractFile, msg):
		go func() {
			m.extractFile()
		}()

	case slices.Contains(common.Hotkeys.CompressFile, msg):
		go func() {
			m.compressSelectedFiles()
		}()

	case slices.Contains(common.Hotkeys.OpenCommandLine, msg):
		m.promptModal.Open(true)
	case slices.Contains(common.Hotkeys.OpenSPFPrompt, msg):
		m.promptModal.Open(false)

	case slices.Contains(common.Hotkeys.OpenHelpMenu, msg):
		m.openHelpMenu()

	case slices.Contains(common.Hotkeys.OpenSortOptionsMenu, msg):
		m.openSortOptionsMenu()

	case slices.Contains(common.Hotkeys.ToggleReverseSort, msg):
		m.toggleReverseSort()

	case slices.Contains(common.Hotkeys.OpenFileWithEditor, msg):
		return m.openFileWithEditor()

	case slices.Contains(common.Hotkeys.OpenCurrentDirectoryWithEditor, msg):
		return m.openDirectoryWithEditor()

	default:
		m.normalAndBrowserModeKey(msg)
	}

	return nil
}

func (m *model) normalAndBrowserModeKey(msg string) {
	// if not focus on the filepanel return
	if m.fileModel.filePanels[m.filePanelFocusIndex].focusType != focus {
		if m.focusPanel == sidebarFocus && slices.Contains(common.Hotkeys.Confirm, msg) {
			m.sidebarSelectDirectory()
		}
		if m.focusPanel == sidebarFocus && slices.Contains(common.Hotkeys.FilePanelItemRename, msg) {
			m.sidebarModel.PinnedItemRename()
		}
		if m.focusPanel == sidebarFocus && slices.Contains(common.Hotkeys.SearchBar, msg) {
			m.sidebarSearchBarFocus()
		}
		return
	}
	// Check if in the select mode and focusOn filepanel
	if m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
		switch {
		case slices.Contains(common.Hotkeys.Confirm, msg):
			m.fileModel.filePanels[m.filePanelFocusIndex].singleItemSelect()
		case slices.Contains(common.Hotkeys.FilePanelSelectModeItemsSelectUp, msg):
			m.fileModel.filePanels[m.filePanelFocusIndex].itemSelectUp(m.mainPanelHeight)
		case slices.Contains(common.Hotkeys.FilePanelSelectModeItemsSelectDown, msg):
			m.fileModel.filePanels[m.filePanelFocusIndex].itemSelectDown(m.mainPanelHeight)
		case slices.Contains(common.Hotkeys.DeleteItems, msg):
			go func() {
				m.deleteItemWarn()
			}()
		case slices.Contains(common.Hotkeys.CopyItems, msg):
			m.copyMultipleItem(false)
		case slices.Contains(common.Hotkeys.CutItems, msg):
			m.copyMultipleItem(true)
		case slices.Contains(common.Hotkeys.FilePanelSelectAllItem, msg):
			m.selectAllItem()
		}
		return
	}

	switch {
	case slices.Contains(common.Hotkeys.Confirm, msg):
		m.enterPanel()
	case slices.Contains(common.Hotkeys.ParentDirectory, msg):
		m.parentDirectory()
	case slices.Contains(common.Hotkeys.DeleteItems, msg):
		go func() {
			m.deleteItemWarn()
		}()
	case slices.Contains(common.Hotkeys.CopyItems, msg):
		m.copySingleItem(false)
	case slices.Contains(common.Hotkeys.CutItems, msg):
		m.copySingleItem(true)
	case slices.Contains(common.Hotkeys.FilePanelItemRename, msg):
		m.panelItemRename()
	case slices.Contains(common.Hotkeys.SearchBar, msg):
		m.searchBarFocus()
	case slices.Contains(common.Hotkeys.CopyPath, msg):
		m.copyPath()
	case slices.Contains(common.Hotkeys.CopyPWD, msg):
		m.copyPWD()
	}
}

// Check the hotkey to cancel operation or create file
func (m *model) typingModalOpenKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		m.typingModal.errorMesssage = ""
		m.cancelTypingModal()
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		m.createItem()
	}
}

// TODO : There is a lot of duplication for these models, each one of them has to handle
// ConfirmTyping and CancelTyping in a similar way. There is a scope of some good refactoring here.
func (m *model) warnModalOpenKey(msg string) tea.Cmd {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg) || slices.Contains(common.Hotkeys.Quit, msg):
		m.cancelWarnModal()
		if m.warnModal.warnType == confirmRenameItem {
			m.cancelRename()
		}
	case slices.Contains(common.Hotkeys.Confirm, msg):
		m.warnModal.open = false
		switch m.warnModal.warnType {
		case confirmDeleteItem:
			return m.getDeleteCmd()
		case confirmRenameItem:
			m.confirmRename()
		}
	}
	return nil
}

func (m *model) notifyModalOpenKey(msg string) {
	//nolint:gocritic // We use switch here because other key logic is also using switch, so it's more consistent.
	switch {
	case slices.Contains(common.Hotkeys.Confirm, msg):
		m.notifyModal.open = false
	}
}

// Handle key input to confirm or cancel and close quiting warn in SPF
func (m *model) confirmToQuitSuperfile(msg string) bool {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg) || slices.Contains(common.Hotkeys.Quit, msg):
		m.modelQuitState = notQuitting
		return false
	case slices.Contains(common.Hotkeys.Confirm, msg):
		return true
	default:
		return false
	}
}

// Handles key inputs inside sort options menu
func (m *model) sortOptionsKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.OpenSortOptionsMenu, msg):
		m.cancelSortOptions()
	case slices.Contains(common.Hotkeys.Quit, msg):
		m.cancelSortOptions()
	case slices.Contains(common.Hotkeys.Confirm, msg):
		m.confirmSortOptions()
	case slices.Contains(common.Hotkeys.ListUp, msg):
		m.sortOptionsListUp()
	case slices.Contains(common.Hotkeys.ListDown, msg):
		m.sortOptionsListDown()
	}
}

func (m *model) renamingKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		m.cancelRename()
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		if m.IsRenamingConflicting() {
			m.warnModalForRenaming()
		} else {
			m.confirmRename()
		}
	}
}

func (m *model) sidebarRenamingKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		m.sidebarModel.CancelSidebarRename()
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		m.sidebarModel.ConfirmSidebarRename()
	}
}

// Check the key input and cancel or confirms the search
func (m *model) focusOnSearchbarKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		m.cancelSearch()
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		m.confirmSearch()
	}
}

// Check hotkey input in help menu. Possible actions are moving up, down
// and quiting the menu
func (m *model) helpMenuKey(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.ListUp, msg):
		m.helpMenuListUp()
	case slices.Contains(common.Hotkeys.ListDown, msg):
		m.helpMenuListDown()
	case slices.Contains(common.Hotkeys.Quit, msg):
		m.quitHelpMenu()
	}
}
