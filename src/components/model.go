package components

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	sideBarWidth      = 20
	projectTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC241D"))
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C287"))
	pinnedLineStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#A4A2A2"))
	pinnedTextStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#CC241D"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E8751A"))
	cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#8EC07C"))
)

var (
	fileIconStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8EC07C"))
	fileLocation  = lipgloss.NewStyle().Foreground(lipgloss.Color("#458588"))
)

func InitialModel() model {
	return model{
		sideBarModel: sideBarModel{
			pinnedModel: pinnedModel{
				folder: getFolder(),
			},
			choice: "default choice",
			state:  selectDisk,
		},
		fileModel: fileModel{
			fileWindows: []fileWindows{
				{location: "~/Documents/code/returnone/backend", fileState: normal},
				{location: "~/Documents/code/returnone/backend", fileState: normal},
				{location: "~/Documents/code/returnone/backend", fileState: normal},
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
		m.mainPanelHeight = msg.Height - 10
		m.fileModel.width = (msg.Width - sideBarWidth - (4 + (len(m.fileModel.fileWindows)-1)*2)) / len(m.fileModel.fileWindows)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.sideBarModel.cursor > 0 {
				m.sideBarModel.cursor--
			} else {
				m.sideBarModel.cursor = len(m.sideBarModel.pinnedModel.folder) - 1
			}
		case "down", "j":
			if m.sideBarModel.cursor < len(m.sideBarModel.pinnedModel.folder)-1 {
				m.sideBarModel.cursor++
			} else {
				m.sideBarModel.cursor = 0
			}
		case "enter", " ":
			m.sideBarModel.pinnedModel.selected = m.sideBarModel.pinnedModel.folder[m.sideBarModel.cursor].location
		}
	}

	return m, nil
}

func (m model) View() string {
	s := projectTitleStyle.Render("    Super Files     ")
	s += "\n"
	for i, folder := range m.sideBarModel.pinnedModel.folder {
		cursor := " "
		if m.sideBarModel.cursor == i {
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
	s = SideBarBoardStyle(m.mainPanelHeight).Render(s)
	f := make([]string, 3)
	for i, file := range m.fileModel.fileWindows {
		f[i] += fileIconStyle.Render("   ") + fileLocation.Render(TruncateTextBeginning(file.location, m.fileModel.width-4)) + "\n"
		f[i] += pinnedLineStyle.Render(repeatString("─", m.fileModel.width))
		f[i] = FilePanelBoardStyle(m.mainPanelHeight, m.fileModel.width).Render(f[i])
	}
	finalRender := lipgloss.JoinHorizontal(lipgloss.Top, s)

	for _, f := range f {
		finalRender = lipgloss.JoinHorizontal(lipgloss.Top, finalRender, f)
	}

	return finalRender
}
