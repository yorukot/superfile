package internal

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/yorukot/superfile/src/internal/common"
)

var (
	lastClickTime time.Time
	lastClickPos  [2]int
)

func (m *model) handleMouseMsg(msg tea.MouseMsg) tea.Cmd {
	mouse := msg.Mouse()
	msgStr := msg.String()

	if msgStr == "wheelup" || msgStr == "wheeldown" {
		m.updateFocusByCoords(mouse.X, mouse.Y)
		wheelMainAction(msgStr, m)
		return nil
	}

	if mouse.Button == tea.MouseLeft {
		return m.handleLeftClick(mouse.X, mouse.Y)
	}

	return nil
}

func (m *model) handleLeftClick(x, y int) tea.Cmd {
	m.updateFocusByCoords(x, y)

	now := time.Now()
	isDoubleClick := now.Sub(lastClickTime) < 400*time.Millisecond && lastClickPos[0] == x && lastClickPos[1] == y
	lastClickTime = now
	lastClickPos = [2]int{x, y}

	if m.focusPanel == sidebarFocus {
		m.handleSidebarClick(y)
	} else if m.focusPanel == nonePanelFocus {
		return m.handleFilePanelClick(x, y, isDoubleClick)
	}

	return nil
}

func (m *model) updateFocusByCoords(x, y int) {
	// Click in footer
	if m.toggleFooter && y >= m.mainPanelHeight {
		processWidth := m.processBarModel.GetWidth()
		if x < processWidth {
			m.focusOnProcessBar()
		} else {
			m.focusOnMetadata()
		}
		return
	}

	// Click in sidebar
	if x < common.Config.SidebarWidth {
		m.focusOnSideBar()
		return
	}

	// Click in main file panels
	m.focusPanel = nonePanelFocus
	relX := x - common.Config.SidebarWidth
	panelCount := m.fileModel.PanelCount()

	accumX := 0
	for i := 0; i < panelCount; i++ {
		pw := m.fileModel.FilePanels[i].GetWidth()
		if relX >= accumX && relX < accumX+pw {
			m.fileModel.SetFocusedPanelIndex(i)
			break
		}
		accumX += pw
	}
}

func (m *model) handleSidebarClick(y int) {
	// Line 0: top border
	// Line 1: superfile title
	// Line 2: blank line
	// Line 3: searchbar (if rendered)
	headerLines := 3
	if m.sidebarModel.SearchBarFocused() {
		headerLines = 4
	}
	targetIndex := y - headerLines
	if targetIndex >= 0 {
		if m.sidebarModel.SetCursor(targetIndex) {
			m.sidebarSelectDirectory()
		}
	}
}

func (m *model) handleFilePanelClick(x, y int, isDoubleClick bool) tea.Cmd {
	panel := m.getFocusedFilePanel()
	// Line 0: top border
	// Line 1: path title
	// Line 2: section divider
	// Line 3: search bar
	topOffset := 4
	if panel.NeedRenderHeaders() {
		topOffset = 5
	}

	targetRow := y - topOffset
	if targetRow >= 0 {
		targetIndex := panel.RenderIndex() + targetRow
		if targetIndex >= 0 && targetIndex < panel.ElemCount() {
			if panel.GetCursor() == targetIndex && isDoubleClick {
				m.enterPanel()
			} else {
				panel.SetCursorPosition(targetIndex)
			}
		}
	}
	return nil
}
