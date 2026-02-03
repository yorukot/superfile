package gotointeractive

import (
	"os"
	"path/filepath"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func (m *Model) Render() string {
	r := ui.PromptRenderer(m.maxHeight, m.width)

	availableWidth := m.width - len(m.headline) - 3 //nolint:mnd // Space for " - " separator
	titlePath := m.currentPath
	if len(titlePath) > availableWidth {
		titlePath = common.TruncateTextBeginning(titlePath, availableWidth, "...")
	}
	r.SetBorderTitle(m.headline + " " + titlePath)

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

		isDir := false
		if result == ".." {
			isDir = true
		} else {
			fullPath := filepath.Join(m.currentPath, result)
			if info, err := os.Stat(fullPath); err == nil {
				isDir = info.IsDir()
			}
		}

		availablePathWidth := m.width - iconColumnWidth
		displayName := common.TruncateTextBeginning(result, availablePathWidth, "...")

		var line string
		if isDir {
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
