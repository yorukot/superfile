package zoxide

import (
	"fmt"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func (m *Model) Render() string {
	r := ui.ZoxideRenderer(m.maxHeight, m.width)
	r.SetBorderTitle(m.headline)

	if m.zClient == nil {
		r.AddSection()
		r.AddLines(" Zoxide not available (check zoxide_support in config)")
		return r.Render()
	}

	r.AddLines(" " + m.textInput.View())
	r.AddSection()
	if len(m.results) > 0 {
		m.renderResultList(r)
	} else {
		r.AddLines(" No zoxide results found")
	}
	return r.Render()
}

func (m *Model) renderResultList(r *rendering.Renderer) {
	// Calculate visible range
	endIndex := m.renderIndex + maxVisibleResults
	if endIndex > len(m.results) {
		endIndex = len(m.results)
	}
	// Show visible results
	m.renderVisibleResults(r, endIndex)

	// Show scroll indicators if needed
	m.renderScrollIndicators(r, endIndex)
}

func (m *Model) renderVisibleResults(r *rendering.Renderer, endIndex int) {
	for i := m.renderIndex; i < endIndex; i++ {
		result := m.results[i]

		// Truncate path if too long (account for score, separator, and padding)
		// Available width: modal width
		// - borders(2) - padding(2) - score(6)
		// - separator(3) = width - 13
		// 0123456789012345678 => 19 width, path gets 6
		// | 9999.9 | <path> |
		availablePathWidth := m.width - ScoreColumnWidth
		path := common.TruncateTextBeginning(result.Path, availablePathWidth, "...")

		line := fmt.Sprintf(" %6.1f | %s", result.Score, path)

		// Highlight the selected item
		if i == m.cursor {
			line = common.ModalCursorStyle.Render(line)
		}
		r.AddLines(line)
	}
}

func (m *Model) renderScrollIndicators(r *rendering.Renderer, endIndex int) {
	if len(m.results) <= maxVisibleResults {
		return
	}

	if m.renderIndex > 0 {
		r.AddSection()
		r.AddLines(" ↑ More results above")
	}
	if endIndex < len(m.results) {
		if m.renderIndex == 0 {
			r.AddSection()
		}
		r.AddLines(" ↓ More results below")
	}
}
