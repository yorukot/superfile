package filepanel

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			common.Config.PageScrollSize = tt.pageScrollSize

			// Create model with elements
			m := Model{
				Element:     make([]Element, tt.totalElements),
				Cursor:      tt.initialCursor,
				RenderIndex: 0,
			}

			if tt.name == "PgUp with custom scroll size" {
				m.PgUp(tt.panelHeight)
			} else {
				m.PgDown(tt.panelHeight)
			}

			assert.Equal(t, tt.expectedCursor, m.Cursor,
				"Cursor position should match expected after PgUp/PgDown")
		})
	}
}
