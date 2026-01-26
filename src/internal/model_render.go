package internal

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"

	"github.com/yorukot/superfile/src/internal/common"

	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
)

func (m *model) sidebarRender() string {
	return m.sidebarModel.Render(m.focusPanel == sidebarFocus,
		m.getFocusedFilePanel().Location)
}

func (m *model) processBarRender() string {
	return m.processBarModel.Render(m.focusPanel == processBarFocus)
}

func (m *model) terminalSizeWarnRender() string {
	fullWidthString := strconv.Itoa(m.fullWidth)
	fullHeightString := strconv.Itoa(m.fullHeight)
	minimumWidthString := strconv.Itoa(common.MinimumWidth)
	minimumHeightString := strconv.Itoa(common.MinimumHeight)
	if m.fullHeight < common.MinimumHeight {
		fullHeightString = common.TerminalTooSmall.Render(fullHeightString)
	}
	if m.fullWidth < common.MinimumWidth {
		fullWidthString = common.TerminalTooSmall.Render(fullWidthString)
	}
	fullHeightString = common.TerminalCorrectSize.Render(fullHeightString)
	fullWidthString = common.TerminalCorrectSize.Render(fullWidthString)

	heightString := common.MainStyle.Render(" Height = ")
	return common.FullScreenStyle(m.fullHeight, m.fullWidth).Render(`Terminal size too small:`+"\n"+
		"Width = "+fullWidthString+
		heightString+fullHeightString+"\n\n"+
		"Needed for current config:"+"\n"+
		"Width = "+common.TerminalCorrectSize.Render(minimumWidthString)+
		heightString+common.TerminalCorrectSize.Render(minimumHeightString)) + filepreview.ClearKittyImages()
}

func (m *model) terminalSizeWarnAfterFirstRender() string {
	minimumWidthInt := common.Config.SidebarWidth + common.FilePanelWidthUnit*len(
		m.fileModel.FilePanels,
	) + common.FilePanelWidthUnit - 1
	minimumWidthString := strconv.Itoa(minimumWidthInt)
	fullWidthString := strconv.Itoa(m.fullWidth)
	fullHeightString := strconv.Itoa(m.fullHeight)
	minimumHeightString := strconv.Itoa(common.MinimumHeight)

	if m.fullHeight < common.MinimumHeight {
		fullHeightString = common.TerminalTooSmall.Render(fullHeightString)
	}
	if m.fullWidth < minimumWidthInt {
		fullWidthString = common.TerminalTooSmall.Render(fullWidthString)
	}
	fullHeightString = common.TerminalCorrectSize.Render(fullHeightString)
	fullWidthString = common.TerminalCorrectSize.Render(fullWidthString)

	heightString := common.MainStyle.Render(" Height = ")
	return common.FullScreenStyle(m.fullHeight, m.fullWidth).Render(`You change your terminal size too small:`+"\n"+
		"Width = "+fullWidthString+
		heightString+fullHeightString+"\n\n"+
		"Needed for current config:"+"\n"+
		"Width = "+common.TerminalCorrectSize.Render(minimumWidthString)+
		heightString+common.TerminalCorrectSize.Render(minimumHeightString)) + filepreview.ClearKittyImages()
}

func (m *model) typineModalRender() string {
	previewPath := filepath.Join(m.typingModal.location, m.typingModal.textInput.Value())

	fileLocation := common.FilePanelTopDirectoryIconStyle.Render(" "+icon.Directory+icon.Space) +
		common.FilePanelTopPathStyle.Render(
			common.TruncateTextBeginning(previewPath, common.ModalWidth-common.InnerPadding, "..."),
		) + "\n"

	confirm := common.ModalConfirm.Render(" (" + common.Hotkeys.ConfirmTyping[0] + ") Create ")
	cancel := common.ModalCancel.Render(" (" + common.Hotkeys.CancelTyping[0] + ") Cancel ")

	tip := confirm +
		lipgloss.NewStyle().Background(common.ModalBGColor).Render("           ") +
		cancel

	var err string
	if m.typingModal.errorMesssage != "" {
		err = "\n\n" + common.ModalErrorStyle.Render(m.typingModal.errorMesssage)
	}
	// TODO : Move this all to rendering package to avoid specifying newlines manually
	return common.ModalBorderStyle(common.ModalHeight, common.ModalWidth).
		Render(fileLocation + "\n" + m.typingModal.textInput.View() + "\n\n" + tip + err)
}

func (m *model) introduceModalRender() string {
	title := common.SidebarTitleStyle.Render(" Thanks for using superfile!!") +
		common.ModalStyle.Render("\n You can read the following information before starting to use it!")
	vimUserWarn := common.ProcessErrorStyle.Render("  ** Very importantly ** If you are a Vim/Nvim user, go to:\n" +
		"  https://superfile.dev/configure/custom-hotkeys/ to change your hotkey settings!")
	subOne := common.SidebarTitleStyle.Render("  (1)") +
		common.ModalStyle.Render(" If this is your first time, make sure you read:\n"+
			"      https://superfile.dev/getting-started/tutorial/")
	subTwo := common.SidebarTitleStyle.Render("  (2)") +
		common.ModalStyle.Render(" If you forget the relevant keys during use,\n"+
			"      you can press \"?\" (shift+/) at any time to query the keys!")
	subThree := common.SidebarTitleStyle.Render("  (3)") +
		common.ModalStyle.Render(" For more customization you can refer to:\n"+
			"      https://superfile.dev/")
	subFour := common.SidebarTitleStyle.Render("  (4)") +
		common.ModalStyle.Render(" Thank you again for using superfile.\n"+
			"      If you have any questions, please feel free to ask at:\n"+
			"      https://github.com/yorukot/superfile\n"+
			"      Of course, you can always open a new issue to share your idea \n"+
			"      or report a bug!")
	return common.FirstUseModal(m.helpMenu.height, m.helpMenu.width).
		Render(title + "\n\n" + vimUserWarn + "\n\n" + subOne + "\n\n" +
			subTwo + "\n\n" + subThree + "\n\n" + subFour + "\n\n")
}

func (m *model) promptModalRender() string {
	return m.promptModal.Render()
}

func (m *model) zoxideModalRender() string {
	return m.zoxideModal.Render()
}

func (m *model) helpMenuRender() string {
	r := ui.HelpMenuRenderer(m.helpMenu.height, m.helpMenu.width)
	r.AddLines(" " + m.helpMenu.searchBar.View())
	r.AddLines("") // one-line separation between searchbar and content

	// TODO : This computation should not happen at render time. Move this to update
	// TODO : Move these computations to a utility function
	maxKeyLength := 0
	for _, data := range m.helpMenu.filteredData {
		totalKeyLen := 0
		for _, key := range data.hotkey {
			totalKeyLen += len(key)
		}

		separatorLen := max(0, (len(data.hotkey)-1)) * common.FooterGroupCols
		if data.subTitle == "" && totalKeyLen+separatorLen > maxKeyLength {
			maxKeyLength = totalKeyLen + separatorLen
		}
	}

	valueLength := m.helpMenu.width - maxKeyLength - common.BorderPadding
	if valueLength < m.helpMenu.width/common.CenterDivisor {
		valueLength = m.helpMenu.width/common.CenterDivisor - common.BorderPadding
	}

	totalTitleCount := 0
	cursorBeenTitleCount := 0

	for i, data := range m.helpMenu.filteredData {
		if data.subTitle != "" {
			if i < m.helpMenu.cursor {
				cursorBeenTitleCount++
			}
			totalTitleCount++
		}
	}

	renderHotkeyLength := m.getRenderHotkeyLengthHelpMenuModal()
	m.getHelpMenuContent(r, renderHotkeyLength, valueLength)

	current := m.helpMenu.cursor + 1 - cursorBeenTitleCount
	if len(m.helpMenu.filteredData) == 0 {
		current = 0
	}
	r.SetBorderInfoItems(fmt.Sprintf("%s/%s",
		strconv.Itoa(current),
		strconv.Itoa(len(m.helpMenu.filteredData)-totalTitleCount)))
	return r.Render()
}

func (m *model) getRenderHotkeyLengthHelpMenuModal() int {
	renderHotkeyLength := 0
	for i := m.helpMenu.renderIndex; i < m.helpMenu.renderIndex+(m.helpMenu.height-common.InnerPadding) && i < len(m.helpMenu.filteredData); i++ {
		if m.helpMenu.filteredData[i].subTitle != "" {
			continue
		}

		hotkey := common.GetHelpMenuHotkeyString(m.helpMenu.filteredData[i].hotkey)

		renderHotkeyLength = max(renderHotkeyLength, len(common.HelpMenuHotkeyStyle.Render(hotkey)))
	}
	return renderHotkeyLength + 1
}

func (m *model) getHelpMenuContent(r *rendering.Renderer, renderHotkeyLength int, valueLength int) {
	for i := m.helpMenu.renderIndex; i < m.helpMenu.renderIndex+(m.helpMenu.height-common.InnerPadding) && i < len(m.helpMenu.filteredData); i++ {
		if m.helpMenu.filteredData[i].subTitle != "" {
			r.AddLines(common.HelpMenuTitleStyle.Render(" " + m.helpMenu.filteredData[i].subTitle))
			continue
		}

		hotkey := common.GetHelpMenuHotkeyString(m.helpMenu.filteredData[i].hotkey)
		description := common.TruncateText(m.helpMenu.filteredData[i].description, valueLength, "...")

		cursor := "  "
		if m.helpMenu.cursor == i {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor + " ")
		}
		r.AddLines(cursor + common.ModalStyle.Render(fmt.Sprintf("%*s%s", renderHotkeyLength,
			common.HelpMenuHotkeyStyle.Render(hotkey+" "), common.ModalStyle.Render(description))))
	}
}
