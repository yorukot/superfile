package components

import (
	"os"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rkoesters/xdg/basedir"
)

const (
	configFolder     string = "/config"
	themeFolder      string = "/theme"
	dataFolder       string = "/data"
	lastCheckVersion string = "/lastCheckVersion"
	pinnedFile       string = "/pinned.json"
	toggleDotFile    string = "/toggleDotFile"
	logFile          string = "/superfile.log"
	configFile       string = "/config.toml"
	hotkeysFile      string = "/hotkeys.toml"
	themeZipName     string = "/theme.zip"
)

var LastTimeCursorMove = [2]int{int(time.Now().UnixMicro()), 0}
var ListeningMessage = true

var forceReloadElement = false

var HomeDir = basedir.Home
var SuperFileMainDir = basedir.ConfigHome + "/superfile"
var SuperFileCacheDir = basedir.CacheHome + "/superfile"
var SuperFileDataDir = basedir.DataHome + "/superfile"

var theme ThemeType
var Config ConfigType
var hotkeys HotkeysType

var logOutput *os.File
var et *exiftool.Exiftool

var channel = make(chan channelMessage, 1000)

func InitialModel(dir string) model {
	toggleDotFileBool, firstFilePanelDir := loadConfigFile(dir)

	return model{
		filePanelFocusIndex: 0,
		focusPanel:          nonePanelFocus,
		processBarModel: processBarModel{
			process: make(map[string]process),
			cursor:  0,
			render:  0,
		},
		sidebarModel: sidebarModel{
			directories: getDirectories(40),
		},
		fileModel: fileModel{
			filePanels: []filePanel{
				{
					render:          0,
					cursor:          0,
					location:        firstFilePanelDir,
					panelMode:       browserMode,
					focusType:       focus,
					directoryRecord: make(map[string]directoryRecord),
					searchBar:       generateSearchBar(),
				},
			},
			width: 10,
		},
		toggleDotFile: toggleDotFileBool,
	}
}

func listenForChannelMessage(msg chan channelMessage) tea.Cmd {
	return func() tea.Msg {
		m := <-msg
		ListeningMessage = false
		return m
	}
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
		if msg.Height < 30 {
			footerHeight = 10
		} else {
			footerHeight = 14
		}
		m.mainPanelHeight = msg.Height - footerHeight + 1
		m.fileModel.width = (msg.Width - sidebarWidth - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
		m.fullHeight = msg.Height
		m.fullWidth = msg.Width
		m.fileModel.maxFilePanel = (msg.Width - 20) / 24
		if m.fileModel.maxFilePanel >= 10 {
			m.fileModel.maxFilePanel = 10
		}
		return m, nil
	case tea.KeyMsg:
		if m.typingModal.open {
			m = typingModalOpenKey(msg.String(), m)
		} else if m.warnModal.open {
			m = warnModalOpenKey(msg.String(), m)
		} else if m.fileModel.renaming {
			m = renamingKey(msg.String(), m)
		} else if panel.searchBar.Focused() {
			m = focusOnSearchbarKey(msg.String(), m)
		} else {
			// return superfile
			if msg.String() == hotkeys.Quit[0] || msg.String() == hotkeys.Quit[1] {
				return m, tea.Quit
			}
			m = mainKey(msg.String(), m)
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
	m.sidebarModel.directories = getDirectories(m.fullHeight)

	if !ListeningMessage {
		cmd = tea.Batch(cmd, listenForChannelMessage(channel))
	}
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
	return m, cmd
}

func (m model) View() string {
	// check is the terminal size enough
	if m.fullHeight < minimumHeight || m.fullWidth < minimumWidth {
		return terminalSizeWarnRender(m)
	} else if m.typingModal.open {
		return typineModalRender(m)
	} else if m.warnModal.open {
		return warnModalRender(m)
	} else {
		sidebar := sidebarRender(m)

		filePanel := filePanelRender(m)

		mainPanel := lipgloss.JoinHorizontal(0, sidebar, filePanel)

		processBar := processBarRender(m)

		metaData := metadataRender(m)

		clipboardBar := clipboardRender(m)

		footer := lipgloss.JoinHorizontal(0, processBar, metaData, clipboardBar)

		// final render
		finalRender := lipgloss.JoinVertical(0, mainPanel, footer)

		return lipgloss.JoinVertical(lipgloss.Top, finalRender)
	}
}
