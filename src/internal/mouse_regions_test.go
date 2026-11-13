package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/mouse"
)

// renderedRows renders a frame and returns it as plain lines, so that a test can
// look at what a given terminal row actually shows.
func renderedRows(t *testing.T, m *model) []string {
	t.Helper()

	out := ansi.Strip(m.viewContent())
	return strings.Split(out, "\n")
}

// rowForItem returns the terminal row at which the registry places a file panel
// entry, or -1 when the entry has no published region.
func rowForItem(m *model, panelIndex, itemIndex int) int {
	panelX := 0
	if panelIndex > 0 {
		panelX = m.fileModel.PanelOriginX(panelIndex)
	}
	panelX += m.fullWidth - m.fileModel.Width

	for y := range m.fullHeight {
		for x := range m.fullWidth {
			if x < panelX {
				continue
			}
			target := m.mouseTargetAt(x, y)
			if target.Kind == mouse.TargetFilePanelItem &&
				target.PanelIndex == panelIndex && target.ItemIndex == itemIndex {
				return y
			}
		}
	}
	return -1
}

// TestMouseRegionsMatchRenderedRows checks the invariant the whole hit-test
// layer rests on: a row reported as holding a file entry really shows that
// entry's name on screen.
//
// The registry supplies the coordinates and the rendered frame supplies the
// expected content, so neither side can validate itself.
func TestMouseRegionsMatchRenderedRows(t *testing.T) {
	dir := t.TempDir()
	names := []string{"alpha.txt", "bravo.txt", "charlie.txt", "delta.txt", "echo.log"}
	for _, name := range names {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte("data"), 0o600))
	}

	m := defaultTestModel(dir)
	rows := renderedRows(t, m)
	panel := &m.fileModel.FilePanels[0]

	checked := 0
	for itemIndex := range panel.ElemCount() {
		element := panel.GetElementAtIdx(itemIndex)
		row := rowForItem(m, 0, itemIndex)
		if row == -1 {
			// The entry scrolled out of view, nothing to cross-check.
			continue
		}

		require.Less(t, row, len(rows), "row for %q is outside the rendered frame", element.Name)
		assert.Contains(t, rows[row], element.Name,
			"region for item %d points at a row showing %q", itemIndex, strings.TrimSpace(rows[row]))
		checked++
	}

	require.Positive(t, checked, "no file entry was published, the registry is empty")
}

// TestMouseRegionsFollowScroll verifies that the published rows track the
// panel's scroll position instead of staying anchored to the top of the panel.
func TestMouseRegionsFollowScroll(t *testing.T) {
	dir := t.TempDir()
	// Enough entries that the panel has to scroll, since the test model is
	// deliberately short.
	count := 200
	for i := range count {
		name := "entry" + strings.Repeat("0", 3-len(itoa(i))) + itoa(i) + ".txt"
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte("data"), 0o600))
	}

	m := defaultTestModel(dir)
	m.fileModel.FilePanels[0].ListUp()

	rows := renderedRows(t, m)
	panel := &m.fileModel.FilePanels[0]

	require.Positive(t, panel.GetCursor(), "cursor should have moved off the first entry")

	for itemIndex := range panel.ElemCount() {
		element := panel.GetElementAtIdx(itemIndex)
		row := rowForItem(m, 0, itemIndex)
		if row == -1 {
			continue
		}
		require.Less(t, row, len(rows))
		assert.Contains(t, rows[row], element.Name,
			"after scrolling, item %d points at a row showing %q", itemIndex,
			strings.TrimSpace(rows[row]))
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var digits []byte
	for i > 0 {
		digits = append([]byte{byte('0' + i%10)}, digits...)
		i /= 10
	}
	return string(digits)
}
