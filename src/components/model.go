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

func InitialModel() model {
	data, err := os.ReadFile("./.superfile/theme/theme.json")
	if err != nil {
		log.Fatalf("Theme file not exist: %v", err)
	}

	err = json.Unmarshal(data, &theme)
	if err != nil {
		log.Fatalf("Error decoding theme json(your json file may be errors): %v", err)
	}
	LoadThemeConfig()
	return model{
		filePanelFocusIndex: 0,
		sideBarFocus:        false,
		sideBarModel: sideBarModel{
			pinnedModel: pinnedModel{
				folder: getFolder(),
			},
			choice: "default choice",
			state:  selectDisk,
		},
		fileModel: fileModel{
			filePanels: []filePanel{
				{
					render:       0,
					cursor:       0,
					location:     HomeDir,
					fileState:    normal,
					focusType:    focus,
					folderRecord: make(map[string]folderRecord),
				},
			},
			width: 10,
		},
		test: "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle("SuperFile")
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.mainPanelHeight = msg.Height - downBarHeight
		m.fileModel.width = (msg.Width - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fullHeight = msg.Height
		m.fullWidth = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		// return superfile
		case "ctrl+c", "q":
			return m, tea.Quit
		/* LIST CONTROLLER START */
		// up list
		case "up", "k":
			if m.sideBarFocus {
				m = ControllerSideBarListUp(m)
			} else {
				m = ControllerFilePanelListUp(m)
			}
		// down list
		case "down", "j":
			if m.sideBarFocus {
				m = ControllerSideBarListDown(m)
			} else {
				m = ControllerFilePanelListDown(m)
			}
		// select file or disk or folder
		case "enter", " ":
			if m.sideBarFocus {
				m = SideBarSelectFolder(m)
			} else {
				m = EnterPanel(m)
			}
		case "b":
			if !m.sideBarFocus {
				m = ParentFolder(m)
			}
		/* LIST CONTROLLER END */

		/* NAVIGATION CONTROLLER START */
		// change file panel
		case "shift+right", "tab":
			m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
			if m.filePanelFocusIndex == (len(m.fileModel.filePanels) - 1) {
				m.filePanelFocusIndex = 0
			} else {
				m.filePanelFocusIndex++
			}

			m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.sideBarFocus)
		// change file panel
		case "shift+left":
			m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
			if m.filePanelFocusIndex == (len(m.fileModel.filePanels) - 1) {
				m.filePanelFocusIndex = 0
			} else {
				m.filePanelFocusIndex++
			}

			m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.sideBarFocus)
		// close file panel
		case "ctrl+w":
			if len(m.fileModel.filePanels) != 1 {
				m.fileModel.filePanels = append(m.fileModel.filePanels[:m.filePanelFocusIndex], m.fileModel.filePanels[m.filePanelFocusIndex+1:]...)

				if m.filePanelFocusIndex != 0 {
					m.filePanelFocusIndex--
				}
				m.fileModel.width = (m.fullWidth - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
				m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.sideBarFocus)
			}
		// create new file panel
		case "ctrl+n":
			if len(m.fileModel.filePanels) != 4 {
				m.fileModel.filePanels = append(m.fileModel.filePanels, filePanel{
					location:     HomeDir,
					fileState:    normal,
					focusType:    secondFocus,
					folderRecord: make(map[string]folderRecord),
				})

				m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
				m.fileModel.filePanels[m.filePanelFocusIndex+1].focusType = returnFocusType(m.sideBarFocus)
				m.fileModel.width = (m.fullWidth - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
				m.filePanelFocusIndex++

			}
		// focus to sidebar or file panel
		case "ctrl+b":
			if m.sideBarFocus {
				m.sideBarFocus = false
				m.fileModel.filePanels[m.filePanelFocusIndex].focusType = focus
			} else {
				m.sideBarFocus = true
				m.fileModel.filePanels[m.filePanelFocusIndex].focusType = secondFocus
			}
			/* NAVIGATION CONTROLLER END */
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
				bottomBorder := GenerateFilePanelBottomBorder("┫0/0┣", m.fileModel.width+3)
				f[i] = FilePanelBoardStyle(m.mainPanelHeight, m.fileModel.width, filePanel.focusType, bottomBorder).Render(f[i])
			} else {
				for h := filePanel.render; h < filePanel.render+PanelElementHeight(m.mainPanelHeight) && h < len(filePanel.element); h++ {
					cursor := " "
					if h == filePanel.cursor {
						cursor = ""
					}
					f[i] += cursorStyle.Render(cursor) + " " + PrettierName(TruncateText(filePanel.element[h].name, m.fileModel.width-5), filePanel.element[h].folder) + "\n"
				}
				cursorPosition := strconv.Itoa(filePanel.cursor + 1)
				totalElement := strconv.Itoa(len(filePanel.element))
				bottomBorder := GenerateFilePanelBottomBorder(fmt.Sprintf("┫%s/%s┣", cursorPosition, totalElement), m.fileModel.width+3)
				f[i] = FilePanelBoardStyle(m.mainPanelHeight, m.fileModel.width, filePanel.focusType, bottomBorder).Render(f[i])
			}
		}
		finalRender := lipgloss.JoinHorizontal(lipgloss.Top, s)

		for _, f := range f {
			finalRender = lipgloss.JoinHorizontal(lipgloss.Top, finalRender, f)
		}
		return lipgloss.JoinVertical(lipgloss.Top, finalRender, m.test)
	}
}
