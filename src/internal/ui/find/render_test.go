package find

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderFdUnavailableHint(t *testing.T) {
	m := setupTestModel()
	m.open = true
	m.fdFound = false

	rendered := m.Render()

	assert.Contains(t, rendered, "fd is not available (install fd and enable find_file_support)")
}

func TestRenderErrMsg(t *testing.T) {
	m := setupTestModelWithResults(3)
	m.open = true
	m.fdFound = true
	m.errMsg = "some error"

	rendered := m.Render()

	assert.Contains(t, rendered, "some error")
	for _, result := range m.results {
		assert.NotContains(t, rendered, result.Path)
	}
}

func TestRenderNoResults(t *testing.T) {
	m := setupTestModel()
	m.open = true
	m.fdFound = true

	rendered := m.Render()

	assert.Contains(t, rendered, "No results found")
}

func TestRenderScrollIndicator(t *testing.T) {
	// The helper model has a 50-row max height (a 43-row visible window).
	// Shrink it to the height that yields a 5-row window so the 10 results
	// overflow it and the scroll indicators render.
	const testWindow = 5
	testdata := []struct {
		name       string
		resultCnt  int
		cursor     int
		expectUp   bool
		expectDown bool
	}{
		{
			name:       "More above",
			resultCnt:  10,
			cursor:     9,
			expectUp:   true,
			expectDown: false,
		},
		{
			name:       "More below",
			resultCnt:  10,
			cursor:     0,
			expectUp:   false,
			expectDown: true,
		},
		{
			name:       "Both directions",
			resultCnt:  10,
			cursor:     5,
			expectUp:   true,
			expectDown: true,
		},
		{
			name:       "No scroll needed",
			resultCnt:  3,
			cursor:     1,
			expectUp:   false,
			expectDown: false,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			m := setupTestModel()
			m.open = true
			m.fdFound = true
			m.width = 50
			m.maxHeight = testWindow + visibleOverhead
			m.results = setupTestModelWithResults(tt.resultCnt).results
			m.cursor = tt.cursor
			m.updateRenderIndex()

			rendered := m.Render()

			if tt.expectUp {
				assert.Contains(t, rendered, scrollUpIndicator)
			} else {
				assert.NotContains(t, rendered, scrollUpIndicator)
			}

			if tt.expectDown {
				assert.Contains(t, rendered, scrollDownIndicator)
			} else {
				assert.NotContains(t, rendered, scrollDownIndicator)
			}
		})
	}
}

func TestGetCursorPath(t *testing.T) {
	testdata := []struct {
		name      string
		open      bool
		resultCnt int
		cursor    int
		expected  string
	}{
		{
			name:      "closed modal with results",
			open:      false,
			resultCnt: 3,
			cursor:    0,
			expected:  "",
		},
		{
			name:      "open modal without results",
			open:      true,
			resultCnt: 0,
			cursor:    0,
			expected:  "",
		},
		{
			name:      "cursor at first result",
			open:      true,
			resultCnt: 3,
			cursor:    0,
			expected:  "/test/path0",
		},
		{
			name:      "cursor at middle result",
			open:      true,
			resultCnt: 3,
			cursor:    1,
			expected:  "/test/path1",
		},
		{
			name:      "cursor at last result",
			open:      true,
			resultCnt: 3,
			cursor:    2,
			expected:  "/test/path2",
		},
		{
			name:      "cursor out of range",
			open:      true,
			resultCnt: 3,
			cursor:    3,
			expected:  "",
		},
		{
			name:      "cursor below zero",
			open:      true,
			resultCnt: 3,
			cursor:    -1,
			expected:  "",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			m := setupTestModelWithResults(tt.resultCnt)
			m.open = tt.open
			m.cursor = tt.cursor

			assert.Equal(t, tt.expected, m.GetCursorPath())
		})
	}
}
