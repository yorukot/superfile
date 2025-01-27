package internal

import tea "github.com/charmbracelet/bubbletea"

func wheelMainAction(msg string, m *model, cmd tea.Cmd) tea.Cmd {
	switch msg {

	case "wheel up":
		if m.focusPanel == sidebarFocus {
			m.controlSideBarListUp(true)
		} else if m.focusPanel == processBarFocus {
			m.controlProcessbarListUp(true)
		} else if m.focusPanel == metadataFocus {
			m.controlMetadataListUp(true)
		} else if m.focusPanel == nonePanelFocus {
			m.controlFilePanelListUp(true)
			m.fileMetaData.renderIndex = 0
			go func() {
				m.returnMetaData()
			}()
		}

	case "wheel down":
		if m.focusPanel == sidebarFocus {
			m.controlSideBarListDown(true)
		} else if m.focusPanel == processBarFocus {
			m.controlProcessbarListDown(true)
		} else if m.focusPanel == metadataFocus {
			m.controlMetadataListDown(true)
		} else if m.focusPanel == nonePanelFocus {
			m.controlFilePanelListDown(true)
			m.fileMetaData.renderIndex = 0
			go func() {
				m.returnMetaData()
			}()
		}
	}
	return cmd
}
