package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

func testModelWithElemCount(cursor int, renderIndex int, height int, elemCount int) Model {
	return testModel(cursor, renderIndex, height, BrowserMode, make([]Element, elemCount))
}

func testModel(cursor int, renderIndex int, height int, mode PanelMode,
	elements []Element) Model {
	return Model{
		Element:     elements,
		Cursor:      cursor,
		RenderIndex: renderIndex,
		height:      height,
		selected:    make(map[string]int),
		PanelMode:   mode,
	}
}

func Test_filePanelUpDown(t *testing.T) {
	testdata := []struct {
		name           string
		panel          Model
		listDown       bool
		expectedCursor int
		expectedRender int
	}{
		{
			name:           "Down movement within renderable range",
			panel:          testModelWithElemCount(0, 0, 12, 10),
			listDown:       true,
			expectedCursor: 1,
			expectedRender: 0,
		},
		{
			name:           "Down movement when cursor is at bottom",
			panel:          testModelWithElemCount(6, 0, 12, 10),
			listDown:       true,
			expectedCursor: 7,
			expectedRender: 1,
		},
		{
			name:           "Down movement causing wrap to top",
			panel:          testModelWithElemCount(9, 3, 12, 10),
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
		},
		{
			name:           "Up movement within renderable range",
			panel:          testModelWithElemCount(2, 0, 12, 10),
			listDown:       false,
			expectedCursor: 1,
			expectedRender: 0,
		},
		{
			name:           "Up movement when cursor is at top",
			panel:          testModelWithElemCount(3, 3, 12, 10),
			listDown:       false,
			expectedCursor: 2,
			expectedRender: 2,
		},
		{
			name:           "Up movement causing wrap to bottom",
			panel:          testModelWithElemCount(0, 0, 12, 10),
			listDown:       false,
			expectedCursor: 9,
			expectedRender: 3,
		},
		{
			name:           "Down movement on empty panel",
			panel:          testModelWithElemCount(0, 0, 12, 0),
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
		},
		{
			name:           "Up movement on empty panel",
			panel:          testModelWithElemCount(0, 0, 12, 0),
			listDown:       false,
			expectedCursor: 0,
			expectedRender: 0,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if tt.listDown {
				tt.panel.ListDown()
			} else {
				tt.panel.ListUp()
			}
			assert.Equal(t, tt.expectedCursor, tt.panel.Cursor)
			assert.Equal(t, tt.expectedRender, tt.panel.RenderIndex)
		})
	}
}

func TestPgUpDown(t *testing.T) {
	testdata := []struct {
		name           string
		panel          Model
		pageDown       bool
		expectedCursor int
		expectedRender int
	}{
		{
			name:           "Page down with full page of items",
			panel:          testModelWithElemCount(0, 0, 12, 20),
			pageDown:       true,
			expectedCursor: 7,
			expectedRender: 1,
		},
		{
			name:           "Page down near end wraps to start",
			panel:          testModelWithElemCount(18, 12, 12, 20),
			pageDown:       true,
			expectedCursor: 5, // (18 + 7) % 20 = 5
			expectedRender: 5,
		},
		{
			name:           "Page up from middle",
			panel:          testModelWithElemCount(10, 4, 12, 20),
			pageDown:       false,
			expectedCursor: 3, // 10 - 7 = 3
			expectedRender: 3,
		},
		{
			name:           "Page up near beginning wraps to end",
			panel:          testModelWithElemCount(2, 0, 12, 20),
			pageDown:       false,
			expectedCursor: 15, // (2 - 7 + 20) % 20 = 15
			expectedRender: 9,
		},
		{
			name:           "Page navigation with small element count",
			panel:          testModelWithElemCount(0, 0, 12, 5),
			pageDown:       true,
			expectedCursor: 2, // (0 + 7) % 5 = 2
			expectedRender: 0,
		},
		{
			name:           "Page down on empty panel",
			panel:          testModelWithElemCount(0, 0, 12, 0),
			pageDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pageDown {
				tt.panel.PgDown()
			} else {
				tt.panel.PgUp()
			}
			assert.Equal(t, tt.expectedCursor, tt.panel.Cursor)
			assert.Equal(t, tt.expectedRender, tt.panel.RenderIndex)
		})
	}
}

func TestItemSelectUpDown(t *testing.T) {
	testdata := []struct {
		name             string
		panel            Model
		panelToSelect    []string
		selectDown       bool
		expectedCursor   int
		expectedRender   int
		expectedSelected map[string]int
	}{
		{
			name: "Select and move down",
			panel: testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
			}),
			panelToSelect:    []string{},
			selectDown:       true,
			expectedCursor:   1,
			expectedRender:   0,
			expectedSelected: map[string]int{"/tmp/file1.txt": 1},
		},
		{
			name: "Select and move up",
			panel: testModel(2, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
			}),
			panelToSelect:    []string{},
			selectDown:       false,
			expectedCursor:   1,
			expectedRender:   0,
			expectedSelected: map[string]int{"/tmp/file3.txt": 1},
		},
		{
			name: "Deselect already selected item",
			panel: testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
			}),
			panelToSelect:    []string{"/tmp/file1.txt"},
			selectDown:       true,
			expectedCursor:   1,
			expectedRender:   0,
			expectedSelected: map[string]int{},
		},
		{
			name: "Selection at boundary with wrap",
			panel: testModel(1, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
			}),
			panelToSelect:    []string{},
			selectDown:       true,
			expectedCursor:   0, // wraps to beginning
			expectedRender:   0,
			expectedSelected: map[string]int{"/tmp/file2.txt": 1},
		},
		{
			name: "Selection persistence across moves",
			panel: testModel(1, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
				{Name: "file3.txt", Location: "/tmp/file3.txt"},
			}),
			panelToSelect:    []string{"/tmp/file1.txt"},
			selectDown:       true,
			expectedCursor:   2,
			expectedRender:   0,
			expectedSelected: map[string]int{"/tmp/file1.txt": 1, "/tmp/file2.txt": 2},
		},
		{
			name:             "Empty panel selection",
			panel:            testModel(0, 0, 12, SelectMode, []Element{}),
			panelToSelect:    []string{},
			selectDown:       true,
			expectedCursor:   0,
			expectedRender:   0,
			expectedSelected: map[string]int{},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.panel.SetSelectedAll(tt.panelToSelect)

			if tt.selectDown {
				tt.panel.ItemSelectDown()
			} else {
				tt.panel.ItemSelectUp()
			}
			assert.Equal(t, tt.expectedCursor, tt.panel.Cursor)
			assert.Equal(t, tt.expectedRender, tt.panel.RenderIndex)
			assert.Equal(t, tt.expectedSelected, tt.panel.selected)
		})
	}
}

func TestScrollToCursor(t *testing.T) {
	testdata := []struct {
		name           string
		panel          Model
		cursorPos      int
		expectedCursor int
		expectedRender int
	}{
		{
			name:           "Jump to visible cursor no change",
			panel:          testModelWithElemCount(5, 3, 12, 20),
			cursorPos:      4,
			expectedCursor: 4,
			expectedRender: 3,
		},
		{
			name:           "Jump above view",
			panel:          testModelWithElemCount(10, 5, 12, 20),
			cursorPos:      2,
			expectedCursor: 2,
			expectedRender: 2,
		},
		{
			name:           "Jump below view",
			panel:          testModelWithElemCount(5, 0, 12, 20),
			cursorPos:      15,
			expectedCursor: 15,
			expectedRender: 9, // 15 - 7 + 1
		},
		{
			name:           "Jump above view with empty space",
			panel:          testModelWithElemCount(19, 18, 12, 20),
			cursorPos:      17,
			expectedCursor: 17,
			expectedRender: 17,
		},
		{
			name:           "Invalid cursor negative",
			panel:          testModelWithElemCount(5, 2, 12, 10),
			cursorPos:      -1,
			expectedCursor: 5, // unchanged
			expectedRender: 2, // unchanged
		},
		{
			name:           "Invalid cursor beyond count",
			panel:          testModelWithElemCount(5, 2, 12, 10),
			cursorPos:      15,
			expectedCursor: 5, // unchanged
			expectedRender: 2, // unchanged
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.panel.scrollToCursor(tt.cursorPos)
			assert.Equal(t, tt.expectedCursor, tt.panel.Cursor)
			assert.Equal(t, tt.expectedRender, tt.panel.RenderIndex)
		})
	}
}

func TestApplyTargetFileCursor(t *testing.T) {
	panel := testModel(0, 0, 8, BrowserMode, []Element{
		{Name: "file1.txt", Location: "/tmp/file1.txt"},
		{Name: "file2.txt", Location: "/tmp/file2.txt"},
		{Name: "file3.txt", Location: "/tmp/file3.txt"},
		{Name: "file4.txt", Location: "/tmp/file4.txt"},
		{Name: "target.txt", Location: "/tmp/target.txt"},
		{Name: "file6.txt", Location: "/tmp/file6.txt"},
	})
	panel.TargetFile = "target.txt"

	expCursor := 4
	expRender := 2

	panel.applyTargetFileCursor()
	assert.Equal(t, expCursor, panel.Cursor)
	assert.Equal(t, expRender, panel.RenderIndex)
	assert.Empty(t, panel.TargetFile)

	// Shouldn't do anything
	panel.applyTargetFileCursor()
	assert.Equal(t, expCursor, panel.Cursor)
	assert.Equal(t, expRender, panel.RenderIndex)
}

func TestPageScrollSizeConfig(t *testing.T) {
	originalPageScrollSize := common.Config.PageScrollSize
	defer func() {
		common.Config.PageScrollSize = originalPageScrollSize
	}()

	tests := []struct {
		name           string
		pageScrollSize int
		totalElements  int
		initialCursor  int
		panelHeight    int
		expectedCursor int
		pgUp           bool
	}{
		{
			name:           "Default full page scroll (PageScrollSize = 0)",
			pageScrollSize: 0,
			totalElements:  30,
			initialCursor:  0,
			panelHeight:    10, // panelElementHeight = 10 - 3 = 7
			expectedCursor: 7,  // Should move by 7 (full page)
		},
		{
			name:           "Custom scroll size 5",
			pageScrollSize: 5,
			totalElements:  30,
			initialCursor:  0,
			panelHeight:    10,
			expectedCursor: 5, // Should move by 5
		},
		{
			name:           "Custom scroll size 10",
			pageScrollSize: 10,
			totalElements:  30,
			initialCursor:  0,
			panelHeight:    10,
			expectedCursor: 10, // Should move by 10
		},
		{
			name:           "PgUp with custom scroll size",
			pageScrollSize: 3,
			totalElements:  30,
			initialCursor:  10,
			panelHeight:    10,
			expectedCursor: 7, // 10 - 3 = 7
			pgUp:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			common.Config.PageScrollSize = tt.pageScrollSize

			// Create model with elements
			m := testModelWithElemCount(tt.initialCursor, 0, tt.panelHeight+2, tt.totalElements)
			if tt.pgUp {
				m.PgUp()
			} else {
				m.PgDown()
			}

			assert.Equal(t, tt.expectedCursor, m.Cursor,
				"Cursor position should match expected after PgUp/PgDown")
		})
	}
}
