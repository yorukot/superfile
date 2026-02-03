package gotointeractive

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func (m *Model) Render() string {
	r := ui.PromptRenderer(m.maxHeight, m.width)

	separator := " "
	headlineWidth := lipgloss.Width(m.headline)
	separatorWidth := lipgloss.Width(separator)
	availableWidth := m.width - headlineWidth - separatorWidth
	titlePath := m.currentPath
	if lipgloss.Width(titlePath) > availableWidth {
		titlePath = common.TruncateTextBeginning(titlePath, availableWidth, "...")
	}
	r.SetBorderTitle(m.headline + separator + titlePath)

	r.AddLines(" " + m.textInput.View())
	r.AddSection()

	if len(m.results) > 0 {
		m.renderResultList(r)
	} else {
		r.AddLines(" No results found")
	}
	return r.Render()
}

func (m *Model) renderResultList(r *rendering.Renderer) {
	endIndex := m.renderIndex + maxVisibleResults
	if endIndex > len(m.results) {
		endIndex = len(m.results)
	}
	m.renderVisibleResults(r, endIndex)
	m.renderScrollIndicators(r, endIndex)
}

func (m *Model) renderVisibleResults(r *rendering.Renderer, endIndex int) {
	for i := m.renderIndex; i < endIndex; i++ {
		result := m.results[i]

		availablePathWidth := m.width - iconColumnWidth
		if availablePathWidth < 0 {
			availablePathWidth = 0
		}
		displayName := common.TruncateTextBeginning(result.Name, availablePathWidth, "...")

		var line string
		if result.IsDir {
			line = icon.Directory + icon.Space + displayName + "/"
		} else {
			line = icon.Icons["file"].Icon + icon.Space + displayName
		}

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
