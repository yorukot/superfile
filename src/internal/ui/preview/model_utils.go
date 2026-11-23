package preview

import (
	"log/slog"
	"strings"

	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

func (m *Model) GetContent() string {
	if m.sixelHidden && m.sixel != nil {
		return m.hiddenContent
	}
	return m.content
}

// GetSixel returns the Sixel image shown by the current content, or nil.
func (m *Model) GetSixel() *filepreview.SixelImage {
	return m.sixel
}

// SetSixelHidden hides or shows the Sixel image. While hidden, the marker
// cells are replaced with plain spaces, so the renderer rewrites them, which
// erases the image from the terminal.
func (m *Model) SetSixelHidden(hidden bool) {
	m.sixelHidden = hidden
}

// setSixel records that content shows a Sixel image. Must be called right
// after setContent.
func (m *Model) setSixel(img *filepreview.SixelImage) {
	m.sixel = img
	// Alternate the marker so the cells under this image always differ from
	// the cells under the previous one: every cell is then rewritten,
	// erasing the previous image, before this one is drawn.
	m.sixelAlt = !m.sixelAlt
	if m.sixelAlt {
		m.content = strings.ReplaceAll(m.content, filepreview.SixelMarker, filepreview.SixelMarkerAlt)
	}
	m.hiddenContent = strings.NewReplacer(
		filepreview.SixelMarker, " ",
		filepreview.SixelMarkerAlt, " ",
	).Replace(m.content)
}

func (m *Model) GetContentWidth() int {
	return m.contentWidth
}

func (m *Model) GetContentHeight() int {
	return m.contentHeight
}

func (m *Model) GetLocation() string {
	return m.location
}

func (m *Model) SetOpen(open bool) {
	m.open = open
}

func (m *Model) SetLocation(location string) {
	m.location = location
}

func (m *Model) SetLoading() {
	m.loading = true
}

// All content change happen via this only, to ensure the sync between
// content and width x height, and the loading variable reset
func (m *Model) setContent(content string, width int, height int, location string) {
	m.content = content
	m.contentWidth = width
	m.contentHeight = height
	m.location = location
	m.loading = false
	m.sixel = nil
	m.hiddenContent = ""
}

func (m *Model) SetEmptyWithDimensions(width int, height int) {
	m.setContent(m.RenderTextWithDimension("", height, width), width, height, "")
}

func (m *Model) IsLoading() bool {
	return m.loading
}

func (m *Model) ToggleOpen() {
	m.open = !m.open
}

func (m *Model) CleanUp() {
	if m.thumbnailGenerator != nil {
		err := m.thumbnailGenerator.CleanUp()
		if err != nil {
			slog.Error("Error While cleaning up TempDirectory", "error", err)
		}
	}
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Open() {
	m.open = true
}

func (m *Model) Close() {
	m.open = false
}
