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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/internal/ui"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/charmbracelet/lipgloss"
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

	// Todo - Add ansitruncate left in renderer and remove truncation here
	r.AddLines(common.FilePanelTopDirectoryIcon + common.FilePanelTopPathStyle.Render(
		common.TruncateTextBeginning(panel.location, filePanelWidth-4, "...")))

	// Todo : Unit test all these if else chains (?)
	// Todo : Better move it out to a separate function
	var sortTypeString string
	var sortTypeStringSmall string
	if panel.sortOptions.data.reversed {
		sortTypeStringSmall = icon.SortDesc
	} else {
		sortTypeStringSmall = icon.SortAsc
	}

	// Todo : Make "Date Modified" a constant, and move this to a utility function
	if panel.sortOptions.data.options[panel.sortOptions.data.selected] == "Date Modified" {
		sortTypeString = "Date"
	} else {
		sortTypeString = panel.sortOptions.data.options[panel.sortOptions.data.selected]
	}

	if common.Config.Nerdfont {
		sortTypeString = sortTypeStringSmall + icon.Space + sortTypeString
	} else {
		// Todo : Figure out if we can set icon.Space to " " if nerdfont is false
		sortTypeString = sortTypeStringSmall + " " + sortTypeString
	}

	var panelModeString string
	var panelModeStringSmall string

	if panel.panelMode == browserMode {
		panelModeStringSmall = icon.Browser
		panelModeString = "Browser"
	} else if panel.panelMode == selectMode {
		panelModeStringSmall = icon.Select
		panelModeString = "Select"
	}

	// Only append Icon in case of nerdfont being true
	if common.Config.Nerdfont {
		panelModeString = panelModeStringSmall + icon.Space + panelModeString
	}

	r.AddSection()
	r.AddLines(" " + panel.searchBar.View())

	cursorNumber := panel.cursor

	// Make 1-indexed only for non zero filePanel len
	if len(panel.element) > 0 {
		cursorNumber++
	}
	cursorNumberString := fmt.Sprintf("%d/%d", cursorNumber, len(panel.element))

	r.SetBorderInfoItems(sortTypeString, panelModeString, cursorNumberString)
	if r.AreInfoItemsTruncated() {
		// Use smaller values
		r.SetBorderInfoItems(sortTypeStringSmall, panelModeStringSmall, cursorNumberString)
	}

	if len(panel.element) == 0 {
		r.AddLines(common.FilePanelNoneText)
	} else {
		for h := panel.render; h < panel.render+panelElementHeight(mainPanelHeight) && h < len(panel.element); h++ {
			cursor := " "
			// Check if the cursor needs to be displayed, if the user is using the search bar, the cursor is not displayed
			if h == panel.cursor && !panel.searchBar.Focused() {
				cursor = icon.Cursor
			}
			isItemSelected := arrayContains(panel.selected, panel.element[h].location)
			if panel.renaming && h == panel.cursor {
				r.AddLines(panel.rename.View())
			} else {
				// Todo (Performance) : Figure out why we are doing this. This will unnecessarily slow down
				// rendering. There should be a way to avoid this at render
				_, err := os.ReadDir(panel.element[h].location)
				r.AddLines(common.FilePanelCursorStyle.Render(cursor+" ") + common.PrettierName(panel.element[h].name, filePanelWidth-5,
					panel.element[h].directory || (err == nil), isItemSelected, common.FilePanelBGColor))
			}
		}
	}
	return r.Render()
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
	// process bar
	if len(m.fileMetaData.metaData) == 0 && len(m.fileModel.filePanels[m.filePanelFocusIndex].element) > 0 && !m.fileModel.renaming {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"", ""})
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{" " + icon.InOperation + "  Loading metadata...", ""})
		// Todo : This needs to be improved, we are updating m.fileMetaData is a separate goroutine
		// while also modifying it here in the function. It could cause issues.
		go func() {
			m.returnMetaData()
		}()
	}

	// Todo : The whole intention of this is to get the comparisonFields come before
	// other fields. Sorting like this is a bad way of achieving that. This can be improved
	sort.Slice(m.fileMetaData.metaData, func(i, j int) bool {
		// Initialising a new slice in each check by sort functions is too ineffinceint.
		// Todo : Fix it
		comparisonFields := []string{"Name", "Size", "Date Modified", "Date Accessed"}

		for _, field := range comparisonFields {
			if m.fileMetaData.metaData[i][0] == field {
				return true
			} else if m.fileMetaData.metaData[j][0] == field {
				return false
			}
		}

		// Default comparison
		return m.fileMetaData.metaData[i][0] < m.fileMetaData.metaData[j][0]
	})

	// Part where actual rendering happens.
	maxKeyLength := 0
	for _, data := range m.fileMetaData.metaData {
		if len(data[0]) > maxKeyLength {
			maxKeyLength = len(data[0])
		}
	}

	// Todo : Too much calculations that are not in a fuctions, are not
	// unit tested, and have no proper explanation. This makes it
	// very hard to maintain and add any changes
	sprintfLength := maxKeyLength + 1
	valueLength := utils.FooterWidth(m.fullWidth) - maxKeyLength - 2
	if valueLength < utils.FooterWidth(m.fullWidth)/2 {
		valueLength = utils.FooterWidth(m.fullWidth)/2 - 2
		sprintfLength = valueLength
	}
	r := ui.MetadataRenderer(m.footerHeight+2, utils.FooterWidth(m.fullWidth)+2, m.focusPanel == metadataFocus)
	// Todo : We can take this info as input in metadata renderer constructor
	renderIndex := m.fileMetaData.renderIndex
	if len(m.fileMetaData.metaData) > 0 {
		renderIndex++
	}
	r.SetBorderInfoItems(fmt.Sprintf("%d/%d", renderIndex, len(m.fileMetaData.metaData)))

	imax := min(m.footerHeight+m.fileMetaData.renderIndex, len(m.fileMetaData.metaData))
	for i := m.fileMetaData.renderIndex; i < imax; i++ {
		data := common.TruncateMiddleText(m.fileMetaData.metaData[i][1], valueLength, "...")
		metadataName := m.fileMetaData.metaData[i][0]
		if utils.FooterWidth(m.fullWidth)-maxKeyLength-3 < utils.FooterWidth(m.fullWidth)/2 {
			metadataName = common.TruncateMiddleText(m.fileMetaData.metaData[i][0], valueLength, "...")
		}
		r.AddLines(fmt.Sprintf("%-*s %s", sprintfLength, metadataName, data))
	}
	return r.Render()
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
	return m.promptModal.Render(m.helpMenu.height, m.helpMenu.width)
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
	// String builder is much better for efficiency
	// See - https://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go/47798475#47798475
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
		if len(line) > maxLineLength {
			line = line[:maxLineLength]
		}
		// This is critical to avoid layout break, removes non Printable ASCII control characters.
		line = common.MakePrintable(line)
		resultBuilder.WriteString(line + "\n")
		lineCount++
		if previewLine > 0 && lineCount >= previewLine {
			break
		}
	}
	// returns the first non-EOF error that was encountered by the [Scanner]
	return resultBuilder.String(), scanner.Err()
}

func (m *model) filePreviewPanelRender() string {
	previewLine := m.mainPanelHeight + 2
	m.fileModel.filePreview.width += m.fullWidth - common.Config.SidebarWidth - m.fileModel.filePreview.width - ((m.fileModel.width + 2) * len(m.fileModel.filePanels)) - 2

	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	box := common.FilePreviewBox(previewLine, m.fileModel.filePreview.width)

	if len(panel.element) == 0 {
		return box.Render("\n --- " + icon.Error + " No content to preview ---")
	}
	// This could create errors if panel.cursor ever becomes negative, or goes out of bounds
	// We should have a panel validation function in our View() function
	// Panel is a full fledged object with own state, its accessed and modified so many times.
	// Ideally we dont should never access data from it via directly accessing its variables
	// Todo : Instead we should have helper functions for panel object and access data that way
	// like panel.GetCurrentSelectedElem() . This abstration of implemetation of panel is needed.
	// Now this lack of abstraction has caused issues ( See PR#730 ) . And now
	// someone needs to scan through the entire codebase to figure out which access of panel
	// data is causing crash.
	itemPath := panel.element[panel.cursor].location

	// Renamed it to info_err to prevent shadowing with err below
	fileInfo, infoErr := os.Stat(itemPath)

	if infoErr != nil {
		slog.Error("Error get file info", "error", infoErr)
		return box.Render("\n --- " + icon.Error + " Error get file info ---")
	}

	ext := filepath.Ext(itemPath)
	// check if the file is unsupported file, cuz pdf will cause error
	if ext == ".pdf" || ext == ".torrent" {
		return box.Render("\n --- " + icon.Error + " Unsupported formats ---")
	}

	if fileInfo.IsDir() {
		directoryContent := ""
		dirPath := itemPath

		files, err := os.ReadDir(dirPath)
		if err != nil {
			slog.Error("Error render directory preview", "error", err)
			return box.Render("\n --- " + icon.Error + " Error render directory preview ---")
		}

		if len(files) == 0 {
			return box.Render("\n --- empty ---")
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

		for i := 0; i < previewLine && i < len(files); i++ {
			file := files[i]
			directoryContent += common.PrettierDirectoryPreviewName(file.Name(), file.IsDir(), common.FilePanelBGColor)
			if i != previewLine-1 && i != len(files)-1 {
				directoryContent += "\n"
			}
		}
		directoryContent = common.CheckAndTruncateLineLengths(directoryContent, m.fileModel.filePreview.width)
		return box.Render(directoryContent)
	}

	if isImageFile(itemPath) {
		if !m.fileModel.filePreview.open {
			// Todo : These variables can be pre rendered for efficiency and less duplicacy
			return box.Render("\n --- Preview panel is closed ---")
		}

		if !common.Config.ShowImagePreview {
			return box.Render("\n --- Image preview is disabled ---")
		}

		ansiRender, err := filepreview.ImagePreview(itemPath, m.fileModel.filePreview.width, previewLine, common.Theme.FilePanelBG)
		if errors.Is(err, image.ErrFormat) {
			return box.Render("\n --- " + icon.Error + " Unsupported image formats ---")
		}

		if err != nil {
			slog.Error("Error covernt image to ansi", "error", err)
			return box.Render("\n --- " + icon.Error + " Error covernt image to ansi ---")
		}

		return box.AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center).Render(ansiRender)
	}

	format := lexers.Match(filepath.Base(itemPath))

	if format == nil {
		isText, err := common.IsTextFile(itemPath)
		if err != nil {
			slog.Error("Error while checking text file", "error", err)
			return box.Render("\n --- " + icon.Error + " Error get file info ---")
		} else if !isText {
			return box.Render("\n --- " + icon.Error + " Unsupported formats ---")
		}
	}

	// At this point either format is not nil, or we can read the file
	fileContent, err := readFileContent(itemPath, m.fileModel.width+20, previewLine)
	if err != nil {
		slog.Error("Error open file", "error", err)
		return box.Render("\n --- " + icon.Error + " Error open file ---")
	}

	if fileContent == "" {
		return box.Render("\n --- empty ---")
	}

	// We know the format of file, and we can apply syntax highlighting
	if format != nil {
		background := ""
		if !common.Config.TransparentBackground {
			background = common.Theme.FilePanelBG
		}
		if common.Config.CodePreviewer == "bat" {
			if batCmd == "" {
				return box.Render("\n --- " + icon.Error + " 'bat' is not installed or not found. ---\n --- Cannot render file preview. ---")
			}
			fileContent, err = getBatSyntaxHighlightedContent(itemPath, previewLine, background)
		} else {
			fileContent, err = ansichroma.HightlightString(fileContent, format.Config().Name, common.Theme.CodeSyntaxHighlightTheme, background)
		}
		if err != nil {
			slog.Error("Error render code highlight", "error", err)
			return box.Render("\n --- " + icon.Error + " Error render code highlight ---")
		}
	}

	fileContent = common.CheckAndTruncateLineLengths(fileContent, m.fileModel.filePreview.width)
	return box.Render(fileContent)
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
