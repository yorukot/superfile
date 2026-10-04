// Tests for the preview pane update path, including the preview
// override used while the find modal is open.
package filemodel

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/preview"
	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

const (
	testModelWidth   = 100
	testModelHeight  = 40
	testPreviewWidth = 30

	runCmdTimeout = 10 * time.Second
)

// testFixtureFiles maps fixture file names to their contents. Each
// content is unique per file and shorter than testPreviewWidth, so
// rendered lines cannot wrap and assertions are unambiguous.
var testFixtureFiles = map[string]string{ //nolint:gochecknoglobals // fixed test fixture data, not runtime state
	"alpha.txt":   "alpha file content AAA",
	"bravo.txt":   "bravo file content BBB",
	"charlie.txt": "charlie file content CCC",
}

// testModelSetup describes the initial state of a model built by
// newTestModel.
type testModelSetup struct {
	// paneOpen opens the preview pane.
	paneOpen bool
	// emptyPanel skips fixture file creation, leaving the panel empty.
	emptyPanel bool
	// selectedName positions the panel cursor on the fixture file with
	// that name.
	selectedName string
	// overrideName sets previewOverride to the fixture file with that
	// name.
	overrideName string
}

// setupCommonGlobals pins the common package globals that rendering
// depends on. It saves and restores the current values, and never
// calls the config-file loaders that can os.Exit on malformed input.
func setupCommonGlobals(t *testing.T) {
	t.Helper()

	oldConfig := common.Config
	oldHotkeys := common.Hotkeys
	t.Cleanup(func() {
		common.SetConfig(oldConfig)
		common.Hotkeys = oldHotkeys //nolint:reassign // no common.Hotkeys setter exists; restores the saved state
	})

	common.SetConfig(common.ConfigType{
		DefaultOpenFilePreview:  false,
		Nerdfont:                false,
		EnableFilePreviewBorder: false,
		FilePreviewWidth:        0,
		SidebarWidth:            0,
		TransparentBackground:   false,
		CodePreviewer:           "",
	})
	// ConfirmTyping and Quit must be non-empty: LoadPrerenderedVariables
	// indexes their first elements.
	common.Hotkeys = common.HotkeysType{ //nolint:reassign // no common.Hotkeys setter exists; pinned per test
		SearchBar:     []string{"/"},
		ConfirmTyping: []string{"\n"},
		Quit:          []string{"q"},
	}

	// Pin the terminal environment so isKittyCapable() (which reads
	// TERM and TERM_PROGRAM) is deterministically false; otherwise in
	// a kitty-capable terminal, the render path would emit a
	// raw-transmit cmd and trip the no-raw-transmit assertions below.
	t.Setenv("TERM", "dumb")
	t.Setenv("TERM_PROGRAM", "")

	common.LoadThemeConfig()
	common.LoadPrerenderedVariables()
}

// newTestModel builds a Model with a single focused directory panel
// backed by real fixture files in a temporary directory, plus a real
// preview.Model. It returns the model and a name → absolute path map
// for the fixture files.
func newTestModel(t *testing.T, setup testModelSetup) (*Model, map[string]string) {
	t.Helper()
	setupCommonGlobals(t)

	dir := t.TempDir()
	paths := make(map[string]string, len(testFixtureFiles))
	if !setup.emptyPanel {
		for name, content := range testFixtureFiles {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
				t.Fatalf("writing fixture %q: %v", name, err)
			}
			paths[name] = filepath.Join(dir, name)
		}
	}

	panel := filepanel.New(dir, true, "", sortmodel.SortByName, false)
	panel.UpdateElementsIfNeeded(true, false)

	if setup.selectedName != "" {
		idx := panel.FindElementIndexByName(setup.selectedName)
		if idx < 0 {
			t.Fatalf("fixture %q not found among panel elements", setup.selectedName)
		}
		panel.SetCursorPosition(idx)
	}

	previewModel := preview.New()
	t.Cleanup(previewModel.CleanUp)

	m := &Model{
		FilePanels:           []filepanel.Model{panel},
		FocusedPanelIndex:    0,
		Width:                testModelWidth,
		Height:               testModelHeight,
		ExpectedPreviewWidth: testPreviewWidth,
		FilePreview:          previewModel,
	}

	if setup.paneOpen {
		m.FilePreview.SetOpen(true)
	}
	if setup.overrideName != "" {
		p, ok := paths[setup.overrideName]
		if !ok {
			// emptyPanel skips fixture creation, so the map is empty and
			// a lookup would silently yield "".
			t.Fatalf("setup.overrideName %q not in fixtures (emptyPanel skips fixture creation)", setup.overrideName)
		}
		m.previewOverride = p
	}

	return m, paths
}

// runTeaCmd runs a tea.Cmd to completion and returns its message. A
// nil command yields a nil message.
func runTeaCmd(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()
	if cmd == nil {
		return nil
	}
	ch := make(chan tea.Msg, 1)
	go func() {
		ch <- cmd()
	}()
	select {
	case msg := <-ch:
		return msg
	case <-time.After(runCmdTimeout):
		t.Fatalf("timed out after %s running tea.Cmd", runCmdTimeout)
		return nil
	}
}

var ansiEscapeRe = regexp.MustCompile(`\x1B\[[0-9;]*[A-Za-z]`)

// stripANSI removes ANSI escape sequences so content assertions are
// not affected by syntax-highlighting colors.
func stripANSI(s string) string {
	return ansiEscapeRe.ReplaceAllString(s, "")
}

// assertContentContains asserts that content contains marker, ignoring
// ANSI escape sequences.
func assertContentContains(t *testing.T, content, marker, what string) {
	t.Helper()
	if !strings.Contains(stripANSI(content), marker) {
		t.Errorf("%s: content does not contain %q (got %q)", what, marker, stripANSI(content))
	}
}

// assertContentExcludes asserts that content does not contain marker,
// ignoring ANSI escape sequences.
func assertContentExcludes(t *testing.T, content, marker, what string) {
	t.Helper()
	if strings.Contains(stripANSI(content), marker) {
		t.Errorf("%s: content should not contain %q (got %q)", what, marker, stripANSI(content))
	}
}

// applyUpdateMsg runs the render message through the production apply
// path (the dimension and location guards, the preview Apply, and the
// raw-transmit handling) and fails the test if it returns a command.
func applyUpdateMsg(t *testing.T, m *Model, msg preview.UpdateMsg) {
	t.Helper()
	if cmd := m.UpdatePreviewPanel(msg); cmd != nil {
		t.Fatalf("expected no raw-transmit cmd, got non-nil")
	}
}

func TestGetFilePreviewCmd(t *testing.T) {
	t.Run("closed pane returns nil", testGetFilePreviewCmdClosedPane)
	t.Run("open pane with empty panel clears it synchronously", testGetFilePreviewCmdEmptyPanelClears)
	t.Run("same location without force returns nil", testGetFilePreviewCmdSameLocation)
	t.Run("different location submits a render that updates the pane", testGetFilePreviewCmdDifferentLocation)
	t.Run("force re-renders the same location", testGetFilePreviewCmdForceSameLocation)
	t.Run("active override without force leaves the pane untouched", testGetFilePreviewCmdOverrideNoForce)
	t.Run("active override with force renders the override", testGetFilePreviewCmdOverrideForce)
	t.Run("closed pane with active override returns nil", testGetFilePreviewCmdClosedOverride)
}

func TestSetPreviewPathCmd(t *testing.T) {
	t.Run("empty path is a no-op", testSetPreviewPathCmdEmpty)
	t.Run("closed pane remembers the override and renders it on reopen", testSetPreviewPathCmdClosedPane)
	t.Run("path equal to current pane location sets override without rendering", testSetPreviewPathCmdSameLocation)
	t.Run("new path submits a render that updates the pane", testSetPreviewPathCmdNewPath)
	t.Run("setting the same path twice does not re-render", testSetPreviewPathCmdSameTwice)
}

func TestClearPreviewOverride(t *testing.T) {
	t.Run("no override is a no-op", testClearPreviewOverrideNone)
	t.Run("closed pane clears without rendering", testClearPreviewOverrideClosedPane)
	t.Run("pane already showing the selection skips the re-render", testClearPreviewOverrideShowsSelection)
	t.Run("pane showing the override re-renders the selection", testClearPreviewOverrideShowsOverride)
	t.Run("empty panel clears synchronously", testClearPreviewOverrideEmptyPanel)
}

// mkMsg builds a preview.UpdateMsg. It is a named function (not a
// closure) so the subtests can share it and gocritic's unlambda
// checker is satisfied. The height is fixed to testModelHeight;
// width and reqID vary per subtest.
func mkMsg(location, content, raw string, width, reqID int) preview.UpdateMsg {
	return preview.NewUpdateMsg(location, content, raw, width, testModelHeight, reqID)
}

func TestUpdatePreviewPanel(t *testing.T) {
	t.Run("matching message is applied", testUpdatePreviewPanelMatching)
	t.Run("stale location without override is dropped", testUpdatePreviewPanelStaleLocation)
	t.Run("nil selection without override is dropped", testUpdatePreviewPanelNilSelection)
	t.Run("dimension mismatch without override is dropped", testUpdatePreviewPanelDimMismatch)
	t.Run("matching message while override is applied", testUpdatePreviewPanelMatchingOverride)
	t.Run("raw transmit while override is forwarded", testUpdatePreviewPanelRawTransmit)
	t.Run("message for the selection while override is active is dropped", testUpdatePreviewPanelSelectionWhileOverride)
	t.Run("dimension mismatch while override is active is dropped", testUpdatePreviewPanelDimMismatchOverride)
}

func TestHasPreviewOverride(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{selectedName: "alpha.txt"})
	if m.HasPreviewOverride() {
		t.Fatal("expected no override on a fresh model")
	}

	if cmd := m.SetPreviewPathCmd(paths["bravo.txt"]); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if !m.HasPreviewOverride() {
		t.Fatal("expected override after SetPreviewPathCmd")
	}

	if cmd := m.ClearPreviewOverride(); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if m.HasPreviewOverride() {
		t.Fatal("expected no override after ClearPreviewOverride")
	}
}

// Subtest bodies for the tests above, grouped by parent test.

// testGetFilePreviewCmdClosedPane tests that a closed pane returns nil.
func testGetFilePreviewCmdClosedPane(t *testing.T) {
	m, _ := newTestModel(t, testModelSetup{selectedName: "alpha.txt"})
	if cmd := m.GetFilePreviewCmd(false); cmd != nil {
		t.Fatalf("expected nil cmd for closed pane, got non-nil")
	}
}

// testGetFilePreviewCmdEmptyPanelClears tests that an open pane with an empty panel clears synchronously.
func testGetFilePreviewCmdEmptyPanelClears(t *testing.T) {
	m, _ := newTestModel(t, testModelSetup{paneOpen: true, emptyPanel: true})
	m.FilePreview.Apply(preview.NewUpdateMsg("/stale", "STALE CONTENT", "", testPreviewWidth, testModelHeight, 0))

	if cmd := m.GetFilePreviewCmd(false); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != "" {
		t.Errorf("expected pane location cleared, got %q", loc)
	}
	assertContentExcludes(t, m.FilePreview.GetContent(), "STALE CONTENT", "cleared pane")
	if m.FilePreview.IsLoading() {
		t.Error("expected pane to not be loading after synchronous clear")
	}
}

// testGetFilePreviewCmdSameLocation tests that the same location without force returns nil.
func testGetFilePreviewCmdSameLocation(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})
	m.FilePreview.SetLocation(paths["alpha.txt"])
	m.FilePreview.SetLoading()

	if cmd := m.GetFilePreviewCmd(false); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if !m.FilePreview.IsLoading() {
		t.Error("loading state must not change when no render is submitted")
	}
}

// testGetFilePreviewCmdDifferentLocation tests that a different location submits a render that updates the pane.
func testGetFilePreviewCmdDifferentLocation(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})

	cmd := m.GetFilePreviewCmd(false)
	if cmd == nil {
		t.Fatal("expected render cmd, got nil")
	}
	if !m.FilePreview.IsLoading() {
		t.Error("expected pane marked loading right after submission")
	}
	if m.ioReqCnt != 1 {
		t.Fatalf("expected ioReqCnt 1 after one submission, got %d", m.ioReqCnt)
	}

	msg, ok := runTeaCmd(t, cmd).(preview.UpdateMsg)
	if !ok {
		t.Fatalf("expected preview.UpdateMsg, got %T", msg)
	}
	if msg.GetLocation() != paths["alpha.txt"] {
		t.Errorf("expected message location %q, got %q", paths["alpha.txt"], msg.GetLocation())
	}
	if msg.GetReqID() != 0 {
		t.Errorf("expected req id 0, got %d", msg.GetReqID())
	}
	if msg.GetContentWidth() != testPreviewWidth || msg.GetContentHeight() != testModelHeight {
		t.Errorf("unexpected message dimensions %dx%d", msg.GetContentWidth(), msg.GetContentHeight())
	}

	applyUpdateMsg(t, m, msg)
	if loc := m.FilePreview.GetLocation(); loc != paths["alpha.txt"] {
		t.Errorf("expected pane location %q, got %q", paths["alpha.txt"], loc)
	}
	assertContentContains(t, m.FilePreview.GetContent(), "alpha file content", "applied content")
	if w := m.FilePreview.GetContentWidth(); w != testPreviewWidth {
		t.Errorf("expected content width %d, got %d", testPreviewWidth, w)
	}
	if h := m.FilePreview.GetContentHeight(); h != testModelHeight {
		t.Errorf("expected content height %d, got %d", testModelHeight, h)
	}
	if m.FilePreview.IsLoading() {
		t.Error("expected loading reset after apply")
	}
}

// testGetFilePreviewCmdForceSameLocation tests that force re-renders the same location.
func testGetFilePreviewCmdForceSameLocation(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})
	m.FilePreview.SetLocation(paths["alpha.txt"])

	first := m.GetFilePreviewCmd(true)
	if first == nil {
		t.Fatal("expected re-render cmd, got nil")
	}
	applyUpdateMsg(t, m, runTeaCmd(t, first).(preview.UpdateMsg))

	second := m.GetFilePreviewCmd(true)
	if second == nil {
		t.Fatal("expected second re-render cmd, got nil")
	}
	msg := runTeaCmd(t, second).(preview.UpdateMsg)
	if msg.GetReqID() != 1 {
		t.Errorf("expected second req id 1, got %d", msg.GetReqID())
	}
	applyUpdateMsg(t, m, msg)
	assertContentContains(t, m.FilePreview.GetContent(), "alpha file content", "re-rendered content")
}

// testGetFilePreviewCmdOverrideNoForce tests that an active override without force leaves the pane untouched.
func testGetFilePreviewCmdOverrideNoForce(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{
		paneOpen:     true,
		selectedName: "alpha.txt",
		overrideName: "bravo.txt",
	})
	m.FilePreview.SetLocation(paths["charlie.txt"])
	m.FilePreview.SetLoading()

	if cmd := m.GetFilePreviewCmd(false); cmd != nil {
		t.Fatalf("expected nil cmd while override active without force, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != paths["charlie.txt"] {
		t.Errorf("pane should keep its location %q, got %q", paths["charlie.txt"], loc)
	}
	if !m.FilePreview.IsLoading() {
		t.Error("loading state must not change when no render is submitted")
	}
}

// testGetFilePreviewCmdOverrideForce tests that an active override with force renders the override.
func testGetFilePreviewCmdOverrideForce(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{
		paneOpen:     true,
		selectedName: "alpha.txt",
		overrideName: "bravo.txt",
	})
	m.FilePreview.SetLocation(paths["charlie.txt"])

	cmd := m.GetFilePreviewCmd(true)
	if cmd == nil {
		t.Fatal("expected render cmd for active override, got nil")
	}
	msg, ok := runTeaCmd(t, cmd).(preview.UpdateMsg)
	if !ok {
		t.Fatalf("expected preview.UpdateMsg, got %T", msg)
	}
	if msg.GetLocation() != paths["bravo.txt"] {
		t.Errorf("expected message location %q, got %q", paths["bravo.txt"], msg.GetLocation())
	}

	applyUpdateMsg(t, m, msg)
	if loc := m.FilePreview.GetLocation(); loc != paths["bravo.txt"] {
		t.Errorf("expected pane location %q, got %q", paths["bravo.txt"], loc)
	}
	assertContentContains(t, m.FilePreview.GetContent(), "bravo file content", "applied override content")
}

// testGetFilePreviewCmdClosedOverride tests that a closed pane with an active override returns nil.
func testGetFilePreviewCmdClosedOverride(t *testing.T) {
	m, _ := newTestModel(t, testModelSetup{
		selectedName: "alpha.txt",
		overrideName: "bravo.txt",
	})
	if cmd := m.GetFilePreviewCmd(true); cmd != nil {
		t.Fatalf("expected nil cmd for closed pane, got non-nil")
	}
}

// testSetPreviewPathCmdEmpty tests that an empty path is a no-op.
func testSetPreviewPathCmdEmpty(t *testing.T) {
	m, _ := newTestModel(t, testModelSetup{selectedName: "alpha.txt"})
	if cmd := m.SetPreviewPathCmd(""); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if m.HasPreviewOverride() {
		t.Error("expected no override to be set")
	}
}

// testSetPreviewPathCmdClosedPane tests that a closed pane remembers the override and renders it on reopen.
func testSetPreviewPathCmdClosedPane(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{selectedName: "alpha.txt"})

	if cmd := m.SetPreviewPathCmd(paths["bravo.txt"]); cmd != nil {
		t.Fatalf("expected nil cmd while pane closed, got non-nil")
	}
	if !m.HasPreviewOverride() {
		t.Fatal("expected override to be remembered while pane closed")
	}
	if m.ioReqCnt != 0 {
		t.Errorf("expected no render submitted while pane closed, ioReqCnt=%d", m.ioReqCnt)
	}

	// Simulate the pane being reopened (ensurePreviewDimensionsSync
	// always re-renders on dimension changes, which includes
	// opening).
	m.FilePreview.SetOpen(true)
	cmd := m.GetFilePreviewCmd(true)
	if cmd == nil {
		t.Fatal("expected render cmd on reopen, got nil")
	}
	msg, ok := runTeaCmd(t, cmd).(preview.UpdateMsg)
	if !ok {
		t.Fatalf("expected preview.UpdateMsg, got %T", msg)
	}
	if msg.GetLocation() != paths["bravo.txt"] {
		t.Errorf("expected message location %q, got %q", paths["bravo.txt"], msg.GetLocation())
	}
	applyUpdateMsg(t, m, msg)
	if loc := m.FilePreview.GetLocation(); loc != paths["bravo.txt"] {
		t.Errorf("expected pane to show the override %q, got %q", paths["bravo.txt"], loc)
	}
	assertContentContains(t, m.FilePreview.GetContent(), "bravo file content", "reopen render")
}

// testSetPreviewPathCmdSameLocation tests that a path equal to the pane location sets the override without rendering.
func testSetPreviewPathCmdSameLocation(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})
	m.FilePreview.SetLocation(paths["alpha.txt"])

	if cmd := m.SetPreviewPathCmd(paths["alpha.txt"]); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if !m.HasPreviewOverride() {
		t.Fatal("expected override to be set")
	}
	if m.ioReqCnt != 0 {
		t.Errorf("expected no render submitted, ioReqCnt=%d", m.ioReqCnt)
	}
}

// testSetPreviewPathCmdNewPath tests that a new path submits a render that updates the pane.
func testSetPreviewPathCmdNewPath(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})

	cmd := m.SetPreviewPathCmd(paths["bravo.txt"])
	if cmd == nil {
		t.Fatal("expected render cmd, got nil")
	}
	if !m.HasPreviewOverride() {
		t.Fatal("expected override to be set")
	}
	if m.ioReqCnt != 1 {
		t.Fatalf("expected ioReqCnt 1 after one submission, got %d", m.ioReqCnt)
	}

	msg, ok := runTeaCmd(t, cmd).(preview.UpdateMsg)
	if !ok {
		t.Fatalf("expected preview.UpdateMsg, got %T", msg)
	}
	if msg.GetLocation() != paths["bravo.txt"] {
		t.Errorf("expected message location %q, got %q", paths["bravo.txt"], msg.GetLocation())
	}
	applyUpdateMsg(t, m, msg)
	assertContentContains(t, m.FilePreview.GetContent(), "bravo file content", "applied override content")
}

// testSetPreviewPathCmdSameTwice tests that setting the same path twice does not re-render.
func testSetPreviewPathCmdSameTwice(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})

	first := m.SetPreviewPathCmd(paths["bravo.txt"])
	if first == nil {
		t.Fatal("expected render cmd for first set, got nil")
	}
	applyUpdateMsg(t, m, runTeaCmd(t, first).(preview.UpdateMsg))

	count := m.ioReqCnt
	if cmd := m.SetPreviewPathCmd(paths["bravo.txt"]); cmd != nil {
		t.Fatalf("expected nil cmd for same path, got non-nil")
	}
	if m.ioReqCnt != count {
		t.Errorf("expected ioReqCnt unchanged (%d), got %d", count, m.ioReqCnt)
	}
}

// testClearPreviewOverrideNone tests that clearing with no override is a no-op.
func testClearPreviewOverrideNone(t *testing.T) {
	m, _ := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})
	if cmd := m.ClearPreviewOverride(); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
}

// testClearPreviewOverrideClosedPane tests that a closed pane clears without rendering.
func testClearPreviewOverrideClosedPane(t *testing.T) {
	m, _ := newTestModel(t, testModelSetup{overrideName: "bravo.txt"})

	if cmd := m.ClearPreviewOverride(); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if m.HasPreviewOverride() {
		t.Error("expected override to be cleared")
	}
}

// testClearPreviewOverrideShowsSelection tests that a pane already showing the selection skips the re-render.
func testClearPreviewOverrideShowsSelection(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{
		paneOpen:     true,
		selectedName: "alpha.txt",
		overrideName: "bravo.txt",
	})
	// The pane still shows the selection: the override was set while
	// a pending render for it had not applied yet (the realistic
	// fast open/close of the find modal).
	m.FilePreview.SetLocation(paths["alpha.txt"])

	if cmd := m.ClearPreviewOverride(); cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if m.HasPreviewOverride() {
		t.Error("expected override to be cleared")
	}
	if loc := m.FilePreview.GetLocation(); loc != paths["alpha.txt"] {
		t.Errorf("pane should keep the selection %q, got %q", paths["alpha.txt"], loc)
	}
}

// testClearPreviewOverrideShowsOverride tests that a pane showing the override re-renders the selection.
func testClearPreviewOverrideShowsOverride(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{
		paneOpen:     true,
		selectedName: "charlie.txt",
		overrideName: "bravo.txt",
	})
	m.FilePreview.SetLocation(paths["bravo.txt"])

	cmd := m.ClearPreviewOverride()
	if cmd == nil {
		t.Fatal("expected render cmd, got nil")
	}
	if m.HasPreviewOverride() {
		t.Error("expected override to be cleared before the re-render")
	}

	msg, ok := runTeaCmd(t, cmd).(preview.UpdateMsg)
	if !ok {
		t.Fatalf("expected preview.UpdateMsg, got %T", msg)
	}
	if msg.GetLocation() != paths["charlie.txt"] {
		t.Errorf("expected message location %q, got %q", paths["charlie.txt"], msg.GetLocation())
	}
	applyUpdateMsg(t, m, msg)
	if loc := m.FilePreview.GetLocation(); loc != paths["charlie.txt"] {
		t.Errorf("expected pane location %q, got %q", paths["charlie.txt"], loc)
	}
	assertContentContains(t, m.FilePreview.GetContent(), "charlie file content", "selection content")
}

// testClearPreviewOverrideEmptyPanel tests that an empty panel clears synchronously.
func testClearPreviewOverrideEmptyPanel(t *testing.T) {
	m, _ := newTestModel(t, testModelSetup{
		paneOpen:   true,
		emptyPanel: true,
	})
	// No fixture files exist in this case, so the override is a
	// literal path instead of a resolved fixture name.
	m.previewOverride = "/find/results/bravo.txt"
	// The pane holds the override's path with no panel selection to
	// fall back to.
	m.FilePreview.SetLocation("/find/results/bravo.txt")

	if cmd := m.ClearPreviewOverride(); cmd != nil {
		t.Fatalf("expected nil cmd (synchronous clear), got non-nil")
	}
	if m.HasPreviewOverride() {
		t.Error("expected override to be cleared")
	}
	if loc := m.FilePreview.GetLocation(); loc != "" {
		t.Errorf("expected pane location cleared, got %q", loc)
	}
	if m.FilePreview.IsLoading() {
		t.Error("expected pane to not be loading after synchronous clear")
	}
}

// testUpdatePreviewPanelMatching tests that a matching message is applied.
func testUpdatePreviewPanelMatching(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})
	m.FilePreview.Apply(mkMsg("/stale", "STALE CONTENT", "", testPreviewWidth, 0))

	msg := mkMsg(paths["alpha.txt"], "alpha file content", "", testPreviewWidth, 1)
	cmd := m.UpdatePreviewPanel(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != paths["alpha.txt"] {
		t.Errorf("expected pane location %q, got %q", paths["alpha.txt"], loc)
	}
	assertContentContains(t, m.FilePreview.GetContent(), "alpha file content", "applied content")
	assertContentExcludes(t, m.FilePreview.GetContent(), "STALE CONTENT", "applied content")
	if w := m.FilePreview.GetContentWidth(); w != testPreviewWidth {
		t.Errorf("expected content width %d, got %d", testPreviewWidth, w)
	}
	if h := m.FilePreview.GetContentHeight(); h != testModelHeight {
		t.Errorf("expected content height %d, got %d", testModelHeight, h)
	}
}

// testUpdatePreviewPanelStaleLocation tests that a stale location without override is dropped.
func testUpdatePreviewPanelStaleLocation(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})
	m.FilePreview.SetLocation(paths["charlie.txt"])

	msg := mkMsg(paths["bravo.txt"], "STALE", "", testPreviewWidth, 1)
	cmd := m.UpdatePreviewPanel(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != paths["charlie.txt"] {
		t.Errorf("pane should keep its location %q, got %q", paths["charlie.txt"], loc)
	}
}

// testUpdatePreviewPanelNilSelection tests that a nil selection without override is dropped.
func testUpdatePreviewPanelNilSelection(t *testing.T) {
	m, _ := newTestModel(t, testModelSetup{paneOpen: true, emptyPanel: true})
	m.FilePreview.SetLocation("/stale")

	msg := mkMsg("/stale", "STALE", "", testPreviewWidth, 1)
	cmd := m.UpdatePreviewPanel(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != "/stale" {
		t.Errorf("pane should keep its location %q, got %q", "/stale", loc)
	}
}

// testUpdatePreviewPanelDimMismatch tests that a dimension mismatch without override is dropped.
func testUpdatePreviewPanelDimMismatch(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{paneOpen: true, selectedName: "alpha.txt"})
	m.FilePreview.SetLocation(paths["alpha.txt"])

	msg := mkMsg(paths["alpha.txt"], "STALE", "", 20, 1)
	cmd := m.UpdatePreviewPanel(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != paths["alpha.txt"] {
		t.Errorf("pane should keep its location %q, got %q", paths["alpha.txt"], loc)
	}
}

// testUpdatePreviewPanelMatchingOverride tests that a matching message is applied while the override is active.
func testUpdatePreviewPanelMatchingOverride(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{
		paneOpen:     true,
		selectedName: "alpha.txt",
		overrideName: "bravo.txt",
	})

	msg := mkMsg(paths["bravo.txt"], "bravo file content", "", testPreviewWidth, 1)
	cmd := m.UpdatePreviewPanel(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != paths["bravo.txt"] {
		t.Errorf("expected pane location %q, got %q", paths["bravo.txt"], loc)
	}
	assertContentContains(t, m.FilePreview.GetContent(), "bravo file content", "applied content")
}

// testUpdatePreviewPanelRawTransmit tests that a raw payload is forwarded while the override is active.
func testUpdatePreviewPanelRawTransmit(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{
		paneOpen:     true,
		selectedName: "alpha.txt",
		overrideName: "bravo.txt",
	})

	msg := mkMsg(paths["bravo.txt"], "image bytes", "RAW-DATA", testPreviewWidth, 1)
	cmd := m.UpdatePreviewPanel(msg)
	if cmd == nil {
		t.Fatal("expected raw-transmit cmd, got nil")
	}
	rawMsg := runTeaCmd(t, cmd)
	raw, ok := rawMsg.(tea.RawMsg)
	if !ok {
		t.Fatalf("expected tea.RawMsg, got %T", rawMsg)
	}
	if raw.Msg != "RAW-DATA" {
		t.Errorf("expected raw payload %q, got %v", "RAW-DATA", raw.Msg)
	}
	assertContentContains(t, m.FilePreview.GetContent(), "image bytes", "applied content")
}

// testUpdatePreviewPanelSelectionWhileOverride tests that a selection message while the override is active is dropped.
func testUpdatePreviewPanelSelectionWhileOverride(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{
		paneOpen:     true,
		selectedName: "alpha.txt",
		overrideName: "bravo.txt",
	})
	m.FilePreview.SetLocation(paths["charlie.txt"])

	msg := mkMsg(paths["alpha.txt"], "STALE", "", testPreviewWidth, 1)
	cmd := m.UpdatePreviewPanel(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != paths["charlie.txt"] {
		t.Errorf("pane should keep its location %q, got %q", paths["charlie.txt"], loc)
	}
}

// testUpdatePreviewPanelDimMismatchOverride tests that a dimension mismatch while the override is active is dropped.
func testUpdatePreviewPanelDimMismatchOverride(t *testing.T) {
	m, paths := newTestModel(t, testModelSetup{
		paneOpen:     true,
		selectedName: "alpha.txt",
		overrideName: "bravo.txt",
	})
	m.FilePreview.SetLocation(paths["bravo.txt"])

	msg := mkMsg(paths["bravo.txt"], "STALE", "", 20, 1)
	cmd := m.UpdatePreviewPanel(msg)
	if cmd != nil {
		t.Fatalf("expected nil cmd, got non-nil")
	}
	if loc := m.FilePreview.GetLocation(); loc != paths["bravo.txt"] {
		t.Errorf("pane should keep its location %q, got %q", paths["bravo.txt"], loc)
	}
}
