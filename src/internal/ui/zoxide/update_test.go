package zoxide

import (
	"testing"

	zoxidelib "github.com/lazysegtree/go-zoxide"
	"github.com/stretchr/testify/assert"
)

func TestApplyWithMatchingQuery(t *testing.T) {
	m := setupTestModel()
	m.textInput.SetValue("test")
	m.cursor = 5
	m.renderIndex = 2

	results := []zoxidelib.Result{
		{Path: "/test/path1", Score: 100},
		{Path: "/test/path2", Score: 90},
		{Path: "/test/path3", Score: 80},
	}
	msg := NewUpdateMsg("test", results, 1)

	cmd := msg.Apply(&m)
	assert.Nil(t, cmd)
	assert.Len(t, m.results, 3, "results should be updated to 3 items")
	assert.Equal(t, 0, m.cursor, "cursor should be reset to 0")
	assert.Equal(t, 0, m.renderIndex, "renderIndex should be reset to 0")
	assert.Equal(t, results, m.results, "results should match the update message")
}

func TestApplyWithStaleQuery(t *testing.T) {
	m := setupTestModel()
	m.textInput.SetValue("new")
	m.cursor = 1
	m.renderIndex = 1
	originalResults := []zoxidelib.Result{
		{Path: "/original/path", Score: 50},
	}
	m.results = originalResults

	staleResults := []zoxidelib.Result{
		{Path: "/test/path1", Score: 100},
		{Path: "/test/path2", Score: 90},
		{Path: "/test/path3", Score: 80},
	}
	msg := NewUpdateMsg("old", staleResults, 1)

	cmd := msg.Apply(&m)
	assert.Nil(t, cmd)
	assert.Equal(t, originalResults, m.results, "results should remain unchanged")
	assert.Equal(t, 1, m.cursor, "cursor should remain unchanged")
	assert.Equal(t, 1, m.renderIndex, "renderIndex should remain unchanged")
}
