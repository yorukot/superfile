package components

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var HomeDir = getHomeDir()

var theme ThemeType
var Config ConfigType

var logOutput *os.File

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
		focusPanel:          nonePanelFocus,
		processBarModel: processBarModel{
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
	return tea.Batch(
		tea.SetWindowTitle("SuperFile"),
		textinput.Blink, // Assuming textinput.Blink is a valid command
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.mainPanelHeight = msg.Height - bottomBarHeight
		m.fileModel.width = (msg.Width - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fullHeight = msg.Height
		m.fullWidth = msg.Width
		return m, nil
	case tea.KeyMsg:
		if m.createNewItem.open {
			switch msg.String() {
			case Config.Cancel[0], Config.Cancel[1]:
				m = CancelModal(m)
			case Config.Confirm[0], Config.Confirm[1]:
				m = CreateItem(m)
			}
		} else if m.fileModel.renaming {
			switch msg.String() {
			case Config.Cancel[0], Config.Cancel[1]:
				m = CancelReanem(m)
			case Config.Confirm[0], Config.Confirm[1]:
				m = ConfirmRename(m)
			}
		} else {
			switch msg.String() {
			// return superfile
			case Config.Quit[0], Config.Quit[1]:
				return m, tea.Quit
			/* LIST CONTROLLER START */
			// up list
			case Config.ListUp[0], Config.ListUp[1]:
				if m.focusPanel == sideBarFocus {
					m = ControllerSideBarListUp(m)
				} else if m.focusPanel == processBarFocus {

				} else {
					m = ControllerFilePanelListUp(m)
				}
			// down list
			case Config.ListDown[0], Config.ListDown[1]:
				if m.focusPanel == sideBarFocus {
					m = ControllerSideBarListDown(m)
				} else if m.focusPanel == processBarFocus {

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
			case Config.FocusOnProcessBar[0], Config.FocusOnProcessBar[0]:
				m = FocusOnProcessBar(m)
			case Config.PasteItem[0], Config.PasteItem[1]:
				go func() {
					m = PasteItem(m)
				}()
			case Config.FilePanelFileCreate[0], Config.FilePanelFileCreate[1]:
				m = PanelCreateNewFile(m)
			case Config.FilePanelFolderCreate[0], Config.FilePanelFolderCreate[1]:
				m = PanelCreateNewFolder(m)
			default:
				// check if it's the select mode
				if m.fileModel.filePanels[m.filePanelFocusIndex].focusType == focus && m.fileModel.filePanels[m.filePanelFocusIndex].panelMode == selectMode {
					switch msg.String() {
					case Config.FilePanelSelectModeItemSingleSelect[0], Config.FilePanelSelectModeItemSingleSelect[1]:
						m = SingleItemSelect(m)
					case Config.FilePanelSelectModeItemSelectUp[0], Config.FilePanelSelectModeItemSelectUp[1]:
						m = ItemSelectUp(m)
					case Config.FilePanelSelectModeItemSelectDown[0], Config.FilePanelSelectModeItemSelectDown[1]:
						m = ItemSelectDown(m)
					case Config.FilePanelSelectModeItemDelete[0], Config.FilePanelSelectModeItemDelete[1]:
						go func() {
							m = DeleteMultipleItem(m)
						}()
					case Config.FilePanelSelectModeItemCopy[0], Config.FilePanelSelectModeItemCopy[1]:
						m = CopyMultipleItem(m)
					case Config.FilePanelSelectModeItemCut[0], Config.FilePanelSelectModeItemCut[1]:
						m = CutMultipleItem(m)
					case Config.FilePanelSelectAllItem[0], Config.FilePanelSelectAllItem[1]:
						m = SelectAllItem(m)
					}
					// else
				} else {
					switch msg.String() {
					case Config.SelectItem[0], Config.SelectItem[1]:
						if m.focusPanel == sideBarFocus {
							m = SideBarSelectFolder(m)
						} else if m.focusPanel == processBarFocus {

						} else {
							m = EnterPanel(m)
						}
					case Config.ParentFolder[0], Config.ParentFolder[1]:
						m = ParentFolder(m)
					case Config.DeleteItem[0], Config.DeleteItem[1]:
						go func() {
							m = DeleteSingleItem(m)
						}()
					case Config.CopySingleItem[0], Config.CopySingleItem[1]:
						m = CopySingleItem(m)
					case Config.CutSingleItem[0], Config.CutSingleItem[1]:
						m = CutSingleItem(m)
					case Config.FilePanelItemRename[0], Config.FilePanelItemRename[1]:
						m = PanelItemRename(m)
					}
				}
			}
		}
	}
	if m.firstTextInput {
		m.firstTextInput = false
	} else if m.fileModel.renaming {
		m.fileModel.filePanels[m.filePanelFocusIndex].rename, cmd = m.fileModel.filePanels[m.filePanelFocusIndex].rename.Update(msg)
	} else {
		m.createNewItem.textInput, cmd = m.createNewItem.textInput.Update(msg)

	}
	return m, cmd
}

func (m model) View() string {
	// check is the terminal size enough
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
	} else if m.createNewItem.open {
		if m.createNewItem.itemType == rename {
			fileLocation := filePanelTopFolderIcon.Render("   ") + filePanelTopPath.Render(TruncateTextBeginning(m.createNewItem.location+"/"+m.createNewItem.textInput.Value(), modalWidth-4)) + "\n"
			confirm := modalConfirm.Render(" (" + Config.Confirm[0] + ") Confirm ")
			cancel := modalCancel.Render(" (" + Config.Cancel[0] + ") Cancel ")
			tip := confirm + "           " + cancel
			return FullScreenStyle(m.fullHeight, m.fullWidth).Render(FocusedModalStyle(modalHeight, modalWidth).Render(fileLocation + "\n" + m.createNewItem.textInput.View() + "\n\n" + tip))

		} else {
			fileLocation := filePanelTopFolderIcon.Render("   ") + filePanelTopPath.Render(TruncateTextBeginning(m.createNewItem.location+"/"+m.createNewItem.textInput.Value(), modalWidth-4)) + "\n"
			confirm := modalConfirm.Render(" (" + Config.Confirm[0] + ") Confirm ")
			cancel := modalCancel.Render(" (" + Config.Cancel[0] + ") Cancel ")
			tip := confirm + "           " + cancel
			return FullScreenStyle(m.fullHeight, m.fullWidth).Render(FocusedModalStyle(modalHeight, modalWidth).Render(fileLocation + "\n" + m.createNewItem.textInput.View() + "\n\n" + tip))
		}
	} else {
		// side bar
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
		filePanelRender := s
		for _, f := range f {
			filePanelRender = lipgloss.JoinHorizontal(lipgloss.Top, filePanelRender, f)
		}

		// process bar
		processRender := ""
		for _, process := range m.processBarModel.process {
			process.progress.Width = m.fullWidth/3 - 3
			symbol := ""
			line := StringColorRender(theme.ProcessBarSideLine).Render("│ ")
			if process.state == failure {
				symbol = StringColorRender(theme.Fail).Render("")
			} else if process.state == successful {
				symbol = StringColorRender(theme.Done).Render("")
			} else {
				symbol = StringColorRender(theme.Cancel).Render("")
			}
			processRender += line + TruncateText(process.name, m.fullWidth/3-7) + " " + symbol + "\n"
			processRender += line + process.progress.ViewAs(1) + "\n\n"
		}
		bottomBorder := GenerateBottomBorder(fmt.Sprintf("%s┣━┫%s/%s", "100sec", "100", "100"), m.fullWidth/3+2)
		processRender = ProcsssBarBoarder(bottomBarHeight-5, m.fullWidth/3, bottomBorder, m.focusPanel).Render(processRender)
		// final render
		finalRender := lipgloss.JoinVertical(0, filePanelRender, processRender)
		return lipgloss.JoinVertical(lipgloss.Top, finalRender)
	}
}
