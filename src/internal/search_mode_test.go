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

// drainSearchResults pumps search progress messages into the model until the
// search completes. It unwraps batch commands produced by tea.Batch.
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

	// Esc cancels the session and restores the panel listing
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

	// Opening a file result navigates to its containing directory and
	// focuses the file.
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

	// Opening a directory result makes it the new focused directory.
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

	// Editing the query clears stale results immediately, before the
	// debounced walk for the new query runs. "deepz" matches nothing, so
	// any retained result would be stale.
	TeaUpdate(m, utils.TeaRuneKeyMsg("z"))
	require.Empty(t, panel.Search.Results, "stale results must be cleared on query change")
	require.False(t, panel.Search.Done)

	// Confirming during the debounce must not open the stale result.
	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	assert.False(t, panel.Search.Active)
	assert.Equal(t, root, panel.Location, "must not navigate to the stale directory result")
}

func TestSearchModeConfirmEmptyQuery(t *testing.T) {
	root := searchTestDir(t)
	m := defaultTestModel(root)
	TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.SearchMode[0]))

	// Confirming with an empty query just cancels the session.
	TeaUpdate(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	panel := m.getFocusedFilePanel()
	assert.False(t, panel.Search.Active)
	assert.Equal(t, root, panel.Location)
}

func TestSearchToggleHiddenHotkeyConfigured(t *testing.T) {
	require.NotEmpty(t, common.Hotkeys.SearchToggleHidden,
		"search_toggle_hidden must have a default binding")
	assert.Equal(t, "alt+.", common.Hotkeys.SearchToggleHidden[0])
	assert.Contains(t, common.Hotkeys.SearchToggleHidden, "ctrl+.")
}

func TestSearchModeToggleHidden(t *testing.T) {
	root := t.TempDir()
	utils.SetupFiles(t,
		filepath.Join(root, "main.go"),
		filepath.Join(root, ".hidden_main.go"),
	)

	// Keep the persisted dotfile toggle out of the real data dir.
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

	// Toggle hidden files without leaving search. Query must be preserved
	// and the search must rewalk with the new visibility.
	toggleCmd := TeaUpdate(m, tea.KeyPressMsg{Code: '.', Mod: tea.ModAlt})
	require.True(t, panel.Search.Active, "toggle must not exit search")
	assert.Equal(t, "main", panel.SearchBar.Value(), "query must be preserved")
	assert.True(t, m.fileModel.DisplayDotFiles)
	drainSearchResults(t, m, toggleCmd)
	require.True(t, panel.Search.Done)
	assert.Equal(t, int64(2), panel.Search.MatchCount)
	assert.Len(t, panel.Search.Results, 2)

	// Typing "." must still edit the query, not toggle visibility.
	TeaUpdate(m, utils.TeaRuneKeyMsg("."))
	assert.Equal(t, "main.", panel.SearchBar.Value())
	assert.True(t, m.fileModel.DisplayDotFiles, "typing . must not toggle")

	// Toggle back via ctrl+., including the leaky terminal variant that
	// reports Text="." alongside the modifier. The dot must not pollute
	// the query: the key is consumed by the search mode handler.
	panel.SearchBar.SetValue("main")
	backCmd := TeaUpdate(m, tea.KeyPressMsg{Code: '.', Mod: tea.ModCtrl, Text: "."})
	require.True(t, panel.Search.Active, "toggle must not exit search")
	assert.Equal(t, "main", panel.SearchBar.Value(), "leaky ctrl+. must not type")
	assert.False(t, m.fileModel.DisplayDotFiles)
	drainSearchResults(t, m, backCmd)
	require.True(t, panel.Search.Done)
	assert.Equal(t, int64(1), panel.Search.MatchCount)
	require.Len(t, panel.Search.Results, 1)
	assert.Equal(t, "main.go", panel.Search.Results[0].Path)

	_ = os.Remove(variable.ToggleDotFile)
}
