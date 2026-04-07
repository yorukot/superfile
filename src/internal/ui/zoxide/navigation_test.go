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
		// FIX #1: Added missing navigateUp normal decrement test (Issue #1130 Mistake #3)
		{
			name:           "navigateUp at position 3 decrements to 2",
			resultCnt:      5,
			startCursor:    3,
			navigateUp:     true,
			expectedCursor: 2,
		},
		{
			name:           "navigateUp at position 1 decrements to 0",
			resultCnt:      5,
			startCursor:    1,
			navigateUp:     true,
			expectedCursor: 0,
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
	// FIX #2 & #4: Added initialRenderIndex field so we can test renderIndex
	// decrease, and replaced magic numbers with maxVisibleResults expressions
	// so tests don't silently break if the constant changes (Issue #1130 Mistakes #1 & #2)
	testdata := []struct {
		name                string
		resultCnt           int
		cursor              int
		initialRenderIndex  int
		expectedRenderIndex int
	}{
		{
			name:                "cursor at 0 has renderIndex 0",
			resultCnt:           10,
			cursor:              0,
			initialRenderIndex:  0,
			expectedRenderIndex: 0,
		},
		{
			name:                "cursor at last visible position has renderIndex 1",
			resultCnt:           10,
			cursor:              maxVisibleResults, // was hardcoded 5
			initialRenderIndex:  0,
			expectedRenderIndex: 1,
		},
		{
			name:                "cursor at last result has renderIndex at last page",
			resultCnt:           10,
			cursor:              9,
			initialRenderIndex:  0,
			expectedRenderIndex: 10 - maxVisibleResults, // was hardcoded 5
		},
		// FIX #3: Added missing renderIndex decrease test (Issue #1130 Mistake #1)
		// Tests the branch: if m.cursor < m.renderIndex { m.renderIndex = m.cursor }
		{
			name:                "cursor above renderIndex causes renderIndex to decrease",
			resultCnt:           10,
			cursor:              2,
			initialRenderIndex:  5,
			expectedRenderIndex: 2,
		},
		{
			name:                "cursor at 0 with high renderIndex snaps renderIndex to 0",
			resultCnt:           10,
			cursor:              0,
			initialRenderIndex:  4,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 0",
			resultCnt:           3,
			cursor:              0,
			initialRenderIndex:  0,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 1",
			resultCnt:           3,
			cursor:              1,
			initialRenderIndex:  0,
			expectedRenderIndex: 0,
		},
		{
			name:                "renderIndex stays 0 with 3 results, cursor at 2",
			resultCnt:           3,
			cursor:              2,
			initialRenderIndex:  0,
			expectedRenderIndex: 0,
		},
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			m := setupTestModelWithResults(td.resultCnt)
			m.cursor = td.cursor
			m.renderIndex = td.initialRenderIndex // FIX #2: pre-set renderIndex
			m.updateRenderIndex()
			assert.Equal(t, td.expectedRenderIndex, m.renderIndex)
		})
	}
}
