package zoxide

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetMaxVisibleResults(t *testing.T) {
	m := setupTestModel()

	assert.Equal(t, DefaultMaxVisibleResults, m.GetConfiguredMaxVisibleResults())

	m.SetMaxVisibleResults(10)
	assert.Equal(t, 10, m.GetConfiguredMaxVisibleResults())

	// Non-positive values fallback to default
	m.SetMaxVisibleResults(0)
	assert.Equal(t, DefaultMaxVisibleResults, m.GetConfiguredMaxVisibleResults())

	m.SetMaxVisibleResults(-5)
	assert.Equal(t, DefaultMaxVisibleResults, m.GetConfiguredMaxVisibleResults())
}

func TestGetMaxVisibleResultsHeightClamping(t *testing.T) {
	m := setupTestModel()
	m.SetMaxVisibleResults(15)

	// Generous maxHeight (50): 50 - 7 = 43 available, so all 15 can be shown
	m.SetMaxHeight(50)
	assert.Equal(t, 15, m.GetMaxVisibleResults())

	// Constrained maxHeight (12): 12 - 7 = 5 available, so clamps to 5
	m.SetMaxHeight(12)
	assert.Equal(t, 5, m.GetMaxVisibleResults())

	// Tight maxHeight below renderOverhead: does not clamp to non-positive number
	m.SetMaxHeight(5)
	assert.Equal(t, 15, m.GetMaxVisibleResults())
}

func TestNavigationWithCustomMaxVisibleResults(t *testing.T) {
	customVisible := 10
	totalResults := 25
	m := setupTestModelWithResults(totalResults)
	m.SetMaxHeight(50)
	m.SetMaxVisibleResults(customVisible)

	assert.Equal(t, customVisible, m.GetMaxVisibleResults())

	// Cursor within first page (0..9) keeps renderIndex at 0
	for cursor := range customVisible {
		m.cursor = cursor
		m.renderIndex = 0
		m.updateRenderIndex()
		assert.Equal(t, 0, m.renderIndex, "renderIndex should remain 0 while cursor is within visible range")
	}

	// Cursor at index 10 should scroll renderIndex to 1
	m.cursor = 10
	m.renderIndex = 0
	m.updateRenderIndex()
	assert.Equal(t, 1, m.renderIndex)

	// Cursor at last item (24) should scroll renderIndex to 24 - 10 + 1 = 15
	m.cursor = totalResults - 1
	m.renderIndex = 0
	m.updateRenderIndex()
	assert.Equal(t, totalResults-customVisible, m.renderIndex)
}

func TestRenderScrollIndicatorsWithCustomMaxVisibleResults(t *testing.T) {
	m := setupTestModelWithClient(t)
	m.width = 60
	m.SetMaxVisibleResults(10)
	m.SetMaxHeight(10 + renderOverhead)

	// With 8 results and max visible 10, no scrolling needed
	m.results = setupTestModelWithResults(8).results
	m.cursor = 0
	m.updateRenderIndex()
	output := m.Render()
	assert.NotContains(t, output, scrollDownIndicator)
	assert.NotContains(t, output, scrollUpIndicator)

	// With 15 results and max visible 10 at top, scroll down indicator should display
	m.results = setupTestModelWithResults(15).results
	m.cursor = 0
	m.updateRenderIndex()
	output = m.Render()
	assert.Contains(t, output, scrollDownIndicator)
	assert.NotContains(t, output, scrollUpIndicator)
}
