package internal

import (
	"path/filepath"
	"reflect"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/ui/preview"
	"github.com/yorukot/superfile/src/pkg/utils"
)

// runSequence runs the commands of a tea.Sequence in order and returns their
// messages, the way the Bubble Tea event loop would receive them.
func runSequence(t *testing.T, cmd tea.Cmd) []tea.Msg {
	t.Helper()
	require.NotNil(t, cmd)
	seq := reflect.ValueOf(cmd())
	require.Equal(t, reflect.Slice, seq.Kind(), "expected a tea.Sequence")
	msgs := make([]tea.Msg, 0, seq.Len())
	for i := range seq.Len() {
		msgs = append(msgs, seq.Index(i).Interface().(tea.Cmd)())
	}
	return msgs
}

func TestKittyPlaceholdersWaitForImageData(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "image.png")
	utils.SetupFilesWithData(t, []byte("x"), file)

	m := defaultTestModelWithFilePreview(dir)
	fm := &m.fileModel
	const placeholders, imageData = "kitty-placeholder-cells", "kitty-image-data"
	msg := preview.NewUpdateMsg(file, placeholders, imageData, fm.ExpectedPreviewWidth, fm.Height, 1)

	cmd := fm.UpdatePreviewPanel(msg)
	assert.NotContains(t, fm.FilePreview.GetContent(), placeholders,
		"placeholders must not be shown before the image data is sent")

	msgs := runSequence(t, cmd)
	require.Len(t, msgs, 2)
	assert.Equal(t, tea.RawMsg{Msg: imageData}, msgs[0], "image data is sent first")

	follow, ok := msgs[1].(preview.UpdateMsg)
	require.True(t, ok, "then the preview content is applied")
	assert.Nil(t, fm.UpdatePreviewPanel(follow), "image data is sent only once")
	assert.Equal(t, placeholders, fm.FilePreview.GetContent())
}
