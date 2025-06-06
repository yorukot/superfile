package internal

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"image"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid"
	"github.com/yorukot/superfile/src/internal/ui"
	"github.com/yorukot/superfile/src/internal/ui/rendering"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/yorukot/ansichroma"
	"github.com/yorukot/superfile/src/config/icon"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

func (m *model) sidebarRender() string {
	return m.sidebarModel.Render(m.mainPanelHeight, m.focusPanel == sidebarFocus, m.fileModel.filePanels[m.filePanelFocusIndex].location)
}

// This also modifies the m.fileModel.filePanels, which it should not
// what modifications we do on this model object are of no consequence.
// Since bubblea passed this 'model' by value in View() function.
func (m *model) filePanelRender() string {
	f := make([]string, len(m.fileModel.filePanels))
	for i, filePanel := range m.fileModel.filePanels {
		// check if cursor or render out of range
		// Todo - instead of this, have a filepanel.validateAndFix(), and log Error
		// This should not ever happen
		if filePanel.cursor > len(filePanel.element)-1 {
			filePanel.cursor = 0
			filePanel.render = 0
		}
		m.fileModel.filePanels[i] = filePanel

		// Todo : Move this to a utility function and clarify the calculation via comments
		// Maybe even write unit tests
		var filePanelWidth int
		if (m.fullWidth-common.Config.SidebarWidth-(4+(len(m.fileModel.filePanels)-1)*2))%len(m.fileModel.filePanels) != 0 && i == len(m.fileModel.filePanels)-1 {
			if m.fileModel.filePreview.open {
				filePanelWidth = m.fileModel.width
			} else {
				filePanelWidth = (m.fileModel.width + (m.fullWidth-common.Config.SidebarWidth-(4+(len(m.fileModel.filePanels)-1)*2))%len(m.fileModel.filePanels))
			}
		} else {
			filePanelWidth = m.fileModel.width
		}

		f[i] = filePanel.Render(m.mainPanelHeight, filePanelWidth, filePanel.focusType != noneFocus)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, f...)
}

func (panel *filePanel) Render(mainPanelHeight int, filePanelWidth int, focussed bool) string {
	r := ui.FilePanelRenderer(mainPanelHeight+2, filePanelWidth+2, focussed)
	// Todo : Unit test for all the functions
	panel.renderTopBar(r, filePanelWidth)
	panel.renderSearchBar(r)
	panel.renderFooter(r)
	panel.renderFileEntries(r, mainPanelHeight, filePanelWidth)

// Todo - Add AnsiTruncateLeft in ui/renderer package and remove truncation here
	panel.renderTopBar(r, filePanelWidth)
	panel.renderSearchBar(r)
	panel.renderFooter(r)
	panel.renderFileEntries(r, mainPanelHeight, filePanelWidth)

	return r.Render()
}

func (panel *filePanel) renderTopBar(r *rendering.Renderer, filePanelWidth int) {
	truncatedPath := common.TruncateTextBeginning(panel.location, filePanelWidth-4, "...")
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(truncatedPath))
}

func (panel *filePanel) renderSearchBar(r *rendering.Renderer) {
	r.AddSection()
	r.AddLines(" " + panel.searchBar.View())
}

func (panel *filePanel) renderFooter(r *rendering.Renderer) {
	sortLabel, sortIcon := panel.getSortInfo()
	modeLabel, modeIcon := panel.getPanelModeInfo()
	cursorStr := panel.getCursorPosition()

	if common.Config.Nerdfont {
		sortLabel = sortIcon + icon.Space + sortLabel
		modeLabel = modeIcon + icon.Space + modeLabel
	} else {
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

	end := panel.render + panelElementHeight(mainPanelHeight)
	if end > len(panel.element) {
		end = len(panel.element)
	}

	for i := panel.render; i < end; i++ {
		isCursor := i == panel.cursor && !panel.searchBar.Focused()
		isSelected := arrayContains(panel.selected, panel.element[i].location)

		if panel.renaming && i == panel.cursor {
			r.AddLines(panel.rename.View())
			continue
		}

		cursor := " "
		if isCursor {
			cursor = icon.Cursor
		}

		// Performance TODO: Remove or cache this if not needed at render time
		// Figure out why we are doing this. This will unnecessarily slow down
		// rendering. There should be a way to avoid this at render
		_, err := os.ReadDir(panel.element[i].location)
		dirExists := err == nil || panel.element[i].directory

		renderedName := common.PrettierName(
			panel.element[i].name,
			filePanelWidth-5,
			dirExists,
			isSelected,
			common.FilePanelBGColor,
		)

		r.AddLines(common.FilePanelCursorStyle.Render(cursor+" ") + renderedName)
	}
}

func (panel *filePanel) getSortInfo() (string, string) {
	opts := panel.sortOptions.data
	selected := opts.options[opts.selected]
	label := selected

	if selected == common.DateModifiedOption {
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

func (panel *filePanel) getCursorPosition() string {
	cursor := panel.cursor
	if len(panel.element) > 0 {
		cursor++ // Convert to 1-based
	}
	return fmt.Sprintf("%d/%d", cursor, len(panel.element))
}

func (m *model) processBarRender() string {
	if !m.processBarModel.isValid(m.footerHeight) {
		slog.Error("processBar in invalid state", "render", m.processBarModel.render,
			"cursor", m.processBarModel.cursor, "footerHeight", m.footerHeight)
	}

	r := ui.ProcessBarRenderer(m.footerHeight+2, utils.FooterWidth(m.fullWidth)+2, m.focusPanel == processBarFocus)

	// cursor's value itself cannot be used as its zero indexed
	cursorNumber := 0
	// Todo : Instead of directly accessing slice, there should be a method .IsEmpty() , or .CntProcess()
	if len(m.processBarModel.processList) != 0 {
		cursorNumber = m.processBarModel.cursor + 1
	}

	r.SetBorderInfoItems(fmt.Sprintf("%d/%d", cursorNumber, len(m.processBarModel.processList)))
	if len(m.processBarModel.processList) == 0 {
		r.AddLines("", " "+common.ProcessBarNoneText)
		return r.Render()
	}

	// save process in the array and sort the process by finished or not,
	// completion percetage, or finish time
	// Todo : This is very inefficient and can be improved.
	// The whole design needs to be changed so that we dont need to recreate the slice
	// and sort on each render. Idea : Maintain two slices - completed, ongoing
	// Processes should be added / removed to the slice on correct time, and we dont
	// need to redo slice formation and sorting on each render.
	var processes []process
	for _, p := range m.processBarModel.process {
		processes = append(processes, p)
	}
	// sort by the process
	sort.Slice(processes, func(i, j int) bool {
		doneI := (processes[i].state == successful)
		doneJ := (processes[j].state == successful)

		// sort by done or not
		if doneI != doneJ {
			return !doneI
		}

		// if both not done
		if !doneI {
			completionI := float64(processes[i].done) / float64(processes[i].total)
			completionJ := float64(processes[j].done) / float64(processes[j].total)
			return completionI < completionJ // Those who finish first will be ranked later.
		}

		// if both done sort by the doneTime
		return processes[j].doneTime.Before(processes[i].doneTime)
	})

	renderedHeight := 0

	for i := m.processBarModel.render; i < len(processes); i++ {
		// We allow rendering of a process if we have at least 2 lines left
		if m.footerHeight < renderedHeight+2 {
			break
		}
		renderedHeight += 3

		curProcess := processes[i]
		curProcess.progress.Width = utils.FooterWidth(m.fullWidth) - 3

		// Todo : get them via a separate function.
		var symbol string
		var cursor string
		if i == m.processBarModel.cursor {
			// Todo : Prerender it.
			cursor = common.FooterCursorStyle.Render("┃ ")
		} else {
			cursor = common.FooterCursorStyle.Render("  ")
		}
		// Todo : Prerender
		switch curProcess.state {
		case failure:
			symbol = common.ProcessErrorStyle.Render(icon.Warn)
		case successful:
			symbol = common.ProcessSuccessfulStyle.Render(icon.Done)
		case inOperation:
			symbol = common.ProcessInOperationStyle.Render(icon.InOperation)
		case cancel:
			symbol = common.ProcessCancelStyle.Render(icon.Error)
		}

		r.AddLines(cursor + common.FooterStyle.Render(common.TruncateText(curProcess.name, utils.FooterWidth(m.fullWidth)-7, "...")+" ") + symbol)
		r.AddLines(cursor+curProcess.progress.ViewAs(float64(curProcess.done)/float64(curProcess.total)), "")
	}

	return r.Render()
}

// This updates m.fileMetaData
func (m *model) metadataRender() string {
	m.ensureMetadataLoaded()

	sortedMeta := sortMetadata(m.fileMetaData.metaData)
	maxKeyLen := getMaxKeyLength(sortedMeta)
	sprintfLen, valLen := computeWidths(m.fullWidth, maxKeyLen)
	totalWidth := utils.FooterWidth(m.fullWidth)

	lines := formatMetadataLines(sortedMeta, m.fileMetaData.renderIndex, m.footerHeight, sprintfLen, totalWidth, valLen)

	r := ui.MetadataRenderer(m.footerHeight+2, utils.FooterWidth(m.fullWidth)+2, m.focusPanel == metadataFocus)
	if len(sortedMeta) > 0 {
		r.SetBorderInfoItems(fmt.Sprintf("%d/%d", m.fileMetaData.renderIndex+1, len(sortedMeta)))
	}
	for _, line := range lines {
		r.AddLines(line)
	}
	return r.Render()
}

func (m *model) ensureMetadataLoaded() {
	if len(m.fileMetaData.metaData) == 0 &&
		len(m.fileModel.filePanels[m.filePanelFocusIndex].element) > 0 &&
		!m.fileModel.renaming {
		loadingMessage := channelMessage{
			messageID:   shortuuid.New(),
			messageType: sendMetadata,
			metadata: [][2]string{
				{"", ""},
				{" " + icon.InOperation + "  Loading metadata...", ""},
			},
		}
		channel <- loadingMessage
		// Todo : This needs to be improved, we are updating m.fileMetaData is a separate goroutine
		// while also modifying it here in the function. It could cause issues.
		go func() {
			m.returnMetaData()
		}()
	}
}

func sortMetadata(meta [][2]string) [][2]string {
	priority := map[string]int{
		"Name":          0,
		"Size":          1,
		"Date Modified": 2,
		"Date Accessed": 3,
	}

	sort.SliceStable(meta, func(i, j int) bool {
		pi, iok := priority[meta[i][0]]
		pj, jok := priority[meta[j][0]]
		switch {
		case iok && jok:
			return pi < pj
		case iok:
			return true
		case jok:
			return false
		default:
			return meta[i][0] < meta[j][0]
		}
	})

	return meta
}

func getMaxKeyLength(meta [][2]string) int {
	maxLen := 0
	for _, pair := range meta {
		if len(pair[0]) > maxLen {
			maxLen = len(pair[0])
		}
	}
	return maxLen
}

func computeWidths(fullWidth, maxKeyLen int) (int, int) {
	totalWidth := utils.FooterWidth(fullWidth)
	valueLen := totalWidth - maxKeyLen - 2
	var sprintfLen int
	if valueLen < totalWidth/2 {
		valueLen = totalWidth/2 - 2
		sprintfLen = valueLen
	} else {
		sprintfLen = maxKeyLen + 1
	}
	return sprintfLen, valueLen
}

func formatMetadataLines(meta [][2]string, startIdx, height, sprintfLen, totalWidth, valueLen int) []string {
	lines := []string{}
	endIdx := min(startIdx+height, len(meta))
	for i := startIdx; i < endIdx; i++ {
		key := meta[i][0]
		value := common.TruncateMiddleText(meta[i][1], valueLen, "...")

		if totalWidth-sprintfLen-3 < totalWidth/2 {
			key = common.TruncateMiddleText(key, valueLen, "...")
		}
		line := fmt.Sprintf("%-*s %s", sprintfLen, key, value)
		lines = append(lines, line)
	}
	return lines
}

func (m *model) clipboardRender() string {
	// render
	var bottomWidth int
	if m.fullWidth%3 != 0 {
		bottomWidth = utils.FooterWidth(m.fullWidth + m.fullWidth%3 + 2)
	} else {
		bottomWidth = utils.FooterWidth(m.fullWidth)
	}
	r := ui.ClipboardRenderer(m.footerHeight+2, bottomWidth+2)
	if len(m.copyItems.items) == 0 {
		// Todo move this to a string
		r.AddLines("", " "+icon.Error+"  No content in clipboard")
	} else {
		for i := 0; i < len(m.copyItems.items) && i < m.footerHeight; i++ {
			if i == m.footerHeight-1 && i != len(m.copyItems.items)-1 {
				// Last Entry we can render, but there are more that one left
				r.AddLines(strconv.Itoa(len(m.copyItems.items)-i) + " item left....")
			} else {
				fileInfo, err := os.Stat(m.copyItems.items[i])
				if err != nil {
					slog.Error("Clipboard render function get item state ", "error", err)
				}
				if !os.IsNotExist(err) {
					// Todo : There is an inconsistency in parameter that is being passed, and its name in ClipboardPrettierName function
					r.AddLines(common.ClipboardPrettierName(m.copyItems.items[i], utils.FooterWidth(m.fullWidth)-3, fileInfo.IsDir(), false))
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
	return common.FullScreenStyle(m.fullHeight, m.fullWidth).Render(`Terminal size too small:` + "\n" +
		"Width = " + fullWidthString +
		heightString + fullHeightString + "\n\n" +
		"Needed for current config:" + "\n" +
		"Width = " + common.TerminalCorrectSize.Render(minimumWidthString) +
		heightString + common.TerminalCorrectSize.Render(minimumHeightString))
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
	return common.FullScreenStyle(m.fullHeight, m.fullWidth).Render(`You change your terminal size too small:` + "\n" +
		"Width = " + fullWidthString +
		heightString + fullHeightString + "\n\n" +
		"Needed for current config:" + "\n" +
		"Width = " + common.TerminalCorrectSize.Render(minimumWidthString) +
		heightString + common.TerminalCorrectSize.Render(minimumHeightString))
}

func (m *model) typineModalRender() string {
	previewPath := filepath.Join(m.typingModal.location, m.typingModal.textInput.Value())

	fileLocation := common.FilePanelTopDirectoryIconStyle.Render(" "+icon.Directory+icon.Space) +
		common.FilePanelTopPathStyle.Render(common.TruncateTextBeginning(previewPath, common.ModalWidth-4, "...")) + "\n"

	confirm := common.ModalConfirm.Render(" (" + common.Hotkeys.ConfirmTyping[0] + ") Create ")
	cancel := common.ModalCancel.Render(" (" + common.Hotkeys.CancelTyping[0] + ") Cancel ")

	tip := confirm +
		lipgloss.NewStyle().Background(common.ModalBGColor).Render("           ") +
		cancel

	return common.ModalBorderStyle(common.ModalHeight, common.ModalWidth).Render(fileLocation + "\n" + m.typingModal.textInput.View() + "\n\n" + tip)
}

func (m *model) introduceModalRender() string {
	title := common.SidebarTitleStyle.Render(" Thanks for using superfile!!") + common.ModalStyle.Render("\n You can read the following information before starting to use it!")
	vimUserWarn := common.ProcessErrorStyle.Render("  ** Very importantly ** If you are a Vim/Nvim user, go to:\n  https://superfile.netlify.app/configure/custom-hotkeys/ to change your hotkey settings!")
	subOne := common.SidebarTitleStyle.Render("  (1)") + common.ModalStyle.Render(" If this is your first time, make sure you read:\n      https://superfile.netlify.app/getting-started/tutorial/")
	subTwo := common.SidebarTitleStyle.Render("  (2)") + common.ModalStyle.Render(" If you forget the relevant keys during use,\n      you can press \"?\" (shift+/) at any time to query the keys!")
	subThree := common.SidebarTitleStyle.Render("  (3)") + common.ModalStyle.Render(" For more customization you can refer to:\n      https://superfile.netlify.app/")
	subFour := common.SidebarTitleStyle.Render("  (4)") + common.ModalStyle.Render(" Thank you again for using superfile.\n      If you have any questions, please feel free to ask at:\n      https://github.com/yorukot/superfile\n      Of course, you can always open a new issue to share your idea \n      or report a bug!")
	return common.FirstUseModal(m.helpMenu.height, m.helpMenu.width).Render(title + "\n\n" + vimUserWarn + "\n\n" + subOne + "\n\n" + subTwo + "\n\n" + subThree + "\n\n" + subFour + "\n\n")
}

func (m *model) warnModalRender() string {
	title := m.warnModal.title
	content := m.warnModal.content
	confirm := common.ModalConfirm.Render(" (" + common.Hotkeys.Confirm[0] + ") Confirm ")
	cancel := common.ModalCancel.Render(" (" + common.Hotkeys.Quit[0] + ") Cancel ")
	tip := confirm + lipgloss.NewStyle().Background(common.ModalBGColor).Render("           ") + cancel
	return common.ModalBorderStyle(common.ModalHeight, common.ModalWidth).Render(title + "\n\n" + content + "\n\n" + tip)
}

func (m *model) promptModalRender() string {
	return m.promptModal.Render()
}

func (m *model) helpMenuRender() string {
	helpMenuContent := ""
	maxKeyLength := 0

	for _, data := range m.helpMenu.data {
		totalKeyLen := 0
		for _, key := range data.hotkey {
			totalKeyLen += len(key)
		}
		saprateLen := len(data.hotkey) - 1*3
		if data.subTitle == "" && totalKeyLen+saprateLen > maxKeyLength {
			maxKeyLength = totalKeyLen + saprateLen
		}
	}

	valueLength := m.helpMenu.width - maxKeyLength - 2
	if valueLength < m.helpMenu.width/2 {
		valueLength = m.helpMenu.width/2 - 2
	}

	renderHotkeyLength := 0
	totalTitleCount := 0
	cursorBeenTitleCount := 0

	for i, data := range m.helpMenu.data {
		if data.subTitle != "" {
			if i < m.helpMenu.cursor {
				cursorBeenTitleCount++
			}
			totalTitleCount++
		}
	}

	for i := m.helpMenu.renderIndex; i < m.helpMenu.height+m.helpMenu.renderIndex && i < len(m.helpMenu.data); i++ {
		hotkey := ""

		if m.helpMenu.data[i].subTitle != "" {
			continue
		}

		for i, key := range m.helpMenu.data[i].hotkey {
			if i != 0 {
				hotkey += " | "
			}
			hotkey += key
		}

		if len(common.HelpMenuHotkeyStyle.Render(hotkey)) > renderHotkeyLength {
			renderHotkeyLength = len(common.HelpMenuHotkeyStyle.Render(hotkey))
		}
	}

	for i := m.helpMenu.renderIndex; i < m.helpMenu.height+m.helpMenu.renderIndex && i < len(m.helpMenu.data); i++ {
		if i != m.helpMenu.renderIndex {
			helpMenuContent += "\n"
		}

		if m.helpMenu.data[i].subTitle != "" {
			helpMenuContent += common.HelpMenuTitleStyle.Render(" " + m.helpMenu.data[i].subTitle)
			continue
		}

		hotkey := ""
		description := common.TruncateText(m.helpMenu.data[i].description, valueLength, "...")

		for i, key := range m.helpMenu.data[i].hotkey {
			if i != 0 {
				hotkey += " | "
			}
			hotkey += key
		}

		cursor := "  "
		if m.helpMenu.cursor == i {
			cursor = common.FilePanelCursorStyle.Render(icon.Cursor + " ")
		}
		helpMenuContent += cursor + common.ModalStyle.Render(fmt.Sprintf("%*s%s", renderHotkeyLength, common.HelpMenuHotkeyStyle.Render(hotkey+" "), common.ModalStyle.Render(description)))
	}

	bottomBorder := common.GenerateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(m.helpMenu.cursor+1-cursorBeenTitleCount), strconv.Itoa(len(m.helpMenu.data)-totalTitleCount)), m.helpMenu.width-2)

	return common.HelpMenuModalBorderStyle(m.helpMenu.height, m.helpMenu.width, bottomBorder).Render(helpMenuContent)
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
	bottomBorder := common.GenerateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(panel.sortOptions.cursor+1), strconv.Itoa(len(panel.sortOptions.data.options))), panel.sortOptions.width-2)

	return common.SortOptionsModalBorderStyle(panel.sortOptions.height, panel.sortOptions.width, bottomBorder).Render(sortOptionsContent)
}

func readFileContent(filepath string, maxLineLength int, previewLine int) (string, error) {
	var resultBuilder strings.Builder
	file, err := os.Open(filepath)
	if err != nil {
		return resultBuilder.String(), err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = ansi.Truncate(line, maxLineLength, "")
		resultBuilder.WriteString(line)
		resultBuilder.WriteRune('\n')
		lineCount++
		if previewLine > 0 && lineCount >= previewLine {
			break
		}
	}
	// returns the first non-EOF error that was encountered by the [Scanner]
	return resultBuilder.String(), scanner.Err()
}

func (m *model) filePreviewPanelRender() string {
	// Todo : This width adjustment must not be done inside render function. It should
	// only be triggered via Update()
	m.fileModel.filePreview.width += m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.width - ((m.fileModel.width + 2) * len(m.fileModel.filePanels)) - 2

	return m.filePreviewPanelRenderWithDimensions(m.mainPanelHeight+2, m.fileModel.filePreview.width)
}

// Helper function to handle empty panel case
func (m *model) renderEmptyFilePreview(r *rendering.Renderer) string {
	clearCmd := filepreview.ClearKittyImages()
	if clearCmd != "" {
		r.AddLines(clearCmd + common.FilePreviewNoContentText)
	} else {
		r.AddLines(common.FilePreviewNoContentText)
	}
	return r.Render()
}

// Helper function to handle file info errors
func (m *model) renderFileInfoError(r *rendering.Renderer, _ lipgloss.Style, err error) string {
	slog.Error("Error get file info", "error", err)
	clearCmd := filepreview.ClearKittyImages()
	if clearCmd != "" {
		r.AddLines(clearCmd + common.FilePreviewNoFileInfoText)
	} else {
		r.AddLines(common.FilePreviewNoFileInfoText)
	}
	return r.Render()
}

// Helper function to handle unsupported formats
func (m *model) renderUnsupportedFormat(r *rendering.Renderer, _ lipgloss.Style) string {
	clearCmd := filepreview.ClearKittyImages()
	if clearCmd != "" {
		r.AddLines(clearCmd + common.FilePreviewUnsupportedFormatText)
	} else {
		r.AddLines(common.FilePreviewUnsupportedFormatText)
	}
	return r.Render()
}

// Helper function to handle directory preview
func (m *model) renderDirectoryPreview(r *rendering.Renderer, itemPath string, previewHeight int) string {
	clearCmd := filepreview.ClearKittyImages()
	files, err := os.ReadDir(itemPath)
	if err != nil {
		slog.Error("Error render directory preview", "error", err)
		if clearCmd != "" {
			r.AddLines(clearCmd + common.FilePreviewDirectoryUnreadableText)
		} else {
			r.AddLines(common.FilePreviewDirectoryUnreadableText)
		}
		return r.Render()
	}

	if len(files) == 0 {
		if clearCmd != "" {
			r.AddLines(clearCmd + common.FilePreviewEmptyText)
		} else {
			r.AddLines(common.FilePreviewEmptyText)
		}
		return r.Render()
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir() && !files[j].IsDir() {
			return true
		}
		if !files[i].IsDir() && files[j].IsDir() {
			return false
		}
		return files[i].Name() < files[j].Name()
	})

	// Add clear command before directory listing
	if clearCmd != "" {
		r.AddLines(clearCmd)
	}

	for i := 0; i < previewHeight && i < len(files); i++ {
		file := files[i]
		style := common.GetElementIcon(file.Name(), file.IsDir(), common.Config.Nerdfont)
		res := lipgloss.NewStyle().Foreground(lipgloss.Color(style.Color)).Background(common.FilePanelBGColor).
			Render(style.Icon+" ") + common.FilePanelStyle.Render(file.Name())
		r.AddLines(res)
	}
	return r.Render()
}

// Helper function to handle image preview
func (m *model) renderImagePreview(box lipgloss.Style, itemPath string, previewWidth, previewHeight int) string {
	if !m.fileModel.filePreview.open {
		clearCmd := filepreview.ClearKittyImages()
		if clearCmd != "" {
			return box.Render(clearCmd + "\n --- Preview panel is closed ---")
		}
		return box.Render("\n --- Preview panel is closed ---")
	}

	if !common.Config.ShowImagePreview {
		clearCmd := filepreview.ClearKittyImages()
		if clearCmd != "" {
			return box.Render(clearCmd + "\n --- Image preview is disabled ---")
		}
		return box.Render("\n --- Image preview is disabled ---")
	}

	imageRender, err := filepreview.ImagePreview(itemPath, previewWidth, previewHeight, common.Theme.FilePanelBG)
	if errors.Is(err, image.ErrFormat) {
		clearCmd := filepreview.ClearKittyImages()
		if clearCmd != "" {
			return box.Render(clearCmd + "\n --- " + icon.Error + " Unsupported image formats ---")
		}
		return box.Render("\n --- " + icon.Error + " Unsupported image formats ---")
	}

	if err != nil {
		slog.Error("Error convert image to ansi", "error", err)
		clearCmd := filepreview.ClearKittyImages()
		if clearCmd != "" {
			return box.Render(clearCmd + "\n --- " + icon.Error + " Error convert image to ansi ---")
		}
		return box.Render("\n --- " + icon.Error + " Error convert image to ansi ---")
	}

	// Check if this looks like Kitty protocol output (starts with escape sequences)
	// For Kitty protocol, avoid using lipgloss alignment to prevent layout drift
	if strings.HasPrefix(imageRender, "\x1b_G") {
		// This is Kitty protocol output - render directly in a simple box
		// without vertical alignment to avoid layout issues
		return common.FilePreviewBox(previewHeight, previewWidth).Render(imageRender)
	}

	// For ANSI output, we can safely use vertical alignment
	return box.AlignVertical(lipgloss.Center).Render(imageRender)
}

// Helper function to handle text file preview
func (m *model) renderTextPreview(r *rendering.Renderer, box lipgloss.Style, itemPath string, previewWidth, previewHeight int) string {
	clearCmd := filepreview.ClearKittyImages()
	format := lexers.Match(filepath.Base(itemPath))

	if format == nil {
		isText, err := common.IsTextFile(itemPath)
		if err != nil {
			slog.Error("Error while checking text file", "error", err)
			if clearCmd != "" {
				return box.Render(clearCmd + "\n --- " + icon.Error + " Error get file info ---")
			}
			return box.Render("\n --- " + icon.Error + " Error get file info ---")
		} else if !isText {
			if clearCmd != "" {
				return box.Render(clearCmd + "\n --- " + icon.Error + " Unsupported formats ---")
			}
			return box.Render("\n --- " + icon.Error + " Unsupported formats ---")
		}
	}

	fileContent, err := readFileContent(itemPath, previewWidth, previewHeight)
	if err != nil {
		slog.Error("Error open file", "error", err)
		if clearCmd != "" {
			return box.Render(clearCmd + "\n --- " + icon.Error + " Error open file ---")
		}
		return box.Render("\n --- " + icon.Error + " Error open file ---")
	}

	if fileContent == "" {
		if clearCmd != "" {
			return box.Render(clearCmd + "\n --- empty ---")
		}
		return box.Render("\n --- empty ---")
	}

	if format != nil {
		background := ""
		if !common.Config.TransparentBackground {
			background = common.Theme.FilePanelBG
		}
		if common.Config.CodePreviewer == "bat" {
			if batCmd == "" {
				if clearCmd != "" {
					return box.Render(clearCmd + "\n --- " + icon.Error + " 'bat' is not installed or not found. ---\n --- Cannot render file preview. ---")
				}
				return box.Render("\n --- " + icon.Error + " 'bat' is not installed or not found. ---\n --- Cannot render file preview. ---")
			}
			fileContent, err = getBatSyntaxHighlightedContent(itemPath, previewHeight, background)
		} else {
			fileContent, err = ansichroma.HightlightString(fileContent, format.Config().Name, common.Theme.CodeSyntaxHighlightTheme, background)
		}
		if err != nil {
			slog.Error("Error render code highlight", "error", err)
			if clearCmd != "" {
				return box.Render(clearCmd + "\n --- " + icon.Error + " Error render code highlight ---")
			}
			return box.Render("\n --- " + icon.Error + " Error render code highlight ---")
		}
	}

	// Add clear command before text content
	if clearCmd != "" {
		r.AddLines(clearCmd)
	}
	r.AddLines(fileContent)
	return r.Render()
}

func (m *model) filePreviewPanelRenderWithDimensions(previewHeight int, previewWidth int) string {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	box := common.FilePreviewBox(previewHeight, previewWidth)
	r := ui.FilePreviewPanelRenderer(previewHeight, previewWidth)

	if len(panel.element) == 0 {
		return m.renderEmptyFilePreview(r)
	}

	itemPath := panel.element[panel.cursor].location
	fileInfo, infoErr := os.Stat(itemPath)

	if infoErr != nil {
		return m.renderFileInfoError(r, box, infoErr)
	}

	ext := filepath.Ext(itemPath)
	if slices.Contains(common.UnsupportedPreviewFormats, ext) {
		return m.renderUnsupportedFormat(r, box)
	}

	if fileInfo.IsDir() {
		return m.renderDirectoryPreview(r, itemPath, previewHeight)
	}

	if isImageFile(itemPath) {
		return m.renderImagePreview(box, itemPath, previewWidth, previewHeight)
	}

	return m.renderTextPreview(r, box, itemPath, previewWidth, previewHeight)
}

func getBatSyntaxHighlightedContent(itemPath string, previewLine int, background string) (string, error) {
	// --plain: use the plain style without line numbers and decorations
	// --force-colorization: force colorization for non-interactive shell output
	// --line-range <:m>: only read from line 1 to line "m"
	batArgs := []string{itemPath, "--plain", "--force-colorization", "--line-range", fmt.Sprintf(":%d", previewLine-1)}

	// set timeout for the external command execution to 500ms max
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	cmd := exec.CommandContext(ctx, batCmd, batArgs...)

	fileContentBytes, err := cmd.Output()
	if err != nil {
		slog.Error("Error render code highlight", "error", err)
		return "", err
	}

	fileContent := string(fileContentBytes)
	if !common.Config.TransparentBackground {
		fileContent = setBatBackground(fileContent, background)
	}
	return fileContent, nil
}

func setBatBackground(input string, background string) string {
	tokens := strings.Split(input, "\x1b[0m")
	backgroundStyle := lipgloss.NewStyle().Background(lipgloss.Color(background))
	for idx, token := range tokens {
		tokens[idx] = backgroundStyle.Render(token)
	}
	return strings.Join(tokens, "\x1b[0m")
}
