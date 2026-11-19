package openwithmodel

import (
	"fmt"
	"strings"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func (m *Model) Render() string {
	var content strings.Builder
	content.WriteString(common.ModalTitleStyle.Render(" Open With") + "\n\n")
	if len(m.apps) == 0 {
		content.WriteString(common.ModalStyle.Render(" No applications configured") + "\n")
	} else {
		for i, application := range m.apps {
			cursor := " "
			if i == m.cursor {
				cursor = common.FilePanelCursorStyle.Render(icon.Cursor)
			}
			label := common.TruncateText(application.Name, m.Width-common.InnerPadding, "...")
			content.WriteString(cursor + common.ModalStyle.Render(" "+label) + "\n")
		}
	}

	position := "0/0"
	if len(m.apps) > 0 {
		position = fmt.Sprintf("%d/%d", m.cursor+1, len(m.apps))
	}
	bottomBorder := common.GenerateFooterBorder(position, m.Width-common.BorderPadding)
	return common.SortOptionsModalBorderStyle(m.Height, m.Width, bottomBorder).Render(content.String())
}
