package internal

import (
	"path/filepath"
	"strconv"

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
	return common.FirstUseModal(m.helpMenu.GetHeight(), m.helpMenu.GetWidth()).
		Render(title + "\n\n" + vimUserWarn + "\n\n" + subOne + "\n\n" +
			subTwo + "\n\n" + subThree + "\n\n" + subFour + "\n\n")
}

func (m *model) promptModalRender() string {
	return m.promptModal.Render()
}

func (m *model) zoxideModalRender() string {
	return m.zoxideModal.Render()
}
