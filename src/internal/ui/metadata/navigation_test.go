package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestPgUpDown(t *testing.T) {
	tenItems := Metadata{
		data: make([][2]string, 10),
	}
	testdata := []struct {
		name                string
		m                   Model
		pageScrollSize      int // value forced into common.Config.PageScrollSize
		pageDown            bool
		expectedRenderIndex int
	}{
		{
			name:                "Page down from top",
			m:                   Model{metadata: tenItems, renderIndex: 0},
			pageScrollSize:      3,
			pageDown:            true,
			expectedRenderIndex: 3,
		},
		{
			name:                "Page down near end wraps around",
			m:                   Model{metadata: tenItems, renderIndex: 8},
			pageScrollSize:      3,
			pageDown:            true,
			expectedRenderIndex: 1, // (8 + 3) % 10
		},
		{
			name:                "Page up from middle",
			m:                   Model{metadata: tenItems, renderIndex: 5},
			pageScrollSize:      3,
			pageDown:            false,
			expectedRenderIndex: 2,
		},
		{
			name:                "Page up near beginning wraps around",
			m:                   Model{metadata: tenItems, renderIndex: 1},
			pageScrollSize:      3,
			pageDown:            false,
			expectedRenderIndex: 8, // (1 - 3 + 10) % 10
		},
		{
			name:                "Page down on empty panel stays put",
			m:                   Model{metadata: Metadata{data: nil}, renderIndex: 0},
			pageScrollSize:      3,
			pageDown:            true,
			expectedRenderIndex: 0,
		},
		{
			name:                "Falls back to visible height when config unset",
			m:                   Model{metadata: tenItems, renderIndex: 0, height: 5},
			pageScrollSize:      0, // borderSize is 2 -> page size becomes 5 - 2 = 3
			pageDown:            true,
			expectedRenderIndex: 3,
		},
	}

	origPageScrollSize := common.Config.PageScrollSize
	t.Cleanup(func() { common.Config.PageScrollSize = origPageScrollSize })

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			common.Config.PageScrollSize = tt.pageScrollSize
			if tt.pageDown {
				tt.m.PgDown()
			} else {
				tt.m.PgUp()
			}
			assert.Equal(t, tt.expectedRenderIndex, tt.m.renderIndex)
		})
	}
}
