package internal

import (
	"log/slog"
	"slices"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/notify"

	tea "github.com/charmbracelet/bubbletea"

	variable "github.com/yorukot/superfile/src/config"
)

// mainKey handles most of key commands in the regular state of the application. For
// keys that performs actions in multiple panels, like going up or down,
// check the state of model m and handle properly.
// TODO: This function has grown too big. It needs to be fixed, via major
// updates and fixes in key handling code
func (m *model) mainKey(msg string) tea.Cmd { //nolint: gocyclo,cyclop,funlen // See above
	switch {
	// If move up Key is pressed, check the current state and executes
	case slices.Contains(common.Hotkeys.ListUp, msg):
		switch m.focusPanel {
		case sidebarFocus:
			m.sidebarModel.ListUp(m.mainPanelHeight)
		case processBarFocus:
			m.processBarModel.ListUp(m.footerHeight)
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
			m.processBarModel.ListDown(m.footerHeight)
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
		return m.toggleFooterController()

	case slices.Contains(common.Hotkeys.ExtractFile, msg):
		return m.getExtractFileCmd()

	case slices.Contains(common.Hotkeys.CompressFile, msg):
		return m.getCompressSelectedFilesCmd()

	case slices.Contains(common.Hotkeys.OpenCommandLine, msg):
		m.promptModal.Open(true)
	case slices.Contains(common.Hotkeys.OpenSPFPrompt, msg):
		m.promptModal.Open(false)
	case slices.Contains(common.Hotkeys.OpenZoxide, msg):
		return m.zoxideModal.Open()

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
		return m.normalAndBrowserModeKey(msg)
	}

	return nil
}

func (m *model) normalAndBrowserModeKey(msg string) tea.Cmd {
	if !m.getFocusedFilePanel().isFocused {
		return m.handleSidebarFocusKeys(msg)
	}

	if m.getFocusedFilePanel().panelMode == selectMode {
		return m.handleSelectModeKeys(msg)
	}

	return m.handleBrowserModeKeys(msg)
}

func (m *model) handleSidebarFocusKeys(msg string) tea.Cmd {
	if m.focusPanel != sidebarFocus {
		return nil
	}

	switch {
	case slices.Contains(common.Hotkeys.Confirm, msg):
		m.sidebarSelectDirectory()
	case slices.Contains(common.Hotkeys.FilePanelItemRename, msg):
		m.sidebarModel.PinnedItemRename()
	case slices.Contains(common.Hotkeys.SearchBar, msg):
		m.sidebarSearchBarFocus()
	}
	return nil
}

func (m *model) handleSelectModeKeys(msg string) tea.Cmd {
	if slices.Contains(common.Hotkeys.Confirm, msg) {
		m.fileModel.filePanels[m.filePanelFocusIndex].singleItemSelect()
		return nil
	}
	if slices.Contains(common.Hotkeys.FilePanelSelectModeItemsSelectUp, msg) {
		m.fileModel.filePanels[m.filePanelFocusIndex].itemSelectUp(m.mainPanelHeight)
		return nil
	}
	if slices.Contains(common.Hotkeys.FilePanelSelectModeItemsSelectDown, msg) {
		m.fileModel.filePanels[m.filePanelFocusIndex].itemSelectDown(m.mainPanelHeight)
		return nil
	}
	if slices.Contains(common.Hotkeys.DeleteItems, msg) {
		return m.getDeleteTriggerCmd(false)
	}
	if slices.Contains(common.Hotkeys.PermanentlyDeleteItems, msg) {
		return m.getDeleteTriggerCmd(true)
	}
	
	return m.handleSelectModeFileOperations(msg)
}

func (m *model) handleSelectModeFileOperations(msg string) tea.Cmd {
	if slices.Contains(common.Hotkeys.CopyItems, msg) {
		m.copyMultipleItem(false)
	} else if slices.Contains(common.Hotkeys.CutItems, msg) {
		m.copyMultipleItem(true)
	} else if slices.Contains(common.Hotkeys.BulkRename, msg) {
		m.panelBulkRename()
	} else if slices.Contains(common.Hotkeys.FilePanelSelectAllItem, msg) {
		m.selectAllItem()
	}
	return nil
}

func (m *model) handleBrowserModeKeys(msg string) tea.Cmd {
	switch {
	case slices.Contains(common.Hotkeys.Confirm, msg):
		m.enterPanel()
	case slices.Contains(common.Hotkeys.ParentDirectory, msg):
		m.parentDirectory()
	case slices.Contains(common.Hotkeys.DeleteItems, msg):
		return m.getDeleteTriggerCmd(false)
	case slices.Contains(common.Hotkeys.PermanentlyDeleteItems, msg):
		return m.getDeleteTriggerCmd(true)
	default:
		return m.handleBrowserModeFileOperations(msg)
	}
	return nil
}

func (m *model) handleBrowserModeFileOperations(msg string) tea.Cmd {
	switch {
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
	return nil
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

func (m *model) notifyModelOpenKey(msg string) tea.Cmd {
	isCancel := slices.Contains(common.Hotkeys.CancelTyping, msg) || slices.Contains(common.Hotkeys.Quit, msg)
	isConfirm := slices.Contains(common.Hotkeys.Confirm, msg)

	if !isCancel && !isConfirm {
		slog.Warn("Invalid keypress in notifyModel", "msg", msg)
		return nil
	}
	m.notifyModel.Close()
	action := m.notifyModel.GetConfirmAction()
	if isCancel {
		return m.handleNotifyModelCancel(action)
	}
	return m.handleNotifyModelConfirm(action)
}

func (m *model) handleNotifyModelCancel(action notify.ConfirmActionType) tea.Cmd {
	switch action {
	case notify.RenameAction:
		m.cancelRename()
	case notify.QuitAction:
		m.modelQuitState = notQuitting
	case notify.DeleteAction, notify.NoAction, notify.PermanentDeleteAction:
		// Do nothing
	default:
		slog.Error("Unknown type of action", "action", action)
	}
	return nil
}

func (m *model) handleNotifyModelConfirm(action notify.ConfirmActionType) tea.Cmd {
	switch action {
	case notify.DeleteAction:
		return m.getDeleteCmd(false)
	case notify.PermanentDeleteAction:
		return m.getDeleteCmd(true)
	case notify.RenameAction:
		m.confirmRename()
	case notify.QuitAction:
		m.modelQuitState = quitConfirmationReceived
	case notify.NoAction:
		// Ignore
	default:
		slog.Error("Unknown type of action", "action", action)
	}
	return nil
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

func (m *model) renamingKey(msg string) tea.Cmd {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		m.cancelRename()
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		if m.IsRenamingConflicting() {
			return m.warnModalForRenaming()
		}
		m.confirmRename()
	}

	return nil
}

// Handle bulk rename modal keys
func (m *model) bulkRenameKey(msg string) tea.Cmd {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		m.cancelBulkRename()
	case slices.Contains(common.Hotkeys.ConfirmTyping, msg):
		m.confirmBulkRename()
	case slices.Contains(common.Hotkeys.ListUp, msg):
		m.adjustBulkRenameValue(-1)
	case slices.Contains(common.Hotkeys.ListDown, msg):
		m.adjustBulkRenameValue(1)
	case slices.Contains(common.Hotkeys.NavBulkRename, msg):
		m.bulkRenameNextType()
	case slices.Contains(common.Hotkeys.RevNavBulkRename, msg):
		m.bulkRenamePrevType()
	}
	return nil
}

func (m *model) adjustBulkRenameValue(delta int) {
	switch m.bulkRenameModal.renameType {
	case 3:
		m.adjustNumberingValue(delta)
	case 4:
		m.adjustCaseType(delta)
	default:
		m.navigateBulkRename(delta)
	}
}

func (m *model) adjustNumberingValue(delta int) {
	newValue := m.bulkRenameModal.startNumber + delta
	if newValue >= 0 {
		m.bulkRenameModal.startNumber = newValue
	}
}

func (m *model) adjustCaseType(delta int) {
	newValue := m.bulkRenameModal.caseType + delta
	if newValue >= 0 && newValue <= 2 {
		m.bulkRenameModal.caseType = newValue
	}
}

func (m *model) navigateBulkRename(delta int) {
	if delta < 0 {
		m.bulkRenameNavigateUp()
	} else {
		m.bulkRenameNavigateDown()
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
	if m.helpMenu.searchBar.Focused() {
		switch {
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg), slices.Contains(common.Hotkeys.CancelTyping, msg):
			m.helpMenu.searchBar.Blur()
		default:
			m.filterHelpMenu(m.helpMenu.searchBar.Value())
		}
	} else {
		m.handleHelpMenuNavKeys(msg)
	}
}

func (m *model) handleHelpMenuNavKeys(msg string) {
	switch {
	case slices.Contains(common.Hotkeys.ListUp, msg):
		m.helpMenuListUp()
	case slices.Contains(common.Hotkeys.ListDown, msg):
		m.helpMenuListDown()
	case slices.Contains(common.Hotkeys.Quit, msg):
		m.quitHelpMenu()
	case slices.Contains(common.Hotkeys.SearchBar, msg):
		m.helpMenu.searchBar.Focus()
	}
}
