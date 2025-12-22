package internal

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
)

func (m *model) sidebarRender() string {
	return m.sidebarModel.Render(m.mainPanelHeight, m.focusPanel == sidebarFocus,
		m.fileModel.filePanels[m.filePanelFocusIndex].Location)
}

// This also modifies the m.fileModel.filePanels, which it should not
// what modifications we do on this model object are of no consequence.
// Since bubblea passed this 'model' by value in View() function.
func (m *model) filePanelRender() string {
	f := make([]string, len(m.fileModel.filePanels))
	for i, filePanel := range m.fileModel.filePanels {
		// check if cursor or render out of range
		// TODO - instead of this, have a filepanel.validateAndFix(), and log Error
		// This should not ever happen
		if filePanel.Cursor > len(filePanel.Element)-1 {
			filePanel.Cursor = 0
			filePanel.RenderIndex = 0
		}
		m.fileModel.filePanels[i] = filePanel

		// TODO : Move this to a utility function and clarify the calculation via comments
		// Maybe even write unit tests
		var filePanelWidth int
		if (m.fullWidth-common.Config.SidebarWidth-(common.InnerPadding+(len(m.fileModel.filePanels)-1)*common.BorderPadding))%len(
			m.fileModel.filePanels,
		) != 0 &&
			i == len(m.fileModel.filePanels)-1 {
			if m.fileModel.filePreview.IsOpen() {
				filePanelWidth = m.fileModel.width
			} else {
				filePanelWidth = (m.fileModel.width + (m.fullWidth-common.Config.SidebarWidth-
					(common.InnerPadding+(len(m.fileModel.filePanels)-1)*common.BorderPadding))%len(m.fileModel.filePanels))
			}
		} else {
			filePanelWidth = m.fileModel.width
		}
		filePanel.UpdateDimensions(filePanelWidth+common.BorderPadding, m.mainPanelHeight+common.BorderPadding)

		f[i] = filePanel.Render(filePanel.IsFocused)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, f...)
}

func (m *model) processBarRender() string {
	return m.processBarModel.Render(m.focusPanel == processBarFocus)
}

func (m *model) clipboardRender() string {
	// render
	var bottomWidth int
	if m.fullWidth%3 != 0 {
		bottomWidth = utils.FooterWidth(m.fullWidth + m.fullWidth%common.FooterGroupCols + common.BorderPadding)
	} else {
		bottomWidth = utils.FooterWidth(m.fullWidth)
	}
	r := ui.ClipboardRenderer(m.footerHeight+common.BorderPadding, bottomWidth+common.BorderPadding)
	if len(m.copyItems.items) == 0 {
		// TODO move this to a string
		r.AddLines("", " "+icon.Error+"  No content in clipboard")
	} else {
		for i := 0; i < len(m.copyItems.items) && i < m.footerHeight; i++ {
			if i == m.footerHeight-1 && i != len(m.copyItems.items)-1 {
				// Last Entry we can render, but there are more that one left
				r.AddLines(strconv.Itoa(len(m.copyItems.items)-i) + " item left....")
			} else {
				fileInfo, err := os.Lstat(m.copyItems.items[i])
				if err != nil {
					slog.Error("Clipboard render function get item state ", "error", err)
				}
				if !os.IsNotExist(err) {
					isLink := fileInfo.Mode()&os.ModeSymlink != 0
					// TODO : There is an inconsistency in parameter that is being passed,
					// and its name in ClipboardPrettierName function
					r.AddLines(common.ClipboardPrettierName(m.copyItems.items[i],
						utils.FooterWidth(m.fullWidth)-common.PanelPadding, fileInfo.IsDir(), isLink, false))
				}
			}
		}
	}
	return r.Render()
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
		m.fileModel.filePanels,
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

	renderHotkeyLength := m.getRenderHotkeyLengthHelpmenuModal()
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

func (m *model) getRenderHotkeyLengthHelpmenuModal() int {
	renderHotkeyLength := 0
	for i := m.helpMenu.renderIndex; i < m.helpMenu.renderIndex+(m.helpMenu.height-common.InnerPadding) && i < len(m.helpMenu.filteredData); i++ {
		hotkey := ""

		if m.helpMenu.filteredData[i].subTitle != "" {
			continue
		}

		for i, key := range m.helpMenu.filteredData[i].hotkey {
			if i != 0 {
				hotkey += " | "
			}
			hotkey += key
		}

		renderHotkeyLength = max(renderHotkeyLength, len(common.HelpMenuHotkeyStyle.Render(hotkey)))
	}
	return renderHotkeyLength
}

func (m *model) getHelpMenuContent(r *rendering.Renderer, renderHotkeyLength int, valueLength int) {
	for i := m.helpMenu.renderIndex; i < m.helpMenu.renderIndex+(m.helpMenu.height-common.InnerPadding) && i < len(m.helpMenu.filteredData); i++ {
		if m.helpMenu.filteredData[i].subTitle != "" {
			r.AddLines(common.HelpMenuTitleStyle.Render(" " + m.helpMenu.filteredData[i].subTitle))
			continue
		}

		hotkey := ""
		description := common.TruncateText(m.helpMenu.filteredData[i].description, valueLength, "...")

		for i, key := range m.helpMenu.filteredData[i].hotkey {
			if i != 0 {
				hotkey += " | "
			}
			hotkey += key
		}

		cursor := "  "
		if m.helpMenu.cursor == i {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor + " ")
		}
		r.AddLines(cursor + common.ModalStyle.Render(fmt.Sprintf("%*s%s", renderHotkeyLength,
			common.HelpMenuHotkeyStyle.Render(hotkey+" "), common.ModalStyle.Render(description))))
	}
}

func (m *model) sortOptionsRender() string {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	sortOptionsContent := common.ModalTitleStyle.Render(" Sort Options") + "\n\n"
	for i, option := range panel.SortOptions.Data.Options {
		cursor := " "
		if i == panel.SortOptions.Cursor {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor)
		}
		sortOptionsContent += cursor + common.ModalStyle.Render(" "+option) + "\n"
	}
	bottomBorder := common.GenerateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(panel.SortOptions.Cursor+1),
		strconv.Itoa(len(panel.SortOptions.Data.Options))), panel.SortOptions.Width-common.BorderPadding)

	return common.SortOptionsModalBorderStyle(panel.SortOptions.Height, panel.SortOptions.Width,
		bottomBorder).Render(sortOptionsContent)
}

func (m *model) filePreviewPanelRender() string {
	return m.fileModel.filePreview.GetContent()
}
