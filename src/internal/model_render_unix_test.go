//go:build !windows

package internal

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestFilePreviewWithInvalidMode(t *testing.T) {
	curTestDir := t.TempDir()
	file := filepath.Join(curTestDir, "testf")

	err := syscall.Mkfifo(file, 0644)
	require.NoError(t, err)

	m := defaultTestModel(curTestDir)

	res := m.filePreviewPanelRenderWithDimensions(10, 100)
	assert.Contains(t, res, common.FilePreviewUnsupportedFileMode)
}
