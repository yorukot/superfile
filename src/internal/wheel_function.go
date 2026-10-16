package internal

import (
	"log/slog"

	"github.com/yorukot/superfile/src/internal/common"
)

func wheelMainAction(msg string, m *model) {
	slog.Debug("wheelMainAction called", "msg", msg, "focusPanel", m.focusPanel)
	var action func()
	switch msg {
	case "wheelup":
		switch m.focusPanel {
		case sidebarFocus:
			action = func() { m.sidebarModel.ListUp() }
		case processBarFocus:
			action = func() { m.processBarModel.ListUp() }
		case metadataFocus:
			action = func() { m.fileMetaData.ListUp() }
		case previewFocus:
			action = func() { m.fileModel.FilePreview.ScrollUp() }
		case nonePanelFocus:
			action = func() { m.getFocusedFilePanel().ListUp() }
		}

	case "wheeldown":
		switch m.focusPanel {
		case sidebarFocus:
			action = func() { m.sidebarModel.ListDown() }
		case processBarFocus:
			action = func() { m.processBarModel.ListDown() }
		case metadataFocus:
			action = func() { m.fileMetaData.ListDown() }
		case previewFocus:
			action = func() { m.fileModel.FilePreview.ScrollDown() }
		case nonePanelFocus:
			action = func() { m.getFocusedFilePanel().ListDown() }
		}
	default:
		slog.Error("Unexpected type of mouse action in wheelMainAction", "msg", msg)
		return
	}

	scrollLines := common.Config.WheelScrollSpeed
	if scrollLines <= 0 {
		return
	}
	for range scrollLines {
		action()
	}
}
