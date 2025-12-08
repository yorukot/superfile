package zoxide

import (
	"testing"

	zoxidelib "github.com/lazysegtree/go-zoxide"
	"github.com/stretchr/testify/assert"
)

func TestRenderWithNilZClient(t *testing.T) {
	m := setupTestModel()
	// Ensure deterministic size
	m.width = 80
	m.maxHeight = 24

	output := m.Render()

	assert.Contains(t, output, "Zoxide not available")
}

func TestRenderWithEmptyResults(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.width = 80
	m.maxHeight = 24

	output := m.Render()

	assert.Contains(t, output, "No zoxide results found")
}

func TestRenderWithResults(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.width = 80
	m.maxHeight = 24
	m.results = []zoxidelib.Result{
		{Path: "/dir1", Score: 100},
		{Path: "/dir2", Score: 90},
		{Path: "/dir3", Score: 80},
	}

	output := m.Render()

	assert.Contains(t, output, "/dir1")
	assert.Contains(t, output, "100.0")
	assert.Contains(t, output, "90.0")
	assert.Contains(t, output, "80.0")
}

func TestRenderWithTextInput(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.width = 80
	m.maxHeight = 24
	m.textInput.SetValue("test query")

	output := m.Render()
	assert.Contains(t, output, "test query")
}

func TestRenderScrollIndicator(t *testing.T) {

	testdata := []struct {
		name       string
		resultCnt  int
		cursor     int
		expectUp   bool
		expectDown bool
	}{
		{name: "More above", resultCnt: 10, cursor: 9, expectUp: true, expectDown: false},
		{name: "More below", resultCnt: 10, cursor: 0, expectUp: false, expectDown: true},
		{name: "Both directions", resultCnt: 10, cursor: 5, expectUp: true, expectDown: true},
		{name: "No scroll needed", resultCnt: 3, cursor: 1, expectUp: false, expectDown: false},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			m := setupTestModelWithClient(t)
			m.width = 80
			m.maxHeight = 24 // ensures maxVisibleResults stays 5 with current renderer logic

			m.results = setupTestModelWithResults(tt.resultCnt).results
			m.cursor = tt.cursor
			m.updateRenderIndex()

			rendered := m.Render()

			if tt.expectUp {
				assert.Contains(t, rendered, scrollIndicatorUp)
			} else {
				assert.NotContains(t, rendered, scrollIndicatorUp)
			}
			if tt.expectDown {
				assert.Contains(t, rendered, scrollIndicatorDown)
			} else {
				assert.NotContains(t, rendered, scrollIndicatorDown)
			}
		})
	}
}
