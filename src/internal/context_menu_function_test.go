package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	variable "github.com/yorukot/superfile/src/config"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/mouse"
	"github.com/yorukot/superfile/src/internal/ui/contextmenu"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/spferror"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// backgroundMenuLabels is what the background context menu offers. The same
// labels are asserted in more than one place, so they are written once.
func backgroundMenuLabels() []string {
	return []string{"Paste", "New file", "New folder"}
}

// rightClickAt builds a right button press at the given coordinates.
func rightClickAt(x, y int) tea.MouseClickMsg {
	return tea.MouseClickMsg{X: x, Y: y, Button: tea.MouseRight}
}

// menuCell finds the coordinates of a context menu entry, so that tests drive
// clicks through the same hit-test lookup a real click uses.
func menuCell(t *testing.T, m *model, itemIndex int) (int, int) {
	t.Helper()

	// The menu publishes its regions during render.
	m.viewContent()
	for y := range m.fullHeight {
		for x := range m.fullWidth {
			target := m.mouseTargetAt(x, y)
			if target.Kind == mouse.TargetContextMenuItem && target.ItemIndex == itemIndex {
				return x, y
			}
		}
	}
	require.FailNow(t, "menu entry not on screen", "index %d", itemIndex)
	return 0, 0
}

// TestRightClickOpensMenuOnFileItem checks the basic open path and that the menu
// is anchored where the pointer was.
func TestRightClickOpensMenuOnFileItem(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 2)
	m.handleMouseMsg(rightClickAt(x, y))

	require.True(t, m.contextMenu.IsOpen())
	assert.Equal(t, x, m.contextMenu.OriginX())
	assert.Equal(t, y, m.contextMenu.OriginY())

	action, ok := m.contextMenu.SelectedAction()
	require.True(t, ok)
	assert.Equal(t, contextmenu.ActionOpen, action, "the first entry is Open")
}

// TestRightClickSelectsTheClickedItem checks that the menu acts on the entry that
// was clicked, even when the cursor was somewhere else.
func TestRightClickSelectsTheClickedItem(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)
	panel := m.getFocusedFilePanel()
	require.Equal(t, 0, panel.GetCursor())

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 3)
	m.handleMouseMsg(rightClickAt(x, y))

	require.True(t, m.contextMenu.IsOpen())
	assert.Equal(t, 3, panel.GetCursor())
	clicked := panel.GetElementAtIdx(3).Location
	assert.True(t, panel.CheckSelected(clicked), "the clicked entry should be selected")
	assert.Equal(t, uint(1), panel.SelectedCount())
	assert.Equal(t, filepanel.SelectMode, panel.PanelMode,
		"the panel has to be in SelectMode for the selection to be used")
}

// TestRightClickInsideSelectionKeepsIt checks that right clicking one of several
// selected entries does not throw the rest of the selection away.
func TestRightClickInsideSelectionKeepsIt(t *testing.T) {
	dir := dirWithFiles(t, 6)
	m := defaultTestModel(dir)
	panel := m.getFocusedFilePanel()

	// Select entries 1 and 2 with the keyboard's own model, then right click one
	// of them.
	for _, idx := range []int{1, 2} {
		x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, idx)
		m.handleMouseMsg(ctrlClickAt(x, y))
	}
	require.Equal(t, uint(2), panel.SelectedCount())

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 2)
	m.handleMouseMsg(rightClickAt(x, y))

	require.True(t, m.contextMenu.IsOpen())
	assert.Equal(t, uint(2), panel.SelectedCount(),
		"right clicking inside the selection must keep it")
}

// TestRightClickOutsideSelectionNarrowsIt checks that right clicking an unselected
// entry reduces the selection to just that entry.
func TestRightClickOutsideSelectionNarrowsIt(t *testing.T) {
	dir := dirWithFiles(t, 6)
	m := defaultTestModel(dir)
	panel := m.getFocusedFilePanel()

	for _, idx := range []int{1, 2} {
		x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, idx)
		m.handleMouseMsg(ctrlClickAt(x, y))
	}
	require.Equal(t, uint(2), panel.SelectedCount())

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 4)
	m.handleMouseMsg(rightClickAt(x, y))

	require.True(t, m.contextMenu.IsOpen())
	assert.Equal(t, uint(1), panel.SelectedCount())
	assert.True(t, panel.CheckSelected(panel.GetElementAtIdx(4).Location))
}

// TestMenuRegionsWinOverPanelRegions checks that a cell covered by the menu
// resolves to the menu rather than to the file row drawn underneath it. Without
// this the menu would be unclickable wherever it overlapped a row.
func TestMenuRegionsWinOverPanelRegions(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)

	// Open the menu low in the panel so it overlaps file rows.
	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
	m.handleMouseMsg(rightClickAt(x, y+6))
	require.True(t, m.contextMenu.IsOpen())

	m.viewContent()

	overlapped := false
	for row := range m.fullHeight {
		for col := range m.fullWidth {
			target := m.mouseTargetAt(col, row)
			if target.Kind == mouse.TargetContextMenuItem {
				overlapped = true
				continue
			}
			// Nothing inside the menu's own bounds may resolve to a panel row.
			if col >= m.contextMenu.OriginX() &&
				col < m.contextMenu.OriginX()+m.contextMenu.Width() &&
				row > m.contextMenu.OriginY() &&
				row < m.contextMenu.OriginY()+m.contextMenu.Height()-1 {
				assert.Equal(t, mouse.TargetContextMenuItem, target.Kind,
					"cell inside the menu resolved to something else")
			}
		}
	}
	assert.True(t, overlapped, "the menu should have published some cells")
}

// TestMenuClickActivatesEntry checks that clicking an entry runs that entry's
// action, and that the menu closes afterwards.
func TestMenuClickActivatesEntry(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)

	var copied string
	m.clipboardWriter = func(text string) error {
		copied = text
		return nil
	}

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 1)
	m.handleMouseMsg(rightClickAt(x, y))
	require.True(t, m.contextMenu.IsOpen())

	// The entry order is Open, Copy, Cut, Copy path, Rename, Delete. Click
	// "Copy path", which has a visible effect on the clipboard.
	const copyPathIndex = 3
	cx, cy := menuCell(t, m, copyPathIndex)
	action, ok := m.contextMenu.ChosenAction(copyPathIndex)
	require.True(t, ok)
	require.Equal(t, contextmenu.ActionCopyPath, action)

	m.handleMouseMsg(leftClickAt(cx, cy))

	assert.False(t, m.contextMenu.IsOpen(), "choosing an entry should close the menu")
	assert.Equal(t, m.getFocusedFilePanel().GetElementAtIdx(1).Location, copied,
		"copy path should have written the clicked entry's path")
}

// TestClickAwayFromMenuDismissesWithoutSideEffect checks that the click which
// closes the menu does not also act on the widget underneath it.
func TestClickAwayFromMenuDismissesWithoutSideEffect(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)
	panel := m.getFocusedFilePanel()

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
	m.handleMouseMsg(rightClickAt(x, y+6))
	require.True(t, m.contextMenu.IsOpen())

	// Click a file row that is not covered by the menu.
	targetX, targetY := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
	m.handleMouseMsg(leftClickAt(targetX, targetY))

	assert.False(t, m.contextMenu.IsOpen())
	assert.Equal(t, 0, panel.GetCursor(),
		"the dismissing click must not also move the panel cursor")
}

// TestKeyboardNavigatesAndActivatesMenu checks the menu is usable without a mouse.
func TestKeyboardNavigatesAndActivatesMenu(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 2)
	m.handleMouseMsg(rightClickAt(x, y))
	require.True(t, m.contextMenu.IsOpen())

	m.contextMenuKey("j")
	action, _ := m.contextMenu.SelectedAction()
	assert.Equal(t, contextmenu.ActionCopy, action)

	m.contextMenuKey("k")
	action, _ = m.contextMenu.SelectedAction()
	assert.Equal(t, contextmenu.ActionOpen, action)

	m.contextMenuKey("esc")
	assert.False(t, m.contextMenu.IsOpen(), "escape should dismiss the menu")
}

// TestRightClickOnSidebarDirectoryDoesNotNavigate checks that opening a menu on a
// sidebar entry only moves the cursor. Navigating behind an open menu would be
// surprising, and the Open entry exists for that.
func TestRightClickOnSidebarDirectoryDoesNotNavigate(t *testing.T) {
	dir := dirWithFiles(t, 3)
	m := defaultTestModel(dir)
	m.sidebarModel.UpdateDirectories()
	m.viewContent()
	locationBefore := m.getFocusedFilePanel().Location

	x, y := mustFindCell(t, m, mouse.TargetSidebarDirectory, 1)
	m.handleMouseMsg(rightClickAt(x, y))

	require.True(t, m.contextMenu.IsOpen())
	assert.Equal(t, locationBefore, m.getFocusedFilePanel().Location)
	assert.Equal(t, sidebarFocus, m.focusPanel)
}

// TestRightClickInFooterOffersNothing checks that the file wide actions are not
// offered where no directory is in view.
func TestRightClickInFooterOffersNothing(t *testing.T) {
	dir := dirWithFiles(t, 3)
	m := defaultTestModelWithFooterAndFilePreview(dir)
	m.viewContent()

	m.handleMouseMsg(rightClickAt(m.fullWidth-1, m.mainPanelHeight+1))
	assert.False(t, m.contextMenu.IsOpen(), "the footer should not open a menu")
}

// TestRightClickOnEmptyPanelOffersCreateAndPaste checks the background menu.
func TestRightClickOnEmptyPanelOffersCreateAndPaste(t *testing.T) {
	dir := dirWithFiles(t, 2)
	m := defaultTestModel(dir)
	m.viewContent()

	// A row below the last entry, still inside the panel area.
	m.handleMouseMsg(rightClickAt(m.fullWidth/2, m.mainPanelHeight-1))
	require.True(t, m.contextMenu.IsOpen())

	labels := menuLabels(m)
	assert.Equal(t, backgroundMenuLabels(), labels)
}

// TestRightClickOnBackgroundFocusesThatPanel checks that the panel the pointer is
// over becomes the focused one before its background menu is offered.
//
// Paste, new file, and new folder all act on the focused panel, so leaving focus
// where the keyboard had it would put the new entry in a directory the user never
// pointed at.
func TestRightClickOnBackgroundFocusesThatPanel(t *testing.T) {
	dir := dirWithFiles(t, 3)
	other := t.TempDir()
	m := defaultTestModel(dir)
	_, err := m.fileModel.CreateNewFilePanel(other)
	require.NoError(t, err)
	require.Equal(t, 2, m.fileModel.PanelCount())
	m.viewContent()

	// Focus is deliberately put on the first panel, which is not the one the
	// pointer will be over.
	m.focusFilePanelOnMouse(0)
	require.Equal(t, 0, m.fileModel.FocusedPanelIndex,
		"focus has to start on a panel other than the second for this to mean anything")

	// The area of the second panel, below the entries.
	x := m.filePanelAreaX() + m.fileModel.PanelOriginX(1) + 2
	y := m.mainPanelHeight - 1
	require.Equal(t, mouse.TargetUnknown, m.mouseTargetAt(x, y).Kind,
		"this coordinate should be the second panel's background")

	m.handleMouseMsg(rightClickAt(x, y))

	require.True(t, m.contextMenu.IsOpen(), "the background menu should be open")
	assert.Equal(t, backgroundMenuLabels(), menuLabels(m))
	assert.Equal(t, 1, m.fileModel.FocusedPanelIndex,
		"the panel the pointer was over should have the focus")
	assert.Equal(t, other, m.getFocusedFilePanel().Location,
		"the menu should act on the panel the pointer was over, not on the old focus")
}

// TestRightClickWithoutApplicableTargetClosesMenu checks a right click with
// nothing to offer does not leave a stale menu on screen.
func TestRightClickWithoutApplicableTargetClosesMenu(t *testing.T) {
	dir := dirWithFiles(t, 3)
	m := defaultTestModel(dir)

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
	m.handleMouseMsg(rightClickAt(x, y))
	require.True(t, m.contextMenu.IsOpen())

	// Right click the sidebar background, which offers nothing.
	m.viewContent()
	m.handleMouseMsg(rightClickAt(m.fullWidth-1, m.mainPanelHeight+1))
	assert.False(t, m.contextMenu.IsOpen())
}

// TestHoverHighlightsMenuEntry checks the pointer moves the menu's cursor.
func TestHoverHighlightsMenuEntry(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
	m.handleMouseMsg(rightClickAt(x, y))
	require.True(t, m.contextMenu.IsOpen())
	m.viewContent()

	hx, hy := menuCell(t, m, 2)
	m.handleMouseMsg(tea.MouseMotionMsg{X: hx, Y: hy})

	action, _ := m.contextMenu.SelectedAction()
	assert.Equal(t, contextmenu.ActionCut, action, "hover should highlight the entry")
}

// TestWheelScrollDoesNotHighlightMenuEntry checks that continued scrolling, which
// arrives as motion with the wheel held, is not mistaken for the pointer moving.
func TestWheelScrollDoesNotHighlightMenuEntry(t *testing.T) {
	dir := dirWithFiles(t, 40)
	m := defaultTestModel(dir)

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
	m.handleMouseMsg(rightClickAt(x, y))
	require.True(t, m.contextMenu.IsOpen())
	m.viewContent()

	before, _ := m.contextMenu.SelectedAction()
	m.handleMouseMsg(tea.MouseMotionMsg{X: x, Y: y, Button: tea.MouseWheelDown})
	after, _ := m.contextMenu.SelectedAction()

	assert.Equal(t, before, after, "a wheel scroll must not move the menu cursor")
}

// TestMenuIsDrawnOverThePanels checks the overlay actually contains the entry
// labels, rather than the menu being computed but never shown.
func TestMenuIsDrawnOverThePanels(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
	m.handleMouseMsg(rightClickAt(x, y))
	require.True(t, m.contextMenu.IsOpen())

	rendered := m.viewContent()
	for _, label := range menuLabels(m) {
		assert.Contains(t, stripANSI(rendered), label,
			"entry %q should be visible in the rendered output", label)
	}
}

// menuLabels returns the labels currently offered by the menu.
func menuLabels(m *model) []string {
	items := m.contextMenu.Items()
	labels := make([]string, 0, len(items))
	for _, item := range items {
		labels = append(labels, item.Label)
	}
	return labels
}

// stripANSI removes escape sequences so that rendered text can be searched.
func stripANSI(s string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range s {
		switch {
		case r == 0x1b:
			inEscape = true
		case inEscape && (r == 'm' || r == 'K' || r == 'H'):
			inEscape = false
		case !inEscape:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// TestNewFolderActionCreatesDirectory checks the trailing separator shortcut for
// creating a directory, which is easy to get wrong because createItem decides
// between a file and a directory purely on that suffix.
func TestNewFolderActionCreatesDirectory(t *testing.T) {
	dir := dirWithFiles(t, 2)
	m := defaultTestModel(dir)

	m.viewContent()
	m.handleMouseMsg(rightClickAt(m.fullWidth/2, m.mainPanelHeight-1))
	require.True(t, m.contextMenu.IsOpen())
	labels := menuLabels(m)
	require.Contains(t, labels, "New folder")

	// Choose "New folder" from the menu.
	newFolderIndex := slices.Index(labels, "New folder")
	cx, cy := menuCell(t, m, newFolderIndex)
	m.handleMouseMsg(leftClickAt(cx, cy))

	assert.False(t, m.contextMenu.IsOpen())
	require.True(t, m.typingModal.open, "creating a folder should open the typing modal")

	// The user types a plain name; the modal already knows a directory is wanted.
	assert.Empty(t, m.typingModal.textInput.Value())
	assert.True(t, m.typingModal.directory)

	m.typingModal.textInput.SetValue("created-by-menu")
	cmd := m.getCreateCmd()
	require.NotNil(t, cmd)
	cmd()
	created := filepath.Join(dir, "created-by-menu")

	assert.Eventually(t, func() bool {
		info, err := os.Lstat(created)
		return err == nil && info.IsDir()
	}, DefaultTestTimeout, DefaultTestTick,
		"the menu's New folder should create a directory, got %q", created)
}

// TestNewFileActionCreatesFile checks the plain file variant, so that the
// directory shortcut above cannot regress both at once.
func TestNewFileActionCreatesFile(t *testing.T) {
	dir := dirWithFiles(t, 2)
	m := defaultTestModel(dir)

	m.viewContent()
	m.handleMouseMsg(rightClickAt(m.fullWidth/2, m.mainPanelHeight-1))
	require.True(t, m.contextMenu.IsOpen())

	labels := menuLabels(m)
	newFileIndex := slices.Index(labels, "New file")
	cx, cy := menuCell(t, m, newFileIndex)
	m.handleMouseMsg(leftClickAt(cx, cy))

	require.True(t, m.typingModal.open)
	assert.Empty(t, m.typingModal.textInput.Value(),
		"creating a file should not pretype a separator")

	m.typingModal.textInput.SetValue("created-by-menu.txt")
	cmd := m.getCreateCmd()
	require.NotNil(t, cmd)
	cmd()
	created := filepath.Join(dir, "created-by-menu.txt")

	assert.Eventually(t, func() bool {
		info, err := os.Lstat(created)
		return err == nil && !info.IsDir()
	}, DefaultTestTimeout, DefaultTestTick,
		"the menu's New file should create a regular file, got %q", created)
}

// TestDeleteActionDeletesClickedItem checks the destructive entry is wired to the
// existing delete flow, including that it follows the confirmation modal.
func TestDeleteActionDeletesClickedItem(t *testing.T) {
	dir := dirWithFiles(t, 3)
	m := defaultTestModel(dir)
	panel := m.getFocusedFilePanel()
	target := panel.GetElementAtIdx(1).Location
	toDelete := filepath.Join(dir, filepath.Base(target))

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 1)
	m.handleMouseMsg(rightClickAt(x, y))
	require.True(t, m.contextMenu.IsOpen())

	labels := menuLabels(m)
	deleteIndex := slices.Index(labels, "Delete")
	cx, cy := menuCell(t, m, deleteIndex)
	m.handleMouseMsg(leftClickAt(cx, cy))

	assert.False(t, m.contextMenu.IsOpen())
	assert.True(t, panel.CheckSelected(target),
		"the entry that was right clicked is the one to delete")
	assert.FileExists(t, toDelete, "nothing may be removed before it is confirmed")
}

// TestPasteActionPastes check the background Paste entry writes through the
// existing paste flow.
func TestPasteActionPastes(t *testing.T) {
	source := dirWithFiles(t, 2)
	dest := t.TempDir()
	m := defaultTestModel(dest)

	// The clipboard belongs to the model, so put the file on this model's copy.
	pasted := filepath.Join(source, "filea.txt")
	m.clipboard.Reset(false)
	m.clipboard.Add(pasted)

	m.viewContent()
	m.handleMouseMsg(rightClickAt(m.fullWidth/2, m.mainPanelHeight-1))
	require.True(t, m.contextMenu.IsOpen())

	labels := menuLabels(m)
	pasteIndex := slices.Index(labels, "Paste")
	require.GreaterOrEqual(t, pasteIndex, 0, "the background menu should offer Paste")

	cx, cy := menuCell(t, m, pasteIndex)
	cmd := m.handleMouseMsg(leftClickAt(cx, cy))

	assert.False(t, m.contextMenu.IsOpen())
	if cmd != nil {
		cmd()
	}
	assert.Eventually(t, func() bool {
		return fileExists(filepath.Join(dest, "filea.txt"))
	}, DefaultTestTimeout, DefaultTestTick,
		"pasting from the menu should copy the clipboard into the panel")
}

// fileExists reports whether a regular file is present at the given path.
func fileExists(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && !info.IsDir()
}

// TestMenuFooterShowsTheWholeCount checks the "index/total" footer is not cut off
// by the box border. Lipgloss gives the bottom border fewer cells than the box is
// wide, so a footer sized for the full width loses its last characters.
func TestMenuFooterShowsTheWholeCount(t *testing.T) {
	dir := dirWithFiles(t, 5)
	m := defaultTestModel(dir)
	m.viewContent()
	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 1)
	m.handleMouseMsg(rightClickAt(x, y))
	require.True(t, m.contextMenu.IsOpen())

	total := len(m.contextMenu.Items())
	for cursor := range total {
		m.contextMenu.HighlightItem(cursor)

		want := fmt.Sprintf("%d/%d", cursor+1, total)
		rendered := stripANSI(m.viewContent())
		assert.Contains(t, rendered, want,
			"the footer should show %q when the cursor is on entry %d", want, cursor)
	}
}

// useTemporaryPinnedFile points the pinned file at a temporary one for the rest of
// the test, so that pinning in a test cannot rewrite the user's real pinned list.
//
// The sidebar captures the path when it is built, so this has to run before the
// model under test is created.
func useTemporaryPinnedFile(t *testing.T) {
	t.Helper()
	previous := variable.PinnedFile
	//nolint:reassign // Test fixture, restored by the cleanup below.
	variable.PinnedFile = filepath.Join(t.TempDir(), "pinned.json")
	t.Cleanup(func() {
		//nolint:reassign // Restoring the pinned file after the test.
		variable.PinnedFile = previous
	})
}

// TestMouseIsIgnoredWhileAModalIsOpen checks that a modal captures the pointer as
// well as the keyboard.
//
// Keys are already routed to whichever modal is open, so a click reaching the
// panels behind one would move the cursor, focus, and selection out from under
// the modal, and a click on a menu entry above one would dispatch its action.
func TestMouseIsIgnoredWhileAModalIsOpen(t *testing.T) {
	modals := []struct {
		name    string
		prepare func(t *testing.T)
		open    func(m *model)
		close   func(m *model)
	}{
		{"typing", nil, func(m *model) {
			m.typingModal.open = true
		}, func(m *model) { m.typingModal.open = false }},
		{"prompt", nil, func(m *model) {
			TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
		}, func(m *model) {
			TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.CancelTyping[0]))
		}},
		{"notify", nil, func(m *model) {
			m.notifyModel.Open()
		}, func(m *model) { m.notifyModel.Close() }},
		{"zoxide", nil, func(m *model) {
			m.zoxideModal.Open()
		}, func(m *model) { m.zoxideModal.Close() }},
		{"sort", nil, func(m *model) {
			m.sortModal.Open(m.getFocusedFilePanel().SortKind)
		}, func(m *model) { m.sortModal.Close() }},
		{"help menu", nil, func(m *model) {
			m.helpMenu.Open()
		}, func(m *model) { m.helpMenu.Close() }},
		{"spf error", nil, func(m *model) {
			m.spfError = spferror.New(true, "Error", "boom",
				&spferror.FileListErrorState{})
		}, func(m *model) { m.spfError = spferror.Model{} }},
		{"file rename", nil, func(m *model) {
			m.fileModel.Renaming = true
		}, func(m *model) { m.fileModel.Renaming = false }},
		{"sidebar rename", useTemporaryPinnedFile, func(m *model) {
			// Only a pinned directory can be renamed, so one is pinned and selected.
			dir := m.getFocusedFilePanel().Location
			require.NoError(t, m.sidebarModel.TogglePinnedDirectory(dir))
			m.sidebarModel.UpdateDirectories()
			// Pinned entries come after the well known ones and behind a divider, so
			// the rows are walked until one of them can be renamed.
			for i := 0; i < 30 && !m.sidebarModel.IsRenaming(); i++ {
				m.sidebarModel.SetCursor(i)
				m.sidebarModel.PinnedItemRename()
			}
			require.True(t, m.sidebarModel.IsRenaming(), "the pinned entry should be renaming")
		}, func(m *model) { m.sidebarModel.CancelSidebarRename() }},
	}

	for _, modal := range modals {
		t.Run(modal.name, func(t *testing.T) {
			if modal.prepare != nil {
				modal.prepare(t)
			}
			m := defaultTestModel(dirWithFiles(t, 5))
			fileRows := m.getFocusedFilePanel().ElemCount()
			cursor := m.getFocusedFilePanel().GetCursor()

			modal.open(m)
			t.Cleanup(func() { modal.close(m) })
			require.True(t, m.modalBlocksMouse(), "%s should capture the mouse", modal.name)

			x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, fileRows-1)
			// The rendered screen covers cursor, focus, and selection at once, so
			// comparing it proves that nothing underneath the modal moved.
			before := m.viewContent()

			m.handleMouseMsg(leftClickAt(x, y))
			require.Empty(t, m.contextMenu.Items(), "a click must not open a menu over a modal")
			require.Equal(t, before, m.viewContent(), "a click must not change the screen")

			m.handleMouseMsg(rightClickAt(x, y))
			require.Empty(t, m.contextMenu.Items(),
				"a right click must not open a menu over a modal either")
			require.Equal(t, before, m.viewContent(), "a right click must not change the screen")
			require.Equal(t, cursor, m.getFocusedFilePanel().GetCursor(),
				"the panel cursor must not move underneath the modal")
		})
	}
}

// TestAModalDismissesAMenuOpenedBeforeIt checks the other ordering: the menu was
// already open when the modal appeared.
//
// The menu renders above most modals, so a click on it would otherwise dispatch
// an action over the modal. Dismissing it means the click does nothing at all.
func TestAModalDismissesAMenuOpenedBeforeIt(t *testing.T) {
	m := defaultTestModel(dirWithFiles(t, 5))
	fileRows := m.getFocusedFilePanel().ElemCount()

	x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, fileRows-1)
	m.handleMouseMsg(rightClickAt(x, y))
	require.NotEmpty(t, m.contextMenu.Items(), "the right click should have opened a menu")

	m.notifyModel.Open()
	t.Cleanup(m.notifyModel.Close)

	m.handleMouseMsg(leftClickAt(x, y))
	require.Empty(t, m.contextMenu.Items(), "a modal should have dismissed the menu")
}
