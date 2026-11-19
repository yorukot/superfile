package find

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestMain(m *testing.M) {
	originalFindFileSupport := common.Config.FindFileSupport
	common.Config.FindFileSupport = true
	defer func() {
		common.Config.FindFileSupport = originalFindFileSupport
	}()
	m.Run()
}

func TestFdErrorMessage(t *testing.T) {
	testdata := []struct {
		name     string
		retCode  int
		output   string
		expected string
	}{
		{
			name:     "zero exit code with empty output",
			retCode:  0,
			output:   "",
			expected: "",
		},
		{
			name:     "zero exit code ignores result lines",
			retCode:  0,
			output:   "a.txt\nb.txt\n",
			expected: "",
		},
		{
			name:     "nonzero exit code returns first non-empty line",
			retCode:  1,
			output:   "[fd error]: regex parse error:\n  detail\n",
			expected: "[fd error]: regex parse error:",
		},
		{
			name:     "nonzero exit code with empty output",
			retCode:  1,
			output:   "",
			expected: "",
		},
		{
			name:     "nonzero exit code trims surrounding whitespace",
			retCode:  2,
			output:   "some error",
			expected: "some error",
		},
		{
			name:     "nonzero exit code skips blank lines",
			retCode:  1,
			output:   "\n\n  spaced error  \n",
			expected: "spaced error",
		},
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			assert.Equal(t, td.expected, fdErrorMessage(td.retCode, td.output))
		})
	}
}

func TestParseFindResults(t *testing.T) {
	t.Run("parses fd output into absolute paths", func(t *testing.T) {
		dir := t.TempDir()
		realDir := filepath.Join(dir, "realdir")
		subDir := filepath.Join(realDir, "sub")
		aFile := filepath.Join(dir, "a.txt")
		symDir := filepath.Join(dir, "symdir")
		brokenLink := filepath.Join(dir, "broken")
		utils.SetupDirectories(t, realDir, subDir)
		utils.SetupFiles(t, aFile)
		require.NoError(t, os.Symlink(realDir, symDir))
		require.NoError(t, os.Symlink(filepath.Join(dir, "nonexistent"), brokenLink))

		output := "realdir/\nrealdir/sub/\na.txt\nsymdir\nbroken\n"
		results := parseFindResults(dir, output)

		require.Len(t, results, 5)
		assert.Equal(t, realDir, results[0].Path, "path should be absolute without trailing slash")
		assert.True(t, results[0].IsDir)
		assert.Equal(t, subDir, results[1].Path, "path should be absolute without trailing slash")
		assert.True(t, results[1].IsDir)
		assert.Equal(t, aFile, results[2].Path)
		assert.False(t, results[2].IsDir)
		assert.Equal(t, symDir, results[3].Path)
		assert.True(t, results[3].IsDir, "symlink to dir should be reported as dir")
		assert.Equal(t, brokenLink, results[4].Path)
		assert.False(t, results[4].IsDir, "broken symlink should not be reported as dir")
	})

	t.Run("empty output returns empty results", func(t *testing.T) {
		assert.Empty(t, parseFindResults(t.TempDir(), ""))
	})
}

func TestHandleConfirm(t *testing.T) {
	t.Run("no results returns NoAction", func(t *testing.T) {
		m := setupTestModel()

		action := m.handleConfirm()

		_, isNoAction := action.(common.NoAction)
		assert.True(t, isNoAction, "action should be NoAction when there are no results")
	})

	t.Run("file result returns GoToPathAction with IsDir false", func(t *testing.T) {
		m := setupTestModel()
		m.results = []FindResult{
			{Path: "/x/y.txt", IsDir: false},
		}

		action := m.handleConfirm()

		goToAction, ok := action.(common.GoToPathAction)
		require.True(t, ok, "action should be GoToPathAction")
		assert.Equal(t, common.GoToPathAction{Path: "/x/y.txt", IsDir: false}, goToAction)
	})

	t.Run("dir result returns GoToPathAction with IsDir true", func(t *testing.T) {
		m := setupTestModel()
		m.results = []FindResult{
			{Path: "/x/somedir", IsDir: true},
		}

		action := m.handleConfirm()

		goToAction, ok := action.(common.GoToPathAction)
		require.True(t, ok, "action should be GoToPathAction")
		assert.Equal(t, common.GoToPathAction{Path: "/x/somedir", IsDir: true}, goToAction)
	})
}

func TestJKKeyHandling(t *testing.T) {
	m := setupTestModel()
	m.Open("/some/dir")
	m.fdFound = true

	originalHotkeys := common.Hotkeys.ListDown
	originalUpHotkeys := common.Hotkeys.ListUp
	common.Hotkeys.ListDown = []string{"j", "down"}
	common.Hotkeys.ListUp = []string{"up", "k"}
	defer func() {
		common.Hotkeys.ListDown = originalHotkeys
		common.Hotkeys.ListUp = originalUpHotkeys
	}()

	action, cmd := m.HandleUpdate(utils.TeaRuneKeyMsg("j"))

	assert.NotNil(t, cmd, "HandleUpdate should return cmd for text input update")
	_, isNoAction := action.(common.NoAction)
	assert.True(t, isNoAction, "action should be NoAction for text input")
	assert.Equal(t, "j", m.textInput.Value(), "'j' should be added to textInput")

	action, cmd = m.HandleUpdate(utils.TeaRuneKeyMsg("k"))
	assert.NotNil(t, cmd, "HandleUpdate should return cmd for text input update")
	_, isNoAction = action.(common.NoAction)
	assert.True(t, isNoAction, "action should be NoAction for text input")
	assert.Equal(t, "jk", m.textInput.Value(), "'k' should be added to textInput")

	m.textInput.SetValue("")
	m.results = []FindResult{
		{Path: "/test/path1"},
		{Path: "/test/path2"},
	}
	m.cursor = 0

	action, cmd = m.HandleUpdate(tea.KeyPressMsg{Code: tea.KeyDown})
	assert.Nil(t, cmd, "HandleUpdate with down arrow should not return cmd")
	_, isNoAction = action.(common.NoAction)
	assert.True(t, isNoAction, "action should be NoAction for navigation")
	assert.Equal(t, 1, m.cursor, "down arrow should navigate down")
	assert.Empty(t, m.textInput.Value(), "down arrow should not add to textInput")

	action, cmd = m.HandleUpdate(tea.KeyPressMsg{Code: tea.KeyUp})
	assert.Nil(t, cmd, "HandleUpdate with up arrow should not return cmd")
	_, isNoAction = action.(common.NoAction)
	assert.True(t, isNoAction, "action should be NoAction for navigation")
	assert.Equal(t, 0, m.cursor, "up arrow should navigate up")
	assert.Empty(t, m.textInput.Value(), "up arrow should not add to textInput")
}

func TestApplyStaleDrop(t *testing.T) {
	m := setupTestModel()
	m.textInput.SetValue("abc")
	m.cursor = 2
	m.renderIndex = 3

	cmd := NewUpdateMsg("ab", []FindResult{{Path: "/stale"}}, "", 1).Apply(&m)
	assert.Nil(t, cmd)
	assert.Empty(t, m.results, "results should remain unchanged for stale query")
	assert.Equal(t, 2, m.cursor, "cursor should remain unchanged for stale query")
	assert.Equal(t, 3, m.renderIndex, "renderIndex should remain unchanged for stale query")

	results := []FindResult{
		{Path: "/x/y.txt", IsDir: false},
	}
	cmd = NewUpdateMsg("abc", results, "some error", 2).Apply(&m)
	assert.Nil(t, cmd)
	assert.Equal(t, results, m.results, "results should be applied for matching query")
	assert.Equal(t, "some error", m.errMsg, "errMsg should be set from the message")
	assert.Equal(t, 0, m.cursor, "cursor should be reset to 0")
	assert.Equal(t, 0, m.renderIndex, "renderIndex should be reset to 0")
}

func TestIsAvailable(t *testing.T) {
	testdata := []struct {
		name            string
		findFileSupport bool
		fdFound         bool
		expected        bool
	}{
		{
			name:            "config off, fd missing",
			findFileSupport: false,
			fdFound:         false,
			expected:        false,
		},
		{
			name:            "config off, fd found",
			findFileSupport: false,
			fdFound:         true,
			expected:        false,
		},
		{
			name:            "config on, fd missing",
			findFileSupport: true,
			fdFound:         false,
			expected:        false,
		},
		{
			name:            "config on, fd found",
			findFileSupport: true,
			fdFound:         true,
			expected:        true,
		},
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			originalFindFileSupport := common.Config.FindFileSupport
			common.Config.FindFileSupport = td.findFileSupport
			defer func() {
				common.Config.FindFileSupport = originalFindFileSupport
			}()

			m := setupTestModel()
			m.fdFound = td.fdFound

			assert.Equal(t, td.expected, m.isAvailable())
		})
	}
}

func TestGetQueryCmdNil(t *testing.T) {
	testdata := []struct {
		name            string
		findFileSupport bool
		fdFound         bool
		expectedNil     bool
	}{
		{
			name:            "config off, fd found",
			findFileSupport: false,
			fdFound:         true,
			expectedNil:     true,
		},
		{
			name:            "config on, fd missing",
			findFileSupport: true,
			fdFound:         false,
			expectedNil:     true,
		},
		{
			name:            "config on, fd found",
			findFileSupport: true,
			fdFound:         true,
			expectedNil:     false,
		},
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			originalFindFileSupport := common.Config.FindFileSupport
			common.Config.FindFileSupport = td.findFileSupport
			defer func() {
				common.Config.FindFileSupport = originalFindFileSupport
			}()

			m := setupTestModel()
			m.fdFound = td.fdFound

			cmd := m.GetQueryCmd("query")
			if td.expectedNil {
				assert.Nil(t, cmd, "GetQueryCmd should return nil")
			} else {
				assert.NotNil(t, cmd, "GetQueryCmd should return query cmd")
			}
		})
	}
}

func TestHandleUpdate(t *testing.T) {
	t.Run("closed modal returns NoAction and makes no changes", func(t *testing.T) {
		m := setupTestModel()
		m.textInput.SetValue("existing")

		action, cmd := m.HandleUpdate(utils.TeaRuneKeyMsg("a"))

		_, isNoAction := action.(common.NoAction)
		assert.True(t, isNoAction, "action should be NoAction when modal is closed")
		assert.Nil(t, cmd, "no cmd should be returned when modal is closed")
		assert.False(t, m.IsOpen(), "modal should remain closed")
		assert.Equal(t, "existing", m.textInput.Value(), "text input should not change when modal is closed")
	})

	t.Run("opening keypress is suppressed, then typing works", func(t *testing.T) {
		m := setupTestModel()
		originalFindFile := common.Hotkeys.FindFile
		common.Hotkeys.FindFile = []string{"ctrl+f"}
		defer func() {
			common.Hotkeys.FindFile = originalFindFile
		}()

		m.textInput.SetValue("old")
		openCmd := m.Open("/some/dir")
		m.fdFound = true
		if _, lookupErr := exec.LookPath("fd"); lookupErr == nil {
			assert.NotNil(t, openCmd, "Open should return async query cmd when fd is available")
		}

		action, cmd := m.HandleUpdate(utils.TeaRuneKeyMsg("ctrl+f"))
		_, isNoAction := action.(common.NoAction)
		assert.True(t, isNoAction, "opening keypress should be ignored")
		assert.Nil(t, cmd, "opening keypress should not trigger a query")
		assert.False(t, m.justOpened, "justOpened should be cleared after opening keypress")
		assert.Empty(t, m.textInput.Value(), "opening keypress should not be added to text input")

		action, cmd = m.HandleUpdate(utils.TeaRuneKeyMsg("a"))
		_, isNoAction = action.(common.NoAction)
		assert.True(t, isNoAction, "normal typing should return NoAction")
		assert.NotNil(t, cmd, "normal typing should trigger a query")
		assert.Equal(t, "a", m.textInput.Value(), "'a' should be added to text input")
	})
}

func TestNavigation(t *testing.T) {
	testdata := []struct {
		name           string
		resultCnt      int
		startCursor    int
		navigateUp     bool
		expectedCursor int
	}{
		{
			name:           "navigateDown from position 0 moves to next position",
			resultCnt:      7,
			startCursor:    0,
			navigateUp:     false,
			expectedCursor: 1,
		},
		{
			name:           "navigateUp at position 0 wraps to last position",
			resultCnt:      7,
			startCursor:    0,
			navigateUp:     true,
			expectedCursor: 6,
		},
		{
			name:           "navigateDown at last position wraps to first position",
			resultCnt:      7,
			startCursor:    6,
			navigateUp:     false,
			expectedCursor: 0,
		},
		{
			name:           "navigateUp at position 3 decrements to 2",
			resultCnt:      7,
			startCursor:    3,
			navigateUp:     true,
			expectedCursor: 2,
		},
		{
			name:           "navigateUp with empty results keeps cursor at 0",
			resultCnt:      0,
			startCursor:    0,
			navigateUp:     true,
			expectedCursor: 0,
		},
		{
			name:           "navigateDown with empty results keeps cursor at 0",
			resultCnt:      0,
			startCursor:    0,
			navigateUp:     false,
			expectedCursor: 0,
		},
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			var m Model
			if td.resultCnt == 0 {
				m = setupTestModel()
			} else {
				m = setupTestModelWithResults(td.resultCnt)
			}
			m.cursor = td.startCursor
			if td.navigateUp {
				m.navigateUp()
			} else {
				m.navigateDown()
			}
			assert.Equal(t, td.expectedCursor, m.cursor)
		})
	}
}

func TestNavigationUpdatesRenderIndex(t *testing.T) {
	// Shrink the 50-row helper modal to the height that yields a
	// 5-row visible window, which these 6 presses of down overflow.
	const testWindow = 5
	m := setupTestModelWithResults(7)
	m.maxHeight = testWindow + visibleOverhead
	for range 6 {
		m.navigateDown()
	}
	assert.Equal(t, 6, m.cursor, "cursor should be at last position after 6 downs")
	assert.Equal(t, 2, m.renderIndex, "renderIndex should follow cursor past the visible window")
}

func TestUpdateRenderIndex(t *testing.T) {
	// The helper model has a 50-row max height (a 43-row visible window).
	// Shrink it to the height that yields a 5-row window so the
	// cursor/renderIndex math these cases target is unchanged.
	const testWindow = 5

	testdata := []struct {
		name                string
		resultCnt           int
		cursor              int
		initialRenderIndex  int
		expectedRenderIndex int
	}{
		{
			name:                "cursor at 0 has renderIndex 0",
			resultCnt:           10,
			cursor:              0,
			initialRenderIndex:  0,
			expectedRenderIndex: 0,
		},
		{
			name:                "cursor at last visible position has renderIndex 1",
			resultCnt:           10,
			cursor:              testWindow,
			initialRenderIndex:  0,
			expectedRenderIndex: 1,
		},
		{
			name:                "cursor at last result has renderIndex at last page",
			resultCnt:           10,
			cursor:              9,
			initialRenderIndex:  0,
			expectedRenderIndex: 10 - testWindow,
		},
		{
			name:                "cursor above renderIndex causes renderIndex to decrease",
			resultCnt:           10,
			cursor:              2,
			initialRenderIndex:  5,
			expectedRenderIndex: 2,
		},
		{
			name:                "cursor at 0 with high renderIndex snaps renderIndex to 0",
			resultCnt:           10,
			cursor:              0,
			initialRenderIndex:  4,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 0",
			resultCnt:           3,
			cursor:              0,
			initialRenderIndex:  0,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 1",
			resultCnt:           3,
			cursor:              1,
			initialRenderIndex:  0,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 2",
			resultCnt:           3,
			cursor:              2,
			initialRenderIndex:  0,
			expectedRenderIndex: 0,
		},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			m := setupTestModelWithResults(td.resultCnt)
			m.maxHeight = testWindow + visibleOverhead
			m.cursor = td.cursor
			m.renderIndex = td.initialRenderIndex
			m.updateRenderIndex()

			assert.Equal(t, td.expectedRenderIndex, m.renderIndex)
		})
	}
}

func TestGetQueryCmdDashPrefixedQuery(t *testing.T) {
	if _, err := exec.LookPath("fd"); err != nil {
		t.Skip("fd is not installed")
	}

	m := setupTestModel()
	m.fdFound = true
	m.searchDir = t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(m.searchDir, "-h.txt"), []byte("x"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(m.searchDir, "normal.txt"), []byte("x"), 0o644))

	// Without the `--` terminator, fd would parse `-h` as a flag and print
	// its help text, which the parser would mistake for match results
	msg := m.GetQueryCmd("-h")()
	update, ok := msg.(UpdateMsg)
	require.True(t, ok)

	assert.Empty(t, update.errMsg)
	assert.Len(t, update.results, 1)
	assert.Equal(t, "-h.txt", filepath.Base(update.results[0].Path))
}

func TestApplyDropsPreviousSessionResults(t *testing.T) {
	m := setupTestModel()
	m.open = true
	m.textInput.SetValue("abc")
	m.openReqID = 5
	m.reqCnt = 6

	stale := []FindResult{{Path: "/old/root/a.txt"}}
	cmd := NewUpdateMsg("abc", stale, "", 3).Apply(&m)
	assert.Nil(t, cmd)
	assert.Empty(t, m.results, "completion from a previous session must be dropped")

	fresh := []FindResult{{Path: "/new/root/b.txt"}}
	cmd = NewUpdateMsg("abc", fresh, "", 5).Apply(&m)
	assert.Nil(t, cmd)
	assert.Equal(t, fresh, m.results, "completion from the current session must apply")
}

func TestApplyDropsOlderRequestResult(t *testing.T) {
	m := setupTestModel()
	m.open = true
	m.textInput.SetValue("a")
	// The a -> ab -> a typing sequence submitted requests 0, 1, 2;
	// the latest submitted request is 2
	m.openReqID = 0
	m.reqCnt = 3

	// The older completion for the same query text must not overwrite
	// the latest results or reset the selection
	cmd := NewUpdateMsg("a", []FindResult{{Path: "/old"}}, "", 0).Apply(&m)
	assert.Nil(t, cmd)
	assert.Empty(t, m.results, "older request completion must be dropped even when the query matches")

	cmd = NewUpdateMsg("a", []FindResult{{Path: "/latest"}}, "", 2).Apply(&m)
	assert.Nil(t, cmd)
	assert.Equal(t, []FindResult{{Path: "/latest"}}, m.results, "latest request completion must apply")
}

func TestOpenStartsFreshSession(t *testing.T) {
	m := setupTestModel()
	m.results = []FindResult{{Path: "/stale/x.txt"}}
	m.errMsg = "old error"
	m.cursor = 2
	m.renderIndex = 3

	m.Open("/tmp")

	assert.True(t, m.IsOpen())
	assert.Empty(t, m.results, "opening a new session should clear stale results")
	assert.Empty(t, m.errMsg, "opening a new session should clear the stale error")
	assert.Equal(t, 0, m.cursor)
	assert.Equal(t, 0, m.renderIndex)
	if m.isAvailable() {
		assert.Equal(t, m.reqCnt-1, m.openReqID, "initial query must carry the session's starting id")
	}
}
