package internal

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func genProcessBarModel(count int, cursor int, render int) processBarModel {
	pSlice := make([]string, count)
	pMap := map[string]process{}
	for i := 0; i < count; i++ {
		pSlice[i] = strconv.Itoa(i)
		pMap[pSlice[i]] = process{
			name: pSlice[i],
		}
	}
	return processBarModel{
		processList: pSlice,
		process:     pMap,
		cursor:      cursor,
		render:      render,
	}
}

func Test_cntRenderableProcess(t *testing.T) {
	assert.Equal(t, cntRenderableProcess(4), 1)
	assert.Equal(t, cntRenderableProcess(5), 2)
	assert.Equal(t, cntRenderableProcess(6), 2)
	assert.Equal(t, cntRenderableProcess(7), 2)
	assert.Equal(t, cntRenderableProcess(8), 3)
	assert.Equal(t, cntRenderableProcess(9), 3)
	assert.Equal(t, cntRenderableProcess(10), 3)
	assert.Equal(t, cntRenderableProcess(11), 4)
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

func Test_fileMetadataUpDown(t *testing.T) {
	testdata := []struct {
		name               string
		fm                 fileMetadata
		listDown           bool // Whether to do listDown or listUp
		expectedRendeIndex int
	}{
		{
			name: "Basic down movement 1",
			fm: fileMetadata{
				metaData:    make([][2]string, 5),
				renderIndex: 0,
			},
			listDown:           true,
			expectedRendeIndex: 1,
		},
		{
			name: "Down wraps to top",
			fm: fileMetadata{
				metaData:    make([][2]string, 5),
				renderIndex: 4,
			},
			listDown:           true,
			expectedRendeIndex: 0,
		},
		{
			name: "Basic up movement 1",
			fm: fileMetadata{
				metaData:    make([][2]string, 5),
				renderIndex: 4,
			},
			listDown:           false,
			expectedRendeIndex: 3,
		},
		{
			name: "Up wraps to top",
			fm: fileMetadata{
				metaData:    make([][2]string, 5),
				renderIndex: 0,
			},
			listDown:           false,
			expectedRendeIndex: 4,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if tt.listDown {
				tt.fm.listDown()
			} else {
				tt.fm.listUp()
			}
			assert.Equal(t, tt.expectedRendeIndex, tt.fm.renderIndex)
		})
	}
}
