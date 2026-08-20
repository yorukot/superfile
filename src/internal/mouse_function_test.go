package internal

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/yorukot/superfile/src/internal/common"
)

func TestMouseLeftClickFocusPanel(t *testing.T) {
	testPaths := []string{t.TempDir()}
	m := defaultModelConfig(true, true, false, testPaths, nil)
	m.fullWidth = 120
	m.fullHeight = 40
	m.setHeightValues()
	m.updateComponentDimensions()
	m.firstLoadingComplete = true

	// Click sidebar header area
	m.updateFocusByCoords(5, 0)
	assert.Equal(t, sidebarFocus, m.focusPanel)

	// Click file panel area (x=common.Config.SidebarWidth + 10, y=5)
	mouseMsg := tea.MouseClickMsg{
		X:      common.Config.SidebarWidth + 10,
		Y:      5,
		Button: tea.MouseLeft,
	}

	_ = m.handleMouseMsg(mouseMsg)
	assert.Equal(t, nonePanelFocus, m.focusPanel)
}
