package clipboard

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/utils"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		utils.SetRootLoggerToStdout(true)
	} else {
		utils.SetRootLoggerToDiscarded()
	}
	m.Run()
}

func TestClipboardRender_Empty(t *testing.T) {
	m := &Model{}
	m.SetDimensions(30, 6)

	out := ansi.Strip(m.Render())

	// Empty message present
	assert.Contains(t, out, "No content in clipboard")
}

func TestClipboardRender_SingleItem(t *testing.T) {
	dir := t.TempDir()
	fpath := filepath.Join(dir, "file.txt")
	err := os.WriteFile(fpath, []byte("data"), 0o644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	m := &Model{}
	// Ensure enough width so base name is not truncated too aggressively
	m.SetDimensions(40, 6)
	m.SetItems([]string{fpath})

	out := ansi.Strip(m.Render())

	assert.NotContains(t, out, "No content in clipboard")
	// We just check basename presence to avoid icon/style coupling
	assert.Contains(t, out, filepath.Base(fpath))
}

func TestClipboardRender_OverflowIndicator(t *testing.T) {
	dir := t.TempDir()
	// Create 5 files; with height=4 we get viewHeight=2 so "4 item left...." should show
	var items []string
	for i := range 5 {
		fp := filepath.Join(dir, "f"+strconv.Itoa(i)+".txt")
		if err := os.WriteFile(fp, []byte("x"), 0o644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		items = append(items, fp)
	}

	m := &Model{}
	m.SetDimensions(40, 4) // viewHeight = 2
	m.SetItems(items)

	out := ansi.Strip(m.Render())
	// Ensure last visible line is the overflow indicator; use Contains to avoid width/padding coupling
	assert.Contains(t, out, "4 item left....", "expected overflow indicator in render, got:\n%s", out)
}
