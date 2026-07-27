package internal

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// Regression test for #966 : page up/down keys must be context aware. When a
// non-file panel (e.g. metadata) is focused, page keys must not scroll the file
// browser panel.
func TestPageKeysAreContextAware(t *testing.T) {
	dir := t.TempDir()
	files := make([]string, 10)
	for i := range files {
		files[i] = filepath.Join(dir, fmt.Sprintf("file%02d.txt", i))
	}
	utils.SetupFiles(t, files...)

	// Fix the page size so PgDown deterministically moves the file panel cursor.
	origPageScrollSize := common.Config.PageScrollSize
	common.Config.PageScrollSize = 3
	t.Cleanup(func() { common.Config.PageScrollSize = origPageScrollSize })

	pageDownKey := common.Hotkeys.PageDown[0]

	t.Run("pages the file panel when a file panel is focused", func(t *testing.T) {
		m := defaultTestModel(dir)
		require.Equal(t, nonePanelFocus, m.focusPanel)
		require.Equal(t, 0, m.getFocusedFilePanel().GetCursor())

		m.mainKey(pageDownKey)

		assert.Equal(t, 3, m.getFocusedFilePanel().GetCursor(),
			"file panel should page down when it is focused")
	})

	t.Run("does not page the file panel when metadata is focused", func(t *testing.T) {
		m := defaultTestModel(dir)
		m.focusPanel = metadataFocus
		m.getFocusedFilePanel().IsFocused = false
		require.Equal(t, 0, m.getFocusedFilePanel().GetCursor())

		m.mainKey(pageDownKey)

		assert.Equal(t, 0, m.getFocusedFilePanel().GetCursor(),
			"file panel must not move when the metadata panel is focused")
	})
}
