package sortmodel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpDownModalSort(t *testing.T) {
	model := New()
	// Top position, UP
	model.ListUp()
	assert.Equal(t, model.Cursor, len(SortOptionsStr)-1)
	// Bottom position, Down
	model.ListDown()
	assert.Equal(t, 0, model.Cursor)
	// Top position, DOWN
	model.ListDown()
	assert.Equal(t, 1, model.Cursor)
	// Middle position, DOWN
	model.ListDown()
	assert.Equal(t, 2, model.Cursor)
	// Middle position, UP
	model.ListUp()
	assert.Equal(t, 1, model.Cursor)

	// prepare to next test
	model.ListUp()
	assert.Equal(t, 0, model.Cursor)
	model.ListUp()
	assert.Equal(t, len(SortOptionsStr)-1, model.Cursor)
	// Bottom position, UP
	model.ListUp()
	assert.Equal(t, len(SortOptionsStr)-2, model.Cursor)
}
