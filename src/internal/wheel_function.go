package internal

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/common"
)

func wheelMainAction(msg string, m *model) {
	slog.Debug("wheelMainAction called", "msg", msg, "focusPanel", m.focusPanel)
	var action func()
	switch msg {
	case "wheel up":
		switch m.focusPanel {
		case sidebarFocus:
			action = func() { m.sidebarModel.ListUp(m.mainPanelHeight) }
		case processBarFocus:
			action = func() { m.processBarModel.listUp(m.footerHeight) }
		case metadataFocus:
			action = func() { m.fileMetaData.ListUp() }
		case nonePanelFocus:
			action = func() { m.fileModel.filePanels[m.filePanelFocusIndex].listUp(m.mainPanelHeight) }
		}

	case "wheel down":
		switch m.focusPanel {
		case sidebarFocus:
			action = func() { m.sidebarModel.ListDown(m.mainPanelHeight) }
		case processBarFocus:
			action = func() { m.processBarModel.listDown(m.footerHeight) }
		case metadataFocus:
			action = func() { m.fileMetaData.ListDown() }
		case nonePanelFocus:
			action = func() { m.fileModel.filePanels[m.filePanelFocusIndex].listDown(m.mainPanelHeight) }
		}
	default:
		slog.Error("Unexpected type of mouse action in wheelMainAction", "msg", msg)
		return
	}

	for range common.WheelRunTime {
		action()
	}
}
