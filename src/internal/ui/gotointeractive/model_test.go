package gotointeractive

import (
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
	msg := NewUpdateMsg("test", []string{"file1", "file2"}, 1)
	assert.Equal(t, "test", msg.query)
	assert.Equal(t, []string{"file1", "file2"}, msg.results)
	assert.Equal(t, 1, msg.GetReqID())
}
