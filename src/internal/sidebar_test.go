package internal

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func dirSlice(count int) []directory {
	res := make([]directory, count)
	for i := range count {
		res[i] = directory{name: "Dir" + strconv.Itoa(i), location: "/a/" + strconv.Itoa(i)}
	}
	return res
}

func fullDirSlice(count int) []directory {
	return formDirctorySlice(dirSlice(count), dirSlice(count), dirSlice(count))
}

// Todo : Use t.Run(tt.name
// Todo : Get rid of global vars, use testdata in each test, even if there is a bit of
// duplication.
// Todo : Add tt.names

func Test_noActualDir(t *testing.T) {
	testcases := []struct {
		name     string
		sidebar  sidebarModel
		expected bool
	}{
		{
			"Empty invalid sidebar should have no actual directories",
			sidebarModel{},
			true,
		},
		{
			"Empty sidebar should have no actual directories",
			sidebarModel{
				directories: fullDirSlice(0),
				renderIndex: 0,
				cursor:      0,
			},
			true,
		},
		{
			"Non-Empty Sidebar with only pinned directories",
			sidebarModel{
				directories: formDirctorySlice(nil, dirSlice(10), nil),
			},
			false,
		},
		{
			"Non-Empty Sidebar with all directories",
			sidebarModel{
				directories: fullDirSlice(10),
			},
			false,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.sidebar.noActualDir())
		})
	}
}

func Test_isCursorInvalid(t *testing.T) {
	testcases := []struct {
		name     string
		sidebar  sidebarModel
		expected bool
	}{
		{
			"Empty invalid sidebar",
			sidebarModel{},
			true,
		},
		{
			"Cursor after all directories",
			sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 0,
				cursor:      32,
			},
			true,
		},
		{
			"Curson points to pinned divider",
			sidebarModel{
				directories: fullDirSlice(10),
				cursor:      10,
			},
			true,
		},
		{
			"Non-Empty Sidebar with all directories",
			sidebarModel{
				directories: fullDirSlice(10),
				cursor:      5,
			},
			false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.sidebar.isCursorInvalid())
		})
	}
}

func Test_resetCursor(t *testing.T) {
	data := []struct {
		name              string
		curSideBar        sidebarModel
		expectedCursorPos int
	}{
		{
			name: "Only Pinned directories",
			curSideBar: sidebarModel{
				directories: formDirctorySlice(nil, dirSlice(10), nil),
			},
			expectedCursorPos: 1, // After pinned divider
		},
		{
			name: "All kind of directories",
			curSideBar: sidebarModel{
				directories: fullDirSlice(10),
			},
			expectedCursorPos: 0, // First home
		},
		{
			name: "Only Disk",
			curSideBar: sidebarModel{
				directories: formDirctorySlice(nil, nil, dirSlice(10)),
			},
			expectedCursorPos: 2, // After pinned and dist divider
		},
		{
			name: "Empty Sidebar",
			curSideBar: sidebarModel{
				directories: fullDirSlice(0),
			},
			expectedCursorPos: 0, // Empty sidebar, cursor should reset to 0
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			tt.curSideBar.resetCursor()
			assert.Equal(t, tt.expectedCursorPos, tt.curSideBar.cursor)
		})
	}
}

func Test_lastRenderIndex(t *testing.T) {
	// Setup test data
	sidebarA := sidebarModel{
		directories: formDirctorySlice(
			dirSlice(10), dirSlice(10), dirSlice(10),
		),
	}
	sidebarB := sidebarModel{
		directories: formDirctorySlice(
			dirSlice(1), nil, dirSlice(5),
		),
	}

	testCases := []struct {
		name              string
		sidebar           sidebarModel
		mainPanelHeight   int
		startIndex        int
		expectedLastIndex int
		explanation       string
	}{
		{
			name:              "Small viewport with home directories",
			sidebar:           sidebarA,
			mainPanelHeight:   10,
			startIndex:        0,
			expectedLastIndex: 6,
			explanation:       "3(initialHeight) + 7 (0-6 home dirs)",
		},
		{
			name:              "Medium viewport showing home and some pinned",
			sidebar:           sidebarA,
			mainPanelHeight:   20,
			startIndex:        0,
			expectedLastIndex: 14,
			explanation:       "3(initialHeight) + 10 (0-9 home dirs) + 3 (10-pinned divider) + 4 (11-14 pinned dirs)",
		},
		{
			name:              "Medium viewport starting from pinned dirs",
			sidebar:           sidebarA,
			mainPanelHeight:   20,
			startIndex:        11,
			expectedLastIndex: 25,
			explanation:       "3(initialHeight) + 10 (11-20 pinned dirs) + 3 (21-disk divider) + 4 (22-25 disk dirs)",
		},
		{
			name:              "Large viewport showing all directories",
			sidebar:           sidebarA,
			mainPanelHeight:   100,
			startIndex:        11,
			expectedLastIndex: 31,
			explanation:       "Last dir index is 31",
		},
		{
			name:              "Start index beyond directory count",
			sidebar:           sidebarA,
			mainPanelHeight:   100,
			startIndex:        32,
			expectedLastIndex: 31,
			explanation:       "When startIndex > len(directories), return last valid index",
		},
		{
			name:              "Asymmetric directory distribution",
			sidebar:           sidebarB,
			mainPanelHeight:   12,
			startIndex:        0,
			expectedLastIndex: 4,
			explanation:       "3(initialHeight) + 1 (0-homedir) + 3(1-pinned divider) + 3 (2-diskdivider) + 2 (3-4 diskdirs)",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sidebar.lastRenderedIndex(tt.mainPanelHeight, tt.startIndex)
			assert.Equal(t, tt.expectedLastIndex, result,
				"lastRenderedIndex failed: %s", tt.explanation)
		})
	}
}

func Test_firstRenderIndex(t *testing.T) {
	sidebarA := sidebarModel{
		directories: fullDirSlice(10),
	}
	sidebarB := sidebarModel{
		directories: formDirctorySlice(
			dirSlice(1), nil, dirSlice(5),
		),
	}
	sidebarC := sidebarModel{
		directories: formDirctorySlice(
			nil, dirSlice(5), dirSlice(5),
		),
	}
	sidebarD := sidebarModel{
		directories: formDirctorySlice(
			nil, nil, dirSlice(3),
		),
	}

	// Empty sidebar with only dividers
	sidebarE := sidebarModel{
		directories: fullDirSlice(0),
	}

	testCases := []struct {
		name               string
		sidebar            sidebarModel
		mainPanelHeight    int
		endIndex           int
		expectedFirstIndex int
		explanation        string
	}{
		{
			name:               "Basic calculation from end index",
			sidebar:            sidebarA,
			mainPanelHeight:    10,
			endIndex:           10,
			expectedFirstIndex: 6,
			explanation:        "3(InitialHeight) + 4 (6-9 homedirs) + 3 (10-pinned divider)",
		},
		{
			name:               "Small panel height",
			sidebar:            sidebarA,
			mainPanelHeight:    5,
			endIndex:           15,
			expectedFirstIndex: 14,
			explanation:        "3(InitialHeight) + 2(14-15 pinned dirs)",
		},
		{
			name:               "End index near beginning",
			sidebar:            sidebarA,
			mainPanelHeight:    20,
			endIndex:           3,
			expectedFirstIndex: 0,
			explanation:        "When end index is near beginning, first index should be 0",
		},
		{
			name:               "End index at disk divider",
			sidebar:            sidebarA,
			mainPanelHeight:    15,
			endIndex:           21, // Disk divider in sidebar_a
			expectedFirstIndex: 12,
			explanation:        "3(InitialHeight) + 9(12-20 pinned dirs) + 3(21-disk divider)",
		},
		{
			name:               "Very large panel height showing all items",
			sidebar:            sidebarA,
			mainPanelHeight:    100,
			endIndex:           31, // Last disk dir in sidebar_a
			expectedFirstIndex: 0,
			explanation:        "Large panel should show all directories from start",
		},
		{
			name:               "Asymetric sidebar with few directories",
			sidebar:            sidebarB,
			mainPanelHeight:    12,
			endIndex:           4, // Last disk dir in sidebar_b
			expectedFirstIndex: 0,
			explanation:        "Small sidebar should fit in panel height",
		},
		{
			name:               "No home directories case",
			sidebar:            sidebarC,
			mainPanelHeight:    10,
			endIndex:           6, // Disk dir in sidebar_c
			expectedFirstIndex: 2, // Pinned divider
			explanation:        "3(InitialHeight) + 4(2-5 pinned dirs) + 3(6-disk divider)",
		},
		{
			name:               "Only disk directories case",
			sidebar:            sidebarD,
			mainPanelHeight:    8,
			endIndex:           4, // Last disk dir
			expectedFirstIndex: 2, // Disk divider
			explanation:        "3(InitialHeight) + 3(2-4 disk dirs)",
		},
		{
			name:               "Empty sidebar case",
			sidebar:            sidebarE,
			mainPanelHeight:    10,
			endIndex:           1, // Disk divider
			expectedFirstIndex: 0, // Pinned divider
			explanation:        "Empty sidebar should show all dividers",
		},
		{
			name:               "End index at the start",
			sidebar:            sidebarA,
			mainPanelHeight:    5,
			endIndex:           0,
			expectedFirstIndex: 0,
			explanation:        "When end index is at start, first index should be the same",
		},
		{
			name:               "End index out of bounds",
			sidebar:            sidebarA,
			mainPanelHeight:    20,
			endIndex:           32, // Out of bounds for sidebar_a
			expectedFirstIndex: 33, // endIndex + 1
			explanation:        "When end index is out of bounds, should return endIndex+1",
		},
		{
			name:               "Very small panel height",
			sidebar:            sidebarA,
			mainPanelHeight:    2, // Too small to fit anything
			endIndex:           10,
			expectedFirstIndex: 11,
			explanation:        "With panel height less than initialHeight, first index is invalid",
		},
		{
			name:               "Panel height exactly matches divider",
			sidebar:            sidebarA,
			mainPanelHeight:    6,  // Just enough for initialHeight + divider
			endIndex:           10, // Pinned divider
			expectedFirstIndex: 10,
			explanation:        "When panel height only fits the divider, start index should be the same",
		},
		{
			name:               "Boundary case between directory types",
			sidebar:            sidebarA,
			mainPanelHeight:    7,
			endIndex:           11, // First pinned dir
			expectedFirstIndex: 10, // Pinned divider
			explanation:        "3(InitialHeight) + 3(10-pinned divider) + 1(11-pinned dir)",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sidebar.firstRenderedIndex(tt.mainPanelHeight, tt.endIndex)
			assert.Equal(t, tt.expectedFirstIndex, result,
				"firstRenderedIndex failed: %s", tt.explanation)
		})
	}
}

func Test_updateRenderIndex(t *testing.T) {
	testCases := []struct {
		name                string
		sidebar             sidebarModel
		mainPanelHeight     int
		initialRenderIndex  int
		initialCursor       int
		expectedRenderIndex int
		explanation         string
	}{
		{
			name: "Case I: Cursor moved above render range",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 10, // Start rendering from pinned divider
				cursor:      5,  // Cursor moved to home directory
			},
			mainPanelHeight:     15,
			expectedRenderIndex: 5,
			explanation:         "When cursor moves above render range, renderIndex should be set to cursor",
		},
		{
			name: "Case II: Cursor within render range",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 5, // Start rendering from a home directory
				cursor:      8, // Cursor within visible range
			},
			mainPanelHeight:     15,
			expectedRenderIndex: 5, // No change expected
			explanation:         "When cursor is within render range, renderIndex should not change",
		},
		{
			name: "Case III: Cursor moved below render range",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 0,  // Start rendering from beginning
				cursor:      20, // Cursor moved to a pinned directory outside visible range
			},
			mainPanelHeight:     10,
			expectedRenderIndex: 14, // Should adjust to make cursor visible
			// 3(Initial height) + 7(14-20 pinned dirs)
			explanation: "When cursor moves below render range, renderIndex should adjust to make cursor visible",
		},
		{
			name: "Edge case: Small panel with cursor at end",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 0,
				cursor:      31, // Last disk directory
			},
			mainPanelHeight:     5,
			expectedRenderIndex: 30, // Should show only the last couple items
			explanation:         "With small panel and cursor at end, should adjust renderIndex to show cursor",
		},
		{
			name: "Edge case: Large panel showing everything",
			sidebar: sidebarModel{
				directories: formDirctorySlice(dirSlice(1), nil, dirSlice(5)),
				renderIndex: 2,
				cursor:      4,
			},
			mainPanelHeight:     50, // Large enough to show all directories
			expectedRenderIndex: 2,  // No change needed as everything is visible
			explanation:         "With large panel showing all items, renderIndex should remain unchanged",
		},
		{
			name: "Edge case: Empty sidebar",
			sidebar: sidebarModel{
				directories: fullDirSlice(0),
				renderIndex: 0,
				cursor:      1,
			},
			mainPanelHeight:     10,
			expectedRenderIndex: 0, // No change needed for empty sidebar
			explanation:         "With empty sidebar, renderIndex should remain at 0",
		},
		{
			name: "Case I and III overlap: Cursor exactly at current renderIndex",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 15,
				cursor:      15,
			},
			mainPanelHeight:     10,
			expectedRenderIndex: 15, // No change needed, Case I takes precedence
			explanation:         "When cursor is exactly at renderIndex, Case I takes precedence and renderIndex remains unchanged",
		},
		{
			name: "Boundary case: Cursor at edge of visible range",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 5,
				cursor:      9, // Just at the edge of what's visible
			},
			mainPanelHeight:     8,
			expectedRenderIndex: 5, // Still visible, no change needed
			explanation:         "When cursor is at the edge of visible range, renderIndex should not change",
		},
		{
			name: "Boundary case: Cursor just beyond visible range",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 5,
				cursor:      11, // Just beyond visible range
			},
			mainPanelHeight:     10,
			expectedRenderIndex: 7, // Adjust to make cursor visible
			explanation:         "When cursor is just beyond visible range, renderIndex should adjust",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the sidebar to avoid modifying the original
			sidebar := tt.sidebar

			// Update render index
			sidebar.updateRenderIndex(tt.mainPanelHeight)

			// Check the result
			assert.Equal(t, tt.expectedRenderIndex, sidebar.renderIndex,
				"updateRenderIndex failed: %s", tt.explanation)
		})
	}
}

func Test_listUp(t *testing.T) {
	testCases := []struct {
		name                string
		sidebar             sidebarModel
		mainPanelHeight     int
		expectedCursor      int
		expectedRenderIndex int
		explanation         string
	}{
		{
			name: "Basic cursor movement from middle position",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 5,
				cursor:      5, // Starting from a home directory
			},
			mainPanelHeight:     15,
			expectedCursor:      4, // Should move up one position
			expectedRenderIndex: 4, // Render index should follow cursor
			explanation:         "When cursor is in the middle, it should move up one position",
		},
		{
			name: "Skip divider when moving up",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 8,
				cursor:      11, // Position just after pinned divider
			},
			mainPanelHeight:     10,
			expectedCursor:      9, // Should skip divider (10) and move to home dir (9)
			expectedRenderIndex: 8,
			explanation:         "When moving up to a divider, cursor should skip it and move to previous item",
		},
		{
			name: "Wrap around from top to bottom",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 0,
				cursor:      0, // At the very top
			},
			mainPanelHeight:     10,
			expectedCursor:      31, // Should wrap to last directory (index 31)
			expectedRenderIndex: 25, // Should adjust render to show cursor
			// 3(Initial Height) + 7(25-31 disk dirs)
			explanation: "When at the top, cursor should wrap to the bottom",
		},
		{
			name: "Skip multiple consecutive dividers",
			sidebar: sidebarModel{
				// Create a sidebar with consecutive dividers for testing
				directories: formDirctorySlice(dirSlice(5), nil, dirSlice(5)),
				renderIndex: 5,
				cursor:      7, // Position after consecutive dividers
			},
			mainPanelHeight:     10,
			expectedCursor:      4, // Should skip all dividers and move to item before dividers
			expectedRenderIndex: 4, // Should adjust render index accordingly
			explanation:         "When encountering multiple consecutive dividers, cursor should skip all of them",
		},
		{
			name: "No actual directories case",
			sidebar: sidebarModel{
				directories: fullDirSlice(0), // Empty sidebar with just dividers
				renderIndex: 0,
				cursor:      0,
			},
			mainPanelHeight:     10,
			expectedCursor:      0, // Should remain unchanged
			expectedRenderIndex: 0, // Should remain unchanged
			explanation:         "When there are no actual directories, cursor should not move",
		},
		{
			name: "Large panel showing all directories",
			sidebar: sidebarModel{
				directories: formDirctorySlice(dirSlice(2), dirSlice(2), dirSlice(2)),
				renderIndex: 0,
				cursor:      3, // Some directory in the middle
			},
			mainPanelHeight:     50, // Large enough to show all directories
			expectedCursor:      1,  // Should move up one position
			expectedRenderIndex: 0,  // No change needed as everything is visible
			explanation:         "With large panel showing all items, cursor should move up and renderIndex remain unchanged",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the sidebar to avoid modifying the original
			sidebar := tt.sidebar

			// Call the function to test
			sidebar.listUp(tt.mainPanelHeight)

			// Check the results
			assert.Equal(t, tt.expectedCursor, sidebar.cursor,
				"listUp cursor position: %s", tt.explanation)
			assert.Equal(t, tt.expectedRenderIndex, sidebar.renderIndex,
				"listUp render index: %s", tt.explanation)
		})
	}
}

func Test_listDown(t *testing.T) {
	testCases := []struct {
		name                string
		sidebar             sidebarModel
		mainPanelHeight     int
		expectedCursor      int
		expectedRenderIndex int
		explanation         string
	}{
		{
			name: "Basic cursor movement from middle position",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 5,
				cursor:      5, // Starting from a home directory
			},
			mainPanelHeight:     15,
			expectedCursor:      6, // Should move down one position
			expectedRenderIndex: 5, // Render index should remain the same as cursor is still visible
			explanation:         "When cursor is in the middle, it should move down one position",
		},
		{
			name: "Skip divider when moving down",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 8,
				cursor:      9, // Position just before pinned divider
			},
			mainPanelHeight:     10,
			expectedCursor:      11, // Should skip divider (10) and move to pinned dir (11)
			expectedRenderIndex: 8,  // Should adjust render index to keep cursor visible
			explanation:         "When moving down to a divider, cursor should skip it and move to next item",
		},
		{
			name: "Wrap around from bottom to top",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 26,
				cursor:      31, // At the very bottom
			},
			mainPanelHeight:     10,
			expectedCursor:      0, // Should wrap to first directory (index 0)
			expectedRenderIndex: 0, // Should adjust render to show cursor
			explanation:         "When at the bottom, cursor should wrap to the top",
		},
		{
			name: "Skip multiple consecutive dividers",
			sidebar: sidebarModel{
				// Create a sidebar with consecutive dividers for testing
				directories: formDirctorySlice(dirSlice(5), nil, dirSlice(5)),
				renderIndex: 0,
				cursor:      4, // Position before consecutive dividers
			},
			mainPanelHeight:     10,
			expectedCursor:      7, // Should skip all dividers and move to item after dividers
			expectedRenderIndex: 5, // Should adjust render index accordingly
			// 3 (Initial Height) 6(5,6 - pinned and disk divider), 1 (7-Disk dir)
			explanation: "When encountering multiple consecutive dividers, cursor should skip all of them",
		},
		{
			name: "No actual directories case",
			sidebar: sidebarModel{
				directories: fullDirSlice(0), // Empty sidebar with just dividers
				renderIndex: 0,
				cursor:      0,
			},
			mainPanelHeight:     10,
			expectedCursor:      0, // Should remain unchanged
			expectedRenderIndex: 0, // Should remain unchanged
			explanation:         "When there are no actual directories, cursor should not move",
		},
		{
			name: "Move down from home to pinned section",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 6,
				cursor:      9, // Last home directory
			},
			mainPanelHeight:     10,
			expectedCursor:      11, // Should move to first pinned directory
			expectedRenderIndex: 7,  // Should adjust render index to show cursor
			explanation:         "When moving down from last home directory, cursor should skip divider and go to first pinned directory",
		},
		{
			name: "Large panel showing all directories",
			sidebar: sidebarModel{
				directories: formDirctorySlice(dirSlice(2), dirSlice(2), dirSlice(2)),
				renderIndex: 0,
				cursor:      3, // Some directory in the middle
			},
			mainPanelHeight:     50, // Large enough to show all directories
			expectedCursor:      4,  // Should move down one position
			expectedRenderIndex: 0,  // No change needed as everything is visible
			explanation:         "With large panel showing all items, cursor should move down and renderIndex remain unchanged",
		},
		{
			name: "Cursor at the end of visible range",
			sidebar: sidebarModel{
				directories: fullDirSlice(10),
				renderIndex: 5,
				cursor:      14, // At the end of visible range
			},
			mainPanelHeight:     15,
			expectedCursor:      15, // Should move down one position
			expectedRenderIndex: 6,  // Should increase render index to keep cursor visible
			explanation:         "When cursor is at the end of visible range, moving down should adjust renderIndex",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of the sidebar to avoid modifying the original
			sidebar := tt.sidebar

			// Call the function to test
			sidebar.listDown(tt.mainPanelHeight)

			// Check the results
			assert.Equal(t, tt.expectedCursor, sidebar.cursor,
				"listDown cursor position: %s", tt.explanation)
			assert.Equal(t, tt.expectedRenderIndex, sidebar.renderIndex,
				"listDown render index: %s", tt.explanation)
		})
	}
}
