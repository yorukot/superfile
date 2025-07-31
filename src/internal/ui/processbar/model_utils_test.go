package processbar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModelUtils(t *testing.T) {
	m := NewModelWithOptions(10, 14)
	assert.Equal(t, 8, m.viewHeight())
	assert.Equal(t, 12, m.viewWidth())
	assert.Equal(t, 0, m.cntProcesses())
	assert.Equal(t, 1, m.newReqCnt())
	assert.Equal(t, 2, m.newReqCnt())
	assert.True(t, m.isValid())

	p1 := NewProcess("1", "test", 10)
	p2 := NewProcess("2", "test2", 11)

	_ = m.AddProcess(p1)
	_ = m.AddProcess(p2)

	assert.Equal(t, 2, m.cntProcesses())
	assert.True(t, m.isValid())

	m.cursor = -1
	assert.False(t, m.isValid())
	m.cursor = 0
	m.renderIndex = 1
	assert.False(t, m.isValid())
}
