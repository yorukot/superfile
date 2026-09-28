package internal

import (
	"os"
	"slices"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
)

func (m *model) handleFolderShortcut(key string) bool {
	for _, shortcut := range common.Hotkeys.FolderShortcuts() {
		if !slices.Contains(shortcut.Keys, key) {
			continue
		}
		path := m.sidebarModel.FolderShortcutLocation(shortcut)
		if path == "" {
			return true
		}
		// Check before invoking navigation so missing pins leave history and focus untouched.
		if info, err := os.Stat(path); err != nil || !info.IsDir() {
			return true
		}
		if err := m.updateCurrentFilePanelDir(path); err == nil {
			m.focusPanel = nonePanelFocus
			m.getFocusedFilePanel().IsFocused = true
		}
		return true
	}
	return false
}

func (m *model) updateFolderHintKey(key tea.Key) {
	// Enhanced keyboard events report the modifier state after each event,
	// including releases and when the other Alt key remains pressed.
	m.folderHintAltHeld = key.Mod.Contains(tea.ModAlt)
	m.folderHintBlurred = false
}

func (m *model) showFolderHotkeyHints() bool {
	if !common.Config.ShowFolderHotkeyHints || m.folderHintBlurred || m.folderShortcutsBlocked() {
		return false
	}
	return !m.folderHintKeyboard || m.folderHintAltHeld
}

func (m *model) folderShortcutsBlocked() bool {
	return m.firstUse || m.spfError.IsOpen() || m.typingModal.open || m.promptModal.IsOpen() ||
		m.zoxideModal.IsOpen() || m.notifyModel.IsOpen() || m.fileModel.Renaming ||
		m.sidebarModel.IsRenaming() || m.getFocusedFilePanel().SearchBar.Focused() ||
		m.sidebarModel.SearchBarFocused() || m.sortModal.IsOpen() || m.helpMenu.IsOpen()
}
