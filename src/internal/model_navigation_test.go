package internal

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"
)

func TestFilePanelNavigation(t *testing.T) {
	/*
		We want to test
		(1) Switching to parent directory
		(2) Switching to parent on being at root "/"
		(3) Entering current directory
		(4) Entering via cd / command
		(5) Cd to itself via cd . command

		Make sure to validate
		- Search bar is cleared
		- The cursor and render values are restored correctly
	*/

	curTestDir := t.TempDir()
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	file1 := filepath.Join(curTestDir, "file1.txt")
	// We >=3 files in dir1 and >=2 files in dir2
	// so that cursor=2, and cursor=1 are valid values.
	file2 := filepath.Join(dir1, "file2.txt")
	file3 := filepath.Join(dir1, "file3.txt")
	file4 := filepath.Join(dir1, "file4.txt")
	file5 := filepath.Join(dir2, "file5.txt")
	file6 := filepath.Join(dir2, "file6.txt")

	rootDir := "/"

	if runtime.GOOS == utils.OsWindows {
		rootDir = "\\"
	}

	utils.SetupDirectories(t, dir1, dir2)
	utils.SetupFiles(t, file1, file2, file3, file4, file5, file6)

	testdata := []struct {
		name           string
		startDir       string
		resultDir      string
		startCursor    int
		keyInput       []string
		searchBarClear bool
	}{
		{
			name:        "Switch to parent",
			startDir:    dir1,
			resultDir:   curTestDir,
			startCursor: 1,
			keyInput: []string{
				common.Hotkeys.ParentDirectory[0],
			},
			searchBarClear: true,
		},
		{
			name:        "Switch to parent when at root",
			startDir:    rootDir,
			resultDir:   rootDir,
			startCursor: 0,
			keyInput: []string{
				common.Hotkeys.ParentDirectory[0],
			},
			searchBarClear: false,
		},
		{
			name:        "Enter current directory",
			startDir:    curTestDir,
			resultDir:   dir2,
			startCursor: 1,
			keyInput: []string{
				common.Hotkeys.Confirm[0],
			},
			searchBarClear: true,
		},
		{
			name:        "Enter via cd command first dir",
			startDir:    curTestDir,
			resultDir:   dir1,
			startCursor: 0,
			keyInput: []string{
				common.Hotkeys.OpenSPFPrompt[0],
				// TODO : Have it quoted, once cd command supports quoted paths
				"cd " + dir1,
				common.Hotkeys.ConfirmTyping[0],
			},
			searchBarClear: true,
		},
		{
			name:        "cd . should be ignored",
			startDir:    curTestDir,
			resultDir:   curTestDir,
			startCursor: 2,
			keyInput: []string{
				common.Hotkeys.OpenSPFPrompt[0],
				// TODO : Have it quoted, once cd command supports quoted paths
				"cd .",
				common.Hotkeys.ConfirmTyping[0],
			},
			searchBarClear: false,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			m := defaultTestModel(tt.startDir)
			for range tt.startCursor {
				m.getFocusedFilePanel().ListDown()
			}
			require.Equal(t, tt.startCursor, m.getFocusedFilePanel().GetCursor())
			originalRenderIndex := m.getFocusedFilePanel().GetRenderIndex()
			for _, s := range tt.keyInput {
				TeaUpdate(m, utils.TeaRuneKeyMsg(s))
			}

			assert.Equal(t, tt.resultDir, m.getFocusedFilePanel().Location)

			if tt.searchBarClear {
				assert.Empty(t, m.getFocusedFilePanel().SearchBar.Value())
			}

			// Go back to original directory

			TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
			TeaUpdate(m, utils.TeaRuneKeyMsg("cd "+tt.startDir))
			TeaUpdate(m, tea.KeyMsg{Type: tea.KeyEnter})

			// Make sure we have original cursor and render
			assert.Equal(t, tt.startCursor, m.getFocusedFilePanel().GetCursor())
			assert.Equal(t, originalRenderIndex, m.getFocusedFilePanel().GetRenderIndex())
		})
	}

	t.Run("Focus on current directory on navigation to parent directory", func(t *testing.T) {
		m := defaultTestModel(dir2)
		p := NewTestTeaProgWithEventLoop(t, m)
		p.SendKey(common.Hotkeys.ParentDirectory[0])

		assert.Eventually(t, func() bool {
			return m.getFocusedFilePanel().GetFocusedItem().Location == dir2 &&
				m.getFocusedFilePanel().GetCursor() == 1
		}, DefaultTestTimeout, DefaultTestTick)
	})
}

func TestCursorOutOfBoundsAfterDirectorySwitch(t *testing.T) {
	// Create two directories with different file counts
	tempDir := t.TempDir()
	dir1 := filepath.Join(tempDir, "dir1")
	dir2 := filepath.Join(tempDir, "dir2")
	utils.SetupDirectories(t, dir1, dir2)

	var files1, files2 []string
	for i := range 10 {
		files1 = append(files1, filepath.Join(dir1, string('a'+rune(i))+".txt"))
	}
	for i := range 5 {
		files2 = append(files2, filepath.Join(dir2, string('a'+rune(i))+".txt"))
	}
	utils.SetupFiles(t, files1...)
	utils.SetupFiles(t, files2...)

	// Start with dir1
	m := defaultTestModel(dir1)
	p := NewTestTeaProgWithEventLoop(t, m)

	// It will immediately load as defaultTestModel does one sync TeaUpdate
	assert.Equal(t, 10, m.getFocusedFilePanel().ElemCount(),
		"Should load 10 files in dir1")

	// Move cursor to position 8 (near end of list)
	panel := m.getFocusedFilePanel()
	for range 8 {
		p.Send(tea.KeyMsg{Type: tea.KeyDown})
	}

	// Verify cursor is at position 8
	assert.Eventually(t, func() bool {
		return m.getFocusedFilePanel().GetCursor() == 8
	}, DefaultTestTimeout, DefaultTestTick, "Cursor should be at position 8")
	t.Logf("Cursor at position %d with %d elements", panel.GetCursor(), panel.ElemCount())

	// Navigate to dir2 (this saves cursor=8 in directoryRecords)
	navigateToTargetDir(t, m, dir1, dir2)

	assert.Equal(t, dir2, m.getFocusedFilePanel().Location, "Should be in dir2")
	assert.Equal(t, 5, m.getFocusedFilePanel().ElemCount())

	for i := 4; i < 10; i++ {
		err := os.Remove(files1[i])
		require.NoError(t, err)
	}
	t.Log("Deleted 6 files from dir1 externally")

	// Navigate back to dir1 (this restores cursor=8 from cache)
	navigateToTargetDir(t, m, dir2, dir1)
	assert.Equal(t, 0, panel.GetCursor(), "Cursor not restored as is from directoryRecords cache")
	assert.NoError(t, panel.ValidateCursorAndRenderIndex(), "panel not valid")
}

func TestCursorRemembersParentPosition(t *testing.T) {
	/*
		We want to test that the cursor remembers its position in the parent directory in 3 different cases
		(1) jump back from child with more elements than parent and near top of list of parent
		(2) jump back from child with less elements than parent and near end of list of parent
		(3) jump back from child with no elements
	*/

	curTestDir := t.TempDir()
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	dir3 := filepath.Join(curTestDir, "dir3")
	dir4 := filepath.Join(curTestDir, "dir4")
	dir5 := filepath.Join(curTestDir, "dir5")
	file1 := filepath.Join(dir2, "file1.txt")
	file2 := filepath.Join(dir2, "file2.txt")
	file3 := filepath.Join(dir2, "file3.txt")
	file4 := filepath.Join(dir2, "file4.txt")
	file5 := filepath.Join(dir2, "file5.txt")
	file6 := filepath.Join(dir2, "file6.txt")
	file7 := filepath.Join(dir2, "file7.txt")
	file8 := filepath.Join(dir4, "file8.txt")
	file9 := filepath.Join(dir4, "file9.txt")

	utils.SetupDirectories(t, dir1, dir2, dir3, dir4, dir5)
	utils.SetupFiles(t, file1, file2, file3, file4, file5, file6, file7, file8, file9)

	cases := []struct {
		name           string
		moveDowns      int
		childDir       string
		expectedCursor int
	}{
		{"case1", 1, dir2, 1},
		{"case2", 3, dir4, 3},
		{"case3", 4, dir5, 4},
	}

	// if a case fails the next case(s) fail also
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := defaultTestModel(curTestDir)

			for range tc.moveDowns {
				m.getFocusedFilePanel().ListDown()
			}

			originalRenderIndex := m.getFocusedFilePanel().GetRenderIndex()

			assert.Eventually(t, func() bool {
				return m.getFocusedFilePanel().GetCursor() == tc.expectedCursor
			}, DefaultTestTimeout, DefaultTestTick, "Cursor should be at correct position")

			// Move into child directory
			TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.Confirm[0]))

			assert.Eventually(t, func() bool {
				return m.getFocusedFilePanel().Location == tc.childDir
			}, DefaultTestTimeout, DefaultTestTick, "Should have stepped into child directory")

			// Go back to original directory
			TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.ParentDirectory[0]))

			assert.Eventually(t, func() bool {
				return m.getFocusedFilePanel().Location == curTestDir
			}, DefaultTestTimeout, DefaultTestTick, "Should have stepped into parent directory curTestDir")

			// Make sure we have original cursor and render
			assert.Equal(
				t,
				tc.expectedCursor,
				m.getFocusedFilePanel().GetCursor(),
				"Should have remembered cursor position in parent",
			)
			assert.Equal(t, originalRenderIndex, m.getFocusedFilePanel().GetRenderIndex())
		})
	}
}
