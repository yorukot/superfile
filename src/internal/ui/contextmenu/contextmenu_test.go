package contextmenu

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testItems() []Item {
	return []Item{
		{Label: "Open", Shortcut: "enter", Action: ActionOpen},
		{Label: "Copy", Shortcut: "ctrl+c", Action: ActionCopy},
		{Label: "Delete", Shortcut: "ctrl+d", Action: ActionDelete},
	}
}

// TestOpenClampsToScreen checks that a menu opened near an edge is moved fully
// back into view, since the stored origin is what both rendering and hit-testing
// read.
func TestOpenClampsToScreen(t *testing.T) {
	var m Model
	items := testItems()

	// A menu that fits where it was asked for is not moved.
	m.Open(20, 8, 60, 20, items)
	assert.Equal(t, 20, m.OriginX())
	assert.Equal(t, 8, m.OriginY(), "a menu that fits is not moved")

	// Far past the right edge.
	m.Open(1000, 30, 60, 20, items)
	assert.Equal(t, 60-m.Width(), m.OriginX())
	assert.GreaterOrEqual(t, m.OriginX()+m.Width(), 0)
	assert.LessOrEqual(t, m.OriginX()+m.Width(), 60)

	// Far past the bottom edge.
	m.Open(10, 1000, 60, 20, items)
	assert.Equal(t, 20-m.Height(), m.OriginY())
	assert.LessOrEqual(t, m.OriginY()+m.Height(), 20)

	// Negative coordinates.
	m.Open(-5, -5, 60, 20, items)
	assert.Equal(t, 0, m.OriginX())
	assert.Equal(t, 0, m.OriginY())
}

// TestOpenWithNoItemsStaysClosed checks that there is never a menu with nothing
// in it to choose.
func TestOpenWithNoItemsStaysClosed(t *testing.T) {
	m := Model{}
	m.Open(5, 5, 80, 24, nil)
	assert.False(t, m.IsOpen())
	assert.False(t, m.HasItems())
}

// TestWidthCoversLongestEntryAndShortcut checks that every row is padded to the
// widest label and shortcut, which is what keeps the drawn rows and the
// published regions in agreement.
func TestWidthCoversLongestEntryAndShortcut(t *testing.T) {
	m := Model{}
	m.Open(0, 0, 80, 24, []Item{
		{Label: "a", Action: ActionOpen},
		{Label: "a much longer entry", Action: ActionCopy},
		{Label: "b", Shortcut: "ctrl+shift+x", Action: ActionDelete},
	})

	// The widest label and the widest shortcut both have to fit in one row.
	longestLabel := "a much longer entry"
	longestShortcut := "ctrl+shift+x"
	assert.Equal(t, len(longestLabel)+len(longestShortcut),
		m.Width()-sideBorders-cursorColumn-cursorPad-shortcutGap)
}

// TestHeightIsEntriesPlusBorders checks the row count the hit-test regions rely on.
func TestHeightIsEntriesPlusBorders(t *testing.T) {
	m := Model{}
	m.Open(0, 0, 80, 24, testItems())
	assert.Equal(t, len(testItems())+borderRows, m.Height())
}

// TestNavigationWraps checks the cursor moves through every entry and wraps.
func TestNavigationWraps(t *testing.T) {
	m := Model{}
	m.Open(0, 0, 80, 24, testItems())

	action, ok := m.SelectedAction()
	require.True(t, ok)
	assert.Equal(t, ActionOpen, action)

	m.ListDown()
	action, _ = m.SelectedAction()
	assert.Equal(t, ActionCopy, action)

	m.ListUp()
	m.ListUp()
	action, _ = m.SelectedAction()
	assert.Equal(t, ActionDelete, action, "moving up from the top wraps to the end")

	m.ListDown()
	action, _ = m.SelectedAction()
	assert.Equal(t, ActionOpen, action)
}

// TestChosenActionRejectsOutOfRange checks that a stale region cannot select an
// entry that is not there.
func TestChosenActionRejectsOutOfRange(t *testing.T) {
	m := Model{}
	m.Open(0, 0, 80, 24, testItems())

	_, ok := m.ChosenAction(-1)
	assert.False(t, ok)
	_, ok = m.ChosenAction(len(testItems()))
	assert.False(t, ok)
	_, ok = m.ChosenAction(99)
	assert.False(t, ok)

	action, ok := m.ChosenAction(1)
	require.True(t, ok)
	assert.Equal(t, ActionCopy, action)
}

// TestHighlightItem checks the pointer and the keyboard share one cursor.
func TestHighlightItem(t *testing.T) {
	m := Model{}
	m.Open(0, 0, 80, 24, testItems())

	require.True(t, m.HighlightItem(2))
	action, _ := m.SelectedAction()
	assert.Equal(t, ActionDelete, action)

	assert.False(t, m.HighlightItem(-1))
	assert.False(t, m.HighlightItem(3))
	// A rejected highlight leaves the cursor where it was.
	action, _ = m.SelectedAction()
	assert.Equal(t, ActionDelete, action)
}

// TestCloseResetsState checks a reopened menu starts clean rather than
// remembering the previous entries.
func TestCloseResetsState(t *testing.T) {
	m := Model{}
	m.Open(0, 0, 80, 24, testItems())
	m.ListDown()
	m.Close()

	assert.False(t, m.IsOpen())
	assert.Empty(t, m.Items())

	_, ok := m.SelectedAction()
	assert.False(t, ok)

	m.Open(0, 0, 80, 24, []Item{{Label: "Paste", Action: ActionPaste}})
	action, ok := m.SelectedAction()
	require.True(t, ok)
	assert.Equal(t, ActionPaste, action)
}

// TestNavigationOnClosedMenuIsSafe checks the menu tolerates being driven while
// closed, which happens whenever key routing changes.
func TestNavigationOnClosedMenuIsSafe(t *testing.T) {
	m := Model{}
	m.ListUp()
	m.ListDown()
	assert.False(t, m.HighlightItem(0))
	assert.Empty(t, m.Render())
}

// TestRenderedRowsShareOneWidth checks every row is padded to the same width, so
// the drawn borders line up and the hit-test rows line up with them.
func TestRenderedRowsShareOneWidth(t *testing.T) {
	m := Model{}
	m.Open(0, 0, 80, 24, testItems())

	lines := strings.Split(stripStyle(m.Render()), "\n")
	require.NotEmpty(t, lines)

	want := runewidth.StringWidth(lines[0])
	for row, line := range lines {
		assert.Equal(t, want, runewidth.StringWidth(line),
			"rendered row %d is not the same width as row 0", row)
	}
	assert.Len(t, lines, m.Height()+1,
		"the rendered menu should have one line per row, plus the trailing newline")
}

// stripStyle removes escape sequences so that rendered width can be measured.
func stripStyle(s string) string {
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
