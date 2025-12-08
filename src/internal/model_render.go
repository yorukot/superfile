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
		m.fileModel.filePanels[m.filePanelFocusIndex].location)
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
		if filePanel.cursor > len(filePanel.element)-1 {
			filePanel.cursor = 0
			filePanel.render = 0
		}
		m.fileModel.filePanels[i] = filePanel

		// TODO : Move this to a utility function and clarify the calculation via comments
		// Maybe even write unit tests
		var filePanelWidth int
		if (m.fullWidth-common.Config.SidebarWidth-(4+(len(m.fileModel.filePanels)-1)*2))%len(
			m.fileModel.filePanels,
		) != 0 &&
			i == len(m.fileModel.filePanels)-1 {
			if m.fileModel.filePreview.IsOpen() {
				filePanelWidth = m.fileModel.width
			} else {
				filePanelWidth = (m.fileModel.width + (m.fullWidth-common.Config.SidebarWidth-
					(4+(len(m.fileModel.filePanels)-1)*2))%len(m.fileModel.filePanels))
			}
		} else {
			filePanelWidth = m.fileModel.width
		}

		f[i] = filePanel.Render(m.mainPanelHeight, filePanelWidth, filePanel.isFocused)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, f...)
}

func (panel *filePanel) Render(mainPanelHeight int, filePanelWidth int, focussed bool) string {
	r := ui.FilePanelRenderer(mainPanelHeight+common.BorderPadding, filePanelWidth+common.BorderPadding, focussed)

	panel.renderTopBar(r, filePanelWidth)
	panel.renderSearchBar(r)
	panel.renderFooter(r)
	panel.renderFileEntries(r, mainPanelHeight, filePanelWidth)

	return r.Render()
}

func (panel *filePanel) renderTopBar(r *rendering.Renderer, filePanelWidth int) {
	// TODO - Add ansitruncate left in renderer and remove truncation here
	truncatedPath := common.TruncateTextBeginning(panel.location, filePanelWidth-common.InnerPadding, "...")
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(truncatedPath))
	r.AddSection()
}

func (panel *filePanel) renderSearchBar(r *rendering.Renderer) {
	r.AddLines(" " + panel.searchBar.View())
}

// TODO : Unit test this
func (panel *filePanel) renderFooter(r *rendering.Renderer) {
	sortLabel, sortIcon := panel.getSortInfo()
	modeLabel, modeIcon := panel.getPanelModeInfo()
	cursorStr := panel.getCursorString()

	if common.Config.Nerdfont {
		sortLabel = sortIcon + icon.Space + sortLabel
		modeLabel = modeIcon + icon.Space + modeLabel
	} else {
		// TODO : Figure out if we can set icon.Space to " " if nerdfont is false
		// That would simplify code
		sortLabel = sortIcon + " " + sortLabel
	}

	if common.Config.ShowPanelFooterInfo {
		r.SetBorderInfoItems(sortLabel, modeLabel, cursorStr)
		if r.AreInfoItemsTruncated() {
			r.SetBorderInfoItems(sortIcon, modeIcon, cursorStr)
		}
	} else {
		r.SetBorderInfoItems(cursorStr)
	}
}

func (panel *filePanel) renderFileEntries(r *rendering.Renderer, mainPanelHeight, filePanelWidth int) {
	if len(panel.element) == 0 {
		r.AddLines(common.FilePanelNoneText)
		return
	}

	end := min(panel.render+panelElementHeight(mainPanelHeight), len(panel.element))

	for i := panel.render; i < end; i++ {
		// TODO : Fix this, this is O(n^2) complexity. Considered a file panel with 200 files, and 100 selected
		// We will be doing a search in 100 item slice for all 200 files.
		isSelected := arrayContains(panel.selected, panel.element[i].location)

		if panel.renaming && i == panel.cursor {
			r.AddLines(panel.rename.View())
			continue
		}

		cursor := " "
		if i == panel.cursor && !panel.searchBar.Focused() {
			cursor = icon.Cursor
		}

		selectBox := panel.renderSelectBox(isSelected)

		// Calculate the actual prefix width for proper alignment
		prefixWidth := lipgloss.Width(cursor+" ") + lipgloss.Width(selectBox)

		isLink := panel.element[i].info.Mode()&os.ModeSymlink != 0
		renderedName := common.PrettierName(
			panel.element[i].name,
			filePanelWidth-prefixWidth,
			panel.element[i].directory,
			isLink,
			isSelected,
			common.FilePanelBGColor,
		)

		r.AddLines(common.FilePanelCursorStyle.Render(cursor+" ") + selectBox + renderedName)
	}
}

func (panel *filePanel) getSortInfo() (string, string) {
	opts := panel.sortOptions.data
	selected := opts.options[opts.selected]
	label := selected
	if selected == string(sortingDateModified) {
		label = "Date"
	}

	iconStr := icon.SortAsc

	if opts.reversed {
		iconStr = icon.SortDesc
	}
	return label, iconStr
}

func (panel *filePanel) getPanelModeInfo() (string, string) {
	switch panel.panelMode {
	case browserMode:
		return "Browser", icon.Browser
	case selectMode:
		return "Select", icon.Select
	default:
		return "", ""
	}
}

func (panel *filePanel) getCursorString() string {
	cursor := panel.cursor
	if len(panel.element) > 0 {
		cursor++ // Convert to 1-based
	}
	return fmt.Sprintf("%d/%d", cursor, len(panel.element))
}

func (panel *filePanel) renderSelectBox(isSelected bool) string {
	if !common.Config.ShowSelectIcons || !common.Config.Nerdfont || panel.panelMode != selectMode {
		return ""
	}

	if panel.isFocused {
		if isSelected {
			return common.CheckboxCheckedFocused
		}
		return common.CheckboxEmptyFocused
	}
	if isSelected {
		return common.CheckboxChecked
	}
	return common.CheckboxEmpty
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
	minimumWidthInt := common.Config.SidebarWidth + 20*len(m.fileModel.filePanels) + 20 - 1
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

		separatorLen := max(0, (len(data.hotkey)-1)) * 3
		if data.subTitle == "" && totalKeyLen+separatorLen > maxKeyLength {
			maxKeyLength = totalKeyLen + separatorLen
		}
	}

	valueLength := m.helpMenu.width - maxKeyLength - 2
	if valueLength < m.helpMenu.width/2 {
		valueLength = m.helpMenu.width/2 - 2
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
	for i := m.helpMenu.renderIndex; i < m.helpMenu.renderIndex+(m.helpMenu.height-4) && i < len(m.helpMenu.filteredData); i++ {
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
	for i := m.helpMenu.renderIndex; i < m.helpMenu.renderIndex+(m.helpMenu.height-4) && i < len(m.helpMenu.filteredData); i++ {
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
	for i, option := range panel.sortOptions.data.options {
		cursor := " "
		if i == panel.sortOptions.cursor {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor)
		}
		sortOptionsContent += cursor + common.ModalStyle.Render(" "+option) + "\n"
	}
	bottomBorder := common.GenerateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(panel.sortOptions.cursor+1),
		strconv.Itoa(len(panel.sortOptions.data.options))), panel.sortOptions.width-common.BorderPadding)

	return common.SortOptionsModalBorderStyle(panel.sortOptions.height, panel.sortOptions.width,
		bottomBorder).Render(sortOptionsContent)
}

func (m *model) filePreviewPanelRender() string {
	return m.fileModel.filePreview.GetContent()
}
