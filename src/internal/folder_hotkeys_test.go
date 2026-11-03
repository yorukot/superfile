package internal

import (
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/adrg/xdg"
	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/pkg/utils"
)

//nolint:reassign // Isolate global settings and restore them with t.Cleanup.
func TestFolderShortcutsNavigateCurrentPanel(
	t *testing.T,
) {
	originalHome, originalDownload := xdg.Home, xdg.UserDirs.Download
	xdg.Home, xdg.UserDirs.Download = t.TempDir(), t.TempDir()
	t.Cleanup(func() { xdg.Home, xdg.UserDirs.Download = originalHome, originalDownload })
	utils.SetupFiles(t, filepath.Join(xdg.Home, "a"), filepath.Join(xdg.Home, "b"))
	other := t.TempDir()
	m := defaultTestModel(other, xdg.Home)
	m.fileModel.FocusedPanelIndex = 1
	m.getFocusedFilePanel().IsFocused = true
	panel := m.getFocusedFilePanel()
	panel.ListDown()
	panel.SearchBar.SetValue("old query")
	panel.PanelMode = filepanel.SelectMode
	m.focusPanel = sidebarFocus
	panel.IsFocused = false

	TeaUpdate(m, tea.KeyPressMsg{Code: 'd', Mod: tea.ModAlt})
	assert.Equal(t, xdg.UserDirs.Download, panel.Location)
	assert.Equal(t, other, m.fileModel.FilePanels[0].Location)
	assert.Equal(t, nonePanelFocus, m.focusPanel)
	assert.True(t, panel.IsFocused)
	assert.Empty(t, panel.SearchBar.Value())
	assert.Equal(t, filepanel.SelectMode, panel.PanelMode)

	TeaUpdate(m, tea.KeyPressMsg{Code: 'h', Mod: tea.ModAlt})
	assert.Equal(t, xdg.Home, panel.Location)
	assert.Equal(t, 1, panel.GetCursor())
}

func TestFolderShortcutsBlockedByInputs(t *testing.T) {
	setups := map[string]func(*model){
		"create":         func(m *model) { m.typingModal.open = true },
		"rename":         func(m *model) { m.fileModel.Renaming = true },
		"file search":    func(m *model) { m.getFocusedFilePanel().SearchBar.Focus() },
		"sidebar search": func(m *model) { m.sidebarModel.SearchBarFocus() },
		"command prompt": func(m *model) { m.promptModal.Open(true) },
		"SPF prompt":     func(m *model) { m.promptModal.Open(false) },
		"help":           func(m *model) { m.helpMenu.Open() },
		"sort":           func(m *model) { m.sortModal.Open(m.getFocusedFilePanel().SortKind) },
	}
	for name, setup := range setups {
		t.Run(name, func(t *testing.T) {
			start := t.TempDir()
			m := defaultTestModel(start)
			setup(m)
			TeaUpdate(m, tea.KeyPressMsg{Code: 'h', Mod: tea.ModAlt})
			assert.Equal(t, start, m.getFocusedFilePanel().Location)
			assert.False(t, m.showFolderHotkeyHints())
		})
	}
}

//nolint:reassign // Isolate global settings and restore them with t.Cleanup.
func TestFolderShortcutMissingAndHiddenPins(
	t *testing.T,
) {
	originalPinned, originalConfig := variable.PinnedFile, common.Config
	variable.PinnedFile = filepath.Join(t.TempDir(), "pinned.json")
	common.Config.SidebarWidth = 0
	t.Cleanup(func() { variable.PinnedFile, common.Config = originalPinned, originalConfig })
	start, target := t.TempDir(), t.TempDir()
	m := defaultTestModel(start)
	require.NoError(t, m.sidebarModel.TogglePinnedDirectory(target))
	TeaUpdate(m, tea.KeyPressMsg{Code: '1', Mod: tea.ModAlt})
	assert.Equal(t, target, m.getFocusedFilePanel().Location)
	TeaUpdate(m, tea.KeyPressMsg{Code: '9', Mod: tea.ModAlt})
	assert.Equal(t, target, m.getFocusedFilePanel().Location)
	require.NoError(t, os.Remove(target))
	require.NoError(t, m.updateCurrentFilePanelDir(start))
	m.focusPanel = processBarFocus
	m.getFocusedFilePanel().IsFocused = false
	TeaUpdate(m, tea.KeyPressMsg{Code: '1', Mod: tea.ModAlt})
	assert.Equal(t, start, m.getFocusedFilePanel().Location)
	assert.Equal(t, processBarFocus, m.focusPanel)
}

//nolint:reassign // Isolate global settings and restore them with t.Cleanup.
func TestFolderHintKeyboardLifecycle(
	t *testing.T,
) {
	originalConfig := common.Config
	common.Config.ShowFolderHotkeyHints = true
	t.Cleanup(func() { common.Config = originalConfig })
	m := defaultTestModel(t.TempDir())
	assert.True(t, m.showFolderHotkeyHints(), "legacy terminals show hints permanently")
	view := m.View()
	assert.True(t, view.KeyboardEnhancements.ReportEventTypes)
	assert.True(t, view.KeyboardEnhancements.ReportAllKeysAsEscapeCodes)
	assert.True(t, view.KeyboardEnhancements.ReportAssociatedText)
	TeaUpdate(m, tea.KeyboardEnhancementsMsg{Flags: ansi.KittyReportEventTypes})
	assert.True(t, m.showFolderHotkeyHints(), "partial support still uses permanent hints")
	TeaUpdate(m, tea.KeyboardEnhancementsMsg{Flags: ansi.KittyReportEventTypes | ansi.KittyReportAllKeysAsEscapeCodes})
	assert.False(t, m.showFolderHotkeyHints())
	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyLeftAlt, Mod: tea.ModAlt})
	assert.True(t, m.showFolderHotkeyHints())
	TeaUpdate(m, tea.KeyReleaseMsg{Code: tea.KeyLeftAlt, Mod: tea.ModAlt})
	assert.True(t, m.showFolderHotkeyHints(), "other Alt key is still held")
	TeaUpdate(m, tea.KeyReleaseMsg{Code: tea.KeyRightAlt})
	assert.False(t, m.showFolderHotkeyHints())
	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyRightAlt, Mod: tea.ModAlt})
	TeaUpdate(m, tea.BlurMsg{})
	assert.False(t, m.showFolderHotkeyHints())
	TeaUpdate(m, tea.FocusMsg{})
	assert.False(t, m.showFolderHotkeyHints(), "focus must not restore stale Alt state")
	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyLeftAlt, Mod: tea.ModAlt})
	common.Config.ShowFolderHotkeyHints = false
	assert.False(t, m.showFolderHotkeyHints())
	assert.False(t, m.View().KeyboardEnhancements.ReportAllKeysAsEscapeCodes)
}

//nolint:reassign // Isolate global settings and restore them with t.Cleanup.
func TestFolderShortcutsPreserveExistingBindings(
	t *testing.T,
) {
	original := common.Hotkeys
	t.Cleanup(func() { common.Hotkeys = original })
	common.Hotkeys.OpenHelpMenu = []string{"alt+d"}
	m := defaultTestModel(t.TempDir())
	start := m.getFocusedFilePanel().Location
	TeaUpdate(m, tea.KeyPressMsg{Code: 'd', Mod: tea.ModAlt})
	assert.True(t, m.helpMenu.IsOpen())
	assert.Equal(t, start, m.getFocusedFilePanel().Location)
}
