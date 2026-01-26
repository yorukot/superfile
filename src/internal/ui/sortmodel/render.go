package sortmodel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

func (m *Model) Render() string {
	var sortOptionsContent strings.Builder
	sortOptionsContent.WriteString(common.ModalTitleStyle.Render(" Sort Options") + "\n\n")
	for i, option := range SortOptionsStr {
		cursor := " "
		if i == m.Cursor {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor)
		}
		sortOptionsContent.WriteString(cursor + common.ModalStyle.Render(" "+option) + "\n")
	}
	bottomBorder := common.GenerateFooterBorder(
		fmt.Sprintf("%s/%s", strconv.Itoa(m.Cursor+1),
			strconv.Itoa(len(SortOptionsStr))), m.Width-common.BorderPadding)

	return common.SortOptionsModalBorderStyle(m.Height, m.Width,
		bottomBorder).Render(sortOptionsContent.String())
}
