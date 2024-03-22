package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var HomeDir = getHomeDir()

func InitialModel() model {
	return model{
		filePanelFocusIndex: 0,
		sideBarFocus:        true,
		sideBarModel: sideBarModel{
			pinnedModel: pinnedModel{
				folder: getFolder(),
			},
			choice: "default choice",
			state:  selectDisk,
		},
		fileModel: fileModel{
			filePanels: []filePanel{
				{location: HomeDir + "/Documents/code/", fileState: normal, focusType: secondFocus},
				{location: HomeDir + "/Documents/code/", fileState: normal, focusType: noneFocus},
				{location: HomeDir + "/Documents/code/", fileState: normal, focusType: noneFocus},
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
		m.mainPanelHeight = msg.Height - downBarHeight
		m.fileModel.width = (msg.Width - sideBarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		// return superfile
		case "ctrl+c", "q":
			return m, tea.Quit
		/* LIST CONTROLLER START */
		// up list
		case "up", "k":
			if m.sideBarModel.cursor > 0 {
				m.sideBarModel.cursor--
			} else {
				m.sideBarModel.cursor = len(m.sideBarModel.pinnedModel.folder) - 1
			}
		// down list
		case "down", "j":
			if m.sideBarModel.cursor < len(m.sideBarModel.pinnedModel.folder)-1 {
				m.sideBarModel.cursor++
			} else {
				m.sideBarModel.cursor = 0
			}
		// select file or disk or folder
		case "enter", " ":
			m.sideBarModel.pinnedModel.selected = m.sideBarModel.pinnedModel.folder[m.sideBarModel.cursor].location
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

				m.fileModel.filePanels[m.filePanelFocusIndex].focusType = returnFocusType(m.sideBarFocus)
			}
		// create new file panel
		case "ctrl+n":
			if len(m.fileModel.filePanels) != 3 {
				m.fileModel.filePanels = append(m.fileModel.filePanels, filePanel{location: HomeDir + "/", fileState: normal, focusType: secondFocus})

				m.fileModel.filePanels[m.filePanelFocusIndex].focusType = noneFocus
				m.fileModel.filePanels[m.filePanelFocusIndex+1].focusType = returnFocusType(m.sideBarFocus)

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
	s := projectTitleStyle.Render("    Super Files     ")
	s += "\n"
	for i, folder := range m.sideBarModel.pinnedModel.folder {
		cursor := " "
		if m.sideBarModel.cursor == i && m.sideBarFocus {
			cursor = ""
		}

		if m.sideBarModel.pinnedModel.selected == folder.location {
			s += cursorStyle.Render(cursor) + " " + selectedItemStyle.Render(TruncateText(folder.name, sideBarWidth-2)) + "" + "\n"
		} else {
			s += cursorStyle.Render(cursor) + " " + itemStyle.Render(TruncateText(folder.name, sideBarWidth-2)) + "" + "\n"
		}

		if i == 4 {
			s += "\n" + pinnedTextStyle.Render("󰐃 Pinned") + pinnedLineStyle.Render(" ───────────") + "\n\n"
		}

		if folder.endPinned {
			s += "\n" + pinnedTextStyle.Render("󱇰 Disk") + pinnedLineStyle.Render(" ─────────────") + "\n\n"
		}
	}

	s = SideBarBoardStyle(m.mainPanelHeight, m.sideBarFocus).Render(s)

	f := make([]string, 3)
	for i, filePanel := range m.fileModel.filePanels {
		filePanel.element = returnFolderElement(filePanel.location)
		f[i] += fileIconStyle.Render("   ") + fileLocation.Render(TruncateTextBeginning(filePanel.location, m.fileModel.width-4)) + "\n"
		f[i] += FilePanelDividerStyle(filePanel.focusType).Render(repeatString("─", m.fileModel.width)) + "\n"
		for _, file := range filePanel.element {
			cursor := " "

			f[i] += cursorStyle.Render(cursor) + " " + itemStyle.Render(TruncateText(file.name, m.fileModel.width-2)) + "" + "\n"
		}
		f[i] = FilePanelBoardStyle(m.mainPanelHeight, m.fileModel.width, filePanel.focusType).Render(f[i])
	}
	finalRender := lipgloss.JoinHorizontal(lipgloss.Top, s)

	for _, f := range f {
		finalRender = lipgloss.JoinHorizontal(lipgloss.Top, finalRender, f)
	}

	return finalRender
}
