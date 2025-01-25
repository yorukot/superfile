package internal

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	variable "github.com/yorukot/superfile/src/config"
	stringfunction "github.com/yorukot/superfile/src/pkg/string_function"
)

var LastTimeCursorMove = [2]int{int(time.Now().UnixMicro()), 0}
var ListeningMessage = true

var firstUse = false
var hasTrash = true

var theme ThemeType
var Config ConfigType
var hotkeys HotkeysType

var et *exiftool.Exiftool

var channel = make(chan channelMessage, 1000)
var progressBarLastRenderTime time.Time = time.Now()

// Initialize and return model with default configs
func InitialModel(dir string, firstUseCheck bool, hasTrashCheck bool) model {
	toggleDotFileBool, toggleFooter, firstFilePanelDir := initialConfig(dir)
	firstUse = firstUseCheck
	hasTrash = hasTrashCheck
	return defaultModelConfig(toggleDotFileBool, toggleFooter, firstFilePanelDir)
}

// Init function to be called by Bubble tea framework, sets windows title,
// cursos blinking and starts message streamming channel
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("SuperFile"),
		textinput.Blink, // Assuming textinput.Blink is a valid command
		listenForChannelMessage(channel),
	)
}

// Update function for bubble tea to provide internal communication to the
// application
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case channelMessage:
		m.handleChannelMessage(msg)
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
	case tea.MouseMsg:
		m, cmd = wheelMainAction(msg.String(), m, cmd)
	case tea.KeyMsg:
		m, cmd = m.handleKeyInput(msg, cmd)
	}

	m.updateFilePanelsState(msg, &cmd)
	m.sidebarModel.directories = getDirectories()

	// check if there already have listening message
	if !ListeningMessage {
		cmd = tea.Batch(cmd, listenForChannelMessage(channel))
	}

	m.getFilePanelItems()

	return m, tea.Batch(cmd)
}

// Handle message exchanging whithin the application
func (m *model) handleChannelMessage(msg channelMessage) {
	switch msg.messageType {
	case sendWarnModal:
		m.warnModal = msg.warnModal
	case sendMetadata:
		m.fileMetaData.metaData = msg.metadata
	default:
		if !arrayContains(m.processBarModel.processList, msg.messageId) {
			m.processBarModel.processList = append(m.processBarModel.processList, msg.messageId)
		}
		m.processBarModel.process[msg.messageId] = msg.processNewState
	}
}

// Adjust window size based on msg information
func (m *model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.fullHeight = msg.Height
	m.fullWidth = msg.Width

	if m.fileModel.filePreview.open {
		// File preview panel width same as file panel
		m.setFilePreviewWidth(msg.Width)
	}

	m.setFilePanelsSize(msg.Width)
	m.setFooterSize(msg.Height)
	m.setHelpMenuSize()

	if m.fileModel.maxFilePanel >= 10 {
		m.fileModel.maxFilePanel = 10
	}
}

// Set file preview panel Widht to width. Assure that
func (m *model) setFilePreviewWidth(width int) {
	if Config.FilePreviewWidth == 0 {
		m.fileModel.filePreview.width = (width - Config.SidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
	} else if Config.FilePreviewWidth > 10 || Config.FilePreviewWidth == 1 {
		LogAndExit("Config file file_preview_width invalidation")
	} else {
		m.fileModel.filePreview.width = (width - Config.SidebarWidth) / Config.FilePreviewWidth
	}
}

// Proper set panels size. Assure that panels do not overlap
func (m *model) setFilePanelsSize(width int) {
	// set each file panel size and max file panel amount
	m.fileModel.width = (width - Config.SidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.fileModel.maxFilePanel = (width - Config.SidebarWidth - m.fileModel.filePreview.width) / 20
	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
}

// Set footer size using height
func (m *model) setFooterSize(height int) {
	if !m.toggleFooter {
		footerHeight = 0
	} else if height < 30 {
		footerHeight = 10
	} else if height < 35 {
		footerHeight = 11
	} else if height < 40 {
		footerHeight = 12
	} else if height < 45 {
		footerHeight = 13
	} else {
		footerHeight = 14
	}

	if m.commandLine.input.Focused() && m.toggleFooter {
		footerHeight--
	}

	if m.toggleFooter {
		m.mainPanelHeight = height - footerHeight + 1
	} else {
		m.mainPanelHeight = height - 2
	}
}

// Set help menu size
func (m *model) setHelpMenuSize() {
	m.helpMenu.height = m.fullHeight - 2
	m.helpMenu.width = m.fullWidth - 2

	if m.fullHeight > 35 {
		m.helpMenu.height = 30
	}

	if m.fullWidth > 95 {
		m.helpMenu.width = 90
	}
}

// Identify the current state of the application m and properly handle the
// msg keybind pressed
func (m model) handleKeyInput(msg tea.KeyMsg, cmd tea.Cmd) (model, tea.Cmd) {
	
	slog.Debug("model.handleKeyInput", "msg", msg, "typestr", msg.Type.String(),
		"runes", msg.Runes, "type", int(msg.Type), "paste", msg.Paste, 
		"alt", msg.Alt)
	slog.Debug("model.handleKeyInput. model info. ",
		"filePanelFocusIndex", m.filePanelFocusIndex, 
		"filePanel.focusType", m.fileModel.filePanels[m.filePanelFocusIndex].focusType,
		"filePanel.panelMode", m.fileModel.filePanels[m.filePanelFocusIndex].panelMode,
		"typingModal.open", m.typingModal.open,
		"warnModal.open", m.warnModal.open,
		"fileModel.renaming", m.fileModel.renaming,
		"searchBar.focussed", m.fileModel.filePanels[m.filePanelFocusIndex].searchBar.Focused(),
		"helpMenu.open", m.helpMenu.open,
		"firstTextInput", m.firstTextInput,
		"focusPanel", m.focusPanel,
	)

	if firstUse {
		firstUse = false
		return m, cmd
	}

	if m.typingModal.open {
		m.typingModalOpenKey(msg.String())
	} else if m.warnModal.open {
		m.warnModalOpenKey(msg.String())
		// If renaming a object
	} else if m.fileModel.renaming {
		m.renamingKey(msg.String())
		// If search bar is open
	} else if m.fileModel.filePanels[m.filePanelFocusIndex].searchBar.Focused() {
		m.focusOnSearchbarKey(msg.String())
		// If sort options menu is open
	} else if m.fileModel.filePanels[m.filePanelFocusIndex].sortOptions.open {
		m.sortOptionsKey(msg.String())
		// If help menu is open
	} else if m.helpMenu.open {
		m.helpMenuKey(msg.String())
		// If command line input is send
	} else if m.commandLine.input.Focused() {
		m.commandLineKey(msg.String())
		// If asking to confirm quiting
	} else if m.confirmToQuit {
		quit := m.confirmToQuitSuperfile(msg.String())
		if quit {
			m.quitSuperfile()
			return m, tea.Quit
		}
		// If quiting input pressed, check if has any runing process and displays a
		// warn. Otherwise just quits application
	} else if msg.String() == containsKey(msg.String(), hotkeys.Quit) {
		if m.hasRunningProcesses() {
			m.warnModalForQuit()
			return m, cmd
		}

		m.quitSuperfile()
		return m, tea.Quit
	} else {
		// Handles general kinds of inputs in the regular state of the application
		cmd = m.mainKey(msg.String(), cmd)
	}
	return m, cmd
}

// Update the file panel state. Change name of renamed files, filter out files
// in search, update typingb bar, etc
func (m *model) updateFilePanelsState(msg tea.Msg, cmd *tea.Cmd) {
	focusPanel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	if m.firstTextInput {
		m.firstTextInput = false
	} else if m.fileModel.renaming {
		focusPanel.rename, *cmd = focusPanel.rename.Update(msg)
	} else if focusPanel.searchBar.Focused() {
		focusPanel.searchBar, *cmd = focusPanel.searchBar.Update(msg)
		for _, hotkey := range hotkeys.SearchBar {
			if hotkey == focusPanel.searchBar.Value() {
				focusPanel.searchBar.SetValue("")
				break
			}
		}
	} else if m.commandLine.input.Focused() {
		m.commandLine.input, *cmd = m.commandLine.input.Update(msg)
	} else if m.typingModal.open {
		m.typingModal.textInput, *cmd = m.typingModal.textInput.Update(msg)
	}

	if focusPanel.cursor < 0 {
		focusPanel.cursor = 0
	}
}

// Check if there's any processes running in background
func (m *model) hasRunningProcesses() bool {
	for _, data := range m.processBarModel.process {
		if data.state == inOperation && data.done != data.total {
			return true
		}
	}
	return false
}

// Triggers a warn for confirm quiting
func (m *model) warnModalForQuit() {
	m.confirmToQuit = true
	m.warnModal.title = "Confirm to quit superfile"
	m.warnModal.content = "You still have files being processed. Are you sure you want to exit?"
}

// Implement View function for bubble tea model to handle visualization.
func (m model) View() string {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	// check is the terminal size enough
	if m.fullHeight < minimumHeight || m.fullWidth < minimumWidth {
		return m.terminalSizeWarnRender()
	}
	if m.fileModel.width < 18 {
		return m.terminalSizeWarnAfterFirstRender()
	}
	sidebar := m.sidebarRender()

	filePanel := m.filePanelRender()

	filePreview := m.filePreviewPanelRender()

	mainPanel := lipgloss.JoinHorizontal(0, sidebar, filePanel, filePreview)

	var footer string

	if m.toggleFooter {
		processBar := m.processBarRender()

		metaData := m.metadataRender()

		clipboardBar := m.clipboardRender()

		footer = lipgloss.JoinHorizontal(0, processBar, metaData, clipboardBar)
	}

	if m.commandLine.input.Focused() {
		commandLine := m.commandLineInputBoxRender()
		footer = lipgloss.JoinVertical(0, footer, commandLine)

	}

	var finalRender string

	if m.toggleFooter {
		finalRender = lipgloss.JoinVertical(0, mainPanel, footer)
	} else {
		finalRender = mainPanel
	}
	// check if need pop up modal
	if m.helpMenu.open {
		helpMenu := m.helpMenuRender()
		overlayX := m.fullWidth/2 - m.helpMenu.width/2
		overlayY := m.fullHeight/2 - m.helpMenu.height/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, helpMenu, finalRender)
	}

	if panel.sortOptions.open {
		sortOptions := m.sortOptionsRender()
		overlayX := m.fullWidth/2 - panel.sortOptions.width/2
		overlayY := m.fullHeight/2 - panel.sortOptions.height/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, sortOptions, finalRender)
	}

	if firstUse {
		introduceModal := m.introduceModalRender()
		overlayX := m.fullWidth/2 - m.helpMenu.width/2
		overlayY := m.fullHeight/2 - m.helpMenu.height/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, introduceModal, finalRender)
	}

	if m.typingModal.open {
		typingModal := m.typineModalRender()
		overlayX := m.fullWidth/2 - modalWidth/2
		overlayY := m.fullHeight/2 - modalHeight/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, typingModal, finalRender)
	}

	if m.warnModal.open {
		warnModal := m.warnModalRender()
		overlayX := m.fullWidth/2 - modalWidth/2
		overlayY := m.fullHeight/2 - modalHeight/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, warnModal, finalRender)
	}

	if m.confirmToQuit {
		warnModal := m.warnModalRender()
		overlayX := m.fullWidth/2 - modalWidth/2
		overlayY := m.fullHeight/2 - modalHeight/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, warnModal, finalRender)
	}

	return finalRender
}

// Returns a tea.cmd responsible for listening messages from msg channel
func listenForChannelMessage(msg chan channelMessage) tea.Cmd {
	return func() tea.Msg {
		for {
			m := <-msg
			if m.messageType != sendProcess {
				ListeningMessage = false
				return m
			}
			if time.Since(progressBarLastRenderTime).Seconds() > 2 || m.processNewState.state == successful || m.processNewState.done < 2 {
				ListeningMessage = false
				progressBarLastRenderTime = time.Now()
				return m
			}
		}
	}
}

// Render and update file panel items. Check for changes and updates in files and
// folders in the current directory.
func (m *model) getFilePanelItems() {
	focusPanel := m.fileModel.filePanels[m.filePanelFocusIndex]
	for i, filePanel := range m.fileModel.filePanels {
		var fileElement []element
		nowTime := time.Now()
		// Check last time each element was updated, if less then 3 seconds ignore
		if filePanel.focusType == noneFocus && nowTime.Sub(filePanel.lastTimeGetElement) < 3*time.Second {
      if !m.updatedToggleDotFile {
        continue
      }
		}

		focusPanelReRender := false

		if len(focusPanel.element) > 0 {
			if filepath.Dir(focusPanel.element[0].location) != focusPanel.location {
				focusPanelReRender = true
			}
		} else {
			focusPanelReRender = true
		}

		reRenderTime := int(float64(len(filePanel.element)) / 100)

		if filePanel.focusType != noneFocus && nowTime.Sub(filePanel.lastTimeGetElement) < time.Duration(reRenderTime)*time.Second && !focusPanelReRender {
			continue
		}

		// Get file names based on search bar filter
		if filePanel.searchBar.Value() != "" {
			fileElement = returnFolderElementBySearchString(filePanel.location, m.toggleDotFile, filePanel.searchBar.Value())
		} else {
			fileElement = returnFolderElement(filePanel.location, m.toggleDotFile, filePanel.sortOptions.data)
		}
		// Update file panel list
		filePanel.element = fileElement
		m.fileModel.filePanels[i].element = fileElement
		m.fileModel.filePanels[i].lastTimeGetElement = nowTime
	}

  m.updatedToggleDotFile = false
}

// Close superfile application. Cd into the curent dir if CdOnQuit on and save
// the path in state direcotory
func (m model) quitSuperfile() {
    // close exiftool session
    if Config.Metadata {
        et.Close();
    }
    // cd on quit
    currentDir := m.fileModel.filePanels[m.filePanelFocusIndex].location
    variable.LastDir = currentDir

    if Config.CdOnQuit {
        // escape single quote
        currentDir = strings.ReplaceAll(currentDir, "'", "'\\''")
        os.WriteFile(variable.SuperFileStateDir+"/lastdir", []byte("cd '"+currentDir+"'"), 0755)
    }
}
