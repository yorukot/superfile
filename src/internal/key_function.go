package internal

import (
	"log/slog"
	"slices"

	"github.com/yorukot/superfile/src/internal/common"

	tea "github.com/charmbracelet/bubbletea"
	variable "github.com/yorukot/superfile/src/config"
)

// TODO : Remove duplication here.
// See Issue #968

// NavigationType defines the type of navigation action that can be performed across panels.
// This enum ensures type safety and makes the navigation system extensible.
type NavigationType int

const (
	NavigateUp NavigationType = iota
	NavigateDown
	NavigatePageUp
	NavigatePageDown
)

// Navigation Strategy Documentation:
//
// The navigation system uses a unified approach where:
// 1. File panels use dedicated navigation methods (listUp/Down, pgUp/Down) for precise control
// 2. Non-file panels (sidebar, metadata, processbar) use ListUp/Down methods as fallbacks
// 3. Page navigation maps to single-step navigation for panels without native page support
// 4. This provides consistent UX while respecting each panel's capabilities

// handleSidebarNavigation handles navigation for the sidebar panel.
// Uses ListUp/ListDown for both regular and page navigation as the sidebar
// doesn't distinguish between single-step and page-based navigation.
func (m *model) handleSidebarNavigation(navType NavigationType) {
	switch navType {
	case NavigateUp, NavigatePageUp:
		m.sidebarModel.ListUp(m.mainPanelHeight)
	case NavigateDown, NavigatePageDown:
		m.sidebarModel.ListDown(m.mainPanelHeight)
	}
}

// handleProcessBarNavigation handles navigation for the process bar panel.
// Similar to sidebar, uses single-step navigation for all navigation types
// since the process bar doesn't have dedicated page navigation methods.
func (m *model) handleProcessBarNavigation(navType NavigationType) {
	switch navType {
	case NavigateUp, NavigatePageUp:
		m.processBarModel.listUp(m.footerHeight)
	case NavigateDown, NavigatePageDown:
		m.processBarModel.listDown(m.footerHeight)
	}
}

// handleMetadataNavigation handles navigation for the metadata panel.
// Uses ListUp/ListDown methods for all navigation types as metadata
// doesn't support page-based navigation.
func (m *model) handleMetadataNavigation(navType NavigationType) {
	switch navType {
	case NavigateUp, NavigatePageUp:
		m.fileMetaData.ListUp()
	case NavigateDown, NavigatePageDown:
		m.fileMetaData.ListDown()
	}
}

// handleFilePanelNavigation handles navigation for file panels.
// This is the only panel type that distinguishes between single-step
// navigation (listUp/Down) and page navigation (pgUp/Down).
// This provides the most precise navigation control for file browsing.
func (m *model) handleFilePanelNavigation(navType NavigationType) {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	switch navType {
	case NavigateUp:
		panel.listUp(m.mainPanelHeight)
	case NavigateDown:
		panel.listDown(m.mainPanelHeight)
	case NavigatePageUp:
		panel.pgUp(m.mainPanelHeight)
	case NavigatePageDown:
		panel.pgDown(m.mainPanelHeight)
	}
}

// handlePanelNavigation dispatches navigation commands to the appropriate panel handler
// based on the currently focused panel. This provides a unified navigation interface
// while respecting each panel's specific navigation capabilities.
//
// Panel Focus Mapping:
// - sidebarFocus -> handleSidebarNavigation
// - processBarFocus -> handleProcessBarNavigation
// - metadataFocus -> handleMetadataNavigation
// - nonePanelFocus -> handleFilePanelNavigation
func (m *model) handlePanelNavigation(navType NavigationType) {
	switch m.focusPanel {
	case sidebarFocus:
		m.handleSidebarNavigation(navType)
	case processBarFocus:
		m.handleProcessBarNavigation(navType)
	case metadataFocus:
		m.handleMetadataNavigation(navType)
	case nonePanelFocus:
		m.handleFilePanelNavigation(navType)
	}
}

// mainKey handles most of key commands in the regular state of the application. For
// keys that performs actions in multiple panels, like going up or down,
// check the state of model m and handle properly.
//
// Navigation Keys:
// - ListUp/ListDown: Single-step navigation (arrow keys)
// - PageUp/PageDown: Page-based navigation (respects focused panel)
func (m *model) mainKey(msg string) tea.Cmd {
	switch {
	case slices.Contains(common.Hotkeys.ListUp, msg):
		m.handlePanelNavigation(NavigateUp)

	case slices.Contains(common.Hotkeys.ListDown, msg):
		m.handlePanelNavigation(NavigateDown)

	case slices.Contains(common.Hotkeys.PageUp, msg):
		m.handlePanelNavigation(NavigatePageUp)

	case slices.Contains(common.Hotkeys.PageDown, msg):
		m.handlePanelNavigation(NavigatePageDown)

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
		return m.getExtractFileCmd()

	case slices.Contains(common.Hotkeys.CompressFile, msg):
		return m.getCompressSelectedFilesCmd()

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

