package preview

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
)

func TestPreviewScrollTextFile(t *testing.T) {
	originalBulk := common.Config.PreviewScrollBulk
	common.Config.PreviewScrollBulk = 1
	t.Cleanup(func() {
		common.Config.PreviewScrollBulk = originalBulk
	})

	curTestDir := t.TempDir()
	filePath := filepath.Join(curTestDir, "scroll.txt")
	content := strings.Join([]string{"line1", "line2", "line3", "line4", "line5"}, "\n") + "\n"
	err := os.WriteFile(filePath, []byte(content), 0o644)
	require.NoError(t, err)

	m := New()
	m.Open()
	render, _ := m.RenderWithPath(filePath, 20, 2, 20)
	assert.Contains(t, ansi.Strip(render), "line1")
	assert.Contains(t, ansi.Strip(render), "line2")
	assert.NotContains(t, ansi.Strip(render), "line3")
	assert.True(t, m.CanScrollDown())

	require.True(t, m.ScrollBulkDown(2))
	render, _ = m.RenderWithPath(filePath, 20, 2, 20)
	stripped := ansi.Strip(render)
	assert.Contains(t, stripped, "line3")
	assert.Contains(t, stripped, "line4")
	assert.NotContains(t, stripped, "line1")

	require.True(t, m.ScrollBulkUp(2))
	render, _ = m.RenderWithPath(filePath, 20, 2, 20)
	stripped = ansi.Strip(render)
	assert.Contains(t, stripped, "line1")
	assert.False(t, m.ScrollBulkUp(2))
}

func TestPreviewScrollDirectory(t *testing.T) {
	originalBulk := common.Config.PreviewScrollBulk
	common.Config.PreviewScrollBulk = 1
	t.Cleanup(func() {
		common.Config.PreviewScrollBulk = originalBulk
	})

	curTestDir := t.TempDir()
	for i := range 5 {
		name := filepath.Join(curTestDir, "file"+string(rune('a'+i))+".txt")
		err := os.WriteFile(name, []byte("x"), 0o644)
		require.NoError(t, err)
	}

	m := New()
	m.Open()
	render, _ := m.RenderWithPath(curTestDir, 20, 2, 20)
	stripped := ansi.Strip(render)
	assert.Contains(t, stripped, "filea.txt")
	assert.Contains(t, stripped, "fileb.txt")
	assert.NotContains(t, stripped, "filec.txt")
	assert.True(t, m.CanScrollDown())

	require.True(t, m.ScrollBulkDown(2))
	render, _ = m.RenderWithPath(curTestDir, 20, 2, 20)
	stripped = ansi.Strip(render)
	assert.Contains(t, stripped, "filec.txt")
}

func TestPreviewScrollResetsOnNewLocation(t *testing.T) {
	curTestDir := t.TempDir()
	file1 := filepath.Join(curTestDir, "one.txt")
	file2 := filepath.Join(curTestDir, "two.txt")
	require.NoError(t, os.WriteFile(file1, []byte("a\nb\nc\nd\n"), 0o644))
	require.NoError(t, os.WriteFile(file2, []byte("w\nx\ny\nz\n"), 0o644))

	m := New()
	m.Open()
	_, _ = m.RenderWithPath(file1, 10, 1, 10)
	require.True(t, m.ScrollLineDown())
	assert.Equal(t, 1, m.scrollOffset)

	m.SetLocation(file2)
	assert.Equal(t, 0, m.scrollOffset)
}

func TestPreviewBulkScrollSizeUsesConfig(t *testing.T) {
	original := common.Config.PreviewScrollBulk
	common.Config.PreviewScrollBulk = 2
	t.Cleanup(func() {
		common.Config.PreviewScrollBulk = original
	})

	curTestDir := t.TempDir()
	filePath := filepath.Join(curTestDir, "scroll.txt")
	content := strings.Join([]string{"l1", "l2", "l3", "l4", "l5", "l6"}, "\n") + "\n"
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o644))

	m := New()
	m.Open()
	_, _ = m.RenderWithPath(filePath, 10, 1, 4)
	require.True(t, m.ScrollBulkDown(4))
	assert.Equal(t, 2, m.scrollOffset)

	require.True(t, m.ScrollTop())
	assert.Equal(t, 0, m.scrollOffset)
}

func TestPreviewBulkScrollWholePage(t *testing.T) {
	original := common.Config.PreviewScrollBulk
	common.Config.PreviewScrollBulk = 1
	t.Cleanup(func() {
		common.Config.PreviewScrollBulk = original
	})

	curTestDir := t.TempDir()
	filePath := filepath.Join(curTestDir, "scroll.txt")
	content := strings.Join([]string{"l1", "l2", "l3", "l4"}, "\n") + "\n"
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0o644))

	m := New()
	m.Open()
	_, _ = m.RenderWithPath(filePath, 10, 1, 2)
	require.True(t, m.ScrollBulkDown(2))
	assert.Equal(t, 2, m.scrollOffset)
}
