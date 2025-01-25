package internal

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/lexers"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/term/ansi"
	"github.com/yorukot/ansichroma"
	"github.com/yorukot/superfile/src/config/icon"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"
)

func (m *model) sidebarRender() string {
	if Config.SidebarWidth == 0 {
		return ""
	}
	superfileTitle := sidebarTitleStyle.Render("    " + icon.SuperfileIcon + " superfile")
	superfileTitle = ansi.Truncate(superfileTitle, Config.SidebarWidth, "")
	s := superfileTitle
	s += "\n"

	pinnedDivider := "\n" + sidebarTitleStyle.Render("󰐃 Pinned") + sidebarDividerStyle.Render(" ───────────") + "\n"
	disksDivider := "\n" + sidebarTitleStyle.Render("󱇰 Disks") + sidebarDividerStyle.Render(" ────────────") + "\n"
	disksDivider = ansi.Truncate(disksDivider, Config.SidebarWidth, "")
	pinnedDivider = ansi.Truncate(pinnedDivider, Config.SidebarWidth, "")

	totalHeight := 2
	for i := m.sidebarModel.renderIndex; i < len(m.sidebarModel.directories); i++ {
		if totalHeight >= m.mainPanelHeight {
			break
		} else {
			s += "\n"
		}

		directory := m.sidebarModel.directories[i]

		if directory.location == "Pinned+-*/=?" {
			s += pinnedDivider
			totalHeight += 3
			continue
		}

		if directory.location == "Disks+-*/=?" {
			if m.mainPanelHeight-totalHeight <= 2 {
				break
			}
			s += disksDivider
			totalHeight += 3
			continue
		}

		totalHeight++
		cursor := " "
		if m.sidebarModel.cursor == i && m.focusPanel == sidebarFocus {
			cursor = icon.Cursor
		}

		if directory.location == m.fileModel.filePanels[m.filePanelFocusIndex].location {
			s += filePanelCursorStyle.Render(cursor+" ") + sidebarSelectedStyle.Render(truncateText(directory.name, Config.SidebarWidth-2, "..."))
		} else {
			s += filePanelCursorStyle.Render(cursor+" ") + sidebarStyle.Render(truncateText(directory.name, Config.SidebarWidth-2, "..."))
		}
	}

	return sideBarBorderStyle(m.mainPanelHeight, m.focusPanel).Render(s)
}

// This also modifies the m.fileModel.filePanels
func (m *model) filePanelRender() string {
	// file panel
	f := make([]string, 10)
	for i, filePanel := range m.fileModel.filePanels {

		// check if cursor or render out of range
		if filePanel.cursor > len(filePanel.element)-1 {
			filePanel.cursor = 0
			filePanel.render = 0
		}
		m.fileModel.filePanels[i] = filePanel

		f[i] += filePanelTopDirectoryIconStyle.Render(" "+icon.Directory+icon.Space) + filePanelTopPathStyle.Render(truncateTextBeginning(filePanel.location, m.fileModel.width-4, "...")) + "\n"
		filePanelWidth := 0
		footerBorderWidth := 0

		if (m.fullWidth-Config.SidebarWidth-(4+(len(m.fileModel.filePanels)-1)*2))%len(m.fileModel.filePanels) != 0 && i == len(m.fileModel.filePanels)-1 {
			if m.fileModel.filePreview.open {
				filePanelWidth = m.fileModel.width
			} else {
				filePanelWidth = (m.fileModel.width + (m.fullWidth-Config.SidebarWidth-(4+(len(m.fileModel.filePanels)-1)*2))%len(m.fileModel.filePanels))
			}
			footerBorderWidth = m.fileModel.width + 15
		} else {
			filePanelWidth = m.fileModel.width
			footerBorderWidth = m.fileModel.width + 15
		}

		sortDirectionString := ""
		if filePanel.sortOptions.data.reversed {
			if Config.Nerdfont {
				sortDirectionString = icon.SortDesc
			} else {
				sortDirectionString = "D"
			}
		} else {
			if Config.Nerdfont {
				sortDirectionString = icon.SortAsc
			} else {
				sortDirectionString = "A"
			}
		}
		sortTypeString := ""
		if filePanelWidth < 23 {
			sortTypeString = sortDirectionString
		} else {
			if filePanel.sortOptions.data.options[filePanel.sortOptions.data.selected] == "Date Modified" {
				sortTypeString = sortDirectionString + icon.Space + "Date"
			} else {
				sortTypeString = sortDirectionString + icon.Space + filePanel.sortOptions.data.options[filePanel.sortOptions.data.selected]
			}
		}

		panelModeString := ""
		if filePanelWidth < 23 {
			if filePanel.panelMode == browserMode {
				if Config.Nerdfont {
					panelModeString = icon.Browser
				} else {
					panelModeString = "B"
				}
			} else if filePanel.panelMode == selectMode {
				if Config.Nerdfont {
					panelModeString = icon.Select
				} else {
					panelModeString = "S"
				}
			}
		} else {
			if filePanel.panelMode == browserMode {
				panelModeString = icon.Browser + icon.Space + "Browser"
			} else if filePanel.panelMode == selectMode {
				panelModeString = icon.Select + icon.Space + "Select"
			}
		}

		f[i] += filePanelDividerStyle(filePanel.focusType).Render(strings.Repeat(Config.BorderTop, filePanelWidth)) + "\n"
		f[i] += " " + filePanel.searchBar.View() + "\n"
		if len(filePanel.element) == 0 {
			f[i] += filePanelStyle.Render(" " + icon.Error + "  No such file or directory")
			bottomBorder := generateFooterBorder(fmt.Sprintf("%s%s%s%s%s", sortTypeString, bottomMiddleBorderSplit, panelModeString, bottomMiddleBorderSplit, "0/0"), footerBorderWidth)
			f[i] = filePanelBorderStyle(m.mainPanelHeight, filePanelWidth, filePanel.focusType, bottomBorder).Render(f[i])
		} else {
			for h := filePanel.render; h < filePanel.render+panelElementHeight(m.mainPanelHeight) && h < len(filePanel.element); h++ {
				endl := "\n"
				if h == filePanel.render+panelElementHeight(m.mainPanelHeight)-1 || h == len(filePanel.element)-1 {
					endl = ""
				}
				cursor := " "
				// Check if the cursor needs to be displayed, if the user is using the search bar, the cursor is not displayed
				if h == filePanel.cursor && !filePanel.searchBar.Focused() {
					cursor = icon.Cursor
				}
				isItemSelected := arrayContains(filePanel.selected, filePanel.element[h].location)
				if filePanel.renaming && h == filePanel.cursor {
					f[i] += filePanel.rename.View() + endl
				} else {
					_, err := os.ReadDir(filePanel.element[h].location);
					f[i] += filePanelCursorStyle.Render(cursor+" ") + prettierName(filePanel.element[h].name, m.fileModel.width-5, filePanel.element[h].directory||(err == nil), isItemSelected, filePanelBGColor) + endl
				}
			}
			cursorPosition := strconv.Itoa(filePanel.cursor + 1)
			totalElement := strconv.Itoa(len(filePanel.element))

			bottomBorder := generateFooterBorder(fmt.Sprintf("%s%s%s%s%s/%s", sortTypeString, bottomMiddleBorderSplit, panelModeString, bottomMiddleBorderSplit, cursorPosition, totalElement), footerBorderWidth)
			f[i] = filePanelBorderStyle(m.mainPanelHeight, filePanelWidth, filePanel.focusType, bottomBorder).Render(f[i])
		}
	}

	// file panel render together
	filePanelRender := ""
	for _, f := range f {
		filePanelRender = lipgloss.JoinHorizontal(lipgloss.Top, filePanelRender, f)
	}
	return filePanelRender
}
func (m *model) processBarRender() string {
	// save process in the array
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

	// render
	processRender := ""
	renderTimes := 0

	for i := m.processBarModel.render; i < len(processes); i++ {
		if footerHeight < 14 && renderTimes == 2 {
			break
		}
		if renderTimes == 3 {
			break
		}
		process := processes[i]
		process.progress.Width = footerWidth(m.fullWidth) - 3
		symbol := ""
		cursor := ""
		if i == m.processBarModel.cursor {
			cursor = footerCursorStyle.Render("┃ ")
		} else {
			cursor = footerCursorStyle.Render("  ")
		}
		switch process.state {
		case failure:
			symbol = processErrorStyle.Render(icon.Warn)
		case successful:
			symbol = processSuccessfulStyle.Render(icon.Done)
		case inOperation:
			symbol = processInOperationStyle.Render(icon.InOperation)
		case cancel:
			symbol = processCancelStyle.Render(icon.Error)
		}

		processRender += cursor + footerStyle.Render(truncateText(process.name, footerWidth(m.fullWidth)-7, "...")+" ") + symbol + "\n"
		if renderTimes == 2 {
			processRender += cursor + process.progress.ViewAs(float64(process.done)/float64(process.total)) + ""
		} else if footerHeight < 14 && renderTimes == 1 {
			processRender += cursor + process.progress.ViewAs(float64(process.done)/float64(process.total))
		} else {
			processRender += cursor + process.progress.ViewAs(float64(process.done)/float64(process.total)) + "\n\n"
		}
		renderTimes++
	}

	if len(processes) == 0 {
		processRender += "\n " + icon.Error + "  No processes running"
	}
	courseNumber := 0
	if len(m.processBarModel.processList) == 0 {
		courseNumber = 0
	} else {
		courseNumber = m.processBarModel.cursor + 1
	}
	bottomBorder := generateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(courseNumber), strconv.Itoa(len(m.processBarModel.processList))), footerWidth(m.fullWidth)-3)
	processRender = procsssBarBoarder(bottomElementHeight(footerHeight), footerWidth(m.fullWidth), bottomBorder, m.focusPanel).Render(processRender)

	return processRender
}

// This updates m.fileMetaData 
func (m *model) metadataRender() string {
	// process bar
	metaDataBar := ""
	if len(m.fileMetaData.metaData) == 0 && len(m.fileModel.filePanels[m.filePanelFocusIndex].element) > 0 && !m.fileModel.renaming {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"", ""})
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{" " + icon.InOperation + "  Loading metadata...", ""})
		go func() {
			m.returnMetaData()
		}()
	}
	maxKeyLength := 0
	sort.Slice(m.fileMetaData.metaData, func(i, j int) bool {
		comparisonFields := []string{"FileName", "FileSize", "FolderName", "FolderSize", "FileModifyDate", "FileAccessDate"}

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
	for _, data := range m.fileMetaData.metaData {
		if len(data[0]) > maxKeyLength {
			maxKeyLength = len(data[0])
		}
	}

	sprintfLength := maxKeyLength + 1
	valueLength := footerWidth(m.fullWidth) - maxKeyLength - 2
	if valueLength < footerWidth(m.fullWidth)/2 {
		valueLength = footerWidth(m.fullWidth)/2 - 2
		sprintfLength = valueLength
	}

	for i := m.fileMetaData.renderIndex; i < bottomElementHeight(footerHeight)+m.fileMetaData.renderIndex && i < len(m.fileMetaData.metaData); i++ {
		if i != m.fileMetaData.renderIndex {
			metaDataBar += "\n"
		}
		data := truncateMiddleText(m.fileMetaData.metaData[i][1], valueLength, "...")
		metadataName := m.fileMetaData.metaData[i][0]
		if footerWidth(m.fullWidth)-maxKeyLength-3 < footerWidth(m.fullWidth)/2 {
			metadataName = truncateMiddleText(m.fileMetaData.metaData[i][0], valueLength, "...")
		}
		metaDataBar += fmt.Sprintf("%-*s %s", sprintfLength, metadataName, data)

	}
	bottomBorder := generateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(m.fileMetaData.renderIndex+1), strconv.Itoa(len(m.fileMetaData.metaData))), footerWidth(m.fullWidth)-3)
	metaDataBar = metadataBoarder(bottomElementHeight(footerHeight), footerWidth(m.fullWidth), bottomBorder, m.focusPanel).Render(metaDataBar)

	return metaDataBar
}

func (m *model) clipboardRender() string {

	// render
	clipboardRender := ""
	if len(m.copyItems.items) == 0 {
		clipboardRender += "\n " + icon.Error + "  No content in clipboard"
	} else {
		for i := 0; i < len(m.copyItems.items) && i < bottomElementHeight(footerHeight); i++ {
			if i == bottomElementHeight(footerHeight)-1 {
				clipboardRender += strconv.Itoa(len(m.copyItems.items)-i+1) + " item left...."
			} else {
				fileInfo, err := os.Stat(m.copyItems.items[i])
				if err != nil {
					outPutLog("Clipboard render function get item state error", err)
				}
				if !os.IsNotExist(err) {
					clipboardRender += clipboardPrettierName(m.copyItems.items[i], footerWidth(m.fullWidth)-3, fileInfo.IsDir(), false) + "\n"
				}
			}
		}
	}
	for i := 0; i < len(m.copyItems.items); i++ {

	}
	bottomWidth := 0

	if m.fullWidth%3 != 0 {
		bottomWidth = footerWidth(m.fullWidth + m.fullWidth%3 + 2)
	} else {
		bottomWidth = footerWidth(m.fullWidth)
	}
	clipboardRender = clipboardBoarder(bottomElementHeight(footerHeight), bottomWidth, Config.BorderBottom).Render(clipboardRender)

	return clipboardRender
}

func (m *model) terminalSizeWarnRender() string {
	fullWidthString := strconv.Itoa(m.fullWidth)
	fullHeightString := strconv.Itoa(m.fullHeight)
	minimumWidthString := strconv.Itoa(minimumWidth)
	minimumHeightString := strconv.Itoa(minimumHeight)
	if m.fullHeight < minimumHeight {
		fullHeightString = terminalTooSmall.Render(fullHeightString)
	}
	if m.fullWidth < minimumWidth {
		fullWidthString = terminalTooSmall.Render(fullWidthString)
	}
	fullHeightString = terminalCorrectSize.Render(fullHeightString)
	fullWidthString = terminalCorrectSize.Render(fullWidthString)

	heightString := mainStyle.Render(" Height = ")
	return fullScreenStyle(m.fullHeight, m.fullWidth).Render(`Terminal size too small:` + "\n" +
		"Width = " + fullWidthString +
		heightString + fullHeightString + "\n\n" +

		"Needed for current config:" + "\n" +
		"Width = " + terminalCorrectSize.Render(minimumWidthString) +
		heightString + terminalCorrectSize.Render(minimumHeightString))
}

func (m *model) terminalSizeWarnAfterFirstRender() string {
	minimumWidthInt := Config.SidebarWidth + 20*len(m.fileModel.filePanels) + 20 - 1
	minimumWidthString := strconv.Itoa(minimumWidthInt)
	fullWidthString := strconv.Itoa(m.fullWidth)
	fullHeightString := strconv.Itoa(m.fullHeight)
	minimumHeightString := strconv.Itoa(minimumHeight)

	if m.fullHeight < minimumHeight {
		fullHeightString = terminalTooSmall.Render(fullHeightString)
	}
	if m.fullWidth < minimumWidthInt {
		fullWidthString = terminalTooSmall.Render(fullWidthString)
	}
	fullHeightString = terminalCorrectSize.Render(fullHeightString)
	fullWidthString = terminalCorrectSize.Render(fullWidthString)

	heightString := mainStyle.Render(" Height = ")
	return fullScreenStyle(m.fullHeight, m.fullWidth).Render(`You change your terminal size too small:` + "\n" +
		"Width = " + fullWidthString +
		heightString + fullHeightString + "\n\n" +

		"Needed for current config:" + "\n" +
		"Width = " + terminalCorrectSize.Render(minimumWidthString) +
		heightString + terminalCorrectSize.Render(minimumHeightString))
}

func (m *model) typineModalRender() string {
	previewPath := m.typingModal.location + "/" + m.typingModal.textInput.Value()

	fileLocation := filePanelTopDirectoryIconStyle.Render(" "+icon.Directory+icon.Space) +
		filePanelTopPathStyle.Render(truncateTextBeginning(previewPath, modalWidth-4, "...")) + "\n"

	confirm := modalConfirm.Render(" (" + hotkeys.ConfirmTyping[0] + ") Create ")
	cancel := modalCancel.Render(" (" + hotkeys.CancelTyping[0] + ") Cancel ")

	tip := confirm +
		lipgloss.NewStyle().Background(modalBGColor).Render("           ") +
		cancel

	return modalBorderStyle(modalHeight, modalWidth).Render(fileLocation + "\n" + m.typingModal.textInput.View() + "\n\n" + tip)
}

func (m *model) introduceModalRender() string {
	title := sidebarTitleStyle.Render(" Thanks for using superfile!!") + modalStyle.Render("\n You can read the following information before starting to use it!")
	vimUserWarn := processErrorStyle.Render("  ** Very importantly ** If you are a Vim/Nvim user, go to:\n  https://superfile.netlify.app/configure/custom-hotkeys/ to change your hotkey settings!")
	subOne := sidebarTitleStyle.Render("  (1)") + modalStyle.Render(" If this is your first time, make sure you read:\n      https://superfile.netlify.app/getting-started/tutorial/")
	subTwo := sidebarTitleStyle.Render("  (2)") + modalStyle.Render(" If you forget the relevant keys during use,\n      you can press \"?\" (shift+/) at any time to query the keys!")
	subThree := sidebarTitleStyle.Render("  (3)") + modalStyle.Render(" For more customization you can refer to:\n      https://superfile.netlify.app/")
	subFour := sidebarTitleStyle.Render("  (4)") + modalStyle.Render(" Thank you again for using superfile.\n      If you have any questions, please feel free to ask at:\n      https://github.com/yorukot/superfile\n      Of course, you can always open a new issue to share your idea \n      or report a bug!")
	return firstUseModal(m.helpMenu.height, m.helpMenu.width).Render(title + "\n\n" + vimUserWarn + "\n\n" + subOne + "\n\n" + subTwo + "\n\n" + subThree + "\n\n" + subFour + "\n\n")
}

func (m *model) warnModalRender() string {
	title := m.warnModal.title
	content := m.warnModal.content
	confirm := modalConfirm.Render(" (" + hotkeys.Confirm[0] + ") Confirm ")
	cancel := modalCancel.Render(" (" + hotkeys.Quit[0] + ") Cancel ")
	tip := confirm + lipgloss.NewStyle().Background(modalBGColor).Render("           ") + cancel
	return modalBorderStyle(modalHeight, modalWidth).Render(title + "\n\n" + content + "\n\n" + tip)
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

		if len(helpMenuHotkeyStyle.Render(hotkey)) > renderHotkeyLength {
			renderHotkeyLength = len(helpMenuHotkeyStyle.Render(hotkey))
		}
	}

	for i := m.helpMenu.renderIndex; i < m.helpMenu.height+m.helpMenu.renderIndex && i < len(m.helpMenu.data); i++ {

		if i != m.helpMenu.renderIndex {
			helpMenuContent += "\n"
		}

		if m.helpMenu.data[i].subTitle != "" {
			helpMenuContent += helpMenuTitleStyle.Render(" " + m.helpMenu.data[i].subTitle)
			continue
		}

		hotkey := ""
		description := truncateText(m.helpMenu.data[i].description, valueLength, "...")

		for i, key := range m.helpMenu.data[i].hotkey {
			if i != 0 {
				hotkey += " | "
			}
			hotkey += key
		}

		cursor := "  "
		if m.helpMenu.cursor == i {
			cursor = filePanelCursorStyle.Render(icon.Cursor + " ")
		}
		helpMenuContent += cursor + modalStyle.Render(fmt.Sprintf("%*s%s", renderHotkeyLength, helpMenuHotkeyStyle.Render(hotkey+" "), modalStyle.Render(description)))
	}

	bottomBorder := generateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(m.helpMenu.cursor+1-cursorBeenTitleCount), strconv.Itoa(len(m.helpMenu.data)-totalTitleCount)), m.helpMenu.width-2)

	return helpMenuModalBorderStyle(m.helpMenu.height, m.helpMenu.width, bottomBorder).Render(helpMenuContent)
}

func (m *model) sortOptionsRender() string {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	sortOptionsContent := modalTitleStyle.Render(" Sort Options") + "\n\n"
	for i, option := range panel.sortOptions.data.options {
		cursor := " "
		if i == panel.sortOptions.cursor {
			cursor = filePanelCursorStyle.Render(icon.Cursor)
		}
		sortOptionsContent += cursor + modalStyle.Render(" "+option) + "\n"
	}
	bottomBorder := generateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(panel.sortOptions.cursor+1), strconv.Itoa(len(panel.sortOptions.data.options))), panel.sortOptions.width-2)

	return sortOptionsModalBorderStyle(panel.sortOptions.height, panel.sortOptions.width, bottomBorder).Render(sortOptionsContent)
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
		line = makePrintable(line)
		resultBuilder.WriteString(line+"\n")
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
	m.fileModel.filePreview.width += m.fullWidth - Config.SidebarWidth - m.fileModel.filePreview.width - ((m.fileModel.width + 2) * len(m.fileModel.filePanels)) - 2

	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	box := filePreviewBox(previewLine, m.fileModel.filePreview.width)

	if len(panel.element) == 0 {
		return box.Render("\n --- " + icon.Error + " No content to preview ---")
	}

	itemPath := panel.element[panel.cursor].location

	fileInfo, err := os.Stat(itemPath)

	if err != nil {
		outPutLog("error get file info", err)
		return box.Render("\n --- " + icon.Error + " Error get file info ---")
	}

	ext := filepath.Ext(itemPath)
	// check if the file is unsipported file, cuz pdf will cause error
	if ext == ".pdf" || ext == ".torrent" {
		return box.Render("\n --- " + icon.Error + " Unsupported formats ---")
	}

	if fileInfo.IsDir() {
		directoryContent := ""
		dirPath := itemPath

		files, err := os.ReadDir(dirPath)
		if err != nil {
			outPutLog("Error render directory preview", err)
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
			directoryContent += prettierDirectoryPreviewName(file.Name(), file.IsDir(), filePanelBGColor)
			if i != previewLine-1 && i != len(files)-1 {
				directoryContent += "\n"
			}
		}
		directoryContent = checkAndTruncateLineLengths(directoryContent, m.fileModel.filePreview.width)
		return box.Render(directoryContent)
	}

	if isImageFile(itemPath) {
		if !m.fileModel.filePreview.open {
            return box.Render("\n --- Preview panel is closed ---")
        }
		ansiRender, err := filepreview.ImagePreview(itemPath, m.fileModel.filePreview.width, previewLine, theme.FilePanelBG)
		if err == image.ErrFormat {
			return box.Render("\n --- " + icon.Error + " Unsupported image formats ---")
		}

		if err != nil {
			outPutLog("Error covernt image to ansi", err)
			return box.Render("\n --- " + icon.Error + " Error covernt image to ansi ---")
		}

		return box.AlignVertical(lipgloss.Center).AlignHorizontal(lipgloss.Center).Render(ansiRender)
	}

	format := lexers.Match(filepath.Base(itemPath))

	if format == nil {
		isText, err := isTextFile(itemPath)
		if err != nil {
			outPutLog("Error while checking text file", err)
			return box.Render("\n --- " + icon.Error + " Error get file info ---")
		} else if !isText {
			return box.Render("\n --- " + icon.Error + " Unsupported formats ---")
		}
	}

	// At this point either format is not nil, or we can read the file
	fileContent , err := readFileContent(itemPath, m.fileModel.width+20, previewLine)
	if err != nil {
		outPutLog(err)
		return box.Render("\n --- " + icon.Error + " Error open file ---")
	}

	// We know the format of file, and we can apply syntax highlighting
	if format != nil {
		background := ""
		if ! Config.TransparentBackground {
			background = theme.FilePanelBG
		}
		fileContent, err = ansichroma.HightlightString(fileContent, format.Config().Name, theme.CodeSyntaxHighlightTheme, background)
		if err != nil {
			outPutLog("Error render code highlight", err)
			return box.Render("\n --- " + icon.Error + " Error render code highlight ---")
		}
	}

	if fileContent == "" {
		return box.Render("\n --- empty ---")
	}
	fileContent = checkAndTruncateLineLengths(fileContent, m.fileModel.filePreview.width)
	return box.Render(fileContent)
}

func (m *model) commandLineInputBoxRender() string {
	return m.commandLine.input.View()
}
