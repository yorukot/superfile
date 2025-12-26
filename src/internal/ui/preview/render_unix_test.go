//go:build !windows

package preview

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

	m := New()
	res := m.RenderWithPath(file, 20, 10, 20)
	assert.Contains(t, res, common.FilePreviewUnsupportedFileMode)
}
