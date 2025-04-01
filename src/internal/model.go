package internal

import (
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/common/utils"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
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
var batCmd = ""

var et *exiftool.Exiftool

var channel = make(chan channelMessage, 1000)
var progressBarLastRenderTime time.Time = time.Now()

// Initialize and return model with default configs
func InitialModel(dir string, firstUseCheck, hasTrashCheck bool) model {
	toggleDotFileBool, toggleFooter, firstFilePanelDir := initialConfig(dir)
	firstUse = firstUseCheck
	hasTrash = hasTrashCheck
	batCmd = checkBatCmd()
	return defaultModelConfig(toggleDotFileBool, toggleFooter, firstFilePanelDir)
}

// Init function to be called by Bubble tea framework, sets windows title,
// cursos blinking and starts message streamming channel
func (m model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("superfile"),
		textinput.Blink, // Assuming textinput.Blink is a valid command
		listenForChannelMessage(channel),
	)
}

// Update function for bubble tea to provide internal communication to the
// application
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("model.Update() called")
	var cmd tea.Cmd

	m.updateSidebarState(msg, &cmd)

	switch msg := msg.(type) {
	case channelMessage:
		m.handleChannelMessage(msg)
	case tea.WindowSizeMsg:
		m.handleWindowResize(msg)
	case tea.MouseMsg:
		msgStr := msg.String()
		if msgStr == "wheel up" || msgStr == "wheel down" {
			wheelMainAction(msgStr, &m)
		} else {
			slog.Debug("Mouse event of type that is not handled", "msg", msgStr)
		}
	case tea.KeyMsg:
		cmd = m.handleKeyInput(msg, cmd)
	default:
		slog.Debug("Message of type that is not handled", "type", reflect.TypeOf(msg))
	}

	m.updateFilePanelsState(msg, &cmd)

	if m.sidebarModel.searchBar.Value() != "" {
		// Todo : All updates of sideBar must be moved to seperate struct functions
		// we have to keep the state of sidebar consistent, and keep values of
		// cursor, directories, renderIndex sane for each update, and it has to
		// take care at one single place, not everywhere we use sideBar
		m.sidebarModel.directories = getFilteredDirectories(m.sidebarModel.searchBar.Value())
		if m.sidebarModel.isCursorInvalid() {
			m.sidebarModel.resetCursor()
		}
	} else {
		m.sidebarModel.directories = getDirectories()
		if m.sidebarModel.isCursorInvalid() {
			m.sidebarModel.resetCursor()
		}
	}

	// check if there already have listening message
	if !ListeningMessage {
		cmd = tea.Batch(cmd, listenForChannelMessage(channel))
	}

	m.getFilePanelItems()
	if !m.firstLoadingComplete {
		m.firstLoadingComplete = true
	}
	return m, tea.Batch(cmd)
}

// Handle message exchanging within the application
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
	m.setHeightValues(msg.Height)
	m.setHelpMenuSize()

	if m.fileModel.maxFilePanel >= 10 {
		m.fileModel.maxFilePanel = 10
	}
}

// Set file preview panel Widht to width. Assure that
func (m *model) setFilePreviewWidth(width int) {
	if common.Config.FilePreviewWidth == 0 {
		m.fileModel.filePreview.width = (width - common.Config.SidebarWidth - (4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
	} else if common.Config.FilePreviewWidth > 10 || common.Config.FilePreviewWidth == 1 {
		utils.LogAndExit("Config file file_preview_width invalidation")
	} else {
		m.fileModel.filePreview.width = (width - common.Config.SidebarWidth) / common.Config.FilePreviewWidth
	}
}

// Proper set panels size. Assure that panels do not overlap
func (m *model) setFilePanelsSize(width int) {
	// set each file panel size and max file panel amount
	m.fileModel.width = (width - common.Config.SidebarWidth - m.fileModel.filePreview.width - (4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.fileModel.maxFilePanel = (width - common.Config.SidebarWidth - m.fileModel.filePreview.width) / 20
	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
}

func (m *model) setHeightValues(height int) {
	if !m.toggleFooter {
		m.footerHeight = 0
	} else if height < 30 {
		m.footerHeight = 6
	} else if height < 35 {
		m.footerHeight = 7
	} else if height < 40 {
		m.footerHeight = 8
	} else if height < 45 {
		m.footerHeight = 9
	} else {
		m.footerHeight = 10
	}
	// Todo : Make it grow even more for bigger screen sizes.
	// Todo : Calculate the value , instead of manually hard coding it.

	// Main panel height = Total terminal height - footer height - 2(footer border) - 2(file panel border)
	m.mainPanelHeight = height - m.footerHeight - 2 - 2
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
func (m *model) handleKeyInput(msg tea.KeyMsg, cmd tea.Cmd) tea.Cmd {

	slog.Debug("model.handleKeyInput", "msg", msg, "typestr", msg.Type.String(),
		"runes", msg.Runes, "type", int(msg.Type), "paste", msg.Paste,
		"alt", msg.Alt)
	slog.Debug("model.handleKeyInput. model info. ",
		"filePanelFocusIndex", m.filePanelFocusIndex,
		"filePanel.focusType", m.fileModel.filePanels[m.filePanelFocusIndex].focusType,
		"filePanel.panelMode", m.fileModel.filePanels[m.filePanelFocusIndex].panelMode,
		"typingModal.open", m.typingModal.open,
		"warnModal.open", m.warnModal.open,
		"promptModal.open", m.promptModal.IsOpen(),
		"fileModel.renaming", m.fileModel.renaming,
		"searchBar.focussed", m.fileModel.filePanels[m.filePanelFocusIndex].searchBar.Focused(),
		"helpMenu.open", m.helpMenu.open,
		"firstTextInput", m.firstTextInput,
		"focusPanel", m.focusPanel,
	)

	if firstUse {
		firstUse = false
		return cmd
	}

	if m.typingModal.open {
		m.typingModalOpenKey(msg.String())

	} else if m.promptModal.IsOpen() {
		// Ignore keypress. It will be handled in Update call via
		// updateFilePanelState

	} else if m.warnModal.open {
		m.warnModalOpenKey(msg.String())
		// If renaming a object
	} else if m.fileModel.renaming {
		m.renamingKey(msg.String())
	} else if m.sidebarModel.renaming {
		m.sidebarRenamingKey(msg.String())
		// If search bar is open
	} else if m.fileModel.filePanels[m.filePanelFocusIndex].searchBar.Focused() {
		m.focusOnSearchbarKey(msg.String())
		// If sort options menu is open
	} else if m.sidebarModel.searchBar.Focused() {
		m.sidebarSearchBarKey(msg.String())
		// If sort options menu is open
	} else if m.fileModel.filePanels[m.filePanelFocusIndex].sortOptions.open {
		m.sortOptionsKey(msg.String())
		// If help menu is open
	} else if m.helpMenu.open {
		m.helpMenuKey(msg.String())
		// If asking to confirm quiting
	} else if m.confirmToQuit {
		quit := m.confirmToQuitSuperfile(msg.String())
		if quit {
			m.quitSuperfile()
			return tea.Quit
		}
		// If quiting input pressed, check if has any running process and displays a
		// warn. Otherwise just quits application
	} else if msg.String() == containsKey(msg.String(), common.Hotkeys.Quit) {
		if m.hasRunningProcesses() {
			m.warnModalForQuit()
			return cmd
		}

		m.quitSuperfile()
		return tea.Quit
	} else {
		// Handles general kinds of inputs in the regular state of the application
		cmd = m.mainKey(msg.String(), cmd)
	}
	return cmd
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
	} else if m.typingModal.open {
		m.typingModal.textInput, *cmd = m.typingModal.textInput.Update(msg)
	} else if m.promptModal.IsOpen() {
		// *cmd is a non-name, and cannot be used on left of :=
		var action common.PromptAction
		// Taking returned cmd is necessary for blinking
		action, *cmd = m.promptModal.HandleMessage(msg)
		m.applyPromptModalAction(action)
	}

	// Todo : This is like duct taping a bigger problem
	// The code should never reach this state.
	if focusPanel.cursor < 0 {
		focusPanel.cursor = 0
	}
}

func (m *model) applyPromptModalAction(action common.PromptAction) {
	switch action.Action {
	case common.NoAction:
		return
	case common.ShellCommandAction:
		if len(action.Args) != 1 {
			slog.Error("Invalid ShellCommandAction without exactly one arg",
				"args", action.Args)
			return
		}
		m.applyShellCommandAction(action.Args[0])
	case common.SplitPanelAction:
		if len(action.Args) != 0 {
			slog.Warn("Invalid SplitPanelAction with extra args. Ignoring.",
				"args", action.Args)
		}
		slog.Debug("SplitPanelAction")
		m.splitPanel()
	}
}

// Todo : Move them around to appropriate places
func (m *model) applyShellCommandAction(shellCommand string) {
	focusPanelDir := ""
	for _, panel := range m.fileModel.filePanels {
		if panel.focusType == focus {
			focusPanelDir = panel.location
		}
	}

	// Linux and Darwin
	baseCmd := "/bin/sh"
	args := []string{"-c", shellCommand}

	if runtime.GOOS == "windows" {
		baseCmd = "powershell.exe"
		args[0] = "-Command"
	}

	retCode, output, err := utils.ExecuteShellCommand(common.DefaultCommandTimeoutMsec, focusPanelDir,
		baseCmd, args...)

	if err != nil {
		slog.Error("Command execution failed", "retCode", retCode,
			"error", err, "output", string(output))
		return
	}
}

func (m *model) splitPanel() {
	m.createNewFilePanel(m.fileModel.filePanels[m.filePanelFocusIndex].location)
}

// Update the sidebar state. Change name of the renaming pinned directory.
func (m *model) updateSidebarState(msg tea.Msg, cmd *tea.Cmd) {
	sidebar := &m.sidebarModel
	if sidebar.renaming {
		sidebar.rename, *cmd = sidebar.rename.Update(msg)
	} else if sidebar.searchBar.Focused() {
		sidebar.searchBar, *cmd = sidebar.searchBar.Update(msg)
	}

	if sidebar.cursor < 0 {
		sidebar.cursor = 0
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
	slog.Debug("model.View() called", "mainPanelHeight", m.mainPanelHeight,
		"footerHeight", m.footerHeight, "fullHeight", m.fullHeight,
		"fullWidth", m.fullWidth)

	if !m.firstLoadingComplete {
		return "Loading..."
	}
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	// check is the terminal size enough
	if m.fullHeight < common.MinimumHeight || m.fullWidth < common.MinimumWidth {
		return m.terminalSizeWarnRender()
	}
	if m.fileModel.width < 18 {
		return m.terminalSizeWarnAfterFirstRender()
	}

	if err := m.validateLayout(); err != nil {
		slog.Error("Invalid layout", "error", err)
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

	if m.promptModal.IsOpen() {
		promptModal := m.promptModalRender()
		overlayX := m.fullWidth/2 - m.helpMenu.width/2
		overlayY := m.fullHeight/2 - m.helpMenu.height/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, promptModal, finalRender)
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
		overlayX := m.fullWidth/2 - common.ModalWidth/2
		overlayY := m.fullHeight/2 - common.ModalHeight/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, typingModal, finalRender)
	}

	if m.warnModal.open {
		warnModal := m.warnModalRender()
		overlayX := m.fullWidth/2 - common.ModalWidth/2
		overlayY := m.fullHeight/2 - common.ModalHeight/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, warnModal, finalRender)
	}

	if m.confirmToQuit {
		warnModal := m.warnModalRender()
		overlayX := m.fullWidth/2 - common.ModalWidth/2
		overlayY := m.fullHeight/2 - common.ModalHeight/2
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
			fileElement = returnDirElementBySearchString(filePanel.location, m.toggleDotFile, filePanel.searchBar.Value())
		} else {
			fileElement = returnDirElement(filePanel.location, m.toggleDotFile, filePanel.sortOptions.data)
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
func (m *model) quitSuperfile() {
	// close exiftool session
	if common.Config.Metadata && et != nil {
		et.Close()
	}
	// cd on quit
	currentDir := m.fileModel.filePanels[m.filePanelFocusIndex].location
	variable.LastDir = currentDir

	if common.Config.CdOnQuit {
		// escape single quote
		currentDir = strings.ReplaceAll(currentDir, "'", "'\\''")
		err := os.WriteFile(variable.LastDirFile, []byte("cd '"+currentDir+"'"), 0755)
		if err != nil {
			slog.Error("Error during writing lastdir file", "error", err)
		}
	}
	slog.Debug("Quitting superfile", "current dir", currentDir)
}

// Check if bat is an executable in PATH and whether to use bat or batcat as command
func checkBatCmd() string {
	if _, err := exec.LookPath("bat"); err == nil {
		return "bat"
	}
	// on ubuntu bat executable is called batcat
	if _, err := exec.LookPath("batcat"); err == nil {
		return "batcat"
	}
	return ""
}
