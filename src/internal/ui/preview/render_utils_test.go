package preview

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCountFileLinesComplete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lines.txt")
	content := strings.Join([]string{"a", "b", "c", "d", "e"}, "\n") + "\n"
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))

	count, complete, err := countFileLines(path)
	require.NoError(t, err)
	assert.True(t, complete)
	assert.Equal(t, 5, count)
}

func TestCountFileLinesRespectsPreviewTimeout(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "many-lines.txt")
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = "line"
	}
	require.NoError(t, os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644))

	count, complete, err := countFileLinesBefore(path, time.Now().Add(-time.Second))
	require.NoError(t, err)
	assert.False(t, complete)
	assert.Positive(t, count)
	assert.Less(t, count, len(lines))
}
