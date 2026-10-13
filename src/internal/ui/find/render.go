package find

import (
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

const (
	scrollUpIndicator   = " ↑ More results above"
	scrollDownIndicator = " ↓ More results below"
)

func (m *Model) Render() string {
	r := ui.FindRenderer(m.maxHeight, m.width)
	r.SetBorderTitle(m.headline)

	if !m.isAvailable() {
		r.AddSection()
		r.AddLines(" fd is not available (install fd and enable find_file_support)")
		return r.Render()
	}

	r.AddLines(" " + m.textInput.View())
	r.AddSection()
	switch {
	case m.errMsg != "":
		r.AddLines(" " + m.errMsg)
	case len(m.results) > 0:
		m.renderResultList(r)
	default:
		r.AddLines(" No results found")
	}
	return r.Render()
}

func (m *Model) resultMarker(isDir bool) string {
	if isDir {
		return icon.Directory
	}
	return icon.Icons["file"].Icon
}

func (m *Model) renderResultList(r *rendering.Renderer) {
	// Calculate visible range
	visible := m.visibleResults()
	endIndex := m.renderIndex + visible
	endIndex = min(endIndex, len(m.results))
	// Show visible results
	m.renderVisibleResults(r, endIndex)

	// Show scroll indicators if needed
	m.renderScrollIndicators(r, endIndex, visible)
}

func (m *Model) renderVisibleResults(r *rendering.Renderer, endIndex int) {
	for i := m.renderIndex; i < endIndex; i++ {
		result := m.results[i]

		// Truncate path if too long (account for marker, separator, and padding)
		// Available width: modal width
		// - borders(2) - padding(2) - marker(1)
		// - separator(3) = width - 8
		availablePathWidth := m.width - markerColumnWidth
		path := common.TruncateTextBeginning(result.Path, availablePathWidth, "...")

		line := " " + m.resultMarker(result.IsDir) + " | " + path

		// Highlight the selected item
		if i == m.cursor {
			line = common.ModalCursorStyle.Render(line)
		}
		r.AddLines(line)
	}
}

func (m *Model) renderScrollIndicators(r *rendering.Renderer, endIndex int, visible int) {
	if len(m.results) <= visible {
		return
	}

	if m.renderIndex > 0 {
		r.AddSection()
		r.AddLines(scrollUpIndicator)
	}
	if endIndex < len(m.results) {
		if m.renderIndex == 0 {
			r.AddSection()
		}
		r.AddLines(scrollDownIndicator)
	}
}
