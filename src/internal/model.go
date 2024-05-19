package internal

import (
	"log"
	"os"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	varibale "github.com/yorukot/superfile/src/config"
	stringfunction "github.com/yorukot/superfile/src/pkg/stringFunction"
)

var LastTimeCursorMove = [2]int{int(time.Now().UnixMicro()), 0}
var ListeningMessage = true

var forceReloadElement = false
var firstUse = false

var theme ThemeType
var Config ConfigType
var hotkeys HotkeysType

var logOutput *os.File
var et *exiftool.Exiftool

var channel = make(chan channelMessage, 1000)

func InitialModel(dir string, firstUseCheck bool) model {
	toggleDotFileBool, firstFilePanelDir := initialConfig(dir)
	firstUse = firstUseCheck
	return defaultModelConfig(toggleDotFileBool, firstFilePanelDir)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("SuperFile"),
		textinput.Blink, // Assuming textinput.Blink is a valid command
		listenForChannelMessage(channel),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	switch msg := msg.(type) {
	case channelMessage:
		if msg.returnWarnModal {
			m.warnModal = msg.warnModal
		} else if msg.loadMetadata {
			m.fileMetaData.metaData = msg.metadata
		} else {
			if !arrayContains(m.processBarModel.processList, msg.messageId) {
				m.processBarModel.processList = append(m.processBarModel.processList, msg.messageId)
			}
			m.processBarModel.process[msg.messageId] = msg.processNewState
		}
		forceReloadElement = true
	case tea.WindowSizeMsg:
		m.fullHeight = msg.Height
		m.fullWidth = msg.Width

		if m.fileModel.filePreview.open {
			// File preview panel width same as file panel
			if Config.FilePreviewWidth == 0 {
				m.fileModel.filePreview.width = (msg.Width - sidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
			} else {
				if Config.FilePreviewWidth > 10 || Config.FilePreviewWidth == 1 {
					log.Fatalln("Config file file_preview_width invalidation")
				}
				m.fileModel.filePreview.width = (msg.Width - sidebarWidth) / Config.FilePreviewWidth
			}
		}

		// set each file panel size and max file panel amount
		m.fileModel.width = (msg.Width - sidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fileModel.maxFilePanel = (msg.Width - sidebarWidth - m.fileModel.filePreview.width) / 20
		for i := range m.fileModel.filePanels {
			m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
		}

		// set footer size
		if msg.Height < 30 {
			footerHeight = 10
		} else if msg.Height < 35 {
			footerHeight = 11
		} else if msg.Height < 40 {
			footerHeight = 12
		} else if msg.Height < 45 {
			footerHeight = 13
		} else {
			footerHeight = 14
		}

		m.mainPanelHeight = msg.Height - footerHeight + 1

		// set help menu size
		m.helpMenu.height = m.fullHeight - 2
		m.helpMenu.width = m.fullWidth - 2

		if m.fullHeight > 35 {
			m.helpMenu.height = 30
		}

		if m.fullWidth > 95 {
			m.helpMenu.width = 90
		}

		if m.fileModel.maxFilePanel >= 10 {
			m.fileModel.maxFilePanel = 10
		}
		return m, nil
	case tea.KeyMsg:

		if firstUse {
			firstUse = false
			return m, cmd
		}

		if m.typingModal.open {
			m = typingModalOpenKey(msg.String(), m)
		} else if m.warnModal.open {
			m = warnModalOpenKey(msg.String(), m)
		} else if m.fileModel.renaming {
			m = renamingKey(msg.String(), m)
		} else if panel.searchBar.Focused() {
			m = focusOnSearchbarKey(msg.String(), m)
		} else if m.helpMenu.open {
			m = helpMenuKey(msg.String(), m)
		} else {
			// return superfile
			if msg.String() == hotkeys.Quit[0] || msg.String() == hotkeys.Quit[1] {
				// cd on quit
				if Config.CdOnQuit {
					currentDir := m.fileModel.filePanels[m.filePanelFocusIndex].location
					if currentDir == varibale.HomeDir {
						return m, tea.Quit
					}
					os.WriteFile(varibale.SuperFileStateDir+"/lastdir", []byte("cd "+currentDir), 0755)
				}
				return m, tea.Quit
			}

			m, cmd = mainKey(msg.String(), m, cmd)

			if m.editorMode {
				m.editorMode = false
				return m, cmd
			}
		}
	}

	if m.firstTextInput {
		m.firstTextInput = false
	} else if m.fileModel.renaming {
		m.fileModel.filePanels[m.filePanelFocusIndex].rename, cmd = m.fileModel.filePanels[m.filePanelFocusIndex].rename.Update(msg)
	} else if panel.searchBar.Focused() {
		m.fileModel.filePanels[m.filePanelFocusIndex].searchBar, cmd = m.fileModel.filePanels[m.filePanelFocusIndex].searchBar.Update(msg)
	} else if m.typingModal.open {
		m.typingModal.textInput, cmd = m.typingModal.textInput.Update(msg)
	}

	if m.fileModel.filePanels[m.filePanelFocusIndex].cursor < 0 {
		m.fileModel.filePanels[m.filePanelFocusIndex].cursor = 0
	}

	cmd = tea.Batch(cmd)
	m.sidebarModel.directories = getDirectories()

	// check if there already have listening message
	if !ListeningMessage {
		cmd = tea.Batch(cmd, listenForChannelMessage(channel))
	}

	m = getFilePanelItems(m)

	return m, tea.Batch(cmd)
}

func (m model) View() string {
	// check is the terminal size enough
	if m.fullHeight < minimumHeight || m.fullWidth < minimumWidth {
		return terminalSizeWarnRender(m)
	}
	if m.fileModel.width < 18 {
		return terminalSizeWarnAfterFirstRender(m)
	}
	sidebar := sidebarRender(m)

	filePanel := filePanelRender(m)

	filePreview := filePreviewPanelRender(m)

	mainPanel := lipgloss.JoinHorizontal(0, sidebar, filePanel, filePreview)

	
	processBar := processBarRender(m)

	metaData := metadataRender(m)

	clipboardBar := clipboardRender(m)

	footer := lipgloss.JoinHorizontal(0, processBar, metaData, clipboardBar)

	finalRender := lipgloss.JoinVertical(0, mainPanel, footer)

	// check if need pop up modal
	if m.helpMenu.open {
		helpMenu := helpMenuRender(m)
		overlayX := m.fullWidth/2 - m.helpMenu.width/2
		overlayY := m.fullHeight/2 - m.helpMenu.height/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, helpMenu, finalRender)
	}

	if firstUse {
		introduceModal := introduceModalRender(m)
		overlayX := m.fullWidth/2 - m.helpMenu.width/2
		overlayY := m.fullHeight/2 - m.helpMenu.height/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, introduceModal, finalRender)
	}

	if m.typingModal.open {
		typingModal := typineModalRender(m)
		overlayX := m.fullWidth/2 - modalWidth/2
		overlayY := m.fullHeight/2 - modalHeight/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, typingModal, finalRender)
	}

	if m.warnModal.open {
		warnModal := warnModalRender(m)
		overlayX := m.fullWidth/2 - modalWidth/2
		overlayY := m.fullHeight/2 - modalHeight/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, warnModal, finalRender)
	}

	return finalRender
}

func listenForChannelMessage(msg chan channelMessage) tea.Cmd {
	return func() tea.Msg {
		m := <-msg
		ListeningMessage = false
		return m
	}
}

func getFilePanelItems(m model) model {
	for i, filePanel := range m.fileModel.filePanels {
		var fileElenent []element
		nowTime := time.Now()
		if filePanel.focusType == noneFocus && nowTime.Sub(filePanel.lastTimeGetElement) < 3*time.Second && !forceReloadElement {
			continue
		}
		if len(filePanel.element) > 500 && (len(filePanel.element) > 500 && (nowTime.Sub(filePanel.lastTimeGetElement) > 3*time.Second)) && !forceReloadElement {
			continue
		}

		if filePanel.searchBar.Value() != "" {
			fileElenent = returnFolderElementBySearchString(filePanel.location, m.toggleDotFile, filePanel.searchBar.Value())
		} else {
			fileElenent = returnFolderElement(filePanel.location, m.toggleDotFile)
		}
		filePanel.element = fileElenent
		m.fileModel.filePanels[i].element = fileElenent
		m.fileModel.filePanels[i].lastTimeGetElement = nowTime
		forceReloadElement = false
	}
	return m
}
