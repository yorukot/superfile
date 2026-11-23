package filemodel

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/preview"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

const (
	testWidth        = 100
	testPreviewWidth = 40
	testHeight       = 20
)

// newSixelTestModel returns a model whose preview shows a Sixel image.
func newSixelTestModel(t *testing.T) *Model {
	t.Helper()
	oldSidebar, oldBorder := common.Config.SidebarWidth, common.Config.EnableFilePreviewBorder
	common.Config.SidebarWidth = 0
	common.Config.EnableFilePreviewBorder = false
	t.Cleanup(func() {
		common.Config.SidebarWidth, common.Config.EnableFilePreviewBorder = oldSidebar, oldBorder
	})

	m := &Model{Width: testWidth, ExpectedPreviewWidth: testPreviewWidth, Height: testHeight}
	m.FilePreview.Open()
	applySixel(m, "/img.png")
	return m
}

func applySixel(m *Model, location string) {
	result := preview.RenderResult{
		Content: strings.Repeat(filepreview.SixelMarker, 4),
		Sixel:   &filepreview.SixelImage{DCS: "<dcs>", OffsetX: 2, OffsetY: 3, Cols: 4, Rows: 1},
	}
	m.FilePreview.Apply(preview.NewUpdateMsg(location, result, testPreviewWidth, testHeight, 0))
	m.InvalidateSixel()
}

// runDraw runs a scheduled draw command (waiting for its delay) and returns
// what HandleSixelDraw makes of the resulting message.
func runDraw(t *testing.T, m *Model, cmd tea.Cmd) tea.Cmd {
	t.Helper()
	require.NotNil(t, cmd)
	msg, ok := cmd().(SixelDrawMsg)
	require.True(t, ok)
	return m.HandleSixelDraw(msg)
}

func TestSyncSixelDrawsOnceAtPreviewPosition(t *testing.T) {
	m := newSixelTestModel(t)

	cmd := m.SyncSixel(false)
	assert.Nil(t, m.SyncSixel(false), "no second draw while the image is still shown")

	draw := runDraw(t, m, cmd)
	require.NotNil(t, draw)
	raw, ok := draw().(tea.RawMsg)
	require.True(t, ok)
	// Preview starts at column 100-40 = 60; offset (2, 3) gives 0-based (62, 3).
	assert.Equal(t, "\x1b7\x1b[4;63H<dcs>\x1b8", raw.Msg)
}

func TestSyncSixelHiddenWhileCovered(t *testing.T) {
	m := newSixelTestModel(t)
	first := m.SyncSixel(false)

	assert.Nil(t, m.SyncSixel(true))
	assert.Equal(t, "    ", m.FilePreview.GetContent(),
		"hidden image's marker cells become spaces, so the renderer erases it")
	assert.Nil(t, runDraw(t, m, first), "a draw scheduled before hiding is dropped")

	again := m.SyncSixel(false)
	assert.True(t, hasMarker(m.FilePreview.GetContent()))
	assert.NotNil(t, runDraw(t, m, again), "the image is drawn again once uncovered")
}

func hasMarker(s string) bool {
	return strings.Contains(s, filepreview.SixelMarker) || strings.Contains(s, filepreview.SixelMarkerAlt)
}

func TestApplySixelAlternatesMarkers(t *testing.T) {
	m := newSixelTestModel(t)
	first := m.FilePreview.GetContent()
	applySixel(m, "/img.png")
	second := m.FilePreview.GetContent()

	assert.True(t, hasMarker(first))
	assert.True(t, hasMarker(second))
	assert.NotEqual(t, first, second,
		"consecutive images use different markers, so every cell is rewritten between them")
}

func TestSyncSixelDropsStaleDraws(t *testing.T) {
	m := newSixelTestModel(t)
	first := m.SyncSixel(false)

	// The user moves on to another image before the first is drawn.
	applySixel(m, "/other.png")
	second := m.SyncSixel(false)

	assert.Nil(t, runDraw(t, m, first))
	assert.NotNil(t, runDraw(t, m, second))
}

func TestSyncSixelRedrawsAfterInvalidate(t *testing.T) {
	m := newSixelTestModel(t)
	require.NotNil(t, m.SyncSixel(false))

	m.InvalidateSixel() // e.g. after a resize repainted the screen
	assert.NotNil(t, m.SyncSixel(false))
}

func TestSyncSixelNothingToDraw(t *testing.T) {
	m := newSixelTestModel(t)

	m.FilePreview.SetLoading()
	assert.Nil(t, m.SyncSixel(false), "not drawn while loading the next preview")

	m.FilePreview.Apply(preview.NewUpdateMsg("/a.txt", preview.RenderResult{Content: "text"},
		testPreviewWidth, testHeight, 1))
	assert.Nil(t, m.SyncSixel(false), "not drawn for non-image content")

	applySixel(m, "/img.png")
	m.FilePreview.Close()
	assert.Nil(t, m.SyncSixel(false), "not drawn while the preview is closed")
}
