package internal

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func dirSlice(count int) []directory {
	res := make([]directory, count)
	for i := 0; i < count; i++ {
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
	sidebar_a := sidebarModel{
		directories: formDirctorySlice(
			dirSlice(10), dirSlice(10), dirSlice(10),
		),
	}
	sidebar_b := sidebarModel{
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
			sidebar:           sidebar_a,
			mainPanelHeight:   10,
			startIndex:        0,
			expectedLastIndex: 6,
			explanation:       "3(initialHeight) + 7 (0-6 home dirs)",
		},
		{
			name:              "Medium viewport showing home and some pinned",
			sidebar:           sidebar_a,
			mainPanelHeight:   20,
			startIndex:        0,
			expectedLastIndex: 14,
			explanation:       "3(initialHeight) + 10 (0-9 home dirs) + 3 (10-pinned divider) + 4 (11-14 pinned dirs)",
		},
		{
			name:              "Medium viewport starting from pinned dirs",
			sidebar:           sidebar_a,
			mainPanelHeight:   20,
			startIndex:        11,
			expectedLastIndex: 25,
			explanation:       "3(initialHeight) + 10 (11-20 pinned dirs) + 3 (21-disk divider) + 4 (22-25 disk dirs)",
		},
		{
			name:              "Large viewport showing all directories",
			sidebar:           sidebar_a,
			mainPanelHeight:   100,
			startIndex:        11,
			expectedLastIndex: 31,
			explanation:       "Last dir index is 31",
		},
		{
			name:              "Start index beyond directory count",
			sidebar:           sidebar_a,
			mainPanelHeight:   100,
			startIndex:        32,
			expectedLastIndex: 31,
			explanation:       "When startIndex > len(directories), return last valid index",
		},
		{
			name:              "Asymmetric directory distribution",
			sidebar:           sidebar_b,
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
	sidebar_a := sidebarModel{
		directories: formDirctorySlice(
			dirSlice(10), dirSlice(10), dirSlice(10),
		),
	}
	sidebar_b := sidebarModel{
		directories: formDirctorySlice(
			dirSlice(1), nil, dirSlice(5),
		),
	}
	sidebar_c := sidebarModel{
		directories: formDirctorySlice(
			nil, dirSlice(5), dirSlice(5),
		),
	}
	sidebar_d := sidebarModel{
		directories: formDirctorySlice(
			nil, nil, dirSlice(3),
		),
	}
	
	// Empty sidebar with only dividers
	sidebar_e := sidebarModel{
		directories: []directory{pinnedDividerDir, diskDividerDir},
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
			sidebar:            sidebar_a,
			mainPanelHeight:    10,
			endIndex:           10,
			expectedFirstIndex: 6,
			explanation:        "3(InitialHeight) + 4 (6-9 homedirs) + 3 (10-pinned divider)",
		},
		{
			name:               "Small panel height",
			sidebar:            sidebar_a,
			mainPanelHeight:    5,
			endIndex:           15,
			expectedFirstIndex: 14,
			explanation:        "3(InitialHeight) + 2(14-15 pinned dirs)",
		},
		{
			name:               "End index near beginning",
			sidebar:            sidebar_a,
			mainPanelHeight:    20,
			endIndex:           3,
			expectedFirstIndex: 0,
			explanation:        "When end index is near beginning, first index should be 0",
		},
		{
			name:               "End index at disk divider",
			sidebar:            sidebar_a,
			mainPanelHeight:    15,
			endIndex:           21, // Disk divider in sidebar_a
			expectedFirstIndex: 12,
			explanation:        "3(InitialHeight) + 9(12-20 pinned dirs) + 3(21-disk divider)",
		},
		{
			name:               "Very large panel height showing all items",
			sidebar:            sidebar_a,
			mainPanelHeight:    100,
			endIndex:           31, // Last disk dir in sidebar_a
			expectedFirstIndex: 0,
			explanation:        "Large panel should show all directories from start",
		},
		{
			name:               "Assymetric sidebar with few directories",
			sidebar:            sidebar_b,
			mainPanelHeight:    12,
			endIndex:           4, // Last disk dir in sidebar_b
			expectedFirstIndex: 0,
			explanation:        "Small sidebar should fit in panel height",
		},
		{
			name:               "No home directories case",
			sidebar:            sidebar_c,
			mainPanelHeight:    10,
			endIndex:           6, // Disk dir in sidebar_c
			expectedFirstIndex: 2, // Pinned divider
			explanation:        "3(InitialHeight) + 4(2-5 pinned dirs) + 3(6-disk divider)",
		},
		{
			name:               "Only disk directories case",
			sidebar:            sidebar_d,
			mainPanelHeight:    8,
			endIndex:           4, // Last disk dir
			expectedFirstIndex: 2, // Disk divider
			explanation:        "3(InitialHeight) + 3(2-4 disk dirs)",
		},
		{
			name:               "Empty sidebar case",
			sidebar:            sidebar_e,
			mainPanelHeight:    10,
			endIndex:           1, // Disk divider
			expectedFirstIndex: 0, // Pinned divider
			explanation:        "Empty sidebar should show all dividers",
		},
		{
			name:               "End index at the start",
			sidebar:            sidebar_a,
			mainPanelHeight:    5,
			endIndex:           0,
			expectedFirstIndex: 0,
			explanation:        "When end index is at start, first index should be the same",
		},
		{
			name:               "End index out of bounds",
			sidebar:            sidebar_a,
			mainPanelHeight:    20,
			endIndex:           32, // Out of bounds for sidebar_a
			expectedFirstIndex: 33, // endIndex + 1
			explanation:        "When end index is out of bounds, should return endIndex+1",
		},
		{
			name:               "Very small panel height",
			sidebar:            sidebar_a,
			mainPanelHeight:    2, // Too small to fit anything
			endIndex:           10,
			expectedFirstIndex: 11,
			explanation:        "With panel height less than initialHeight, first index is invalid",
		},
		{
			name:               "Panel height exactly matches divider",
			sidebar:            sidebar_a,
			mainPanelHeight:    6, // Just enough for initialHeight + divider
			endIndex:           10, // Pinned divider
			expectedFirstIndex: 10,
			explanation:        "When panel height only fits the divider, start index should be the same",
		},
		{
			name:               "Boundary case between directory types",
			sidebar:            sidebar_a,
			mainPanelHeight:    7,
			endIndex:           11, // First pinned dir
			expectedFirstIndex: 10, // Pinned divider
			explanation:        "3(InitialHeight) + 3(10-pinned divider) + 1(11-pinned dir)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.sidebar.firstRenderedIndex(tc.mainPanelHeight, tc.endIndex)
			assert.Equal(t, tc.expectedFirstIndex, result,
				"firstRenderedIndex failed: %s", tc.explanation)
		})
	}
}