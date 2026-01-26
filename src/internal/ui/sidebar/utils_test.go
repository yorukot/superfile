package sidebar

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

func setupTestConfig() {
	common.Config.SidebarOrder = []string{"home", "pinned", "disks"}
	common.Config.SidebarShowHomeDirs = true
	common.Config.SidebarShowPinned = true
	common.Config.SidebarShowDisks = true
}

func dirSlice(count int) []directory {
	res := make([]directory, count)
	for i := range count {
		res[i] = directory{Name: "Dir" + strconv.Itoa(i), Location: "/a/" + strconv.Itoa(i)}
	}
	return res
}

func fullDirSlice(count int) []directory {
	return formDirctorySlice(dirSlice(count), dirSlice(count), dirSlice(count))
}

// TODO : Use t.Run(tt.name
// TODO : Get rid of global vars, use testdata in each test, even if there is a bit of
// duplication.
// TODO : Add tt.names

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
			Model{
				directories: fullDirSlice(0),
				renderIndex: 0,
				cursor:      0,
			},
			true,
		},
		{
			"Non-Empty Sidebar with only pinned directories",
			Model{
				directories: formDirctorySlice(nil, dirSlice(10), nil),
			},
			false,
		},
		{
			"Non-Empty Sidebar with all directories",
			Model{
				directories: fullDirSlice(10),
			},
			false,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			setupTestConfig()
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
			Model{
				directories: fullDirSlice(10),
				renderIndex: 0,
				cursor:      32,
			},
			true,
		},
		{
			"Curson points to pinned divider",
			Model{
				directories: fullDirSlice(10),
				cursor:      10,
			},
			true,
		},
		{
			"Non-Empty Sidebar with all directories",
			Model{
				directories: fullDirSlice(10),
				cursor:      5,
			},
			false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			setupTestConfig()
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
			curSideBar: Model{
				directories: formDirctorySlice(nil, dirSlice(10), nil),
			},
			expectedCursorPos: 1, // After pinned divider
		},
		{
			name: "All kind of directories",
			curSideBar: Model{
				directories: fullDirSlice(10),
			},
			expectedCursorPos: 0, // First home
		},
		{
			name: "Only Disk",
			curSideBar: Model{
				directories: formDirctorySlice(nil, nil, dirSlice(10)),
			},
			expectedCursorPos: 2, // After pinned and dist divider
		},
		{
			name: "Empty Sidebar",
			curSideBar: Model{
				directories: fullDirSlice(0),
			},
			expectedCursorPos: 0, // Empty sidebar, cursor should reset to 0
		},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			setupTestConfig()
			tt.curSideBar.resetCursor()
			assert.Equal(t, tt.expectedCursorPos, tt.curSideBar.cursor)
		})
	}
}
