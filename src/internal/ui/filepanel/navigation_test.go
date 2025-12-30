package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

func testModel(cursor int, renderIndex int, height int, elemCount int) Model {
	return Model{
		Element:     make([]Element, elemCount),
		Cursor:      cursor,
		RenderIndex: renderIndex,
		height:      height,
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
			panel:          testModel(0, 0, 12, 10),
			listDown:       true,
			expectedCursor: 1,
			expectedRender: 0,
		},
		{
			name:           "Down movement when cursor is at bottom",
			panel:          testModel(6, 0, 12, 10),
			listDown:       true,
			expectedCursor: 7,
			expectedRender: 1,
		},
		{
			name:           "Down movement causing wrap to top",
			panel:          testModel(9, 3, 12, 10),
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
		},
		{
			name:           "Up movement within renderable range",
			panel:          testModel(2, 0, 12, 10),
			listDown:       false,
			expectedCursor: 1,
			expectedRender: 0,
		},
		{
			name:           "Up movement when cursor is at top",
			panel:          testModel(3, 3, 12, 10),
			listDown:       false,
			expectedCursor: 2,
			expectedRender: 2,
		},
		{
			name:           "Up movement causing wrap to bottom",
			panel:          testModel(0, 0, 12, 10),
			listDown:       false,
			expectedCursor: 9,
			expectedRender: 3,
		},
		{
			name:           "Down movement on empty panel",
			panel:          testModel(0, 0, 12, 0),
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
		},
		{
			name:           "Up movement on empty panel",
			panel:          testModel(0, 0, 12, 0),
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
		name            string
		panel           Model
		pageDown        bool
		mainPanelHeight int
		expectedCursor  int
		expectedRender  int
	}{
		{
			name: "Page down with full page of items",
			panel: Model{
				Element:     make([]Element, 20),
				Cursor:      0,
				RenderIndex: 0,
			},
			pageDown:        true,
			mainPanelHeight: 10, // panelElementHeight = 10 - 3 = 7
			expectedCursor:  7,
			expectedRender:  1,
		},
		{
			name: "Page down near end wraps to start",
			panel: Model{
				Element:     make([]Element, 20),
				Cursor:      18,
				RenderIndex: 12,
			},
			pageDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  5, // (18 + 7) % 20 = 5
			expectedRender:  5,
		},
		{
			name: "Page up from middle",
			panel: Model{
				Element:     make([]Element, 20),
				Cursor:      10,
				RenderIndex: 4,
			},
			pageDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  3, // 10 - 7 = 3
			expectedRender:  3,
		},
		{
			name: "Page up near beginning wraps to end",
			panel: Model{
				Element:     make([]Element, 20),
				Cursor:      2,
				RenderIndex: 0,
			},
			pageDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  15, // (2 - 7 + 20) % 20 = 15
			expectedRender:  9,
		},
		{
			name: "Page navigation with small element count",
			panel: Model{
				Element:     make([]Element, 5),
				Cursor:      0,
				RenderIndex: 0,
			},
			pageDown:        true,
			mainPanelHeight: 10, // panelElementHeight = 7, but only 5 elements
			expectedCursor:  2,  // (0 + 7) % 5 = 2
			expectedRender:  0,
		},
		{
			name: "Page down on empty panel",
			panel: Model{
				Element:     make([]Element, 0),
				Cursor:      0,
				RenderIndex: 0,
			},
			pageDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.panel.SetHeight(tt.mainPanelHeight + 2)

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
		mainPanelHeight  int
		expectedCursor   int
		expectedRender   int
		expectedSelected map[string]int
	}{
		{
			name: "Select and move down",
			panel: Model{
				Element: []Element{
					{Name: "file1.txt", Location: "/tmp/file1.txt"},
					{Name: "file2.txt", Location: "/tmp/file2.txt"},
					{Name: "file3.txt", Location: "/tmp/file3.txt"},
				},
				Cursor:      0,
				RenderIndex: 0,
				PanelMode:   SelectMode,
			},
			panelToSelect:    []string{},
			selectDown:       true,
			mainPanelHeight:  10,
			expectedCursor:   1,
			expectedRender:   0,
			expectedSelected: map[string]int{"/tmp/file1.txt": 1},
		},
		{
			name: "Select and move up",
			panel: Model{
				Element: []Element{
					{Name: "file1.txt", Location: "/tmp/file1.txt"},
					{Name: "file2.txt", Location: "/tmp/file2.txt"},
					{Name: "file3.txt", Location: "/tmp/file3.txt"},
				},
				Cursor:      2,
				RenderIndex: 0,
				PanelMode:   SelectMode,
			},
			panelToSelect:    []string{},
			selectDown:       false,
			mainPanelHeight:  10,
			expectedCursor:   1,
			expectedRender:   0,
			expectedSelected: map[string]int{"/tmp/file3.txt": 1},
		},
		{
			name: "Deselect already selected item",
			panel: Model{
				Element: []Element{
					{Name: "file1.txt", Location: "/tmp/file1.txt"},
					{Name: "file2.txt", Location: "/tmp/file2.txt"},
				},
				Cursor:      0,
				RenderIndex: 0,
				PanelMode:   SelectMode,
			},
			panelToSelect:    []string{"/tmp/file1.txt"},
			selectDown:       true,
			mainPanelHeight:  10,
			expectedCursor:   1,
			expectedRender:   0,
			expectedSelected: map[string]int{},
		},
		{
			name: "Selection at boundary with wrap",
			panel: Model{
				Element: []Element{
					{Name: "file1.txt", Location: "/tmp/file1.txt"},
					{Name: "file2.txt", Location: "/tmp/file2.txt"},
				},
				Cursor:      1,
				RenderIndex: 0,
				PanelMode:   SelectMode,
			},
			panelToSelect:    []string{},
			selectDown:       true,
			mainPanelHeight:  10,
			expectedCursor:   0, // wraps to beginning
			expectedRender:   0,
			expectedSelected: map[string]int{"/tmp/file2.txt": 1},
		},
		{
			name: "Selection persistence across moves",
			panel: Model{
				Element: []Element{
					{Name: "file1.txt", Location: "/tmp/file1.txt"},
					{Name: "file2.txt", Location: "/tmp/file2.txt"},
					{Name: "file3.txt", Location: "/tmp/file3.txt"},
				},
				Cursor:      1,
				RenderIndex: 0,
				PanelMode:   SelectMode,
			},
			panelToSelect:    []string{"/tmp/file1.txt"},
			selectDown:       true,
			mainPanelHeight:  10,
			expectedCursor:   2,
			expectedRender:   0,
			expectedSelected: map[string]int{"/tmp/file1.txt": 1, "/tmp/file2.txt": 2},
		},
		{
			name: "Empty panel selection",
			panel: Model{
				Element:     []Element{},
				Cursor:      0,
				RenderIndex: 0,
				PanelMode:   SelectMode,
			},
			panelToSelect:    []string{},
			selectDown:       true,
			mainPanelHeight:  10,
			expectedCursor:   0,
			expectedRender:   0,
			expectedSelected: map[string]int{},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.panel.SetSelectedAll(tt.panelToSelect, true)
			tt.panel.SetHeight(tt.mainPanelHeight + 2)

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
		name            string
		panel           Model
		cursorPos       int
		mainPanelHeight int
		expectedCursor  int
		expectedRender  int
	}{
		{
			name: "Jump to visible cursor no change",
			panel: Model{
				Element:     make([]Element, 20),
				Cursor:      5,
				RenderIndex: 3,
			},
			cursorPos:       4,
			mainPanelHeight: 10, // visible range [3, 9]
			expectedCursor:  4,
			expectedRender:  3,
		},
		{
			name: "Jump above view",
			panel: Model{
				Element:     make([]Element, 20),
				Cursor:      10,
				RenderIndex: 5,
			},
			cursorPos:       2,
			mainPanelHeight: 10, // panelElementHeight = 7
			expectedCursor:  2,
			expectedRender:  2,
		},
		{
			name: "Jump below view",
			panel: Model{
				Element:     make([]Element, 20),
				Cursor:      5,
				RenderIndex: 0,
			},
			cursorPos:       15,
			mainPanelHeight: 10, // visible range [0, 6]
			expectedCursor:  15,
			expectedRender:  9, // 15 - 7 + 1
		},
		{
			name: "Jump above view with empty space",
			panel: Model{
				Element:     make([]Element, 20),
				Cursor:      19,
				RenderIndex: 18,
			},
			cursorPos:       17,
			mainPanelHeight: 10, // visible range [0, 6]
			expectedCursor:  17,
			expectedRender:  17,
		},
		{
			name: "Invalid cursor negative",
			panel: Model{
				Element:     make([]Element, 10),
				Cursor:      5,
				RenderIndex: 2,
			},
			cursorPos:       -1,
			mainPanelHeight: 10,
			expectedCursor:  5, // unchanged
			expectedRender:  2, // unchanged
		},
		{
			name: "Invalid cursor beyond count",
			panel: Model{
				Element:     make([]Element, 10),
				Cursor:      5,
				RenderIndex: 2,
			},
			cursorPos:       15,
			mainPanelHeight: 10,
			expectedCursor:  5, // unchanged
			expectedRender:  2, // unchanged
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.panel.SetHeight(tt.mainPanelHeight + 2)
			tt.panel.scrollToCursor(tt.cursorPos)
			assert.Equal(t, tt.expectedCursor, tt.panel.Cursor)
			assert.Equal(t, tt.expectedRender, tt.panel.RenderIndex)
		})
	}
}

func TestApplyTargetFileCursor(t *testing.T) {
	panel := Model{
		Element: []Element{
			{Name: "file1.txt", Location: "/tmp/file1.txt"},
			{Name: "file2.txt", Location: "/tmp/file2.txt"},
			{Name: "file3.txt", Location: "/tmp/file3.txt"},
			{Name: "file4.txt", Location: "/tmp/file4.txt"},
			{Name: "target.txt", Location: "/tmp/target.txt"},
			{Name: "file6.txt", Location: "/tmp/file6.txt"},
		},
		Cursor:      0,
		RenderIndex: 0,
		TargetFile:  "target.txt",
		height:      8,
	}
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
			m := testModel(tt.initialCursor, 0, tt.panelHeight+2, tt.totalElements)
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
