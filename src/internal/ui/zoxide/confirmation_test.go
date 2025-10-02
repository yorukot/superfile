package zoxide

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestHandleConfirmWithValidSelection(t *testing.T) {
	m := setupTestModelWithResults(3)
	m.cursor = 1

	action := m.handleConfirm()

	cdAction, ok := action.(common.CDCurrentPanelAction)
	require.True(t, ok, "action should be CDCurrentPanelAction")
	assert.Equal(t, m.results[1].Path, cdAction.Location, "action should navigate to results[1].Path")
}

func TestHandleConfirmWithNoResults(t *testing.T) {
	m := setupTestModel()

	action := m.handleConfirm()

	_, ok := action.(common.NoAction)
	assert.True(t, ok, "action should be NoAction when there are no results")
}

func TestHandleConfirmWithInvalidCursor(t *testing.T) {
	m := setupTestModelWithResults(3)
	m.cursor = 5

	action := m.handleConfirm()

	_, ok := action.(common.NoAction)
	assert.True(t, ok, "action should be NoAction when cursor is out of bounds")
}
