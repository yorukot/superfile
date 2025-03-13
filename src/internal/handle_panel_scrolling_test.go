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

// Control processbar panel list down
func Test_processBarModel(t *testing.T) {
	testdata := []struct {
		name           string
		pModel         processBarModel
		listDown       bool // Whether to do listDown or listUp
		expectedCursor int
		expectedRender int
		footerHeight   int
		explanation    string
	}{
		{
			name:           "Basic down movement 1",
			pModel:         genProcessBarModel(10, 0, 0),
			listDown:       true,
			expectedCursor: 1,
			expectedRender: 0,
			footerHeight:   10,
			explanation:    "Basic movements",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			if tt.listDown {
				tt.pModel.listDown(footerHeight)
			} else {
				tt.pModel.listUp(footerHeight)
			}

			assert.Equal(t, tt.expectedCursor, tt.pModel.cursor)
			assert.Equal(t, tt.expectedRender, tt.pModel.render)
		})
	}
}
