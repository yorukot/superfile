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
		// Existing good cases
		{name: "navigateUp at 0 wraps to last", resultCnt: 5, startCursor: 0, navigateUp: true, expectedCursor: 4},
		{name: "navigateDown at 0 moves to 1", resultCnt: 5, startCursor: 0, navigateUp: false, expectedCursor: 1},
		{name: "navigateDown at last wraps to 0", resultCnt: 5, startCursor: 4, navigateUp: false, expectedCursor: 0},
		// NEW: normal navigateUp decrement (fixes mistake #3)
		{name: "navigateUp normal decrement 3→2", resultCnt: 5, startCursor: 3, navigateUp: true, expectedCursor: 2},
		{name: "navigateUp normal decrement 1→0", resultCnt: 5, startCursor: 1, navigateUp: true, expectedCursor: 0},
		// Edge cases
		{name: "navigateUp empty results", resultCnt: 0, startCursor: 0, navigateUp: true, expectedCursor: 0},
		{name: "navigateDown empty results", resultCnt: 0, startCursor: 0, navigateUp: false, expectedCursor: 0},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			var m Model
			if td.resultCnt == 0 {
				m = setupTestModel()
			} else {
				m = setupTestModelWithResults(td.resultCnt)
			}
			// Ensure deterministic visible count
			m.width = 80
			m.maxHeight = 24

			m.cursor = td.startCursor

			if td.navigateUp {
				m.navigateUp()
			} else {
				m.navigateDown()
			}

			assert.Equal(t, td.expectedCursor, m.cursor, "cursor position mismatch")
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
		{name: "cursor 0 → renderIndex 0", resultCnt: 10, cursor: 0, expectedRenderIndex: 0},
		{name: "cursor 5 → renderIndex 1 (at bottom of page)", resultCnt: 10, cursor: 5, expectedRenderIndex: 1},
		{name: "cursor 9 → renderIndex 5 (last page)", resultCnt: 10, cursor: 9, expectedRenderIndex: 5},
		{name: "few results → renderIndex always 0", resultCnt: 3, cursor: 2, expectedRenderIndex: 0},
		// NEW: test that renderIndex decreases when moving cursor up (fixes mistake #1)
		{
			name:                "moving cursor from 9 → 0 decreases renderIndex",
			resultCnt:           10,
			cursor:              0,
			expectedRenderIndex: 0,
		},
	}

	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			m := setupTestModelWithResults(td.resultCnt)
			m.width = 80
			m.maxHeight = 24
			m.cursor = td.cursor
			m.updateRenderIndex()

			assert.Equal(t, td.expectedRenderIndex, m.renderIndex,
				"renderIndex wrong for cursor=%d, results=%d", td.cursor, td.resultCnt)
		})
	}
}
