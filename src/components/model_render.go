package components

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func sidebarRender(m model) string {
	s := sidebarTitleStyle.Render("     Super File")
	s += "\n\n"

	// Ugly shit workaround from hell code, made by @lescx
	amountWellKnownDirectories := len(getWellKnownDirectories())
	amountPinnedDirectories := len(getPinnedDirectories())
	pinnedRendered := false
	externalRendered := false

	pinnedDivider := "\n" + sidebarTitleStyle.Render("󰐃 Pinned") + sidebarDividerStyle.Render(" ───────────") + "\n\n"
	disksDivider := "\n" + sidebarTitleStyle.Render("󱇰 Disks") + sidebarDividerStyle.Render(" ────────────") + "\n\n"

	for i, directory := range m.sidebarModel.directories {
		if i == amountWellKnownDirectories && !pinnedRendered {
			s += pinnedDivider
			pinnedRendered = true
		}

		if i == amountPinnedDirectories+amountWellKnownDirectories && !externalRendered {
			s += disksDivider
			externalRendered = true
		}
		cursor := " "
		if m.sidebarModel.cursor == i && m.focusPanel == sidebarFocus {
			cursor = ""
		}
		if directory.location == m.fileModel.filePanels[m.filePanelFocusIndex].location {
			s += filePanelCursorStyle.Render(cursor) + sidebarSelectedStyle.Render(" "+truncateText(directory.name, sidebarWidth-2)) + "\n"
		} else {
			s += filePanelCursorStyle.Render(cursor) + sidebarStyle.Render(" "+truncateText(directory.name, sidebarWidth-2)) + "\n"
		}
	}

	// In case no pinned directories or external drives are pinned,
	// list menu item at the bottom
	if m.fullHeight > 30 {
		if !pinnedRendered {
			s += pinnedDivider
		}
		if !externalRendered {
			s += disksDivider
		}
	}

	return sideBarBorderStyle(m.mainPanelHeight, m.focusPanel).Render(s)
}

func filePanelRender(m model) string {
	// file panel
	f := make([]string, 10)
	for i, filePanel := range m.fileModel.filePanels {


		// check if cursor or render out of range
		if filePanel.cursor > len(filePanel.element)-1 {
			filePanel.cursor = 0
			filePanel.render = 0
		}
		m.fileModel.filePanels[i] = filePanel

		f[i] += filePanelTopDirectoryIconStyle.Render("   ") + filePanelTopPathStyle.Render(truncateTextBeginning(filePanel.location, m.fileModel.width-4)) + "\n"
		filePanelWidth := 0
		bottomBorderWidth := 0

		if (m.fullWidth-sidebarWidth-(4+(len(m.fileModel.filePanels)-1)*2))%len(m.fileModel.filePanels) != 0 && i == len(m.fileModel.filePanels)-1 {
			filePanelWidth = (m.fileModel.width + (m.fullWidth-sidebarWidth-(4+(len(m.fileModel.filePanels)-1)*2))%len(m.fileModel.filePanels))
			bottomBorderWidth = m.fileModel.width + 7
		} else {
			filePanelWidth = m.fileModel.width
			bottomBorderWidth = m.fileModel.width + 6
		}

		f[i] += filePanelDividerStyle(filePanel.focusType).Render(strings.Repeat("━", filePanelWidth)) + "\n"
		f[i] += " " + filePanel.searchBar.View() + "\n"
		if len(filePanel.element) == 0 {
			f[i] += filePanelStyle.Render("   No such file or directory")
			bottomBorder := generateFooterBorder("0/0", m.fileModel.width+5)
			f[i] = filePanelBorderStyle(m.mainPanelHeight, m.fileModel.width, filePanel.focusType, bottomBorder).Render(f[i])
		} else {
			for h := filePanel.render; h < filePanel.render+panelElementHeight(m.mainPanelHeight) && h < len(filePanel.element); h++ {
				endl := "\n"
				if h == filePanel.render+panelElementHeight(m.mainPanelHeight)-1 || h == len(filePanel.element)-1 {
					endl = ""
				}
				cursor := " "
				// Check if the cursor needs to be displayed, if the user is using the search bar, the cursor is not displayed
				if h == filePanel.cursor && !filePanel.searchBar.Focused() {
					cursor = ""
				}
				isItemSelected := arrayContains(filePanel.selected, filePanel.element[h].location)
				if filePanel.renaming && h == filePanel.cursor {
					f[i] += filePanel.rename.View() + endl
				} else {
					f[i] += filePanelCursorStyle.Render(cursor+" ") + prettierName(filePanel.element[h].name, m.fileModel.width-5, filePanel.element[h].directory, isItemSelected, filePanelBGColor) + endl
				}
			}
			cursorPosition := strconv.Itoa(filePanel.cursor + 1)
			totalElement := strconv.Itoa(len(filePanel.element))
			panelModeString := ""
			if filePanel.panelMode == browserMode {
				panelModeString = "󰈈 Browser"
			} else if filePanel.panelMode == selectMode {
				panelModeString = "󰆽 Select"
			}
			bottomBorder := generateFooterBorder(fmt.Sprintf("%s┣━┫%s/%s", panelModeString, cursorPosition, totalElement), bottomBorderWidth)
			f[i] = filePanelBorderStyle(m.mainPanelHeight, filePanelWidth, filePanel.focusType, bottomBorder).Render(f[i])
		}
	}

	// file panel render togther
	filePanelRender := ""
	for _, f := range f {
		filePanelRender = lipgloss.JoinHorizontal(lipgloss.Top, filePanelRender, f)
	}
	return filePanelRender
}

func processBarRender(m model) string {
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
			symbol = processErrorStyle.Render("")
		case successful:
			symbol = processSuccessfulStyle.Render("")
		case inOperation:
			symbol = processInOperationStyle.Render("󰥔")
		case cancel:
			symbol = processCancelStyle.Render("")
		}

		processRender += cursor + footerStyle.Render(truncateText(process.name, footerWidth(m.fullWidth)-7)+" ") + symbol + "\n"
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
		processRender += "\n   No processes running"
	}
	courseNumber := 0
	if len(m.processBarModel.processList) == 0 {
		courseNumber = 0
	} else {
		courseNumber = m.processBarModel.cursor + 1
	}
	bottomBorder := generateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(courseNumber), strconv.Itoa(len(m.processBarModel.processList))), footerWidth(m.fullWidth)-3)
	processRender = procsssBarBoarder(bottomElementHight(footerHeight), footerWidth(m.fullWidth), bottomBorder, m.focusPanel).Render(processRender)

	return processRender
}

func metadataRender(m model) string {
	// process bar
	metaDataBar := ""
	if len(m.fileMetaData.metaData) == 0 && len(m.fileModel.filePanels[m.filePanelFocusIndex].element) > 0 && !m.fileModel.renaming {
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{"", ""})
		m.fileMetaData.metaData = append(m.fileMetaData.metaData, [2]string{" 󰥔  Loading metadata...", ""})
		go func() {
			m = returnMetaData(m)
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
	vauleLength := footerWidth(m.fullWidth) - maxKeyLength - 2
	if vauleLength < footerWidth(m.fullWidth)/2 {
		vauleLength = footerWidth(m.fullWidth)/2 - 2
		sprintfLength = vauleLength
	}

	for i := m.fileMetaData.renderIndex; i < bottomElementHight(footerHeight)+m.fileMetaData.renderIndex && i < len(m.fileMetaData.metaData); i++ {
		if i != m.fileMetaData.renderIndex {
			metaDataBar += "\n"
		}
		data := truncateMiddleText(m.fileMetaData.metaData[i][1], vauleLength)
		metadataName := m.fileMetaData.metaData[i][0]
		if footerWidth(m.fullWidth)-maxKeyLength-3 < footerWidth(m.fullWidth)/2 {
			metadataName = truncateMiddleText(m.fileMetaData.metaData[i][0], vauleLength)
		}
		metaDataBar += fmt.Sprintf("%-*s %s", sprintfLength, metadataName, data)

	}
	bottomBorder := generateFooterBorder(fmt.Sprintf("%s/%s", strconv.Itoa(m.fileMetaData.renderIndex+1), strconv.Itoa(len(m.fileMetaData.metaData))), footerWidth(m.fullWidth)-3)
	metaDataBar = metadataBoarder(bottomElementHight(footerHeight), footerWidth(m.fullWidth), bottomBorder, m.focusPanel).Render(metaDataBar)

	return metaDataBar
}

func clipboardRender(m model) string {

	// render
	clipboardRender := ""
	if len(m.copyItems.items) == 0 {
		clipboardRender += "\n   No content in clipboard"
	} else {
		for i := 0; i < len(m.copyItems.items) && i < bottomElementHight(footerHeight); i++ {
			if i == bottomElementHight(footerHeight)-1 {
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
	clipboardRender = clipboardBoarder(bottomElementHight(footerHeight), bottomWidth, "━").Render(clipboardRender)

	return clipboardRender
}

func terminalSizeWarnRender(m model) string {
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

func typineModalRender(m model) string {
	if m.typingModal.itemType == newFile {
		fileLocation := filePanelTopDirectoryIconStyle.Render("   ") + filePanelTopPathStyle.Render(truncateTextBeginning(m.typingModal.location+"/"+m.typingModal.textInput.Value(), modalWidth-4)) + "\n"
		confirm := modalConfirm.Render(" (" + hotkeys.Confirm[0] + ") New File ")
		cancel := modalCancel.Render(" (" + hotkeys.Cancel[0] + ") Cancel ")
		tip := confirm + lipgloss.NewStyle().Background(modalBGColor).Render("           ") + cancel
		return fullScreenStyle(m.fullHeight, m.fullWidth).Render(modalBorderStyle(modalHeight, modalWidth).Render(fileLocation + "\n" + m.typingModal.textInput.View() + "\n\n" + tip))
	} else {
		fileLocation := filePanelTopDirectoryIconStyle.Render("   ") + filePanelTopPathStyle.Render(truncateTextBeginning(m.typingModal.location+"/"+m.typingModal.textInput.Value(), modalWidth-4)) + "\n"
		confirm := modalConfirm.Render(" (" + hotkeys.Confirm[0] + ") New Folder ")
		cancel := modalCancel.Render(" (" + hotkeys.Cancel[0] + ") Cancel ")
		tip := confirm + lipgloss.NewStyle().Background(modalBGColor).Render("           ") + cancel
		return fullScreenStyle(m.fullHeight, m.fullWidth).Render(modalBorderStyle(modalHeight, modalWidth).Render(fileLocation + "\n" + m.typingModal.textInput.View() + "\n\n" + tip))
	}
}

func warnModalRender(m model) string {
	title := m.warnModal.title
	content := m.warnModal.content
	confirm := modalCancel.Render(" (" + hotkeys.Confirm[0] + ") Confirm ")
	cancel := modalCancel.Render(" (" + hotkeys.Cancel[0] + ") Cancel ")
	tip := confirm + lipgloss.NewStyle().Background(modalBGColor).Render("           ") + cancel
	return fullScreenStyle(m.fullHeight, m.fullWidth).Render(modalBorderStyle(modalHeight, modalWidth).Render(title + "\n\n" + content + "\n\n" + tip))
}
