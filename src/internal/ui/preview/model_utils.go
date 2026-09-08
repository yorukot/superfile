package preview

import (
	"log/slog"
	"strings"
)

func (m *Model) GetContent() string {
	return m.content
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
	m.lines = nil
	m.renderIndex = 0
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
	m.resetScroll()
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
func (m *Model) IsTextPreview() bool {
	return m.lines != nil
}

func (m *Model) ScrollUp() {
	if m.lines == nil || m.renderIndex <= 0 {
		return
	}
	m.renderIndex--
}

func (m *Model) ScrollDown() {
	if m.lines == nil {
		return
	}
	maxIndex := len(m.lines) - m.contentHeight
	if maxIndex < 0 {
		maxIndex = 0
	}
	if m.renderIndex < maxIndex {
		m.renderIndex++
	}
}

func (m *Model) PgUp() {
	if m.lines == nil {
		return
	}
	m.renderIndex -= m.contentHeight
	if m.renderIndex < 0 {
		m.renderIndex = 0
	}
}

func (m *Model) PgDown() {
	if m.lines == nil {
		return
	}
	maxIndex := len(m.lines) - m.contentHeight
	if maxIndex < 0 {
		maxIndex = 0
	}
	m.renderIndex += m.contentHeight
	if m.renderIndex > maxIndex {
		m.renderIndex = maxIndex
	}
}

func (m *Model) resetScroll() {
	m.renderIndex = 0
}

// GetScrollRender returns the currently visible window of a scrolled text preview.
// For non-text previews it falls back to the pre-rendered content.
func (m *Model) GetScrollRender() string {
	if m.lines == nil {
		return m.content
	}
	start := m.renderIndex
	end := start + m.contentHeight
	if end > len(m.lines) {
		end = len(m.lines)
	}
	window := strings.Join(m.lines[start:end], "\n")
	return m.RenderTextWithDimension(window, m.contentHeight, m.contentWidth)
}
