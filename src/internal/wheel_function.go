package internal

import tea "github.com/charmbracelet/bubbletea"

func wheelMainAction(msg string, m model, cmd tea.Cmd) (model, tea.Cmd) {
	switch msg {

	case "wheel up":
		if m.focusPanel == sidebarFocus {
			m = controlSideBarListUp(m, true)
		} else if m.focusPanel == processBarFocus {
			m = controlProcessbarListUp(m, true)
		} else if m.focusPanel == metadataFocus {
			m = controlMetadataListUp(m, true)
		} else if m.focusPanel == nonePanelFocus {
			m = controlFilePanelListUp(m, true)
			m.fileMetaData.renderIndex = 0
			go func() {
				m = returnMetaData(m)
			}()
		}

	case "wheel down":
		if m.focusPanel == sidebarFocus {
			m = controlSideBarListDown(m, true)
		} else if m.focusPanel == processBarFocus {
			m = controlProcessbarListDown(m, true)
		} else if m.focusPanel == metadataFocus {
			m = controlMetadataListDown(m, true)
		} else if m.focusPanel == nonePanelFocus {
			m = controlFilePanelListDown(m, true)
			m.fileMetaData.renderIndex = 0
			go func() {
				m = returnMetaData(m)
			}()
		}
	}
	return m, cmd
}
