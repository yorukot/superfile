package helpmenu

import (
	"fmt"
	"strconv"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
)

func (m *Model) Render() string {
	r := ui.HelpMenuRenderer(m.height, m.width)
	r.AddLines(" " + m.searchBar.View())
	r.AddLines("") // one-line separation between searchbar and content

	// TODO : This computation should not happen at render time. Move this to update
	// TODO : Move these computations to a utility function
	maxKeyLength := 0
	for _, data := range m.filteredData {
		totalKeyLen := 0
		for _, key := range data.hotkey {
			totalKeyLen += len(key)
		}

		separatorLen := max(0, (len(data.hotkey)-1)) * common.FooterGroupCols
		if data.subTitle == "" && totalKeyLen+separatorLen > maxKeyLength {
			maxKeyLength = totalKeyLen + separatorLen
		}
	}

	valueLength := m.width - maxKeyLength - common.BorderPadding
	if valueLength < m.width/common.CenterDivisor {
		valueLength = m.width/common.CenterDivisor - common.BorderPadding
	}

	totalTitleCount := 0
	cursorBeenTitleCount := 0

	for i, data := range m.filteredData {
		if data.subTitle != "" {
			if i < m.cursor {
				cursorBeenTitleCount++
			}
			totalTitleCount++
		}
	}

	renderHotkeyLength := m.getRenderHotkeyLength()
	m.getContent(r, renderHotkeyLength, valueLength)

	current := m.cursor + 1 - cursorBeenTitleCount
	if len(m.filteredData) == 0 {
		current = 0
	}
	r.SetBorderInfoItems(fmt.Sprintf("%s/%s",
		strconv.Itoa(current),
		strconv.Itoa(len(m.filteredData)-totalTitleCount)))
	return r.Render()
}

func (m *Model) getRenderHotkeyLength() int {
	renderHotkeyLength := 0
	for i := m.renderIndex; i < m.renderIndex+(m.height-common.InnerPadding) && i < len(m.filteredData); i++ {
		if m.filteredData[i].subTitle != "" {
			continue
		}

		hotkey := common.GetHelpMenuHotkeyString(m.filteredData[i].hotkey)

		renderHotkeyLength = max(renderHotkeyLength, len(common.HelpMenuHotkeyStyle.Render(hotkey)))
	}
	return renderHotkeyLength + 1
}

func (m *Model) getContent(r *rendering.Renderer, renderHotkeyLength int, valueLength int) {
	for i := m.renderIndex; i < m.renderIndex+(m.height-common.InnerPadding) && i < len(m.filteredData); i++ {
		if m.filteredData[i].subTitle != "" {
			r.AddLines(common.HelpMenuTitleStyle.Render(" " + m.filteredData[i].subTitle))
			continue
		}

		hotkey := common.GetHelpMenuHotkeyString(m.filteredData[i].hotkey)
		description := common.TruncateText(m.filteredData[i].description, valueLength, "...")

		cursor := "  "
		if m.cursor == i {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor + " ")
		}
		r.AddLines(cursor + common.ModalStyle.Render(fmt.Sprintf("%*s%s", renderHotkeyLength,
			common.HelpMenuHotkeyStyle.Render(hotkey+" "), common.ModalStyle.Render(description))))
	}
}
