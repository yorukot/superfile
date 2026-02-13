package gotointeractive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")

	assert.False(t, m.IsOpen())
	assert.Equal(t, 80, m.GetWidth())
	assert.Equal(t, 20, m.GetMaxHeight())
	assert.Equal(t, "/tmp", m.GetCurrentPath())
}

func TestOpenClose(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")

	cmd := m.Open()
	assert.NotNil(t, cmd)
	assert.True(t, m.IsOpen())
	assert.Empty(t, m.GetTextInputValue())

	m.Close()
	assert.False(t, m.IsOpen())
	assert.Empty(t, m.GetTextInputValue())
	assert.Empty(t, m.GetResults())
}

func TestSetCurrentPath(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")

	m.SetCurrentPath("/home")
	assert.Equal(t, "/home", m.GetCurrentPath())
}

func TestSetDimensions(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")

	m.SetWidth(100)
	assert.Equal(t, 100, m.GetWidth())

	m.SetMaxHeight(30)
	assert.Equal(t, 30, m.GetMaxHeight())
}

func TestSetDimensionsMinValues(t *testing.T) {
	m := DefaultModel(20, 40, "/tmp")

	m.SetWidth(10)
	assert.Equal(t, GotoMinWidth, m.GetWidth())

	m.SetMaxHeight(5)
	assert.Equal(t, GotoMinHeight, m.GetMaxHeight())
}

func TestUpdateMsg(t *testing.T) {
	results := []Result{
		{Name: "file1", IsDir: false},
		{Name: "file2", IsDir: true},
	}
	msg := NewUpdateMsg("test", results, 1, "/tmp")
	assert.Equal(t, "test", msg.query)
	assert.Equal(t, results, msg.results)
	assert.Equal(t, 1, msg.GetReqID())
	assert.Equal(t, "/tmp", msg.path)
}

func TestNavigateUp(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")
	m.results = []Result{
		{Name: "file1", IsDir: false},
		{Name: "file2", IsDir: true},
		{Name: "file3", IsDir: false},
	}
	m.cursor = 1

	m.navigateUp()
	assert.Equal(t, 0, m.cursor)

	m.navigateUp()
	assert.Equal(t, 2, m.cursor)
}

func TestNavigateDown(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")
	m.results = []Result{
		{Name: "file1", IsDir: false},
		{Name: "file2", IsDir: true},
		{Name: "file3", IsDir: false},
	}
	m.cursor = 1

	m.navigateDown()
	assert.Equal(t, 2, m.cursor)

	m.navigateDown()
	assert.Equal(t, 0, m.cursor)
}

func TestNavigatePageUp(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")
	m.results = make([]Result, 20)
	for i := 0; i < 20; i++ {
		m.results[i] = Result{Name: "file" + string(rune('0'+i)), IsDir: i%2 == 0}
	}
	m.cursor = 19
	m.renderIndex = 5

	m.navigatePageUp()
	assert.Equal(t, 5, m.cursor)
	assert.Equal(t, 5, m.renderIndex)
}

func TestNavigatePageDown(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")
	m.results = make([]Result, 20)
	for i := 0; i < 20; i++ {
		m.results[i] = Result{Name: "file" + string(rune('0'+i)), IsDir: false}
	}
	m.cursor = 0

	m.navigatePageDown()
	assert.Equal(t, 14, m.cursor)
	assert.Equal(t, 0, m.renderIndex)
}

func TestApplyWithMatchingQuery(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")
	m.textInput.SetValue("test")
	m.cursor = 5
	m.renderIndex = 2

	results := []Result{
		{Name: "file1", IsDir: false},
		{Name: "file2", IsDir: true},
		{Name: "file3", IsDir: false},
	}
	msg := NewUpdateMsg("test", results, 1, "/tmp")

	cmd := msg.Apply(&m)
	assert.Nil(t, cmd)
	assert.Len(t, m.results, 3, "results should be updated to 3 items")
	assert.Equal(t, 0, m.cursor, "cursor should be reset to 0")
	assert.Equal(t, 0, m.renderIndex, "renderIndex should be reset to 0")
	assert.Equal(t, results, m.results, "results should match the update message")
}

func TestApplyWithStaleQuery(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")
	m.textInput.SetValue("new")
	m.cursor = 1
	m.renderIndex = 1
	originalResults := []Result{
		{Name: "original", IsDir: false},
	}
	m.results = originalResults

	staleResults := []Result{
		{Name: "file1", IsDir: false},
		{Name: "file2", IsDir: true},
		{Name: "file3", IsDir: false},
	}
	msg := NewUpdateMsg("old", staleResults, 1, "/tmp")

	cmd := msg.Apply(&m)
	assert.Nil(t, cmd)
	assert.Equal(t, originalResults, m.results, "results should remain unchanged")
	assert.Equal(t, 1, m.cursor, "cursor should remain unchanged")
	assert.Equal(t, 1, m.renderIndex, "renderIndex should remain unchanged")
}

func TestApplyWithStalePath(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")
	m.textInput.SetValue("test")
	m.cursor = 1
	m.renderIndex = 1
	m.currentPath = "/tmp"
	originalResults := []Result{
		{Name: "original", IsDir: false},
	}
	m.results = originalResults

	staleResults := []Result{
		{Name: "file1", IsDir: false},
	}
	msg := NewUpdateMsg("test", staleResults, 1, "/other")

	cmd := msg.Apply(&m)
	assert.Nil(t, cmd)
	assert.Equal(t, originalResults, m.results, "results should remain unchanged")
	assert.Equal(t, 1, m.cursor, "cursor should remain unchanged")
	assert.Equal(t, 1, m.renderIndex, "renderIndex should remain unchanged")
}

func TestApplyWithStaleQueryAndPath(t *testing.T) {
	m := DefaultModel(20, 80, "/tmp")
	m.textInput.SetValue("new")
	m.cursor = 1
	m.renderIndex = 1
	m.reqCnt = 3
	m.currentPath = "/tmp"
	originalResults := []Result{
		{Name: "original", IsDir: false},
	}
	m.results = originalResults

	staleResults := []Result{
		{Name: "file1", IsDir: false},
	}
	msg := NewUpdateMsg("old", staleResults, 1, "/tmp")

	cmd := msg.Apply(&m)
	assert.Nil(t, cmd)
	assert.Equal(t, originalResults, m.results, "results should remain unchanged")
	assert.Equal(t, 1, m.cursor, "cursor should remain unchanged")
	assert.Equal(t, 1, m.renderIndex, "renderIndex should remain unchanged")
}

func TestFilterResultsExcludesParentOnNonMatchingQuery(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("test"), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("test"), 0644)

	results := filterResults("abc", tmpDir)
	hasParent := false
	for _, r := range results {
		if r.Name == ".." {
			hasParent = true
		}
	}
	assert.False(t, hasParent, ".. should not be in results when query 'abc' doesn't match")
}

func TestFilterResultsIncludesParentOnEmptyQuery(t *testing.T) {
	path := "/home/user/test"
	results := filterResults("", path)
	hasParent := false
	for _, r := range results {
		if r.Name == ".." {
			hasParent = true
		}
	}
	assert.True(t, hasParent, ".. should be in results when query is empty")
}

func TestFilterResultsIncludesParentOnMatchingQuery(t *testing.T) {
	path := "/home/user/test"
	results := filterResults("..", path)
	hasParent := false
	for _, r := range results {
		if r.Name == ".." {
			hasParent = true
		}
	}
	assert.True(t, hasParent, ".. should be in results when query matches")
}
