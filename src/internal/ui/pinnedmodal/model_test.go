package pinnedmodal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	m := DefaultModel(20, 80)

	assert.False(t, m.IsOpen())
	assert.Equal(t, 80, m.GetWidth())
	assert.Equal(t, 20, m.GetMaxHeight())
}

func TestOpenClose(t *testing.T) {
	m := DefaultModel(20, 80)

	cmd := m.Open()
	assert.NotNil(t, cmd)
	assert.True(t, m.IsOpen())
	assert.Empty(t, m.GetTextInputValue())

	m.Close()
	assert.False(t, m.IsOpen())
	assert.Empty(t, m.GetTextInputValue())
	assert.Empty(t, m.GetResults())
}

func TestSetDimensions(t *testing.T) {
	m := DefaultModel(20, 80)

	m.SetWidth(100)
	assert.Equal(t, 100, m.GetWidth())

	m.SetMaxHeight(30)
	assert.Equal(t, 30, m.GetMaxHeight())
}

func TestSetDimensionsMinValues(t *testing.T) {
	m := DefaultModel(20, 40)

	m.SetWidth(10)
	assert.Equal(t, PinnedModalMinWidth, m.GetWidth())

	m.SetMaxHeight(2)
	assert.Equal(t, PinnedModalMinHeight, m.GetMaxHeight())
}

func TestUpdateMsg(t *testing.T) {
	results := []Directory{
		{Name: "dir1", Location: "/path/dir1"},
		{Name: "dir2", Location: "/path/dir2"},
	}
	msg := NewUpdateMsg("test", results, 1, "")
	assert.Equal(t, "test", msg.query)
	assert.Equal(t, results, msg.results)
	assert.Equal(t, 1, msg.reqID)
	assert.Empty(t, msg.path)
}

func TestNavigateUp(t *testing.T) {
	m := DefaultModel(20, 80)
	m.results = []Directory{
		{Name: "dir1", Location: "/path/dir1"},
		{Name: "dir2", Location: "/path/dir2"},
		{Name: "dir3", Location: "/path/dir3"},
	}
	m.cursor = 1

	m.navigateUp()
	assert.Equal(t, 0, m.cursor)

	m.navigateUp()
	assert.Equal(t, 2, m.cursor)
}

func TestNavigateDown(t *testing.T) {
	m := DefaultModel(20, 80)
	m.results = []Directory{
		{Name: "dir1", Location: "/path/dir1"},
		{Name: "dir2", Location: "/path/dir2"},
		{Name: "dir3", Location: "/path/dir3"},
	}
	m.cursor = 1

	m.navigateDown()
	assert.Equal(t, 2, m.cursor)

	m.navigateDown()
	assert.Equal(t, 0, m.cursor)
}

func TestNavigatePageUp(t *testing.T) {
	m := DefaultModel(20, 80)
	m.results = make([]Directory, 20)
	for i := 0; i < 20; i++ {
		m.results[i] = Directory{Name: "dir" + string(rune('0'+i)), Location: "/path/dir" + string(rune('0'+i))}
	}
	m.cursor = 19
	m.renderIndex = 5

	m.navigatePageUp()
	assert.Equal(t, 15, m.cursor)
	assert.Equal(t, 5, m.renderIndex)
}

func TestNavigatePageDown(t *testing.T) {
	m := DefaultModel(20, 80)
	m.results = make([]Directory, 20)
	for i := 0; i < 20; i++ {
		m.results[i] = Directory{Name: "dir" + string(rune('0'+i)), Location: "/path/dir" + string(rune('0'+i))}
	}
	m.cursor = 0

	m.navigatePageDown()
	assert.Equal(t, 4, m.cursor)
	assert.Equal(t, 0, m.renderIndex)
}

func TestApplyWithMatchingQuery(t *testing.T) {
	m := DefaultModel(20, 80)
	m.textInput.SetValue("test")
	m.cursor = 5
	m.renderIndex = 2

	results := []Directory{
		{Name: "dir1", Location: "/path/dir1"},
		{Name: "dir2", Location: "/path/dir2"},
		{Name: "dir3", Location: "/path/dir3"},
	}
	msg := NewUpdateMsg("test", results, 1, "")

	cmd := msg.Apply(&m)
	assert.Nil(t, cmd)
	assert.Len(t, m.results, 3, "results should be updated to 3 items")
	assert.Equal(t, 0, m.cursor, "cursor should be reset to 0")
	assert.Equal(t, 0, m.renderIndex, "renderIndex should be reset to 0")
	assert.Equal(t, results, m.results, "results should match the update message")
}

func TestApplyWithStaleQuery(t *testing.T) {
	m := DefaultModel(20, 80)
	m.textInput.SetValue("new")
	m.cursor = 1
	m.renderIndex = 1
	originalResults := []Directory{
		{Name: "original", Location: "/path/original"},
	}
	m.results = originalResults

	staleResults := []Directory{
		{Name: "dir1", Location: "/path/dir1"},
		{Name: "dir2", Location: "/path/dir2"},
		{Name: "dir3", Location: "/path/dir3"},
	}
	msg := NewUpdateMsg("old", staleResults, 1, "")

	cmd := msg.Apply(&m)
	assert.Nil(t, cmd)
	assert.Equal(t, originalResults, m.results, "results should remain unchanged")
	assert.Equal(t, 1, m.cursor, "cursor should remain unchanged")
	assert.Equal(t, 1, m.renderIndex, "renderIndex should remain unchanged")
}

func TestLoadPinnedDirs(t *testing.T) {
	m := DefaultModel(20, 80)
	dirs := []Directory{
		{Name: "dir1", Location: "/path/dir1"},
		{Name: "dir2", Location: "/path/dir2"},
		{Name: "dir3", Location: "/path/dir3"},
	}

	m.LoadPinnedDirs(dirs)
	assert.Equal(t, dirs, m.allDirs, "allDirs should be set")
	assert.Equal(t, dirs, m.results, "results should be set")
	assert.Equal(t, 0, m.cursor, "cursor should be reset")
	assert.Equal(t, 0, m.renderIndex, "renderIndex should be reset")
}
