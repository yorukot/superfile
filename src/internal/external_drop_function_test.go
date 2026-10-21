package internal

import (
	"os"
	"path/filepath"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/mouse"
	"github.com/yorukot/superfile/src/internal/ui/spferror"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// TestParseDroppedPaths covers every shape a file manager is known to hand a
// terminal, and the prose that has to be left alone.
//
// A drop arrives as pasted text, so the parser is the whole of the safety: a
// paste that is not a list of files must never become an import.
func TestParseDroppedPaths(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "plain.txt")
	spaced := filepath.Join(dir, "my file.txt")
	other := filepath.Join(dir, "other.txt")
	missing := filepath.Join(dir, "missing.txt")
	for _, path := range []string{plain, spaced, other} {
		require.NoError(t, os.WriteFile(path, []byte("data"), 0o600))
	}

	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{"bare path", plain, []string{plain}},
		{"file uri", "file://" + plain, []string{plain}},
		{"file uri with localhost", "file://localhost" + plain, []string{plain}},
		{"file uri percent encoded", "file://" + filepath.ToSlash(filepath.Join(dir, "my%20file.txt")),
			[]string{spaced}},
		{"path with spaces", spaced, []string{spaced}},
		{"shell escaped path", filepath.ToSlash(dir) + `/my\ file.txt`, []string{spaced}},
		{"double quoted path", `"` + spaced + `"`, []string{spaced}},
		{"single quoted path", `'` + spaced + `'`, []string{spaced}},
		{"two paths on one line", plain + " " + other, []string{plain, other}},
		{"uri list", "file://" + plain + "\nfile://" + other, []string{plain, other}},
		{"uri list with crlf", "file://" + plain + "\r\nfile://" + other + "\r\n",
			[]string{plain, other}},
		{"uri list with comments", "# dragged from a file manager\r\nfile://" + plain + "\r\n",
			[]string{plain}},
		{"blank lines are skipped", "\n\n" + plain + "\n\n", []string{plain}},
		{"duplicate paths collapse", plain + "\n" + plain, []string{plain}},

		{"empty", "", nil},
		{"whitespace only", " \t\n ", nil},
		{"comment only", "# nothing here", nil},
		{"prose", "just some words", nil},
		{"prose with a path", "see " + plain + " for details", nil},
		{"relative path", "plain.txt", nil},
		{"nonexistent path", missing, nil},
		{"one missing among many", plain + "\n" + missing, nil},
		{"one prose line among many", plain + "\nnot a path", nil},
		{"network share uri", "file://server/share/file.txt", nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, parseDroppedPaths(tc.content))
		})
	}
}

// TestPasteMsgImportsDroppedFile is the main path: a paste naming a file
// elsewhere on disk copies it into the directory the focused panel shows.
//
// It insists the source survives, since a drop from another window is an import
// and not a move, and that the clipboard stays empty, since a drop is not the
// item the user asked to paste.
func TestPasteMsgImportsDroppedFile(t *testing.T) {
	dest := t.TempDir()
	source := t.TempDir()
	dropped := filepath.Join(source, "dropped.txt")
	require.NoError(t, os.WriteFile(dropped, []byte("data"), 0o600))

	m := defaultTestModel(dest)
	require.Equal(t, dest, m.getFocusedFilePanel().Location)

	cmd := m.handlePasteMsg(tea.PasteMsg{Content: "file://" + dropped})
	require.NotNil(t, cmd, "a dropped file should start an import")
	runCmd(cmd)

	assert.Eventually(t, func() bool {
		return fileExists(filepath.Join(dest, "dropped.txt"))
	}, DefaultTestTimeout, DefaultTestTick,
		"the dropped file should be copied into the focused panel")
	assert.FileExists(t, dropped, "a drop copies, so the source stays where it was")
	assert.Zero(t, m.clipboard.Len(), "a drop is not a paste, so the clipboard stays empty")
	assert.False(t, m.clipboard.IsCut(), "a drop must not turn the clipboard into a cut")
}

// TestPasteMsgDropsIntoSecondPanel checks that the destination is the directory
// the focused panel shows rather than the first panel, so that a file follows
// focus the same way a key press would.
func TestPasteMsgDropsIntoSecondPanel(t *testing.T) {
	first := t.TempDir()
	second := t.TempDir()
	source := t.TempDir()
	dropped := filepath.Join(source, "dropped.txt")
	require.NoError(t, os.WriteFile(dropped, []byte("data"), 0o600))

	m := defaultTestModel(first, second)
	m.fileModel.SetFocusedPanelIndex(1)
	require.Equal(t, second, m.getFocusedFilePanel().Location)

	cmd := m.handlePasteMsg(tea.PasteMsg{Content: dropped})
	require.NotNil(t, cmd, "a dropped file should start an import")
	runCmd(cmd)

	assert.Eventually(t, func() bool {
		return fileExists(filepath.Join(second, "dropped.txt"))
	}, DefaultTestTimeout, DefaultTestTick,
		"the dropped file should follow the focused panel")
	assert.NoFileExists(t, filepath.Join(first, "dropped.txt"),
		"the first panel is not focused, so nothing should land there")
}

// TestPasteMsgYieldsToFocusedSearchBar checks that a path pasted while a field
// has the keyboard is typed into that field rather than imported.
//
// This is the whole of the promise that pasting into superfile still pastes:
// the drop only happens when nothing else would take the text.
func TestPasteMsgYieldsToFocusedSearchBar(t *testing.T) {
	dest := t.TempDir()
	source := t.TempDir()
	dropped := filepath.Join(source, "dropped.txt")
	require.NoError(t, os.WriteFile(dropped, []byte("data"), 0o600))

	m := defaultTestModel(dest)
	m.searchBarFocus()
	// Focusing the bar sets firstTextInput, which swallows the first message so
	// the key that focused it is not typed into it. A later paste is not that
	// message, so the flag is cleared the way the first update would clear it.
	m.firstTextInput = false

	require.Nil(t, m.handlePasteMsg(tea.PasteMsg{Content: dropped}),
		"a paste belongs to the field that has focus")
	TeaUpdate(m, tea.PasteMsg{Content: dropped})

	assert.Equal(t, dropped, m.getFocusedFilePanel().SearchBar.Value(),
		"the path should be typed into the search bar")
	assert.NoFileExists(t, filepath.Join(dest, "dropped.txt"),
		"nothing should be imported while a field has the keyboard")
}

// TestPasteMsgIgnoresProse checks that pasting ordinary text with nothing
// focused does nothing at all, which is what keeps a stray Ctrl+V from filling
// the panel.
func TestPasteMsgIgnoresProse(t *testing.T) {
	dir := dirWithFiles(t, 3)
	m := defaultTestModel(dir)
	before := m.getFocusedFilePanel().ElemCount()

	for _, content := range []string{"", "hello world", "/tmp/does-not-exist-here"} {
		assert.Nil(t, m.handlePasteMsg(tea.PasteMsg{Content: content}),
			"pasting %q should not start an import", content)
	}
	assert.Equal(t, before, m.getFocusedFilePanel().ElemCount(),
		"nothing should have been added to the panel")
}

// TestPasteMsgSkipsFileAlreadyInDestination checks that dropping a file onto
// the directory it is already in does nothing, rather than making a renamed
// duplicate of it.
func TestPasteMsgSkipsFileAlreadyInDestination(t *testing.T) {
	dir := dirWithFiles(t, 3)
	m := defaultTestModel(dir)
	before := m.getFocusedFilePanel().ElemCount()
	existing := filepath.Join(dir, "filea.txt")
	require.FileExists(t, existing)

	assert.Nil(t, m.handlePasteMsg(tea.PasteMsg{Content: existing}),
		"there is nothing to copy into the directory the file is already in")
	assert.Equal(t, before, m.getFocusedFilePanel().ElemCount(),
		"no duplicate should have been made")
}

// TestPasteIsIgnoredWhileAModalIsOpen checks that a modal captures an incoming
// paste as well as the keys and the pointer.
//
// Files arriving behind a modal would be imported somewhere the user is not
// looking, and a text modal would otherwise get the paths as its own text.
func TestPasteIsIgnoredWhileAModalIsOpen(t *testing.T) {
	modals := []struct {
		name string
		open func(m *model)
	}{
		{"typing", func(m *model) { m.typingModal.open = true }},
		{"prompt", func(m *model) {
			TeaUpdate(m, utils.TeaRuneKeyMsg(common.Hotkeys.OpenSPFPrompt[0]))
		}},
		{"notify", func(m *model) { m.notifyModel.Open() }},
		{"zoxide", func(m *model) { m.zoxideModal.Open() }},
		{"sort", func(m *model) { m.sortModal.Open(m.getFocusedFilePanel().SortKind) }},
		{"help menu", func(m *model) { m.helpMenu.Open() }},
		{"spf error", func(m *model) {
			m.spfError = spferror.New(true, "Error", "boom", &spferror.FileListErrorState{})
		}},
		{"file rename", func(m *model) { m.fileModel.Renaming = true }},
		{"context menu", func(m *model) {
			x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
			m.handleMouseMsg(rightClickAt(x, y))
			require.True(t, m.contextMenu.IsOpen())
		}},
	}

	source := t.TempDir()
	dropped := filepath.Join(source, "dropped.txt")
	require.NoError(t, os.WriteFile(dropped, []byte("data"), 0o600))

	for _, modal := range modals {
		t.Run(modal.name, func(t *testing.T) {
			dir := dirWithFiles(t, 3)
			m := defaultTestModel(dir)
			modal.open(m)

			assert.Nil(t, m.handlePasteMsg(tea.PasteMsg{Content: dropped}),
				"a paste should be left to the %s", modal.name)
			assert.NoFileExists(t, filepath.Join(dir, "dropped.txt"),
				"nothing should be imported while the %s is open", modal.name)
		})
	}
}

// TestPasteIsIgnoredWhileAFieldHasFocus checks the fields that are not modals:
// a search bar is open but only typing into itself, and a drag owns the pointer
// until it ends.
func TestPasteIsIgnoredWhileAFieldHasFocus(t *testing.T) {
	source := t.TempDir()
	dropped := filepath.Join(source, "dropped.txt")
	require.NoError(t, os.WriteFile(dropped, []byte("data"), 0o600))

	t.Run("file panel search bar", func(t *testing.T) {
		dir := dirWithFiles(t, 3)
		m := defaultTestModel(dir)
		m.searchBarFocus()
		m.firstTextInput = false

		assert.Nil(t, m.handlePasteMsg(tea.PasteMsg{Content: dropped}))
		assert.NoFileExists(t, filepath.Join(dir, "dropped.txt"))
	})

	t.Run("sidebar search bar", func(t *testing.T) {
		dir := dirWithFiles(t, 3)
		m := defaultTestModel(dir)
		m.sidebarModel.SearchBarFocus()

		assert.Nil(t, m.handlePasteMsg(tea.PasteMsg{Content: dropped}))
		assert.NoFileExists(t, filepath.Join(dir, "dropped.txt"))
	})

	t.Run("drag in progress", func(t *testing.T) {
		dir := dirWithFiles(t, 3)
		m := defaultTestModel(dir)
		m.viewContent()
		x, y := mustFindCell(t, m, mouse.TargetFilePanelItem, 0)
		m.handleMouseMsg(leftClickAt(x, y))
		m.handleMouseMsg(mouseMotionAt(x+5, y))
		require.True(t, m.drag.pending || m.drag.active, "the press should have started a drag")

		assert.Nil(t, m.handlePasteMsg(tea.PasteMsg{Content: dropped}),
			"a drag owns the pointer until it ends")
		assert.NoFileExists(t, filepath.Join(dir, "dropped.txt"))
	})
}

// TestGetDroppedItemsCmdRefusesImpossibleDestinations checks that a drop onto a
// directory that is not there reports it instead of failing silently, the way a
// paste onto a bad destination does.
func TestGetDroppedItemsCmdRefusesImpossibleDestinations(t *testing.T) {
	dir := dirWithFiles(t, 3)
	m := defaultTestModel(dir)
	gone := filepath.Join(dir, "not-here")

	cmd := m.getDroppedItemsCmd(gone, []string{filepath.Join(dir, "filea.txt")})
	require.NotNil(t, cmd, "an impossible destination should say so rather than do nothing")
	msg := cmd()
	require.IsType(t, NotifyModalUpdateMsg{}, msg, "the drop should open a notification")
}

// TestValidateDropOperation checks the rules a drop has to satisfy, which are
// the rules a paste already follows.
func TestValidateDropOperation(t *testing.T) {
	dir := dirWithFiles(t, 2)
	file := filepath.Join(dir, "filea.txt")
	subdir := filepath.Join(dir, "subdir")
	nested := filepath.Join(subdir, "deep")
	require.NoError(t, os.MkdirAll(nested, 0o755))

	t.Run("file into its directory", func(t *testing.T) {
		assert.NoError(t, validateDropOperation(dir, []string{file}))
	})
	t.Run("file into another directory", func(t *testing.T) {
		assert.NoError(t, validateDropOperation(t.TempDir(), []string{file}))
	})
	t.Run("destination is gone", func(t *testing.T) {
		require.Error(t, validateDropOperation(filepath.Join(dir, "missing"), []string{file}))
	})
	t.Run("destination is a file", func(t *testing.T) {
		err := validateDropOperation(file, []string{file})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})
	t.Run("directory into itself", func(t *testing.T) {
		require.Error(t, validateDropOperation(nested, []string{subdir}))
	})
}

// TestDroppableInto checks what a drop leaves out: a path already in the
// destination, and a path the paste listed twice.
func TestDroppableInto(t *testing.T) {
	dir := t.TempDir()
	here := filepath.Join(dir, "already.txt")
	away := filepath.Join(t.TempDir(), "away.txt")

	got := droppableInto(dir, []string{here, away, away, filepath.Join(dir, "another.txt")})
	assert.Equal(t, []string{away}, got,
		"paths already in the destination and repeated paths should be dropped")
}
