package internal

import (
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Pumps progress until done, unwrapping batches.
func drainSearchResults(t *testing.T, m *model, cmd tea.Cmd) {
	t.Helper()
	cmds := []tea.Cmd{cmd}
	for len(cmds) > 0 {
		cur := cmds[0]
		cmds = cmds[1:]
		msg := ExecuteTeaCmdWithTimeout(cur, DefaultTestTimeout)
		require.NotNil(t, msg, "search progress did not arrive in time")
		if batch, ok := msg.(tea.BatchMsg); ok {
			cmds = append(cmds, []tea.Cmd(batch)...)
			continue
		}
		nextCmd := TeaUpdate(m, msg)
		if sp, ok := msg.(SearchProgressMsg); ok {
			if sp.progress.Done {
				return
			}
			cmds = append(cmds, nextCmd)
		}
	}
}

func searchTestDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	sub := filepath.Join(root, "src", "deep")
	utils.SetupDirectories(t, sub)
	utils.SetupFiles(t,
		filepath.Join(root, "main.go"),
		filepath.Join(root, "other.txt"),
		filepath.Join(root, "src", "utils.go"),
		filepath.Join(root, "src", "deep", "main.go"),
	)
	return root
}

func typeSearchQuery(t *testing.T, m *model, query string) tea.Cmd {
	t.Helper()
	var cmd tea.Cmd
	for _, r := range query {
		cmd = TeaUpdate(m, utils.TeaRuneKeyMsg(string(r)))
	}
	require.Equal(t, query, m.getFocusedFilePanel().SearchBar.Value())
	return cmd
}

func TestSearchModeEnterExit(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))

	panel := m.getFocusedFilePanel()
	require.True(t, panel.Search.Active)
	require.Equal(t, root, panel.Search.Root)
	require.True(t, panel.SearchBar.Focused())
	require.Empty(t, panel.SearchBar.Value())

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyEsc})
	require.False(t, panel.Search.Active)
	require.Empty(t, panel.SearchBar.Value())
	require.False(t, panel.SearchBar.Focused())
	require.NotZero(t, panel.ElemCount(), "panel listing should be restored")
}

func TestSearchModeResults(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	cmd := typeSearchQuery(t, m, "main")
	drainSearchResults(t, m, cmd)

	panel := m.getFocusedFilePanel()
	require.True(t, panel.Search.Done)
	assert.Equal(t, int64(2), panel.Search.MatchCount, "main.go files should match 'main'")
	require.Len(t, panel.Search.Results, 2)
	assert.Equal(t, "main.go", panel.Search.Results[0].Path)
	assert.Equal(t, "src/deep/main.go", panel.Search.Results[1].Path)
}

func TestSearchModeNoResults(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	cmd := typeSearchQuery(t, m, "zzznomatch")
	drainSearchResults(t, m, cmd)

	panel := m.getFocusedFilePanel()
	require.True(t, panel.Search.Done)
	assert.Zero(t, panel.Search.MatchCount)
	assert.Empty(t, panel.Search.Results)
}

func TestSearchModeNavigation(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	cmd := typeSearchQuery(t, m, "main")
	drainSearchResults(t, m, cmd)

	panel := m.getFocusedFilePanel()
	require.Len(t, panel.Search.Results, 2)
	assert.Equal(t, 0, panel.Search.Cursor)

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Equal(t, 1, panel.Search.Cursor)

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Equal(t, 0, panel.Search.Cursor)
}

func TestSearchModeConfirmFile(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	cmd := typeSearchQuery(t, m, "main")
	drainSearchResults(t, m, cmd)

	panel := m.getFocusedFilePanel()
	require.Len(t, panel.Search.Results, 2)
	require.False(t, panel.Search.Results[0].Dir)

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyEnter})

	assert.False(t, panel.Search.Active, "search mode should end on confirm")
	assert.Equal(t, root, panel.Location)
	idx := panel.FindElementIndexByName("main.go")
	require.NotEqual(t, -1, idx, "main.go should be in the panel listing")
	assert.Equal(t, idx, panel.GetCursor(), "cursor should focus the opened file")
	assert.Empty(t, panel.TargetFile, "target file should be consumed by the reload")
	assert.Empty(t, panel.SearchBar.Value())
}

func TestSearchModeConfirmDir(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	cmd := typeSearchQuery(t, m, "deep")
	drainSearchResults(t, m, cmd)

	panel := m.getFocusedFilePanel()
	require.NotEmpty(t, panel.Search.Results, "deep should match its directory")
	require.True(t, panel.Search.Results[0].Dir, "the directory should rank above the file inside it")

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyEnter})

	assert.False(t, panel.Search.Active)
	assert.Equal(t, filepath.Join(root, "src", "deep"), panel.Location)
}

func TestSearchModeConfirmDuringDebounce(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	cmd := typeSearchQuery(t, m, "deep")
	drainSearchResults(t, m, cmd)

	panel := m.getFocusedFilePanel()
	require.NotEmpty(t, panel.Search.Results, "deep should match before editing")
	require.True(t, panel.Search.Results[0].Dir)

	// Query change clears stale results before debounced walk.
	TeaUpdate(m, utils.TeaRuneKeyMsg("z"))
	require.Empty(t, panel.Search.Results, "stale results must be cleared on query change")
	require.False(t, panel.Search.Done)

	// Confirm during debounce must not open stale result.
	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	assert.False(t, panel.Search.Active)
	assert.Equal(t, root, panel.Location, "must not navigate to the stale directory result")
}

func TestSearchModeConfirmEmptyQuery(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	panel := m.getFocusedFilePanel()
	assert.False(t, panel.Search.Active)
	assert.Equal(t, root, panel.Location)
}

func TestSearchModeConfigurableNav(t *testing.T) {
	root := searchTestDir(t)

	// Vim-style modifiers for test. Bare letters type.
	oldListUp, oldListDown := common.Hotkeys.ListUp, common.Hotkeys.ListDown
	common.Hotkeys.ListUp = []string{"up", "k", "ctrl+p"}
	common.Hotkeys.ListDown = []string{"down", "j", "ctrl+n"}
	t.Cleanup(func() {
		common.Hotkeys.ListUp, common.Hotkeys.ListDown = oldListUp, oldListDown
	})

	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	cmd := typeSearchQuery(t, m, "main")
	drainSearchResults(t, m, cmd)

	panel := m.getFocusedFilePanel()
	require.Len(t, panel.Search.Results, 2)
	require.Equal(t, 0, panel.Search.Cursor)

	TeaUpdate(m, tea.KeyPressMsg{Code: 'n', Mod: tea.ModCtrl})
	assert.Equal(t, 1, panel.Search.Cursor)
	assert.Equal(t, "main", panel.SearchBar.Value())
	TeaUpdate(m, tea.KeyPressMsg{Code: 'p', Mod: tea.ModCtrl})
	assert.Equal(t, 0, panel.Search.Cursor)
	assert.Equal(t, "main", panel.SearchBar.Value())

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Equal(t, 1, panel.Search.Cursor)
	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Equal(t, 0, panel.Search.Cursor)

	// Bare j/k must type, not move.
	TeaUpdate(m, utils.TeaRuneKeyMsg("j"))
	assert.Equal(t, "mainj", panel.SearchBar.Value(), "bare j must type, not move")
	require.True(t, panel.Search.Active)
}

func TestSearchToggleHiddenHotkeyConfigured(t *testing.T) {
	require.NotEmpty(t, common.Hotkeys.SearchToggleHidden,
		"search_toggle_hidden must have a default binding")
	assert.Equal(t, "alt+.", common.Hotkeys.SearchToggleHidden[0])
	assert.NotContains(t, common.Hotkeys.SearchToggleHidden, "ctrl+.",
		"ctrl+. is not deliverable by most terminals")
}

func TestSearchModeToggleHidden(t *testing.T) {
	root := t.TempDir()
	utils.SetupFiles(t,
		filepath.Join(root, "main.go"),
		filepath.Join(root, ".hidden_main.go"),
	)

	// Isolate persisted toggle.
	oldToggleFile := variable.ToggleDotFile
	variable.ToggleDotFile = filepath.Join(t.TempDir(), "toggleDotFile")
	t.Cleanup(func() { variable.ToggleDotFile = oldToggleFile })

	m := defaultTestModel(root)
	require.False(t, m.fileModel.DisplayDotFiles)

	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	panel := m.getFocusedFilePanel()
	require.True(t, panel.Search.Active)

	cmd := typeSearchQuery(t, m, "main")
	drainSearchResults(t, m, cmd)
	require.True(t, panel.Search.Done)
	assert.Equal(t, int64(1), panel.Search.MatchCount)
	require.Len(t, panel.Search.Results, 1)
	assert.Equal(t, "main.go", panel.Search.Results[0].Path)

	// Toggle preserves query and rewalks.
	toggleCmd := TeaUpdate(m, tea.KeyPressMsg{Code: '.', Mod: tea.ModAlt})
	require.True(t, panel.Search.Active, "toggle must not exit search")
	assert.Equal(t, "main", panel.SearchBar.Value(), "query must be preserved")
	assert.True(t, m.fileModel.DisplayDotFiles)
	drainSearchResults(t, m, toggleCmd)
	require.True(t, panel.Search.Done)
	assert.Equal(t, int64(2), panel.Search.MatchCount)
	assert.Len(t, panel.Search.Results, 2)

	TeaUpdate(m, utils.TeaRuneKeyMsg("."))
	assert.Equal(t, "main.", panel.SearchBar.Value())
	assert.True(t, m.fileModel.DisplayDotFiles, "typing . must not toggle")

	panel.SearchBar.SetValue("main")
	backCmd := TeaUpdate(m, tea.KeyPressMsg{Code: '.', Mod: tea.ModAlt})
	require.True(t, panel.Search.Active, "toggle must not exit search")
	assert.Equal(t, "main", panel.SearchBar.Value(), "query must be preserved")
	assert.False(t, m.fileModel.DisplayDotFiles)
	drainSearchResults(t, m, backCmd)
	require.True(t, panel.Search.Done)
	assert.Equal(t, int64(1), panel.Search.MatchCount)
	require.Len(t, panel.Search.Results, 1)
	assert.Equal(t, "main.go", panel.Search.Results[0].Path)

	// ctrl+. with Text degrades to typing.
	TeaUpdate(m, tea.KeyPressMsg{Code: '.', Mod: tea.ModCtrl, Text: "."})
	assert.Equal(t, "main.", panel.SearchBar.Value())
	assert.False(t, m.fileModel.DisplayDotFiles, "must not toggle")

	_ = os.Remove(variable.ToggleDotFile)
}

func TestSearchModeToggleHiddenPreservesCursor(t *testing.T) {
	root := t.TempDir()
	utils.SetupFiles(t,
		filepath.Join(root, "main.go"),
		filepath.Join(root, "main_test.go"),
		filepath.Join(root, ".hidden_main.go"),
	)

	// Isolate persisted toggle.
	oldToggleFile := variable.ToggleDotFile
	variable.ToggleDotFile = filepath.Join(t.TempDir(), "toggleDotFile")
	t.Cleanup(func() { variable.ToggleDotFile = oldToggleFile })

	m := defaultTestModel(root)
	require.False(t, m.fileModel.DisplayDotFiles)

	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	panel := m.getFocusedFilePanel()
	cmd := typeSearchQuery(t, m, "main")
	drainSearchResults(t, m, cmd)
	require.True(t, panel.Search.Done)
	require.Len(t, panel.Search.Results, 2)

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyDown})
	require.Equal(t, 1, panel.Search.Cursor)
	selected := panel.Search.Results[1].Path

	// Toggle rewalks; the selected path is still present and stays selected.
	toggleCmd := TeaUpdate(m, tea.KeyPressMsg{Code: '.', Mod: tea.ModAlt})
	assert.Equal(t, "main", panel.SearchBar.Value(), "query must be preserved")
	drainSearchResults(t, m, toggleCmd)
	require.True(t, panel.Search.Done)
	require.NotNil(t, panel.GetSearchCursorResult())
	assert.Equal(t, selected, panel.GetSearchCursorResult().Path, "cursor must follow selection")

	_ = os.Remove(variable.ToggleDotFile)
}

func TestSearchModeNavKeysKeepQuery(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))
	cmd := typeSearchQuery(t, m, "main")
	drainSearchResults(t, m, cmd)

	panel := m.getFocusedFilePanel()
	require.Len(t, panel.Search.Results, 2)

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Equal(t, "main", panel.SearchBar.Value(), "nav must not edit query")
	assert.Equal(t, 1, panel.Search.Cursor)

	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Equal(t, "main", panel.SearchBar.Value(), "nav must not edit query")
	assert.Equal(t, 0, panel.Search.Cursor)
}
