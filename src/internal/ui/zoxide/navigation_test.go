package zoxide

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNavigateUpDownWrapping(t *testing.T) {
	m := setupTestModelWithResults(5)
	m.cursor = 0

	m.navigateUp()
	assert.Equal(t, 4, m.cursor, "navigateUp at position 0 should wrap to position 4")

	m.cursor = 0
	m.navigateDown()
	assert.Equal(t, 1, m.cursor, "navigateDown at position 0 should move to position 1")

	m.cursor = 4
	m.navigateDown()
	assert.Equal(t, 0, m.cursor, "navigateDown at position 4 should wrap to position 0")
}

func TestNavigationWithEmptyResults(t *testing.T) {
	m := setupTestModel()
	m.cursor = 0

	m.navigateUp()
	assert.Equal(t, 0, m.cursor, "navigateUp with empty results should keep cursor at 0")

	m.navigateDown()
	assert.Equal(t, 0, m.cursor, "navigateDown with empty results should keep cursor at 0")
}

func TestUpdateRenderIndexScrolling(t *testing.T) {
	m := setupTestModelWithResults(10)

	m.cursor = 0
	m.updateRenderIndex()
	assert.Equal(t, 0, m.renderIndex, "cursor at 0 should have renderIndex = 0")

	m.cursor = 5
	m.updateRenderIndex()
	assert.Equal(t, 1, m.renderIndex, "cursor at 5 should have renderIndex = 1 (cursor visible at bottom)")

	m.cursor = 9
	m.updateRenderIndex()
	assert.Equal(t, 5, m.renderIndex, "cursor at 9 should have renderIndex = 5 (last page)")

	m.cursor = 0
	m.updateRenderIndex()
	assert.Equal(t, 0, m.renderIndex, "cursor at 0 should scroll back up to renderIndex = 0")
}

func TestUpdateRenderIndexWithFewResults(t *testing.T) {
	m := setupTestModelWithResults(3)

	m.cursor = 0
	m.updateRenderIndex()
	assert.Equal(t, 0, m.renderIndex, "renderIndex should stay 0 with 3 results, cursor at 0")

	m.cursor = 1
	m.updateRenderIndex()
	assert.Equal(t, 0, m.renderIndex, "renderIndex should stay 0 with 3 results, cursor at 1")

	m.cursor = 2
	m.updateRenderIndex()
	assert.Equal(t, 0, m.renderIndex, "renderIndex should stay 0 with 3 results, cursor at 2")
}
