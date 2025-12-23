package processbar

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func genProcessBarModel(count int, cursor int, render int, viewHeight int) Model {
	pMap := map[string]Process{}
	for i := range count {
		pID := strconv.Itoa(i)
		pMap[pID] = Process{
			ID:   pID,
			Name: pID,
		}
	}
	return Model{
		processes:   pMap,
		cursor:      cursor,
		renderIndex: render,
		width:       minWidth,
		height:      viewHeight + 2,
	}
}

func Test_processBarModelUpDown(t *testing.T) {
	testdata := []struct {
		name           string
		processCnt     int
		cursor         int
		render         int
		listDown       bool // Whether to do ListDown or ListUp
		expectedCursor int
		expectedRender int
		footerHeight   int
	}{
		{
			name:           "Basic down movement 1",
			processCnt:     10,
			cursor:         0,
			render:         0,
			listDown:       true,
			expectedCursor: 1,
			expectedRender: 0,
			footerHeight:   10,
		},
		{
			name:           "Down at the last process - Footer height is plenty",
			processCnt:     3,
			cursor:         2,
			render:         0,
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
			footerHeight:   10,
		},
		{
			name:           "Down at the last process - Footer height just enough",
			processCnt:     3,
			cursor:         2,
			render:         0,
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
			footerHeight:   8,
		},
		{
			name:           "Down at the last process - Footer height is small",
			processCnt:     10,
			cursor:         9,
			render:         7,
			listDown:       true,
			expectedCursor: 0,
			expectedRender: 0,
			footerHeight:   8,
		},
		{
			name:           "Down at the process causing render index to move",
			processCnt:     10,
			cursor:         3,
			render:         0,
			listDown:       true,
			expectedCursor: 4,
			expectedRender: 1,
			footerHeight:   11, // Can hold 4 processes
		},
		{
			name:           "Basic up movement 1",
			processCnt:     10,
			cursor:         1,
			render:         0,
			listDown:       false,
			expectedCursor: 0,
			expectedRender: 0,
			footerHeight:   10,
		},
		{
			name:           "Up at top wraps to last and adjusts render",
			processCnt:     10,
			cursor:         0,
			render:         0,
			listDown:       false,
			expectedCursor: 9,
			expectedRender: 6, // 10 processes , 4 renderable
			footerHeight:   11,
		},
		{
			name:           "Up causes render index decrement",
			processCnt:     10,
			cursor:         3,
			render:         3,
			listDown:       false,
			expectedCursor: 2,
			expectedRender: 2, // Cursor moved above render start
			footerHeight:   8, // Renders 3 processes
		},
		{
			name:           "Up on short list wraps correctly",
			processCnt:     3,
			cursor:         0,
			render:         0,
			listDown:       false,
			expectedCursor: 2,
			expectedRender: 0, // 3 processes, 3 renderable
			footerHeight:   11,
		},
		{
			name:           "Up within render window maintains position",
			processCnt:     8,
			cursor:         5,
			render:         3,
			listDown:       false,
			expectedCursor: 4,
			expectedRender: 3, // Remain in render window
			footerHeight:   11,
		},
		{
			name:           "Up with minimal footer height",
			processCnt:     5,
			cursor:         0,
			render:         0,
			listDown:       false,
			expectedCursor: 4,
			expectedRender: 3,
			footerHeight:   5,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			pModel := genProcessBarModel(tt.processCnt, tt.cursor, tt.render, tt.footerHeight)
			assert.True(t, pModel.isValid())
			if tt.listDown {
				pModel.ListDown()
			} else {
				pModel.ListUp()
			}

			assert.Equal(t, tt.expectedCursor, pModel.cursor)
			assert.Equal(t, tt.expectedRender, pModel.renderIndex)
		})
	}
}
