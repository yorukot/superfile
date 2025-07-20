package internal

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func genProcessBarModel(count int, cursor int, render int) processBarModel {
	pList := make([]string, count)
	pMap := map[string]process{}
	for i := range count {
		pList[i] = strconv.Itoa(i)
		pMap[pList[i]] = process{
			name: pList[i],
		}
	}
	return processBarModel{
		processList: pList,
		process:     pMap,
		cursor:      cursor,
		render:      render,
	}
}

func Test_cntRenderableProcess(t *testing.T) {
	assert.Equal(t, 1, cntRenderableProcess(4))
	assert.Equal(t, 2, cntRenderableProcess(5))
	assert.Equal(t, 2, cntRenderableProcess(6))
	assert.Equal(t, 2, cntRenderableProcess(7))
	assert.Equal(t, 3, cntRenderableProcess(8))
	assert.Equal(t, 3, cntRenderableProcess(9))
	assert.Equal(t, 3, cntRenderableProcess(10))
	assert.Equal(t, 4, cntRenderableProcess(11))
}

func Test_processBarModelUpDown(t *testing.T) {
	testdata := []struct {
		name           string
		pModel         processBarModel
		listDown       bool // Whether to do listDown or listUp
		expectedCursor int
		expectedRender int
		footerHeight   int
	}{
		{
			name:           "Basic down movement 1",
			pModel:         genProcessBarModel(10, 0, 0),
			listDown:       true,
			expectedCursor: 1,
			expectedRender: 0,
			footerHeight:   10,
		},
		{
			name:           "Down at the last process - Footer height is plenty",
			pModel:         genProcessBarModel(3, 2, 0),
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
			footerHeight:   10,
		},
		{
			name:           "Down at the last process - Footer height just enough",
			pModel:         genProcessBarModel(3, 2, 0),
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
			footerHeight:   8,
		},
		{
			name:           "Down at the last process - Footer height is small",
			pModel:         genProcessBarModel(10, 9, 7),
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
			footerHeight:   8,
		},
		{
			name:           "Down at the process causing render index to move",
			pModel:         genProcessBarModel(10, 3, 0),
			listDown:       true,
			expectedCursor: 4,
			expectedRender: 1,
			footerHeight:   11, // Can hold 4 processes
		},
		{
			name:           "Basic up movement 1",
			pModel:         genProcessBarModel(10, 1, 0),
			listDown:       false,
			expectedCursor: 0,
			expectedRender: 0,
			footerHeight:   10,
		},
		{
			name:           "Up at top wraps to last and adjusts render",
			pModel:         genProcessBarModel(10, 0, 0),
			listDown:       false,
			expectedCursor: 9,
			expectedRender: 6, // 10 processes , 4 renderable
			footerHeight:   11,
		},
		{
			name:           "Up causes render index decrement",
			pModel:         genProcessBarModel(10, 3, 3),
			listDown:       false,
			expectedCursor: 2,
			expectedRender: 2, // Cursor moved above render start
			footerHeight:   8, // Renders 3 processes
		},
		{
			name:           "Up on short list wraps correctly",
			pModel:         genProcessBarModel(3, 0, 0),
			listDown:       false,
			expectedCursor: 2,
			expectedRender: 0, // 3 processes, 3 renderable
			footerHeight:   11,
		},
		{
			name:           "Up within render window maintains position",
			pModel:         genProcessBarModel(8, 5, 3),
			listDown:       false,
			expectedCursor: 4,
			expectedRender: 3, // Remain in render window
			footerHeight:   11,
		},
		{
			name:           "Up with minimal footer height",
			pModel:         genProcessBarModel(5, 0, 0),
			listDown:       false,
			expectedCursor: 4,
			expectedRender: 3,
			footerHeight:   5,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.pModel.isValid(tt.footerHeight))
			if tt.listDown {
				tt.pModel.listDown(tt.footerHeight)
			} else {
				tt.pModel.listUp(tt.footerHeight)
			}

			assert.Equal(t, tt.expectedCursor, tt.pModel.cursor)
			assert.Equal(t, tt.expectedRender, tt.pModel.render)
		})
	}
}

func Test_filePanelUpDown(t *testing.T) {
	testdata := []struct {
		name            string
		panel           filePanel
		listDown        bool
		mainPanelHeight int
		expectedCursor  int
		expectedRender  int
	}{
		{
			name: "Down movement within renderable range",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  0,
				render:  0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  1,
			expectedRender:  0,
		},
		{
			name: "Down movement when cursor is at bottom",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  6, // 3 - Header lines + 7(0-6 files)
				render:  0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  7,
			expectedRender:  1,
		},
		{
			name: "Down movement causing wrap to top",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  9, // 3 - Header lines + 7(3-9 files)
				render:  3,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
		{
			name: "Up movement within renderable range",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  2,
				render:  0,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  1,
			expectedRender:  0,
		},
		{
			name: "Up movement when cursor is at top",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  3,
				render:  3,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  2,
			expectedRender:  2,
		},
		{
			name: "Up movement causing wrap to bottom",
			panel: filePanel{
				element: make([]element, 10),
				cursor:  0,
				render:  0,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  9,
			expectedRender:  3,
		},
		{
			name: "Down movement on empty panel",
			panel: filePanel{
				element: make([]element, 0),
				cursor:  0,
				render:  0,
			},
			listDown:        true,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
		{
			name: "Up movement on empty panel",
			panel: filePanel{
				element: make([]element, 0),
				cursor:  0,
				render:  0,
			},
			listDown:        false,
			mainPanelHeight: 10,
			expectedCursor:  0,
			expectedRender:  0,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if tt.listDown {
				tt.panel.listDown(tt.mainPanelHeight)
			} else {
				tt.panel.listUp(tt.mainPanelHeight)
			}
			assert.Equal(t, tt.expectedCursor, tt.panel.cursor)
			assert.Equal(t, tt.expectedRender, tt.panel.render)
		})
	}
}

// TODO : Write tests for File Panel pgUp and pgDown and itemSelectUp/itemSelectDown
