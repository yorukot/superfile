package internal

import (
	"log/slog"
	"slices"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/mouse"
	"github.com/yorukot/superfile/src/internal/ui/contextmenu"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

// openContextMenu shows the context menu for whatever was right clicked at the
// given terminal coordinates.
func (m *model) openContextMenu(target mouse.Target, x, y int) tea.Cmd {
	items := m.contextMenuItems(target, x, y)
	if len(items) == 0 {
		// Nothing to offer there, so a menu must not be left hanging open.
		m.contextMenu.Close()
		return nil
	}

	m.contextMenu.Open(x, y, m.fullWidth, m.fullHeight, items)
	return nil
}

// contextMenuItems focuses the clicked target and returns the entries that apply
// to it. An empty result means there is nothing to offer at that spot.
func (m *model) contextMenuItems(target mouse.Target, x, y int) []contextmenu.Item {
	switch target.Kind {
	case mouse.TargetFilePanelItem:
		return m.fileItemContextItems(target)
	case mouse.TargetSidebarDirectory:
		return m.sidebarDirectoryContextItems(target)
	case mouse.TargetUnknown:
		if !m.inFilePanelArea(x, y) {
			return nil
		}
		// Paste, new file, and new folder all act on the focused panel, so the
		// panel the pointer is over has to be the one focused first. A right
		// click does not go through the click handlers that would do it, and
		// leaving focus where it was would put the new entry in a directory the
		// user never pointed at.
		if panelIndex := m.filePanelIndexAtX(x); panelIndex >= 0 {
			m.focusFilePanelOnMouse(panelIndex)
		}
		return backgroundContextItems()
	case mouse.TargetProcessBarItem, mouse.TargetMetadataItem, mouse.TargetContextMenuItem:
		// The process bar and the metadata panel have no actions to offer yet,
		// and a right click that landed on the open menu itself only moves it.
		return nil
	}
	return nil
}

// inFilePanelArea reports whether a point that belongs to no widget is inside
// the directory area, which is where the file wide actions belong.
//
// Right clicking the footer or the area outside the panels must not offer
// paste and create, since those act on a directory and no directory is in view
// there.
func (m *model) inFilePanelArea(x, y int) bool {
	if y >= m.mainPanelHeight {
		return false
	}
	if x < m.filePanelAreaX() {
		return false
	}
	return x < m.fullWidth
}

// backgroundContextItems returns the entries for the directory area around the
// file panels.
func backgroundContextItems() []contextmenu.Item {
	return []contextmenu.Item{
		{
			Label: "Paste", Shortcut: shortcutOf(common.Hotkeys.PasteItems),
			Action: contextmenu.ActionPaste,
		},
		{
			Label: "New file", Shortcut: shortcutOf(common.Hotkeys.FilePanelItemCreate),
			Action: contextmenu.ActionNewFile,
		},
		{Label: "New folder", Action: contextmenu.ActionNewFolder},
	}
}

// fileItemContextItems returns the entries for a file or directory row, after
// focusing its panel and narrowing the selection to it where appropriate.
func (m *model) fileItemContextItems(target mouse.Target) []contextmenu.Item {
	m.focusFilePanelOnMouse(target.PanelIndex)
	panel := m.getFocusedFilePanel()

	if !panel.SetCursorToIndex(target.ItemIndex) {
		return nil
	}

	location := panel.GetElementAtIdx(target.ItemIndex).Location
	if !panel.CheckSelected(location) {
		// Right clicking outside the current selection narrows it to the clicked
		// entry, so that copy, cut and delete apply to what was clicked. Right
		// clicking inside the selection keeps every earlier entry selected.
		panel.ResetSelected()
		panel.SetSelected(location)
		panel.SetPanelMode(filepanel.SelectMode)
	}

	return []contextmenu.Item{
		{Label: "Open", Shortcut: shortcutOf(common.Hotkeys.Confirm), Action: contextmenu.ActionOpen},
		{Label: "Copy", Shortcut: shortcutOf(common.Hotkeys.CopyItems), Action: contextmenu.ActionCopy},
		{Label: "Cut", Shortcut: shortcutOf(common.Hotkeys.CutItems), Action: contextmenu.ActionCut},
		{
			Label: "Copy path", Shortcut: shortcutOf(common.Hotkeys.CopyPath),
			Action: contextmenu.ActionCopyPath,
		},
		{
			Label: "Rename", Shortcut: shortcutOf(common.Hotkeys.FilePanelItemRename),
			Action: contextmenu.ActionRename,
		},
		{Label: "Delete", Shortcut: shortcutOf(common.Hotkeys.DeleteItems), Action: contextmenu.ActionDelete},
	}
}

// sidebarDirectoryContextItems returns the entries for a sidebar directory.
//
// A right click only moves the sidebar cursor. Navigating is left to the Open
// entry, so that opening a menu does not change directory behind the user.
func (m *model) sidebarDirectoryContextItems(target mouse.Target) []contextmenu.Item {
	m.focusSidebarOnMouse()
	if !m.sidebarModel.SetCursor(target.ItemIndex) {
		return nil
	}

	return []contextmenu.Item{
		{Label: "Open", Shortcut: shortcutOf(common.Hotkeys.Confirm), Action: contextmenu.ActionOpen},
		{
			Label: "Copy path", Shortcut: shortcutOf(common.Hotkeys.CopyPath),
			Action: contextmenu.ActionCopyPath,
		},
	}
}

// runContextMenuAction performs the chosen entry's action.
//
// Focus cannot move while the menu is open, because the menu captures input, so
// the actions can rely on the panel or sidebar that was focused when it opened.
func (m *model) runContextMenuAction(action contextmenu.Action) tea.Cmd {
	switch action {
	case contextmenu.ActionOpen:
		if m.focusPanel == sidebarFocus {
			m.sidebarSelectDirectory()
			return nil
		}
		m.enterPanel()
	case contextmenu.ActionCopy:
		m.copyContextMenuTarget(false)
	case contextmenu.ActionCut:
		m.copyContextMenuTarget(true)
	case contextmenu.ActionPaste:
		return m.getPasteItemCmd()
	case contextmenu.ActionRename:
		m.panelItemRename()
	case contextmenu.ActionDelete:
		return m.getDeleteTriggerCmd(false)
	case contextmenu.ActionCopyPath:
		m.copyContextMenuTargetPath()
	case contextmenu.ActionNewFile:
		m.panelCreateNewFile()
	case contextmenu.ActionNewFolder:
		m.panelCreateNewFolder()
	case contextmenu.ActionNone:
		// An entry without an action does nothing, so a partially built menu can
		// never panic.
	}
	return nil
}

// copyContextMenuTarget copies either the whole selection or the single focused
// entry, matching what the keyboard does in each mode.
func (m *model) copyContextMenuTarget(cut bool) {
	panel := m.getFocusedFilePanel()
	if panel.PanelMode == filepanel.SelectMode && panel.SelectedCount() > 0 {
		m.copyMultipleItem(cut)
		return
	}
	m.copySingleItem(cut)
}

// copyContextMenuTargetPath copies the path of whatever the menu was opened on.
func (m *model) copyContextMenuTargetPath() {
	if m.focusPanel == sidebarFocus {
		// The sidebar was right clicked, so the path of interest is the sidebar
		// directory rather than the focused file panel's entry.
		location := m.sidebarModel.GetCurrentDirectoryLocation()
		if location == "" {
			return
		}
		if err := m.writeClipboard(location); err != nil {
			slog.Error("Error while copying directory path", "error", err)
		}
		return
	}

	m.copyPath()
}

// panelCreateNewFolder starts creating a directory in the focused panel.
//
// It is the same flow as creating a file, with the typing modal told that the
// name it is given belongs to a directory.
func (m *model) panelCreateNewFolder() {
	m.panelCreateNewFile()
	m.typingModal.directory = true
}

// closeContextMenu dismisses the menu, if it is open.
//
// It is nil when there was nothing to dismiss, so a caller can pass its result
// straight back as the command for the message that triggered it.
func (m *model) closeContextMenu() tea.Cmd {
	if !m.contextMenu.IsOpen() {
		return nil
	}
	m.contextMenu.Close()
	return nil
}

// contextMenuKey handles keys while the context menu is open.
func (m *model) contextMenuKey(msg string) tea.Cmd {
	switch {
	case slices.Contains(common.Hotkeys.CancelTyping, msg):
		m.contextMenu.Close()
	case slices.Contains(common.Hotkeys.ListUp, msg):
		m.contextMenu.ListUp()
	case slices.Contains(common.Hotkeys.ListDown, msg):
		m.contextMenu.ListDown()
	case slices.Contains(common.Hotkeys.Confirm, msg):
		action, ok := m.contextMenu.SelectedAction()
		m.contextMenu.Close()
		if !ok {
			return nil
		}
		return m.runContextMenuAction(action)
	}
	return nil
}

// shortcutOf returns the hotkey to show next to a menu entry.
//
// Hotkey lists are padded with empty strings so a second key can be added without
// editing the line, so the first non empty entry is the one to display.
func shortcutOf(keys []string) string {
	for _, key := range keys {
		if key != "" {
			return key
		}
	}
	return ""
}
