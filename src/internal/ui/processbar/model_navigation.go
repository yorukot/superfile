package processbar

import (
	"github.com/yorukot/superfile/src/internal/common"
)

// Control processbar panel list up
// There is a shadowing happening here, but it will be removed
// Once we make footerHeight part of model struct
func (m *Model) ListUp() {
	cntP := m.cntProcesses()
	if cntP == 0 {
		return
	}
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.renderIndex {
			m.renderIndex--
		}
	} else {
		m.cursor = cntP - 1
		// Either start from beginning or
		// from a process so that we could render last one
		m.renderIndex = max(0, cntP-m.cntRenderableProcess())
	}
}

// Control processbar panel list down
func (m *Model) ListDown() {
	cntP := m.cntProcesses()
	if cntP == 0 {
		return
	}
	if m.cursor < cntP-1 {
		m.cursor++
		if m.cursor > m.renderIndex+m.cntRenderableProcess()-1 {
			m.renderIndex++
		}
	} else {
		m.renderIndex = 0
		m.cursor = 0
	}
}

func (m *Model) cntRenderableProcess() int {
	footerHeight := m.height - common.BorderPadding
	return cntRenderableProcess(footerHeight)
}

func cntRenderableProcess(footerHeight int) int {
	// We can render one process in three lines
	// And last process in two or three lines ( with/without a line separtor)
	return (footerHeight + 1) / linesPerProcess
}
