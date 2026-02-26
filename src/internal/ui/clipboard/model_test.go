package clipboard

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func TestMain(m *testing.M) {
	//nolint:reassign // Needed to tests
	common.ClipboardNoneText = " " + icon.Error + icon.Space + " No content in clipboard"
	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}
	m.Run()
}

func TestClipboardRender_Empty(t *testing.T) {
	dir := t.TempDir()
	var items []string
	for i := range 5 {
		fp := filepath.Join(dir, "f"+strconv.Itoa(i)+".txt")
		items = append(items, fp)
	}
	m := &Model{}
	m.SetDimensions(15+len(items[0]), 6)
	t.Run("Empty", func(t *testing.T) {
		out := ansi.Strip(m.Render())
		assert.Contains(t, out, common.ClipboardNoneText)
	})

	utils.CreateFiles(items[0])
	t.Run("Single Item", func(t *testing.T) {
		m.SetItems([]string{items[0]})
		out := ansi.Strip(m.Render())
		assert.NotContains(t, out, common.ClipboardNoneText)
		assert.Contains(t, out, items[0])
		assert.NotContains(t, out, items[1])
	})

	utils.CreateFiles(items[1])
	t.Run("Only two items exist, rest don't", func(t *testing.T) {
		m.SetItems(items)
		out := ansi.Strip(m.Render())
		assert.NotContains(t, out, common.ClipboardNoneText)
		assert.Contains(t, out, items[0])
		assert.Contains(t, out, items[1])
		for i := 2; i < 5; i++ {
			assert.NotContains(t, out, items[i])
		}
	})

	utils.CreateFiles(items[2:]...)
	t.Run("Overflow", func(t *testing.T) {
		m.SetItems(items)
		out := ansi.Strip(m.Render())
		assert.NotContains(t, out, common.ClipboardNoneText)
		for i := range 3 {
			assert.Contains(t, out, items[i])
		}
		assert.Contains(t, out, "2 items left....", "expected overflow indicator in render")
	})
}

func TestPruneInaccessibleItemsAndGet(t *testing.T) {
	dir := t.TempDir()
	files := []string{filepath.Join(dir, "f1"), filepath.Join(dir, "f2")}
	utils.SetupFiles(t, files...)

	m := &Model{}
	m.SetItems(files)
	assert.Equal(t, files, m.PruneInaccessibleItemsAndGet())
	require.NoError(t, os.Remove(files[1]))
	assert.Equal(t, []string{files[0]}, m.PruneInaccessibleItemsAndGet())
}
