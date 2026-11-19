package openwithmodel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestNavigationAndSelection(t *testing.T) {
	model := New()
	applications := []common.OpenWithApplication{
		{Name: "First", Command: "first"},
		{Name: "Second", Command: "second"},
	}
	model.Open("/tmp/file with spaces.txt", applications)

	assert.True(t, model.IsOpen())
	assert.Equal(t, "/tmp/file with spaces.txt", model.Path())
	model.ListUp()
	assert.Equal(t, 1, model.Cursor())
	selected, ok := model.Selected()
	require.True(t, ok)
	assert.Equal(t, applications[1], selected)
	model.ListDown()
	assert.Equal(t, 0, model.Cursor())

	model.Close()
	assert.False(t, model.IsOpen())
	assert.Empty(t, model.Path())
	_, ok = model.Selected()
	assert.False(t, ok)
}

func TestEmptyApplicationList(t *testing.T) {
	model := New()
	model.Open("/tmp/file.txt", nil)

	model.ListUp()
	model.ListDown()
	_, ok := model.Selected()

	assert.False(t, ok)
	assert.Contains(t, model.Render(), "No applications configured")
}
