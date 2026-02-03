package pinnedmodal

import (
	"fmt"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func (m *Model) Render() string {
	r := ui.ZoxideRenderer(m.maxHeight, m.width)
	r.SetBorderTitle(m.headline)

	r.AddLines(" " + m.textInput.View())
	r.AddSection()
	if len(m.results) > 0 {
		m.renderResultList(r)
	} else {
		r.AddLines(" No pinned directories found")
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
		dir := m.results[i]

		availablePathWidth := m.width - columnWidth
		name := common.TruncateTextBeginning(dir.Name, availablePathWidth/2, "...")
		path := common.TruncateTextBeginning(dir.Location, availablePathWidth/2, "...")

		line := fmt.Sprintf(" %-15s | %s", name, path)

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
