package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestNewOpenWithCommandPreservesArgumentsAndPath(t *testing.T) {
	application := common.OpenWithApplication{
		Name:    "Viewer",
		Command: "viewer",
		Args:    []string{"--new-window", "literal;argument"},
	}
	filePath := "/tmp/report $(touch unsafe); final.pdf"

	cmd, err := newOpenWithCommand(application, filePath)

	require.NoError(t, err)
	assert.Equal(t, []string{"viewer", "--new-window", "literal;argument", filePath}, cmd.Args)
}

func TestNewOpenWithCommandRejectsInvalidChoice(t *testing.T) {
	_, err := newOpenWithCommand(common.OpenWithApplication{}, "/tmp/report.pdf")
	assert.ErrorIs(t, err, errInvalidOpenWithChoice)

	_, err = newOpenWithCommand(common.OpenWithApplication{Command: "viewer"}, "")
	assert.ErrorIs(t, err, errInvalidOpenWithChoice)
}
