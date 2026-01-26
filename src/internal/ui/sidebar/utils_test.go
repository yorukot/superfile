package sidebar

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testModel(cursor int, renderIndex int, height int, sections []string,
	directories []directory) Model {
	return Model{
		directories: directories,
		cursor:      cursor,
		renderIndex: renderIndex,
		height:      height,
		sections:    sections,
	}
}

func dirSlice(count int) []directory {
	res := make([]directory, count)
	for i := range count {
		res[i] = directory{Name: "Dir" + strconv.Itoa(i), Location: "/a/" + strconv.Itoa(i)}
	}
	return res
}

func fullDirSlice(count int) []directory {
	sections := []string{"home", "pinned", "disks"}
	return formDirctorySlice(dirSlice(count), dirSlice(count), dirSlice(count), sections)
}

func Test_noActualDir(t *testing.T) {
	testcases := []struct {
		name     string
		sidebar  Model
		expected bool
	}{
		{
			"Empty invalid sidebar should have no actual directories",
			Model{},
			true,
		},
		{
			"Empty sidebar should have no actual directories",
			testModel(0, 0, 10, []string{"home", "pinned", "disks"}, fullDirSlice(0)),
			true,
		},
		{
			"Non-Empty Sidebar with only pinned directories",
			testModel(0, 0, 10, []string{"home", "pinned", "disks"},
				formDirctorySlice(nil, dirSlice(10), nil, []string{"home", "pinned", "disks"})),
			false,
		},
		{
			"Non-Empty Sidebar with all directories",
			testModel(0, 0, 10, []string{"home", "pinned", "disks"}, fullDirSlice(10)),
			false,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.sidebar.NoActualDir())
		})
	}
}

func Test_isCursorInvalid(t *testing.T) {
	testcases := []struct {
		name     string
		sidebar  Model
		expected bool
	}{
		{
			"Empty invalid sidebar",
			Model{},
			true,
		},
		{
			"Cursor after all directories",
			testModel(32, 0, 10, []string{"home", "pinned", "disks"}, fullDirSlice(10)),
			true,
		},
		{
			"Curson points to pinned divider",
			testModel(10, 0, 10, []string{"home", "pinned", "disks"}, fullDirSlice(10)),
			true,
		},
		{
			"Non-Empty Sidebar with all directories",
			testModel(5, 0, 10, []string{"home", "pinned", "disks"}, fullDirSlice(10)),
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
		curSideBar        Model
		expectedCursorPos int
	}{
		{
			name: "Only Pinned directories",
			curSideBar: testModel(0, 0, 10, []string{"home", "pinned", "disks"},
				formDirctorySlice(nil, dirSlice(10), nil, []string{"home", "pinned", "disks"})),
			expectedCursorPos: 1, // After pinned divider
		},
		{
			name: "All kind of directories",
			curSideBar: testModel(0, 0, 10, []string{"home", "pinned", "disks"},
				fullDirSlice(10)),
			expectedCursorPos: 0, // First home
		},
		{
			name: "Only Disk",
			curSideBar: testModel(0, 0, 10, []string{"home", "pinned", "disks"},
				formDirctorySlice(nil, nil, dirSlice(10), []string{"home", "pinned", "disks"})),
			expectedCursorPos: 2, // After pinned and dist divider
		},
		{
			name: "Empty Sidebar",
			curSideBar: testModel(0, 0, 10, []string{"home", "pinned", "disks"},
				fullDirSlice(0)),
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

func TestSidebarSectionsVisibility(t *testing.T) {
	testcases := []struct {
		name          string
		sections      []string
		homeDirs      int
		pinnedDirs    int
		diskDirs      int
		expectedLen   int
		expectHomeDiv bool
	}{
		{
			name:        "Only one section (pinned)",
			sections:    []string{"pinned"},
			pinnedDirs:  5,
			expectedLen: 6, // divider + 5 dirs
		},
		{
			name:        "No sections",
			sections:    []string{},
			expectedLen: 0,
		},
		{
			name:          "Reordered sections (pinned, home)",
			sections:      []string{"pinned", "home"},
			homeDirs:      3,
			pinnedDirs:    3,
			expectedLen:   1 + 3 + 1 + 3, // pinned divider + 3 pinned + home divider + 3 home
			expectHomeDiv: true,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			dirs := formDirctorySlice(
				dirSlice(tt.homeDirs),
				dirSlice(tt.pinnedDirs),
				dirSlice(tt.diskDirs),
				tt.sections,
			)
			assert.Len(t, dirs, tt.expectedLen)
			if tt.expectHomeDiv {
				foundHomeDiv := false
				for _, d := range dirs {
					if d == homeDividerDir {
						foundHomeDiv = true
						break
					}
				}
				assert.True(t, foundHomeDiv, "Expected home divider to be present")
			}
		})
	}
}
