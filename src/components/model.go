package components

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"strconv"
)

var HomeDir = getHomeDir()

var theme ThemeType
var Config ConfigType

var logOutput *os.File

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

func InitialModel() model {
	var err error
	logOutput, err = os.OpenFile("superfile.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	data, err := os.ReadFile("./.superfile/config/config.json")
	if err != nil {
		log.Fatalf("HotKey file not exist: %v", err)
	}

	err = json.Unmarshal(data, &Config)
	if err != nil {
		log.Fatalf("Error decoding HotKey json(your config  file may be errors): %v", err)
	}

	data, err = os.ReadFile("./.superfile/theme/theme.json")
	if err != nil {
		log.Fatalf("Theme file not exist: %v", err)
	}

	err = json.Unmarshal(data, &theme)
	if err != nil {
		log.Fatalf("Error decoding theme json(your theme file may be errors): %v", err)
	}
	LoadThemeConfig()
	return model{
		filePanelFocusIndex: 0,
		sideBarFocus:        false,
		procsssBarFocus:     false,
		processBar: processBar{
			process: []process{},
			cursor:  0,
		},
		sideBarModel: sideBarModel{
			pinnedModel: pinnedModel{
				folder: getFolder(),
			},
		},
		fileModel: fileModel{
			filePanels: []filePanel{
				{
					render:       0,
					cursor:       0,
					location:     HomeDir,
					panelMode:    browserMode,
					focusType:    focus,
					folderRecord: make(map[string]folderRecord),
				},
			},
			width: 10,
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("SuperFile")
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.mainPanelHeight = msg.Height - bottomBarHeight
		m.fileModel.width = (msg.Width - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fullHeight = msg.Height
		m.fullWidth = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		// return superfile
		case Config.Quit[0], Config.Quit[1]:
			return m, tea.Quit
		/* LIST CONTROLLER START */
		// up list
		case Config.ListUp[0], Config.ListUp[1]:
			if m.sideBarFocus {
				m = ControllerSideBarListUp(m)
			} else {
				m = ControllerFilePanelListUp(m)
			}
		// down list
		case Config.ListDown[0], Config.ListDown[1]:
			if m.sideBarFocus {
				m = ControllerSideBarListDown(m)
			} else {
				m = ControllerFilePanelListDown(m)
			}
		/* LIST CONTROLLER END */
		case Config.ChangePanelMode[0], Config.ChangePanelMode[1]:
			m = SelectedMode(m)
		/* NAVIGATION CONTROLLER START */
		// change file panel
		case Config.NextFilePanel[0], Config.NextFilePanel[1]:
			m = NextFilePanel(m)
		// change file panel
		case Config.PreviousFilePanel[0], Config.PreviousFilePanel[1]:
			m = PreviousFilePanel(m)
		// close file panel
		case Config.CloseFilePanel[0], Config.CloseFilePanel[1]:
			m = CloseFilePanel(m)
		// create new file panel
		case Config.CreateNewFilePanel[0], Config.CreateNewFilePanel[1]:
			m = CreateNewFilePanel(m)
		// focus to sidebar or file panel
		case Config.FocusOnSideBar[0], Config.FocusOnSideBar[1]:
			m = FocusOnSideBar(m)
			/* NAVIGATION CONTROLLER END */
		default:
			if m.fileModel.filePanels[m.filePanelFocusIndex].focusType == focus && m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
				switch msg.String() {
				case Config.FilePanelSelectModeItemSingleSelect[0], Config.FilePanelSelectModeItemSingleSelect[1]:
					m = SingleItemSelect(m)
				case Config.FilePanelSelectModeItemSelectUp[0], Config.FilePanelSelectModeItemSelectUp[1]:
					m = ItemSelectUp(m)
				case Config.FilePanelSelectModeItemSelectDown[0], Config.FilePanelSelectModeItemSelectDown[1]:
					m = ItemSelectDown(m)
				}
			} else {
				switch msg.String() {
				// select file or disk or folder
				case Config.SelectItem[0], Config.SelectItem[1]:
					if m.sideBarFocus {
						m = SideBarSelectFolder(m)
					} else {
						m = EnterPanel(m)
					}
				/* LIST CONTROLLER END */
				case Config.ParentFolder[0], Config.ParentFolder[1]:
					m = ParentFolder(m)
				case Config.DeleteItem[0], Config.DeleteItem[1]:
					m = DeleteSingleItem(m)
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.fullHeight < minimumHeight || m.fullWidth < minimumWidth {
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
	} else {

		s := sideBarTitle.Render("    Super Files     ")
		s += "\n"
		for i, folder := range m.sideBarModel.pinnedModel.folder {
			cursor := " "
			if m.sideBarModel.cursor == i && m.sideBarFocus {
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

		s = SideBarBoardStyle(m.mainPanelHeight, m.sideBarFocus).Render(s)

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
					f[i] += cursorStyle.Render(cursor) + " " + PrettierName(TruncateText(filePanel.element[h].name, m.fileModel.width-5), filePanel.element[h].folder, isItemSelected) + "\n"
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
		bottomBorder := GenerateBottomBorder(fmt.Sprintf("%s┣━┫%s/%s", "100sec", "100", "100"), m.fileModel.width+6)
		finalRender := s
		for _, f := range f {
			finalRender = lipgloss.JoinHorizontal(lipgloss.Top, finalRender, f)
		}
		processRender := ""
		for _, process := range m.processBar.process {
			process.progress.Width = m.fullWidth/3 - 3
			symbol := ""
			line := ""
			if process.state == failure {
				symbol = StringColorRender("#FFAE00").Render("")
				line = StringColorRender("#FFAE00").Render("│ ")
			} else if process.state == successful {
				symbol = StringColorRender("#77FF00").Render("")
				line = StringColorRender("#77FF00").Render("│ ")
			} else {
				symbol = StringColorRender("#A1A1A1").Render("")
				line = StringColorRender("#A1A1A1").Render("│ ")
			}
			processRender += line + TruncateText(process.name, m.fullWidth/3-7) + " " + symbol + "\n"
			processRender += line + process.progress.ViewAs(1) + "\n\n"
		}
		finalRender = lipgloss.JoinVertical(0, finalRender, ProcsssBarBoarder(bottomBarHeight-5, m.fullWidth/3, m.procsssBarFocus, bottomBorder).Render(processRender))
		return lipgloss.JoinVertical(lipgloss.Top, finalRender)
	}
}
