package components

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func SideBarRender(m model) string {
	s := sideBarTitle.Render("    Super Files     ")
	s += "\n"
	for i, folder := range m.sideBarModel.pinnedModel.folder {
		cursor := " "
		if m.sideBarModel.cursor == i && m.focusPanel == sideBarFocus {
			cursor = ""
		}
		if folder.location == m.fileModel.filePanels[m.filePanelFocusIndex].location {
			s += cursorStyle.Render(cursor) + " " + sideBarSelected.Render(TruncateText(folder.name, sideBarWidth-2)) + "" + "\n"
		} else {
			s += cursorStyle.Render(cursor) + " " + sideBarItem.Render(TruncateText(folder.name, sideBarWidth-2)) + "" + "\n"
		}
		if i == 4 {
			s += "\n" + sideBarTitle.Render("󰐃 Pinned") + borderStyle.Render(" ───────────") + "\n\n"
		}
		if folder.endPinned {
			s += "\n" + sideBarTitle.Render("󱇰 Disk") + borderStyle.Render(" ─────────────") + "\n\n"
		}
	}

	s = SideBarBoardStyle(m.mainPanelHeight, m.focusPanel).Render(s)

	return s
}

func FilePanelRender(m model, sideBar string) string {
	// file panel
	f := make([]string, 4)
	for i, filePanel := range m.fileModel.filePanels {
		fileElenent := returnFolderElement(filePanel.location)
		filePanel.element = fileElenent
		m.fileModel.filePanels[i].element = fileElenent
		f[i] += filePanelTopFolderIcon.Render("   ") + filePanelTopPath.Render(TruncateTextBeginning(filePanel.location, m.fileModel.width-4)) + "\n"
		f[i] += FilePanelDividerStyle(filePanel.focusType).Render(repeatString("━", m.fileModel.width)) + "\n"
		if len(filePanel.element) == 0 {
			f[i] += "   No any file or folder"
			bottomBorder := GenerateBottomBorder("0/0", m.fileModel.width+5)
			f[i] = FilePanelBoardStyle(m.mainPanelHeight, m.fileModel.width, filePanel.focusType, bottomBorder).Render(f[i])
		} else {
			for h := filePanel.render; h < filePanel.render+PanelElementHeight(m.mainPanelHeight) && h < len(filePanel.element); h++ {
				cursor := " "
				if h == filePanel.cursor {
					cursor = ""
				}
				isItemSelected := ArrayContains(filePanel.selected, filePanel.element[h].location)
				if filePanel.renaming && h == filePanel.cursor {
					f[i] += filePanel.rename.View() + "\n"
				} else {
					f[i] += cursorStyle.Render(cursor) + " " + PrettierName(filePanel.element[h].name, m.fileModel.width-5, filePanel.element[h].folder, isItemSelected) + "\n"
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
			bottomBorder := GenerateBottomBorder(fmt.Sprintf("%s┣━┫%s/%s", panelModeString, cursorPosition, totalElement), m.fileModel.width+6)
			f[i] = FilePanelBoardStyle(m.mainPanelHeight, m.fileModel.width, filePanel.focusType, bottomBorder).Render(f[i])
		}
	}

	// file panel render togther
	filePanelRender := sideBar
	for _, f := range f {
		filePanelRender = lipgloss.JoinHorizontal(lipgloss.Top, filePanelRender, f)
	}
	return filePanelRender
}

func ProcessBarRender(m model) string {
	// process bar
	processRender := ""

	for _, process := range m.processBarModel.process {
		process.progress.Width = m.fullWidth/3 - 3
		symbol := ""
		line := StringColorRender(theme.ProcessBarSideLine).Render("│ ")
		switch process.state {
		case failure:
			symbol = StringColorRender(theme.Fail).Render("")
		case successful:
			symbol = StringColorRender(theme.Done).Render("")
		case inOperation:
			symbol = StringColorRender(theme.Cancel).Render("󰥔")
		case cancel:
			symbol = StringColorRender(theme.Cancel).Render("")
		}
		processRender += line + TruncateText(process.name, m.fullWidth/3-7) + " " + symbol + "\n"
		processRender += line + process.progress.ViewAs(float64(process.done)/float64(process.total)) + "\n\n"
	}

	bottomBorder := GenerateBottomBorder(fmt.Sprintf("%s/%s", "100", "100"), m.fullWidth/3-3)
	processRender = ProcsssBarBoarder(bottomBarHeight-5, m.fullWidth/3, bottomBorder, m.focusPanel).Render(processRender)

	return processRender
}

func MetaDataRender(m model) string {
	// process bar
	metaDataBar := ""

	maxKeyLength := 0
	sort.Slice(m.fileMetaData.metaData, func(i, j int) bool {

		if m.fileMetaData.metaData[i][0] == "FileName" {
			return true
		} else if m.fileMetaData.metaData[j][0] == "FileName" {
			return false
		} else if m.fileMetaData.metaData[i][0] == "FileSize" {
			return true
		} else if m.fileMetaData.metaData[j][0] == "FileSize" {
			return false
		} else if m.fileMetaData.metaData[i][0] == "FolderName" {
			return true
		} else if m.fileMetaData.metaData[j][0] == "FolderName" {
			return false
		} else if m.fileMetaData.metaData[i][0] == "FolderSize" {
			return true
		} else if m.fileMetaData.metaData[j][0] == "FolderSize" {
			return false
		} else if m.fileMetaData.metaData[i][0] == "FileModifyDate" {
			return true
		} else if m.fileMetaData.metaData[j][0] == "FileModifyDate" {
			return false
		} else if m.fileMetaData.metaData[i][0] == "FileAccessDate" {
			return true
		} else if m.fileMetaData.metaData[j][0] == "FileAccessDate" {
			return false
		} else {
			return m.fileMetaData.metaData[i][0] < m.fileMetaData.metaData[j][0]
		}
	})

	for _, data := range m.fileMetaData.metaData {
		if len(data[0]) > maxKeyLength {
			maxKeyLength = len(data[0])
		}
	}
	for i := m.fileMetaData.renderIndex; i < bottomBarHeight-5+m.fileMetaData.renderIndex && i < len(m.fileMetaData.metaData); i++ {
		if i != m.fileMetaData.renderIndex {
			metaDataBar += "\n"
		}
		data := TruncateMiddleText(m.fileMetaData.metaData[i][1], (m.fullWidth/3)-maxKeyLength-3)
		metaDataBar += fmt.Sprintf("%-*s %s", maxKeyLength+1, m.fileMetaData.metaData[i][0], data)
	}
	bottomBorder := GenerateBottomBorder(fmt.Sprintf("%s/%s", strconv.Itoa(m.fileMetaData.renderIndex+1), strconv.Itoa(len(m.fileMetaData.metaData))), m.fullWidth/3-3)
	metaDataBar = MetaDataBoarder(bottomBarHeight-5, m.fullWidth/3, bottomBorder, m.focusPanel).Render(metaDataBar)

	return metaDataBar
}

func TerminalSizeWarnRender(m model) string {
	focusedModelStyle := lipgloss.NewStyle().
		Height(m.fullHeight).
		Width(m.fullWidth).
		Align(lipgloss.Center, lipgloss.Center).
		BorderForeground(lipgloss.Color("69"))
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
	fullHeightString = terminalMinimumSize.Render(fullHeightString)
	fullWidthString = terminalMinimumSize.Render(fullWidthString)

	return focusedModelStyle.Render(`Terminal size too small:` + "\n" +
		"Width = " + fullWidthString +
		" Height = " + fullHeightString + "\n\n" +

		"Needed for current config:" + "\n" +
		"Width = " + terminalMinimumSize.Render(minimumWidthString) +
		" Height = " + terminalMinimumSize.Render(minimumHeightString))
}

func ModalRender(m model) string {
	if m.createNewItem.itemType == newFile {
		fileLocation := filePanelTopFolderIcon.Render("   ") + filePanelTopPath.Render(TruncateTextBeginning(m.createNewItem.location+"/"+m.createNewItem.textInput.Value(), modalWidth-4)) + "\n"
		confirm := modalConfirm.Render(" (" + Config.Confirm[0] + ") New File ")
		cancel := modalCancel.Render(" (" + Config.Cancel[0] + ") Cancel ")
		tip := confirm + "           " + cancel
		return FullScreenStyle(m.fullHeight, m.fullWidth).Render(FocusedModalStyle(modalHeight, modalWidth).Render(fileLocation + "\n" + m.createNewItem.textInput.View() + "\n\n" + tip))
	} else {
		fileLocation := filePanelTopFolderIcon.Render("   ") + filePanelTopPath.Render(TruncateTextBeginning(m.createNewItem.location+"/"+m.createNewItem.textInput.Value(), modalWidth-4)) + "\n"
		confirm := modalConfirm.Render(" (" + Config.Confirm[0] + ") New Folder ")
		cancel := modalCancel.Render(" (" + Config.Cancel[0] + ") Cancel ")
		tip := confirm + "           " + cancel
		return FullScreenStyle(m.fullHeight, m.fullWidth).Render(FocusedModalStyle(modalHeight, modalWidth).Render(fileLocation + "\n" + m.createNewItem.textInput.View() + "\n\n" + tip))
	}
}
