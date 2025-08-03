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

// KeyMessage represents a keyboard input message for type-safe key handling.
// This reduces the percentage of string-heavy function arguments by providing
// a structured approach to keyboard input processing.
type KeyMessage string

// NewKeyMessage creates a new KeyMessage from a string input.
func NewKeyMessage(msg string) KeyMessage {
	return KeyMessage(msg)
}

// String returns the string representation of the key message.
func (k KeyMessage) String() string {
	return string(k)
}

// CancelConfirmHandler defines the interface for handling cancel/confirm operations.
// This eliminates code duplication across modal and typing handlers.
type CancelConfirmHandler interface {
	HandleCancel()
	HandleConfirm()
}

// handleCancelConfirmKey provides a unified way to handle cancel/confirm key patterns.
// This reduces code duplication across multiple similar functions.
func (m *model) handleCancelConfirmKey(msg KeyMessage, handler CancelConfirmHandler) {
	msgStr := msg.String()
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msgStr):
		handler.HandleCancel()
	case slices.Contains(common.Hotkeys.ConfirmTyping, msgStr):
		handler.HandleConfirm()
	}
}

// typingModalHandler implements CancelConfirmHandler for typing modal operations.
type typingModalHandler struct {
	model *model
}

func (h *typingModalHandler) HandleCancel() {
	h.model.typingModal.errorMesssage = ""
	h.model.cancelTypingModal()
}

func (h *typingModalHandler) HandleConfirm() {
	h.model.createItem()
}

// searchbarHandler implements CancelConfirmHandler for searchbar operations.
type searchbarHandler struct {
	model *model
}

func (h *searchbarHandler) HandleCancel() {
	h.model.cancelSearch()
}

func (h *searchbarHandler) HandleConfirm() {
	h.model.confirmSearch()
}

// sidebarRenamingHandler implements CancelConfirmHandler for sidebar renaming operations.
type sidebarRenamingHandler struct {
	model *model
}

func (h *sidebarRenamingHandler) HandleCancel() {
	h.model.sidebarModel.CancelSidebarRename()
}

func (h *sidebarRenamingHandler) HandleConfirm() {
	h.model.sidebarModel.ConfirmSidebarRename()
}

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
func (m *model) mainKey(msg KeyMessage) tea.Cmd {
	msgStr := msg.String()
	switch {
	case slices.Contains(common.Hotkeys.ListUp, msgStr):
		m.handlePanelNavigation(NavigateUp)

	case slices.Contains(common.Hotkeys.ListDown, msgStr):
		m.handlePanelNavigation(NavigateDown)

	case slices.Contains(common.Hotkeys.PageUp, msgStr):
		m.handlePanelNavigation(NavigatePageUp)

	case slices.Contains(common.Hotkeys.PageDown, msgStr):
		m.handlePanelNavigation(NavigatePageDown)

	case slices.Contains(common.Hotkeys.ChangePanelMode, msgStr):
		m.getFocusedFilePanel().changeFilePanelMode()

	case slices.Contains(common.Hotkeys.NextFilePanel, msgStr):
		m.nextFilePanel()

	case slices.Contains(common.Hotkeys.PreviousFilePanel, msgStr):
		m.previousFilePanel()

	case slices.Contains(common.Hotkeys.CloseFilePanel, msgStr):
		m.closeFilePanel()

	case slices.Contains(common.Hotkeys.CreateNewFilePanel, msgStr):
		err := m.createNewFilePanel(variable.HomeDir)
		if err != nil {
			slog.Error("error while creating new panel", "error", err)
		}
	case slices.Contains(common.Hotkeys.ToggleFilePreviewPanel, msgStr):
		m.toggleFilePreviewPanel()

	case slices.Contains(common.Hotkeys.FocusOnSidebar, msgStr):
		m.focusOnSideBar()

	case slices.Contains(common.Hotkeys.FocusOnProcessBar, msgStr):
		m.focusOnProcessBar()

	case slices.Contains(common.Hotkeys.FocusOnMetaData, msgStr):
		m.focusOnMetadata()

	case slices.Contains(common.Hotkeys.PasteItems, msgStr):
		return m.getPasteItemCmd()

	case slices.Contains(common.Hotkeys.FilePanelItemCreate, msgStr):
		m.panelCreateNewFile()
	case slices.Contains(common.Hotkeys.PinnedDirectory, msgStr):
		m.pinnedDirectory()

	case slices.Contains(common.Hotkeys.ToggleDotFile, msgStr):
		m.toggleDotFileController()

	case slices.Contains(common.Hotkeys.ToggleFooter, msgStr):
		m.toggleFooterController()

	case slices.Contains(common.Hotkeys.ExtractFile, msgStr):
		return m.getExtractFileCmd()

	case slices.Contains(common.Hotkeys.CompressFile, msgStr):
		return m.getCompressSelectedFilesCmd()

	case slices.Contains(common.Hotkeys.OpenCommandLine, msgStr):
		m.promptModal.Open(true)
	case slices.Contains(common.Hotkeys.OpenSPFPrompt, msgStr):
		m.promptModal.Open(false)

	case slices.Contains(common.Hotkeys.OpenHelpMenu, msgStr):
		m.openHelpMenu()

	case slices.Contains(common.Hotkeys.OpenSortOptionsMenu, msgStr):
		m.openSortOptionsMenu()

	case slices.Contains(common.Hotkeys.ToggleReverseSort, msgStr):
		m.toggleReverseSort()

	case slices.Contains(common.Hotkeys.OpenFileWithEditor, msgStr):
		return m.openFileWithEditor()

	case slices.Contains(common.Hotkeys.OpenCurrentDirectoryWithEditor, msgStr):
		return m.openDirectoryWithEditor()

	default:
		m.normalAndBrowserModeKey(msg)
	}

	return nil
}

func (m *model) normalAndBrowserModeKey(msg KeyMessage) {
	msgStr := msg.String()
	// if not focus on the filepanel return
	if m.fileModel.filePanels[m.filePanelFocusIndex].focusType != focus {
		if m.focusPanel == sidebarFocus && slices.Contains(common.Hotkeys.Confirm, msgStr) {
			m.sidebarSelectDirectory()
		}
		if m.focusPanel == sidebarFocus && slices.Contains(common.Hotkeys.FilePanelItemRename, msgStr) {
			m.sidebarModel.PinnedItemRename()
		}
		if m.focusPanel == sidebarFocus && slices.Contains(common.Hotkeys.SearchBar, msgStr) {
			m.sidebarSearchBarFocus()
		}
		return
	}
	// Check if in the select mode and focusOn filepanel
	if m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
		switch {
		case slices.Contains(common.Hotkeys.Confirm, msgStr):
			m.fileModel.filePanels[m.filePanelFocusIndex].singleItemSelect()
		case slices.Contains(common.Hotkeys.FilePanelSelectModeItemsSelectUp, msgStr):
			m.fileModel.filePanels[m.filePanelFocusIndex].itemSelectUp(m.mainPanelHeight)
		case slices.Contains(common.Hotkeys.FilePanelSelectModeItemsSelectDown, msgStr):
			m.fileModel.filePanels[m.filePanelFocusIndex].itemSelectDown(m.mainPanelHeight)
		case slices.Contains(common.Hotkeys.DeleteItems, msgStr):
			go func() {
				m.deleteItemWarn()
			}()
		case slices.Contains(common.Hotkeys.CopyItems, msgStr):
			m.copyMultipleItem(false)
		case slices.Contains(common.Hotkeys.CutItems, msgStr):
			m.copyMultipleItem(true)
		case slices.Contains(common.Hotkeys.FilePanelSelectAllItem, msgStr):
			m.selectAllItem()
		}
		return
	}

	switch {
	case slices.Contains(common.Hotkeys.Confirm, msgStr):
		m.enterPanel()
	case slices.Contains(common.Hotkeys.ParentDirectory, msgStr):
		m.parentDirectory()
	case slices.Contains(common.Hotkeys.DeleteItems, msgStr):
		go func() {
			m.deleteItemWarn()
		}()
	case slices.Contains(common.Hotkeys.CopyItems, msgStr):
		m.copySingleItem(false)
	case slices.Contains(common.Hotkeys.CutItems, msgStr):
		m.copySingleItem(true)
	case slices.Contains(common.Hotkeys.FilePanelItemRename, msgStr):
		m.panelItemRename()
	case slices.Contains(common.Hotkeys.SearchBar, msgStr):
		m.searchBarFocus()
	case slices.Contains(common.Hotkeys.CopyPath, msgStr):
		m.copyPath()
	case slices.Contains(common.Hotkeys.CopyPWD, msgStr):
		m.copyPWD()
	}
}

// Check the hotkey to cancel operation or create file
func (m *model) typingModalOpenKey(msg KeyMessage) {
	handler := &typingModalHandler{model: m}
	m.handleCancelConfirmKey(msg, handler)
}

// TODO : There is a lot of duplication for these models, each one of them has to handle
// ConfirmTyping and CancelTyping in a similar way. There is a scope of some good refactoring here.
func (m *model) warnModalOpenKey(msg KeyMessage) tea.Cmd {
	msgStr := msg.String()
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msgStr) || slices.Contains(common.Hotkeys.Quit, msgStr):
		m.cancelWarnModal()
		if m.warnModal.warnType == confirmRenameItem {
			m.cancelRename()
		}
	case slices.Contains(common.Hotkeys.Confirm, msgStr):
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

func (m *model) notifyModalOpenKey(msg KeyMessage) {
	msgStr := msg.String()
	//nolint:gocritic // We use switch here because other key logic is also using switch, so it's more consistent.
	switch {
	case slices.Contains(common.Hotkeys.Confirm, msgStr):
		m.notifyModal.open = false
	}
}

// Handle key input to confirm or cancel and close quiting warn in SPF
func (m *model) confirmToQuitSuperfile(msg KeyMessage) bool {
	msgStr := msg.String()
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msgStr) || slices.Contains(common.Hotkeys.Quit, msgStr):
		m.modelQuitState = notQuitting
		return false
	case slices.Contains(common.Hotkeys.Confirm, msgStr):
		return true
	default:
		return false
	}
}

// Handles key inputs inside sort options menu
func (m *model) sortOptionsKey(msg KeyMessage) {
	msgStr := msg.String()
	switch {
	case slices.Contains(common.Hotkeys.OpenSortOptionsMenu, msgStr):
		m.cancelSortOptions()
	case slices.Contains(common.Hotkeys.Quit, msgStr):
		m.cancelSortOptions()
	case slices.Contains(common.Hotkeys.Confirm, msgStr):
		m.confirmSortOptions()
	case slices.Contains(common.Hotkeys.ListUp, msgStr):
		m.sortOptionsListUp()
	case slices.Contains(common.Hotkeys.ListDown, msgStr):
		m.sortOptionsListDown()
	}
}

func (m *model) renamingKey(msg KeyMessage) {
	msgStr := msg.String()
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msgStr):
		m.cancelRename()
	case slices.Contains(common.Hotkeys.ConfirmTyping, msgStr):
		if m.IsRenamingConflicting() {
			m.warnModalForRenaming()
		} else {
			m.confirmRename()
		}
	}
}

func (m *model) sidebarRenamingKey(msg KeyMessage) {
	handler := &sidebarRenamingHandler{model: m}
	m.handleCancelConfirmKey(msg, handler)
}

// Check the key input and cancel or confirms the search
func (m *model) focusOnSearchbarKey(msg KeyMessage) {
	handler := &searchbarHandler{model: m}
	m.handleCancelConfirmKey(msg, handler)
}

// Check hotkey input in help menu. Possible actions are moving up, down
// and quiting the menu
func (m *model) helpMenuKey(msg KeyMessage) {
	msgStr := msg.String()
	switch {
	case slices.Contains(common.Hotkeys.ListUp, msgStr):
		m.helpMenuListUp()
	case slices.Contains(common.Hotkeys.ListDown, msgStr):
		m.helpMenuListDown()
	case slices.Contains(common.Hotkeys.Quit, msgStr):
		m.quitHelpMenu()
	}
}
