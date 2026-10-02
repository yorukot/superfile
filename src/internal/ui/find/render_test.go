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
