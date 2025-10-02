package zoxide

import (
	"testing"

	zoxidelib "github.com/lazysegtree/go-zoxide"
	"github.com/stretchr/testify/assert"
)

func TestRenderWithNilZClient(t *testing.T) {
	m := setupTestModel()

	output := m.Render()

	assert.Contains(t, output, "Zoxide not available", "output should contain 'Zoxide not available'")
}

func TestRenderWithEmptyResults(t *testing.T) {
	m := setupTestModelWithClient(t)

	output := m.Render()

	assert.Contains(t, output, "No zoxide results found", "output should contain 'No zoxide results found'")
}

func TestRenderWithResults(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.results = []zoxidelib.Result{
		{Path: "/dir1", Score: 100},
		{Path: "/dir2", Score: 90},
		{Path: "/dir3", Score: 80},
	}

	output := m.Render()

	assert.Contains(t, output, "/dir1", "output should contain /dir1")
	assert.Contains(t, output, "/dir2", "output should contain /dir2")
	assert.Contains(t, output, "/dir3", "output should contain /dir3")
	assert.Contains(t, output, "100.0", "output should contain score 100.0")
	assert.Contains(t, output, "90.0", "output should contain score 90.0")
	assert.Contains(t, output, "80.0", "output should contain score 80.0")
}

func TestRenderWithTextInput(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.textInput.SetValue("test query")

	output := m.Render()

	assert.Contains(t, output, "test query", "output should contain text input value")
}

func TestRenderScrollIndicatorMoreAbove(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.results = setupTestModelWithResults(10).results
	m.renderIndex = 3
	m.cursor = 5

	output := m.Render()

	assert.Contains(t, output, "↑", "output should contain '↑' indicator when there are results above")
}

func TestRenderScrollIndicatorMoreBelow(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.results = setupTestModelWithResults(10).results
	m.renderIndex = 0
	m.cursor = 0

	output := m.Render()

	assert.Contains(t, output, "↓", "output should contain '↓' indicator when there are results below")
}

func TestRenderScrollIndicatorsBothDirections(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.results = setupTestModelWithResults(10).results
	m.renderIndex = 3
	m.cursor = 5

	output := m.Render()

	assert.Contains(t, output, "↑", "output should contain '↑' indicator when there are results above")
	assert.Contains(t, output, "↓", "output should contain '↓' indicator when there are results below")
}

func TestRenderScrollIndicatorsNoneNeeded(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.results = setupTestModelWithResults(3).results

	output := m.Render()

	assert.NotContains(t, output, "↑", "output should not contain '↑' indicator with <= maxVisibleResults")
	assert.NotContains(t, output, "↓", "output should not contain '↓' indicator with <= maxVisibleResults")
}
