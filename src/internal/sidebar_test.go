package internal

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

var home_dir_list_1 = []directory{
	{name: "Home1", location: "/a/b"},
	{name: "Home2", location: "/a/b"},
}

var pinned_dir_list_1 = []directory{
	{name: "Pinned1", location: "/a/b"},
	{name: "Pinned2", location: "/a/b"},
}

var disk_dir_list_1 = []directory{
	{name: "Disk1", location: "/a/b"},
	{name: "Disk2", location: "/a/b"},
}

var directory_list_1 = formDirctorySlice(home_dir_list_1, pinned_dir_list_1, disk_dir_list_1)

// Its still invalid
var sidebar1_empty_invalid = sidebarModel{}

var sidebar2_empty = sidebarModel{
	directories: []directory{pinnedDividerDir, diskDividerDir},
	renderIndex: 0,
	cursor:      0,
}

var sidebar3_only_pinned = sidebarModel{
	directories: formDirctorySlice(nil, pinned_dir_list_1, nil),
}

var sidebar4_all = sidebarModel{
	directories: directory_list_1,
	renderIndex: 0,
	cursor:      0,
}

var sidebar5_invalid = sidebarModel{
	directories: directory_list_1,
	renderIndex: 0,
	cursor:      len(directory_list_1) + 1,
}

var sidebar6_invalid = sidebarModel{
	directories: formDirctorySlice(home_dir_list_1, nil, nil),
	renderIndex: 0,
	// Cursor now points to pinned divider, which is invalid
	cursor: len(home_dir_list_1),
}

var sidebar7_only_disks = sidebarModel{
	directories: formDirctorySlice(nil, nil, disk_dir_list_1),
}

func generate_dir_list(count int) []directory {
	res := make([]directory, count)
	for i := 0; i < count; i++ {
		res[i] = directory{name: "Dir" + strconv.Itoa(i), location: "/a/" + strconv.Itoa(i)}
	}
	return res
}

func Test_noActualDir(t *testing.T) {
	assert.True(t, sidebar1_empty_invalid.noActualDir(), "Empty sidebar should have no actual directories")
	assert.True(t, sidebar2_empty.noActualDir(), "Empty sidebar should have no actual directories")
	assert.False(t, sidebar3_only_pinned.noActualDir(), "Non-Empty Sidebar should have actual directories")
	assert.False(t, sidebar4_all.noActualDir(), "Non-Empty Sidebar should have actual directories")
}

func Test_isCursorInvalid(t *testing.T) {
	assert.True(t, sidebar1_empty_invalid.isCursorInvalid(), "Expected cursor to be invalid")
	assert.True(t, sidebar5_invalid.isCursorInvalid(), "Expected cursor to be invalid")
	assert.True(t, sidebar6_invalid.isCursorInvalid(), "Expected cursor to be invalid")
	assert.False(t, sidebar4_all.isCursorInvalid(), "Expected cursor to be valid")
}

func Test_resetCursor(t *testing.T) {
	data := []struct {
		curSideBar        sidebarModel
		expectedCursorPos int
	}{
		{
			curSideBar:        sidebar3_only_pinned,
			expectedCursorPos: 1, // After pinned divider
		},
		{
			curSideBar:        sidebar4_all,
			expectedCursorPos: 0, // First home
		},
		{
			curSideBar:        sidebar7_only_disks,
			expectedCursorPos: 2, // After pinned and dist divider
		},
		{
			curSideBar:        sidebar2_empty,
			expectedCursorPos: 0, // Empty sidebar, cursor should reset to 0
		},
	}

	for _, tt := range data {
		tt.curSideBar.resetCursor()
		assert.Equal(t, tt.expectedCursorPos, tt.curSideBar.cursor)
	}
}

// Todo : we can add more tests
func Test_renderIndex(t *testing.T) {

	sidebar_a := sidebarModel{
		directories: formDirctorySlice(
			generate_dir_list(10), generate_dir_list(10), generate_dir_list(10),
		),
	}
	sidebar_b := sidebarModel{
		directories: formDirctorySlice(
			generate_dir_list(1), nil, generate_dir_list(5),
		),
	}

	lastRenderedIndex_data := []struct {
		curSideBar        sidebarModel
		mainPanelHeight   int
		startIndex        int
		expectedLastIndex int
	}{
		// 3(initialHeight) , 7 (0-6 home dirs)
		{sidebar_a, 10, 0, 6},

		// 3(initialHeight) , 10 (0-9 home dirs), 3 (10-pinned divider)
		// 4(11-14 pinned dirs)
		{sidebar_a, 20, 0, 14},

		// 3(initialHeight) , 10 (11-20 pinned dirs), 3 (21-disk divider)
		// 4(22-25 disk dirs)
		{sidebar_a, 20, 11, 25},

		// Last dir - 31
		{sidebar_a, 100, 11, 31},

		// startIndex is more then len(directories), and startIndex-1 is returned
		{sidebar_a, 100, 32, 31},

		// 3(initialHeight), 1 (0-homedir), 6(1-pinned divider, 2-diskdivider),
		// 2(3-4 diskdirs)
		{sidebar_b, 12, 0, 4},
	}

	for _, tt := range lastRenderedIndex_data {
		assert.Equal(t, tt.curSideBar.lastRenderedIndex(tt.mainPanelHeight, tt.startIndex),
			tt.expectedLastIndex)
	}

	firstRenderedIndex_data := []struct {
		curSideBar         sidebarModel
		mainPanelHeight    int
		endIndex           int
		expectedFirstIndex int
	}{
		// 3(InitialHeight), 4(6-9 homedirs), 3(10-pinned divider)
		{sidebar_a, 10, 10, 6},
	}

	for _, tt := range firstRenderedIndex_data {
		assert.Equal(t, tt.curSideBar.firstRenderedIndex(tt.mainPanelHeight, tt.endIndex),
			tt.expectedFirstIndex)
	}
}
