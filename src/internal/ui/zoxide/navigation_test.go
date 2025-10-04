package zoxide

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNavigation(t *testing.T) {
	testdata := []struct {
		name           string
		resultCnt      int
		startCursor    int
		navigateUp     bool
		expectedCursor int
	}{
		{
			name:           "navigateUp at position 0 wraps to last position",
			resultCnt:      5,
			startCursor:    0,
			navigateUp:     true,
			expectedCursor: 4,
		},
		{
			name:           "navigateDown at position 0 moves to next position",
			resultCnt:      5,
			startCursor:    0,
			navigateUp:     false,
			expectedCursor: 1,
		},
		{
			name:           "navigateDown at last position wraps to first position",
			resultCnt:      5,
			startCursor:    4,
			navigateUp:     false,
			expectedCursor: 0,
		},
		{
			name:           "navigateUp with empty results keeps cursor at 0",
			resultCnt:      0,
			startCursor:    0,
			navigateUp:     true,
			expectedCursor: 0,
		},
		{
			name:           "navigateDown with empty results keeps cursor at 0",
			resultCnt:      0,
			startCursor:    0,
			navigateUp:     false,
			expectedCursor: 0,
		},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			var m Model
			if td.resultCnt == 0 {
				m = setupTestModel()
			} else {
				m = setupTestModelWithResults(td.resultCnt)
			}
			m.cursor = td.startCursor

			if td.navigateUp {
				m.navigateUp()
			} else {
				m.navigateDown()
			}

			assert.Equal(t, td.expectedCursor, m.cursor)
		})
	}
}

func TestUpdateRenderIndex(t *testing.T) {
	testdata := []struct {
		name                string
		resultCnt           int
		cursor              int
		expectedRenderIndex int
	}{
		{
			name:                "cursor at 0 has renderIndex 0",
			resultCnt:           10,
			cursor:              0,
			expectedRenderIndex: 0,
		},
		{
			name:                "cursor at 5 has renderIndex 1 (visible at bottom)",
			resultCnt:           10,
			cursor:              5,
			expectedRenderIndex: 1,
		},
		{
			name:                "cursor at 9 has renderIndex 5 (last page)",
			resultCnt:           10,
			cursor:              9,
			expectedRenderIndex: 5,
		},
		{
			name:                "cursor back at 0 scrolls back up to renderIndex 0",
			resultCnt:           10,
			cursor:              0,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 0",
			resultCnt:           3,
			cursor:              0,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 1",
			resultCnt:           3,
			cursor:              1,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 2",
			resultCnt:           3,
			cursor:              2,
			expectedRenderIndex: 0,
		},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			m := setupTestModelWithResults(td.resultCnt)
			m.cursor = td.cursor
			m.updateRenderIndex()
			assert.Equal(t, td.expectedRenderIndex, m.renderIndex)
		})
	}
}
