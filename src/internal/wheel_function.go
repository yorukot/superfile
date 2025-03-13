package internal

import tea "github.com/charmbracelet/bubbletea"

func wheelMainAction(msg string, m *model, cmd tea.Cmd) tea.Cmd {
	var action func()
	switch msg {

	case "wheel up":
		if m.focusPanel == sidebarFocus {
			action = func() { m.sidebarModel.listUp(m.mainPanelHeight) }
		} else if m.focusPanel == processBarFocus {
			action = func() { m.processBarModel.listUp(footerHeight) }
		} else if m.focusPanel == metadataFocus {
			action = func() { m.fileMetaData.listUp() }
		} else if m.focusPanel == nonePanelFocus {
			action = func() { m.fileModel.filePanels[m.filePanelFocusIndex].listUp(m.mainPanelHeight) }
		}

	case "wheel down":
		if m.focusPanel == sidebarFocus {
			action = func() { m.sidebarModel.listDown(m.mainPanelHeight) }
		} else if m.focusPanel == processBarFocus {
			action = func() { m.processBarModel.listDown() }
		} else if m.focusPanel == metadataFocus {
			action = func() { m.fileMetaData.listDown() }
		} else if m.focusPanel == nonePanelFocus {
			action = func() { m.fileModel.filePanels[m.filePanelFocusIndex].listDown(m.mainPanelHeight) }
		}
	}

	for i := 0; i < wheelRunTime; i++ {
		action()
	}

	if m.focusPanel == nonePanelFocus {
		m.fileMetaData.renderIndex = 0
		go func() {
			m.returnMetaData()
		}()
	}

	return cmd
}
