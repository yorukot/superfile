package internal

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	cursor "charm.land/bubbles/v2/cursor"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/preview"
)

// Named values for TestFindPreview, to keep magic numbers out of the asserts

const findPreviewAllResults = 6
const findPreviewPaneWidth = 49
const findPreviewMaxQueue = 200
const findPreviewBigCut = 200
const findPreviewPaneLeft = 71
const findPreviewModalTop = 12
const findPreviewModalBottom = 35
const findPreviewMinRows = 36
const findPreviewHeadline = "Find File/Folder"

const findResizeWidth = 80
const findResizeHeight = 30
const findResizePreviewWidth = 29
const findResizePreviewHeight = 30
const findResizeModalWidth = 40
const findResizeModalMaxHeight = 15

// findPreviewDir returns the root of the TestFindPreview fixture, resolved at
// call time so the package testDir can change between test runs
func findPreviewDir() string {
	return filepath.Join(testDir, "TestFindPreview")
}

func setupProgForFind(t *testing.T, dir string) *TeaProg {
	t.Helper()
	common.Config.FindFileSupport = true
	m := defaultTestModel(dir)
	return NewTestTeaProgWithEventLoop(t, m)
}

func openFind(t *testing.T, p *TeaProg) {
	t.Helper()
	p.SendKey(common.Hotkeys.FindFile[0])
	assert.Eventually(t, func() bool {
		return p.getModel().findModal.IsOpen()
	}, DefaultTestTimeout, DefaultTestTick, "Find modal should open")
}

func TestFind(t *testing.T) {
	if _, err := exec.LookPath("fd"); err != nil {
		t.Skipf("fd not installed")
	}

	originalFindFileSupport := common.Config.FindFileSupport
	defer func() {
		common.Config.FindFileSupport = originalFindFileSupport
	}()

	curTestDir := filepath.Join(testDir, "TestFind")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	nestedDir := filepath.Join(dir2, "nested")
	utils.SetupDirectories(t, curTestDir, dir1, dir2, nestedDir)
	utils.SetupFiles(t,
		filepath.Join(dir1, "file1.txt"),
		filepath.Join(dir2, "topfile.txt"),
		filepath.Join(nestedDir, "target_file.txt"))
	t.Cleanup(func() {
		os.RemoveAll(curTestDir)
	})

	t.Run("Find file navigation", func(t *testing.T) {
		p := setupProgForFind(t, curTestDir)
		openFind(t, p)

		p.SendKey("target")
		targetFile := filepath.Join(nestedDir, "target_file.txt")
		assert.Eventually(t, func() bool {
			results := p.getModel().findModal.GetResults()
			return len(results) == 1 && results[0].Path == targetFile && !results[0].IsDir
		}, DefaultTestTimeout, DefaultTestTick, "target_file.txt should be found by find")

		// Press enter to navigate to the file's directory
		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			panel := p.getModel().getFocusedFilePanel()
			return !p.getModel().findModal.IsOpen() &&
				panel.Location == nestedDir &&
				panel.TargetFile == "" &&
				panel.GetFocusedItem().Name == "target_file.txt"
		}, DefaultTestTimeout, DefaultTestTick,
			"Find modal should close and navigate to %s (current location: %s)",
			nestedDir, p.getModel().getFocusedFilePanel().Location)
	})

	t.Run("Find folder navigation", func(t *testing.T) {
		p := setupProgForFind(t, curTestDir)
		openFind(t, p)

		p.SendKey("nested")
		assert.Eventually(t, func() bool {
			results := p.getModel().findModal.GetResults()
			return len(results) == 1 && results[0].Path == nestedDir && results[0].IsDir
		}, DefaultTestTimeout, DefaultTestTick, "nested dir should be found by find")

		// Press enter to navigate to the directory
		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			return !p.getModel().findModal.IsOpen() &&
				p.getModel().getFocusedFilePanel().Location == nestedDir
		}, DefaultTestTimeout, DefaultTestTick,
			"Find modal should close and navigate to %s (current location: %s)",
			nestedDir, p.getModel().getFocusedFilePanel().Location)
	})

	t.Run("Find disabled shows no results", func(t *testing.T) {
		common.Config.FindFileSupport = false
		m := defaultTestModel(dir1)

		TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.FindFile[0]))
		assert.True(t, m.findModal.IsOpen(), "Find modal should open even when FindFileSupport is disabled")

		results := m.findModal.GetResults()
		assert.Empty(t, results, "Find modal should show no results when FindFileSupport is disabled")
	})

	t.Run("Find 'ctrl+f' key suppression on open", func(t *testing.T) {
		p := setupProgForFind(t, curTestDir)
		openFind(t, p)
		assert.Empty(t, p.getModel().findModal.GetTextInputValue(),
			"The 'ctrl+f' key should not be added to textInput")
		p.SendKeyDirectly("abc")
		assert.Equal(t, "abc", p.getModel().findModal.GetTextInputValue())
	})

	t.Run("Search filter cleared on file selection", func(t *testing.T) {
		target2File := filepath.Join(dir2, "target2.txt")
		utils.SetupFiles(t, target2File)
		t.Cleanup(func() {
			os.Remove(target2File)
		})

		p := setupProgForFind(t, dir2)
		p.getModel().getFocusedFilePanel().SearchBar.SetValue("topfile")

		openFind(t, p)

		p.SendKey("target2")
		assert.Eventually(t, func() bool {
			results := p.getModel().findModal.GetResults()
			return len(results) == 1 && results[0].Path == target2File && !results[0].IsDir
		}, DefaultTestTimeout, DefaultTestTick, "target2.txt should be found by find")

		p.SendKey(common.Hotkeys.ConfirmTyping[0])
		assert.Eventually(t, func() bool {
			panel := p.getModel().getFocusedFilePanel()
			return !p.getModel().findModal.IsOpen() &&
				panel.Location == dir2 &&
				panel.SearchBar.Value() == "" &&
				panel.TargetFile == "" &&
				panel.GetFocusedItem().Name == "target2.txt"
		}, DefaultTestTimeout, DefaultTestTick,
			"Search filter should be cleared and cursor should move to %s (search bar: %s)",
			target2File, p.getModel().getFocusedFilePanel().SearchBar.Value())
	})
}

// runSettled processes the command tree produced by a single model update and
// returns how many file preview updates it applies along the way
func runSettled(t *testing.T, m *model, cmd tea.Cmd) int {
	t.Helper()

	queue := []tea.Cmd{}
	if cmd != nil {
		queue = append(queue, cmd)
	}
	applied := 0
	for len(queue) > 0 {
		if len(queue) > findPreviewMaxQueue {
			t.Fatalf("command queue exceeded %d commands, the update chain is not settling", findPreviewMaxQueue)
		}
		leaf := queue[0]
		queue = queue[1:]

		msg := ExecuteTeaCmdWithTimeout(leaf, DefaultTestTimeout)
		if msg == nil {
			continue
		}

		switch typed := msg.(type) {
		case cursor.BlinkMsg:
			// Blinking re-arms itself forever, so the loop must not process it
			continue
		case tea.BatchMsg:
			queue = append(queue, typed...)
			continue
		case preview.UpdateMsg:
			applied++
		}

		// The update may finish without pending work, in which case the
		// command is nil and the chain is done
		if next := TeaUpdate(m, msg); next != nil {
			queue = append(queue, next)
		}
	}
	return applied
}

// primeFindPreviewModel builds a test model with the file preview pane open,
// settles the initial preview render, and focuses the first item in the panel
func primeFindPreviewModel(t *testing.T, dir string) *model {
	t.Helper()
	m := defaultTestModelWithFilePreview(dir)
	runSettled(t, m, TeaUpdate(m, tea.WindowSizeMsg{
		Width: DefaultTestModelWidth, Height: DefaultTestModelHeight,
	}))
	// The model setup discards the initial render command, but the pane's
	// location was set when that command was created. Force one render so
	// the pane is fully settled with content before the test starts
	settled := runSettled(t, m, m.fileModel.GetFilePreviewCmd(true))
	assert.Equal(t, 1, settled, "the forced prime render should apply exactly once")
	setFilePanelSelectedItemByName(t, m.getFocusedFilePanel(), "dir1")
	return m
}

// openFindDirect opens the find modal directly, bypassing key routing, and
// settles the initial query
func openFindDirect(t *testing.T, m *model) int {
	t.Helper()
	return runSettled(t, m, m.findModal.Open(m.getFocusedFilePanel().Location))
}

// sendFindKey sends a key to the model, the same way the find modal receives
// it, and settles the commands it triggers
func sendFindKey(t *testing.T, m *model, key string) int {
	t.Helper()
	return runSettled(t, m, TeaUpdate(m, utils.TeaRuneKeyMsg(key)))
}

// sendFindArrow sends a raw arrow key press, then settles the commands it
// triggers. Raw presses are required because the single-letter aliases are
// alphanumeric and would be typed into the find input
func sendFindArrow(t *testing.T, m *model, code rune) int {
	t.Helper()
	return runSettled(t, m, TeaUpdate(m, tea.KeyPressMsg{Code: code}))
}

// frameRows returns the rendered frame of the model, one line per row, with
// the trailing newline stripped
func frameRows(t *testing.T, m *model) []string {
	t.Helper()
	rows := strings.Split(m.View().Content, "\n")
	if len(rows) > 0 && rows[len(rows)-1] == "" {
		rows = rows[:len(rows)-1]
	}
	return rows
}

// TestFindPreview covers the file preview pane interaction with the find
// modal: while find is open, the pane previews the result under the cursor,
// and the modal is placed clear of the pane
func TestFindPreview(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping find preview test in short mode")
	}
	if _, err := exec.LookPath("fd"); err != nil {
		t.Skipf("fd not installed")
	}

	originalFindFileSupport := common.Config.FindFileSupport
	defer func() {
		common.Config.FindFileSupport = originalFindFileSupport
	}()
	common.Config.FindFileSupport = true

	curTestDir := filepath.Join(testDir, "TestFindPreview")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	nestedDir := filepath.Join(dir2, "nested")
	utils.SetupDirectories(t, curTestDir, dir1, dir2, nestedDir)
	utils.SetupFiles(t,
		filepath.Join(dir1, "file1.txt"),
		filepath.Join(dir2, "topfile.txt"),
		filepath.Join(nestedDir, "target_file.txt"))
	t.Cleanup(func() {
		os.RemoveAll(curTestDir)
	})

	t.Run("Navigation before find", func(t *testing.T) {
		findPreviewNavigateBeforeFind(t)
	})

	t.Run("Open find keeps pane on the cursor item", func(t *testing.T) {
		findPreviewOpenKeepsPane(t)
	})

	t.Run("Arrows move the cursor and the pane follows", func(t *testing.T) {
		findPreviewArrowsMoveCursor(t)
	})

	t.Run("Typing filters results and the pane follows", func(t *testing.T) {
		findPreviewTypingFiltersResults(t)
	})

	t.Run("Typing one key at a time", func(t *testing.T) {
		findPreviewTypingOneKeyAtATime(t)
	})

	t.Run("Clearing the query releases the pane", func(t *testing.T) {
		findPreviewClearingQueryReleasesPane(t)
	})

	t.Run("Cancel typing releases the pane", func(t *testing.T) {
		findPreviewCancelTypingReleasesPane(t)
	})

	t.Run("Cancel typing right after opening", func(t *testing.T) {
		findPreviewCancelTypingRightAfterOpening(t)
	})

	t.Run("Confirm navigates to the found file", func(t *testing.T) {
		findPreviewConfirmNavigatesToFile(t)
	})

	t.Run("Confirm navigates to the found directory", func(t *testing.T) {
		findPreviewConfirmNavigatesToDirectory(t)
	})

	t.Run("Toggle the pane while find is open", func(t *testing.T) {
		findPreviewTogglePaneWhileFindOpen(t)
	})

	t.Run("Resize while find is open", func(t *testing.T) {
		findPreviewResizeWhileFindOpen(t)
	})

	t.Run("Modal stays clear of the preview pane", func(t *testing.T) {
		findPreviewModalStaysClearOfPane(t)
	})
}

// findPreviewNavigateBeforeFind verifies that panel navigation re-renders the
// pane normally while find is closed
func findPreviewNavigateBeforeFind(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())

	count := sendFindArrow(t, m, tea.KeyDown)
	assert.Equal(t, 1, count, "navigation while find is closed should re-render the pane once")

	panel := m.getFocusedFilePanel()
	assert.Equal(t, "dir2", panel.GetFocusedItem().Name, "cursor should move to dir2")
	assert.Equal(t, filepath.Join(findPreviewDir(), "dir2"), m.fileModel.FilePreview.GetLocation(),
		"pane should preview the new cursor item")
	assert.False(t, m.fileModel.HasPreviewOverride(), "no preview override expected outside find")
}

// findPreviewOpenKeepsPane verifies that opening find claims the pane for the
// result under the initial cursor, without re-rendering it
func findPreviewOpenKeepsPane(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())

	count := openFindDirect(t, m)
	assert.Equal(t, 0, count, "opening find must not re-render the pane that already shows the cursor item")

	dir1 := filepath.Join(findPreviewDir(), "dir1")
	assert.True(t, m.findModal.IsOpen())
	results := m.findModal.GetResults()
	assert.Len(t, results, findPreviewAllResults, "an empty query should list the whole fixture tree")
	assert.Equal(t, dir1, m.findModal.GetCursorPath(), "the cursor should start on the first result")
	assert.True(t, m.fileModel.HasPreviewOverride(), "opening find should claim the pane for the cursor item")
	assert.Equal(t, dir1, m.fileModel.FilePreview.GetLocation())
}

// findPreviewArrowsMoveCursor verifies that arrow keys move the find cursor
// and the pane follows each new result
func findPreviewArrowsMoveCursor(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)

	var counts []int
	for _, key := range []rune{tea.KeyDown, tea.KeyDown, tea.KeyUp} {
		counts = append(counts, sendFindArrow(t, m, key))
	}
	assert.Equal(t, []int{1, 1, 1}, counts, "each arrow should re-render the pane for the new cursor item")

	file1 := filepath.Join(findPreviewDir(), "dir1", "file1.txt")
	assert.Equal(t, file1, m.findModal.GetCursorPath(), "the cursor should end on dir1/file1.txt")
	assert.Equal(t, file1, m.fileModel.FilePreview.GetLocation(), "the pane should follow the cursor")
	assert.True(t, m.fileModel.HasPreviewOverride())
}

// findPreviewTypingFiltersResults verifies that typed keys filter the results
// and the pane follows the result under the cursor
func findPreviewTypingFiltersResults(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)
	// Move the cursor the same way as the arrows subtest, so the pane
	// shows a file before typing starts
	sendFindArrow(t, m, tea.KeyDown)
	sendFindArrow(t, m, tea.KeyDown)
	sendFindArrow(t, m, tea.KeyUp)

	dir1 := filepath.Join(findPreviewDir(), "dir1")
	count1 := sendFindKey(t, m, "d")
	assert.Equal(t, 1, count1, "'d' should narrow the results and move the pane to the new cursor item")
	assert.Equal(t, dir1, m.findModal.GetCursorPath(), "the 'd' results should start on dir1")

	count2 := sendFindKey(t, m, "a")
	assert.Equal(t, 0, count2, "'da' should match nothing and release the pane back to the focused item")

	assert.Empty(t, m.findModal.GetResults(), "'da' should match nothing in the fixture")
	assert.False(t, m.fileModel.HasPreviewOverride(), "empty results should release the pane override")
	assert.Equal(t, dir1, m.fileModel.FilePreview.GetLocation())
}

// findPreviewTypingOneKeyAtATime verifies that each typed key re-renders the
// pane once, until the cursor item stops changing
func findPreviewTypingOneKeyAtATime(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)

	targetFile := filepath.Join(findPreviewDir(), "dir2", "nested", "target_file.txt")
	var counts []int
	for _, key := range []string{"t", "a", "r", "g", "e", "t"} {
		counts = append(counts, sendFindKey(t, m, key))
	}
	assert.Equal(t, []int{1, 1, 0, 0, 0, 0}, counts,
		"each query should re-render the pane once until the cursor item stops changing")

	assert.Equal(t, "target", m.findModal.GetTextInputValue())
	assert.Equal(t, targetFile, m.findModal.GetCursorPath())
	assert.Equal(t, targetFile, m.fileModel.FilePreview.GetLocation())
	assert.True(t, m.fileModel.HasPreviewOverride())
}

// findPreviewClearingQueryReleasesPane verifies that filtering to no results
// releases the pane back to the focused item, once
func findPreviewClearingQueryReleasesPane(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)
	for _, key := range []string{"t", "a", "r", "g", "e", "t"} {
		sendFindKey(t, m, key)
	}

	dir1 := filepath.Join(findPreviewDir(), "dir1")
	count1 := sendFindKey(t, m, "z")
	assert.Equal(t, 1, count1, "'targetz' should match nothing and return the pane to the focused item")
	assert.False(t, m.fileModel.HasPreviewOverride())
	assert.Equal(t, dir1, m.fileModel.FilePreview.GetLocation())

	count2 := sendFindKey(t, m, "x")
	assert.Equal(t, 0, count2, "typing with no results should not re-render the pane")
	assert.False(t, m.fileModel.HasPreviewOverride())
	assert.Equal(t, dir1, m.fileModel.FilePreview.GetLocation())
}

// findPreviewCancelTypingReleasesPane verifies that cancelling a non-empty
// query closes find and returns the pane to the focused item, once
func findPreviewCancelTypingReleasesPane(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)
	for _, key := range []string{"t", "a", "r", "g", "e", "t"} {
		sendFindKey(t, m, key)
	}

	dir1 := filepath.Join(findPreviewDir(), "dir1")
	count := sendFindKey(t, m, common.Hotkeys.CancelTyping[1])
	assert.Equal(t, 1, count, "closing find should return the pane to the focused item")
	assert.False(t, m.findModal.IsOpen())
	assert.False(t, m.fileModel.HasPreviewOverride())
	assert.Equal(t, dir1, m.fileModel.FilePreview.GetLocation())
}

// findPreviewCancelTypingRightAfterOpening verifies that cancelling before
// any cursor move does not re-render the pane
func findPreviewCancelTypingRightAfterOpening(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)

	dir1 := filepath.Join(findPreviewDir(), "dir1")
	count := sendFindKey(t, m, common.Hotkeys.CancelTyping[1])
	assert.Equal(t, 0, count, "closing find without a cursor move should not re-render the pane")
	assert.False(t, m.findModal.IsOpen())
	assert.False(t, m.fileModel.HasPreviewOverride())
	assert.Equal(t, dir1, m.fileModel.FilePreview.GetLocation())
}

// findPreviewConfirmNavigatesToFile verifies that confirming a file closes
// find, navigates to its directory, and keeps the pane on that file
func findPreviewConfirmNavigatesToFile(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)
	for _, key := range []string{"t", "a", "r", "g", "e", "t"} {
		sendFindKey(t, m, key)
	}

	nestedDir := filepath.Join(findPreviewDir(), "dir2", "nested")
	targetFile := filepath.Join(nestedDir, "target_file.txt")
	count := sendFindKey(t, m, common.Hotkeys.ConfirmTyping[0])
	assert.Equal(t, 0, count, "confirming should navigate without an extra render, the pane already shows the file")

	panel := m.getFocusedFilePanel()
	assert.False(t, m.findModal.IsOpen())
	assert.False(t, m.fileModel.HasPreviewOverride())
	assert.Equal(t, nestedDir, panel.Location)
	assert.Equal(t, "target_file.txt", panel.GetFocusedItem().Name)
	assert.Empty(t, panel.TargetFile)
	assert.Equal(t, targetFile, m.fileModel.FilePreview.GetLocation())
	assert.Contains(t, ansi.Strip(m.fileModel.FilePreview.GetContent()), "This is sample")
}

// findPreviewConfirmNavigatesToDirectory verifies that confirming a directory
// closes find, navigates into it, and re-renders the pane for the focused item
func findPreviewConfirmNavigatesToDirectory(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)

	count1 := sendFindKey(t, m, "dir1")
	assert.Equal(t, 0, count1, "typing 'dir1' in one key should not re-render the pane that already shows dir1")

	count2 := sendFindKey(t, m, common.Hotkeys.ConfirmTyping[0])
	assert.Equal(t, 1, count2, "navigating into dir1 should re-render the pane for the focused file")

	dir1 := filepath.Join(findPreviewDir(), "dir1")
	panel := m.getFocusedFilePanel()
	assert.False(t, m.findModal.IsOpen())
	assert.False(t, m.fileModel.HasPreviewOverride())
	assert.Equal(t, dir1, panel.Location)
	assert.Equal(t, "file1.txt", panel.GetFocusedItem().Name)

	file1 := filepath.Join(dir1, "file1.txt")
	assert.Equal(t, file1, m.fileModel.FilePreview.GetLocation())
	assert.Contains(t, ansi.Strip(m.fileModel.FilePreview.GetContent()), "This is sample")
}

// findPreviewTogglePaneWhileFindOpen verifies that toggling the pane while
// find is open reuses the existing content without re-rendering
func findPreviewTogglePaneWhileFindOpen(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)

	dir1 := filepath.Join(findPreviewDir(), "dir1")
	// 'f' is alphanumeric and would be typed into the find input, so the
	// toggle is driven directly
	closeCmd := m.fileModel.ToggleFilePreviewPanel()
	assert.Nil(t, closeCmd, "closing the pane should keep the content at the current size")
	assert.False(t, m.fileModel.FilePreview.IsOpen())
	assert.Equal(t, findPreviewPaneWidth, m.fileModel.ExpectedPreviewWidth)
	assert.Equal(t, dir1, m.fileModel.FilePreview.GetLocation())
	assert.Nil(t, m.updateFilePreview(), "the override still points at the focused item, so no render is needed")

	reopenCmd := m.fileModel.ToggleFilePreviewPanel()
	assert.Nil(t, reopenCmd, "reopening the pane should reuse the content at the current size")
	assert.True(t, m.fileModel.FilePreview.IsOpen())
	assert.Equal(t, findPreviewPaneWidth, m.fileModel.FilePreview.GetContentWidth())
	assert.Equal(t, DefaultTestModelHeight, m.fileModel.FilePreview.GetContentHeight())
	assert.Equal(t, dir1, m.fileModel.FilePreview.GetLocation())
	assert.True(t, m.fileModel.HasPreviewOverride())
	assert.Nil(t, m.updateFilePreview(), "the pane already shows the cursor item after reopening")
}

// findPreviewResizeWhileFindOpen verifies that resizing the window re-renders
// the pane once at the new dimensions and keeps the modal sized to fit
func findPreviewResizeWhileFindOpen(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)

	count := runSettled(t, m, TeaUpdate(m, tea.WindowSizeMsg{
		Width: findResizeWidth, Height: findResizeHeight,
	}))
	assert.Equal(t, 1, count, "resizing should re-render the pane once at the new dimensions")

	dir1 := filepath.Join(findPreviewDir(), "dir1")
	pane := m.fileModel.FilePreview
	assert.Equal(t, findResizePreviewWidth, pane.GetContentWidth())
	assert.Equal(t, findResizePreviewHeight, pane.GetContentHeight())
	assert.Equal(t, dir1, pane.GetLocation(), "resizing should not change which item the pane shows")
	assert.Equal(t, findResizePreviewWidth, m.fileModel.ExpectedPreviewWidth)
	assert.Equal(t, findResizeModalWidth, m.findModal.GetWidth())
	assert.Equal(t, findResizeModalMaxHeight, m.findModal.GetMaxHeight())
}

// findPreviewModalStaysClearOfPane verifies that the open modal frame differs
// from the closed one only inside the modal band of the left region, never in
// the preview pane region, and that the headline stays out of the pane
func findPreviewModalStaysClearOfPane(t *testing.T) {
	t.Helper()
	m := primeFindPreviewModel(t, findPreviewDir())
	openFindDirect(t, m)

	openFrame := frameRows(t, m)
	sendFindKey(t, m, common.Hotkeys.CancelTyping[1])
	closedFrame := frameRows(t, m)

	assert.Len(t, openFrame, len(closedFrame), "the frames should have the same row count")
	assert.GreaterOrEqual(t, len(openFrame), findPreviewMinRows, "the frame should span most of the test window")

	for i, openRow := range openFrame {
		closedRow := closedFrame[i]
		assert.Equal(t, ansi.Strip(ansi.Cut(closedRow, findPreviewPaneLeft, findPreviewBigCut)),
			ansi.Strip(ansi.Cut(openRow, findPreviewPaneLeft, findPreviewBigCut)),
			"row %d: the find modal should not touch the preview pane region", i)
	}

	var leftOutsideEqual, bandDiffers, headlineRows int
	for i := range openFrame {
		openLeft := ansi.Strip(ansi.Cut(openFrame[i], 0, findPreviewPaneLeft))
		closedLeft := ansi.Strip(ansi.Cut(closedFrame[i], 0, findPreviewPaneLeft))
		if i < findPreviewModalTop || i > findPreviewModalBottom {
			assert.Equal(t, closedLeft, openLeft,
				"row %d: rows outside the modal band should be identical in the left region", i)
			leftOutsideEqual++
		} else if openLeft != closedLeft {
			bandDiffers++
		}
		if strings.Contains(openLeft, findPreviewHeadline) {
			headlineRows++
		}
		paneLeft := ansi.Strip(ansi.Cut(openFrame[i], findPreviewPaneLeft, findPreviewBigCut))
		if strings.Contains(paneLeft, findPreviewHeadline) {
			t.Fatalf("row %d: the find modal headline leaked into the preview pane region", i)
		}
	}
	assert.Positive(t, leftOutsideEqual, "there should be rows outside the modal band")
	assert.Positive(t, bandDiffers, "the find modal should occupy rows of the left region")
	assert.GreaterOrEqual(t, headlineRows, 1, "the find headline should be visible in the open frame")
}
