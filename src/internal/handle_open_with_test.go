package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestEnterPanelReturnsCommandForTerminalApplication(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "settings.conf")
	require.NoError(t, os.WriteFile(filePath, nil, 0o600))

	previousOpenWith := common.Config.OpenWith
	common.Config.OpenWith = map[string]any{
		"conf": common.OpenWithApplication{Command: "nvim", Terminal: true},
	}
	t.Cleanup(func() {
		common.Config.OpenWith = previousOpenWith
	})

	m := defaultTestModel(dir)
	setFilePanelSelectedItemByName(t, m.getFocusedFilePanel(), filepath.Base(filePath))

	assert.NotNil(t, m.enterPanel())
}

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
