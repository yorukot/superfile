package zoxide

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestOpenResetsState(t *testing.T) {
	m := setupTestModelWithClient(t)
	common.Config.ZoxideSupport = true
	m.textInput.SetValue("old")
	m.cursor = 2
	m.renderIndex = 1

	cmd := m.Open()

	assert.True(t, m.open, "open should be true after Open()")
	assert.True(t, m.justOpened, "justOpened should be true after Open()")
	assert.Empty(t, m.textInput.Value(), "textInput should be empty after Open()")
	assert.NotNil(t, cmd, "Open() should return non-nil Cmd for async query")
}

func TestCloseClearsState(t *testing.T) {
	m := setupTestModelWithResults(5)
	m.open = true
	m.cursor = 2
	m.renderIndex = 1
	m.textInput.SetValue("test")

	m.Close()

	assert.False(t, m.open, "open should be false after Close()")
	assert.Empty(t, m.results, "results should be empty after Close()")
	assert.Equal(t, 0, m.cursor, "cursor should be 0 after Close()")
	assert.Equal(t, 0, m.renderIndex, "renderIndex should be 0 after Close()")
	assert.Empty(t, m.textInput.Value(), "textInput should be empty after Close()")
}

func TestGetResultsReturnsCopy(t *testing.T) {
	m := setupTestModelWithResults(3)
	originalPath := m.results[0].Path

	results := m.GetResults()
	results[0].Path = "/modified/path"

	assert.Equal(
		t,
		originalPath,
		m.results[0].Path,
		"modifying returned results should not affect original model.results",
	)
}

func TestSetWidthBoundsChecking(t *testing.T) {
	m := setupTestModel()

	m.SetWidth(5)
	assert.Equal(t, ZoxideMinWidth, m.width, "width should be set to ZoxideMinWidth when value < ZoxideMinWidth")

	m.SetWidth(100)
	assert.Equal(t, 100, m.width, "width should be set to provided value when >= ZoxideMinWidth")
}

func TestSetMaxHeightBoundsChecking(t *testing.T) {
	m := setupTestModel()

	m.SetMaxHeight(1)
	assert.Equal(
		t,
		ZoxideMinHeight,
		m.maxHeight,
		"maxHeight should be set to ZoxideMinHeight when value < ZoxideMinHeight",
	)

	m.SetMaxHeight(50)
	assert.Equal(t, 50, m.maxHeight, "maxHeight should be set to provided value when >= ZoxideMinHeight")
}
