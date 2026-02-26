package internal

import (
	"errors"
	"log/slog"
	"os"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/preview"

	variable "github.com/yorukot/superfile/src/config"
	zoxideui "github.com/yorukot/superfile/src/internal/ui/zoxide"
	stringfunction "github.com/yorukot/superfile/src/pkg/string_function"
)

// These represent model's state information, its not a global preperty
var (
	LastTimeCursorMove = [2]int{int(time.Now().UnixMicro()), 0} //nolint: gochecknoglobals // TODO: Move to model struct
	et                 *exiftool.Exiftool                       //nolint: gochecknoglobals // TODO: Move to model struct
)

// Initialize and return model with default configs
// It returns only tea.Model because when it used in main, the return value
// is passed to tea.NewProgram() which accepts tea.Model
// Either way type 'model' is not exported, so there is not way main package can
// be aware of it, and use it directly
func InitialModel(firstPanelPaths []string, firstUseCheck bool) tea.Model {
	toggleDotFile, toggleFooter, zClient := initialConfig(firstPanelPaths)
	return defaultModelConfig(toggleDotFile, toggleFooter, firstUseCheck, firstPanelPaths, zClient)
}

// Init function to be called by Bubble tea framework, sets windows title,
// cursos blinking and starts message streamming channel
// Note : What init should do, for example read file panel data, read sidebar directories, and
// disk, is being done in at the creation of model of object. Right now creation of model object
// and its initialization isn't well separated.
func (m *model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("superfile"),
		textinput.Blink, // Assuming textinput.Blink is a valid command
		processCmdToTeaCmd(m.processBarModel.GetListenCmd()),
	)
}

// Update function for bubble tea to provide internal communication to the
// application
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Debug("model.Update() called", "msgType", reflect.TypeOf(msg))

	var sidebarCmd, inputCmd, updateCmd, panelCmd,
		metadataCmd, filePreviewCmd, helpMenuCmd, resizeCmd tea.Cmd

	// These are above the key message handing to prevent issues with firstKeyInput
	// if someone presses `/` to focus to searchBar, searchBar will otherwise
	// get `/` input too.
	sidebarCmd = m.sidebarModel.UpdateState(msg)
	// Necessary for blinking. Can't do this in HandleKey, as we only pass KeyMsg there
	helpMenuCmd = m.helpMenu.HandleTeaMsg(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		resizeCmd = m.handleWindowResize(msg)
	case tea.MouseMsg:
		m.handleMouseMsg(msg)
	case tea.KeyMsg:
		inputCmd = m.handleKeyInput(msg)

	// Has to handle zoxide messages separately as they could be generated via
	// zoxide update commands, or batched commands from textinput
	// Cannot do it like processbar messages
	case zoxideui.UpdateMsg:
		slog.Debug("Got ModelUpdate message", "id", msg.GetReqID())
		updateCmd = msg.Apply(&m.zoxideModal)

	// Its a pain to interconvert commands like processBar
	case preview.UpdateMsg:
		slog.Debug("Got ModelUpdate message", "id", msg.GetReqID())
		m.fileModel.UpdatePreviewPanel(msg)
	case ModelUpdateMessage:
		slog.Debug("Got ModelUpdate message", "id", msg.GetReqID())
		updateCmd = msg.ApplyToModel(m)

	default:
		slog.Debug("Message of type that is not explicitly handled")
	}

	// This is needed for blink, etc to work
	panelCmd = m.updateComponentState(msg)

	m.updateModelStateAfterMsg()
	filePreviewCmd = m.fileModel.GetFilePreviewCmd(false)

	metadataCmd = m.getMetadataCmd()

	return m, tea.Batch(sidebarCmd, helpMenuCmd, inputCmd, updateCmd,
		panelCmd, metadataCmd, filePreviewCmd, resizeCmd)
}

func (m *model) handleMouseMsg(msg tea.MouseMsg) {
	msgStr := msg.String()
	if msgStr == "wheel up" || msgStr == "wheel down" {
		wheelMainAction(msgStr, m)
	} else {
		slog.Debug("Mouse event of type that is not handled", "msg", msgStr)
	}
}

func (m *model) updateModelStateAfterMsg() {
	m.sidebarModel.UpdateDirectories()
	m.fileModel.UpdateFilePanelsIfNeeded(false)
	// TODO: Move to utility
	if m.focusPanel != metadataFocus {
		m.fileMetaData.ResetRender()
	}
	// TODO: Entirely remove the need of this variable, and handle first loading via Init()
	// Init() should return a basic model object with all IO waiting via a tea.Cmd
	if !m.firstLoadingComplete {
		m.firstLoadingComplete = true
	}
}

// Note : Maybe we should not trigger metadata fetch for updates
// that dont change the currently selected file panel element
// TODO : At least dont trigger metadata fetch when user is scrolling
// through the metadata panel
func (m *model) getMetadataCmd() tea.Cmd {
	if m.disableMetadata {
		return nil
	}
	if m.getFocusedFilePanel().EmptyOrInvalid() {
		m.fileMetaData.SetBlank()
		return nil
	}
	selectedItem := m.getFocusedFilePanel().GetFocusedItem()
	metadataFocused := m.focusPanel == metadataFocus
	// Note : This will cause metadata not being refreshed there is any file update events.
	// We can have a cache with TTL or watch filesystem changes to fix this
	if selectedItem.Location == m.fileMetaData.GetMetadataLocation() &&
		metadataFocused == m.fileMetaData.GetMetadataExpectedFocused() {
		return nil
	}
	if m.fileMetaData.UpdateMetadataIfExistsInCache(selectedItem.Location, metadataFocused) {
		return nil
	}

	m.fileMetaData.SetMetadataLocationAndFocused(selectedItem.Location, metadataFocused)

	if m.fileMetaData.IsBlank() {
		m.fileMetaData.SetInfoMsg(icon.InOperation + icon.Space + "Loading metadata...")
	}

	reqCnt := m.ioReqCnt
	m.ioReqCnt++
	// If there are too many metadata fetches, we need to have a cache with path as a key
	// and timeout based eviction
	slog.Debug("Submitting metadata fetch request", "id", reqCnt, "path", selectedItem.Location)
	return func() tea.Msg {
		return NewMetadataMsg(
			metadata.GetMetadata(selectedItem.Location, metadataFocused, et), metadataFocused, reqCnt)
	}
}

// Adjust window size based on msg information
func (m *model) handleWindowResize(msg tea.WindowSizeMsg) tea.Cmd {
	m.fullHeight = msg.Height
	m.fullWidth = msg.Width
	m.setHeightValues()
	return m.updateComponentDimensions()
}

func (m *model) setHeightValues() {
	//nolint: gocritic // This is to be separated out to a function, and made better later. No need to refactor here
	if !m.toggleFooter {
		m.footerHeight = 0
	} else if m.fullHeight < common.HeightBreakA {
		m.footerHeight = 6
	} else if m.fullHeight < common.HeightBreakB {
		m.footerHeight = 7
	} else if m.fullHeight < common.HeightBreakC {
		m.footerHeight = 8
	} else if m.fullHeight < common.HeightBreakD {
		m.footerHeight = 9
	} else {
		m.footerHeight = 10
	}
	// TODO : Make it grow even more for bigger screen sizes.
	// TODO : Calculate the value , instead of manually hard coding it.

	// Main panel height = Total terminal height- 2(file panel border) - footer height
	m.mainPanelHeight = m.fullHeight - common.BorderPadding - utils.FullFooterHeight(m.footerHeight, m.toggleFooter)
}

func (m *model) updateComponentDimensions() tea.Cmd {
	m.setHelpMenuSize()
	m.setPromptModelSize()
	m.setZoxideModelSize()
	m.setFooterComponentSize()

	// File preview panel requires explicit height update, unlike sidebar/file panels
	// which receive height as render parameters and update automatically on each frame
	// Force re-render of preview content with new dimensions
	return m.setMainModelDimensions()
}

func (m *model) setMainModelDimensions() tea.Cmd {
	fileModelWidth := m.fullWidth
	if common.Config.SidebarWidth != 0 {
		fileModelWidth -= common.Config.SidebarWidth + common.BorderPadding
	}
	m.sidebarModel.SetHeight(m.mainPanelHeight + common.BorderPadding)
	return m.fileModel.SetDimensions(fileModelWidth, m.mainPanelHeight+common.BorderPadding)
}

// Set help menu size
func (m *model) setHelpMenuSize() {
	height := m.fullHeight - common.BorderPadding
	width := m.fullWidth - common.BorderPadding
	if m.fullHeight > common.HeightBreakB {
		height = 30
	}
	if m.fullWidth > common.ResponsiveWidthThreshold {
		width = 90
	}
	m.helpMenu.SetDimensions(width, height)
}

func (m *model) setPromptModelSize() {
	// Scale prompt model's maxHeight - 33% of total height
	m.promptModal.SetMaxHeight(m.fullHeight / 3) //nolint:mnd // modal uses third height for layout

	// Scale prompt model's maxHeight - 50% of total height
	m.promptModal.SetWidth(m.fullWidth / 2) //nolint:mnd // modal uses half width for layout
}

func (m *model) setZoxideModelSize() {
	// Scale zoxide model's maxHeight - 50% of total height to accommodate scroll indicators
	m.zoxideModal.SetMaxHeight(m.fullHeight / 2) //nolint:mnd // modal uses half height for layout

	// Scale zoxide model's width - 50% of total width
	m.zoxideModal.SetWidth(m.fullWidth / 2) //nolint:mnd // modal uses half width for layout
}

func (m *model) setFooterComponentSize() {
	var width, clipBoardwidth, height int
	height = m.footerHeight + common.BorderPadding
	width = m.fullWidth / utils.CntFooterPanels
	clipBoardwidth = width + m.fullWidth%utils.CntFooterPanels
	m.fileMetaData.SetDimensions(width, height)
	m.processBarModel.SetDimensions(width, height)
	m.clipboard.SetDimensions(clipBoardwidth, height)
}

// Identify the current state of the application m and properly handle the
// msg keybind pressed
func (m *model) handleKeyInput(msg tea.KeyMsg) tea.Cmd {
	slog.Debug("model.handleKeyInput", "msg", msg, "typestr", msg.Type.String(),
		"runes", msg.Runes, "type", int(msg.Type), "paste", msg.Paste,
		"alt", msg.Alt)
	slog.Debug("model.handleKeyInput. model info. ",
		"fileModel.FocusedPanelIndex", m.fileModel.FocusedPanelIndex,
		"filePanel.isFocused", m.getFocusedFilePanel().IsFocused,
		"filePanel.panelMode", m.getFocusedFilePanel().PanelMode,
		"typingModal.open", m.typingModal.open,
		"notifyModel.open", m.notifyModel.IsOpen(),
		"promptModal.open", m.promptModal.IsOpen(),
		"fileModel.renaming", m.fileModel.Renaming,
		"searchBar.focused", m.getFocusedFilePanel().SearchBar.Focused(),
		"helpMenu.open", m.helpMenu.IsOpen(),
		"firstTextInput", m.firstTextInput,
		"focusPanel", m.focusPanel,
	)
	if m.firstUse {
		m.firstUse = false
		return nil
	}
	var cmd tea.Cmd
	cdOnQuit := common.Config.CdOnQuit
	switch {
	case m.typingModal.open:
		m.typingModalOpenKey(msg.String())
	case m.promptModal.IsOpen():
		// Ignore keypress. It will be handled in Update call via
		// updateFilePanelState
		// TODO: Convert that to async via tea.Cmd
	case m.zoxideModal.IsOpen():
		// Ignore keypress. It will be handled in Update call via
		// updateFilePanelState

	// Handles all warn models except the warn model for confirming to quit
	case m.notifyModel.IsOpen():
		cmd = m.notifyModelOpenKey(msg.String())

	// If renaming a object
	case m.fileModel.Renaming:
		cmd = m.renamingKey(msg.String())
	case m.sidebarModel.IsRenaming():
		m.sidebarRenamingKey(msg.String())
	// If search bar is open
	case m.getFocusedFilePanel().SearchBar.Focused():
		m.focusOnSearchbarKey(msg.String())
	// If sort options menu is open
	case m.sidebarModel.SearchBarFocused():
		m.sidebarModel.HandleSearchBarKey(msg.String())
	case m.sortModal.IsOpen():
		m.sortOptionsKey(msg.String())
	// If help menu is open
	case m.helpMenu.IsOpen():
		m.helpMenu.HandleKey(msg.String())

	case slices.Contains(common.Hotkeys.Quit, msg.String()):
		m.modelQuitState = quitInitiated

	case slices.Contains(common.Hotkeys.CdQuit, msg.String()):
		m.modelQuitState = quitInitiated
		cdOnQuit = true

	default:
		// Handles general kinds of inputs in the regular state of the application
		cmd = m.mainKey(msg.String())
	}

	// If quiting input pressed, check if has any running process and displays a
	// warn. Otherwise just quits application
	if m.modelQuitState == quitInitiated {
		if m.processBarModel.HasRunningProcesses() {
			// Dont quit now, get a confirmation first.
			m.modelQuitState = quitConfirmationInitiated
			m.warnModalForQuit()
			return cmd
		}
		m.modelQuitState = quitConfirmationReceived
	}
	if m.modelQuitState == quitConfirmationReceived {
		m.quitSuperfile(cdOnQuit)
		return tea.Quit
	}
	return cmd
}

// Update the file panel state. Change name of renamed files, filter out files
// in search, update typingb bar, etc
func (m *model) updateComponentState(msg tea.Msg) tea.Cmd {
	focusPanel := m.getFocusedFilePanel()
	var cmd tea.Cmd
	var action common.ModelAction
	switch {
	case m.firstTextInput:
		m.firstTextInput = false
	case m.fileModel.Renaming:
		focusPanel.Rename, cmd = focusPanel.Rename.Update(msg)
	case focusPanel.SearchBar.Focused():
		focusPanel.SearchBar, cmd = focusPanel.SearchBar.Update(msg)
	case m.typingModal.open:
		m.typingModal.textInput, cmd = m.typingModal.textInput.Update(msg)
	case m.promptModal.IsOpen():
		// TODO : Separate this to a utility
		cwdLocation := m.getFocusedFilePanel().Location
		action, cmd = m.promptModal.HandleUpdate(msg, cwdLocation)
		cmd = tea.Batch(cmd, m.applyPromptModalAction(action))
	case m.zoxideModal.IsOpen():
		action, cmd = m.zoxideModal.HandleUpdate(msg)
		cmd = tea.Batch(cmd, m.applyZoxideModalAction(action))
	}
	return cmd
}

// Apply the Action and notify the promptModal
func (m *model) applyPromptModalAction(action common.ModelAction) tea.Cmd {
	successMsg, cmd, actionErr := m.logAndExecuteAction(action)
	if actionErr != nil {
		m.promptModal.HandleSPFActionResults(false, actionErr.Error())
	} else if successMsg != "" {
		m.promptModal.HandleSPFActionResults(true, successMsg)
	}
	return cmd
}

// Utility function to log and execute actions, reducing duplication
func (m *model) logAndExecuteAction(action common.ModelAction) (string, tea.Cmd, error) {
	// Only log actions that aren't NoAction to reduce debug noise
	if _, ok := action.(common.NoAction); !ok {
		slog.Debug("Applying model action", "action", action)
	}

	switch action := action.(type) {
	case common.NoAction:
		return "", nil, nil
	case common.ShellCommandAction:
		// Shell commands are handled separately and don't return here
		m.applyShellCommandAction(action.Command)
		return "", nil, nil
	case common.SplitPanelAction:
		cmd, err := m.splitPanel()
		return "Panel successfully split", cmd, err
	case common.CDCurrentPanelAction:
		return "Panel directory changed", nil, m.updateCurrentFilePanelDir(action.Location)
	case common.OpenPanelAction:
		cmd, err := m.createNewFilePanelRelativeToCurrent(action.Location)
		return "New panel opened", cmd, err
	default:
		return "", nil, errors.New("unhandled action type")
	}
}

// Apply the Action for zoxide modal (no result notifications needed)
func (m *model) applyZoxideModalAction(action common.ModelAction) tea.Cmd {
	_, cmd, _ := m.logAndExecuteAction(action)
	return cmd
}

// TODO : Move them around to appropriate places
func (m *model) applyShellCommandAction(shellCommand string) {
	focusPanelDir := m.getFocusedFilePanel().Location

	retCode, output, err := utils.ExecuteCommandInShell(common.DefaultCommandTimeout, focusPanelDir, shellCommand)

	m.promptModal.HandleShellCommandResults(retCode, output)

	if err != nil {
		slog.Error("Command execution failed", "retCode", retCode,
			"error", err, "output", output)
		return
	}
}

func (m *model) splitPanel() (tea.Cmd, error) {
	return m.fileModel.CreateNewFilePanel(m.getFocusedFilePanel().Location)
}

func (m *model) createNewFilePanelRelativeToCurrent(path string) (tea.Cmd, error) {
	currentDir := m.getFocusedFilePanel().Location
	return m.fileModel.CreateNewFilePanel(utils.ResolveAbsPath(currentDir, path))
}

// simulates a 'cd' action
func (m *model) updateCurrentFilePanelDir(path string) error {
	panel := m.getFocusedFilePanel()
	err := panel.UpdateCurrentFilePanelDir(path)
	if err == nil {
		// Track the directory change with zoxide
		m.trackDirectoryWithZoxide(panel.Location)
	}
	return err
}

// trackDirectoryWithZoxide adds the directory to zoxide database if zoxide is available and enabled
func (m *model) trackDirectoryWithZoxide(path string) {
	if !common.Config.ZoxideSupport || m.zClient == nil {
		return
	}

	err := m.zClient.Add(path)
	if err != nil {
		slog.Debug("Failed to add directory to zoxide", "path", path, "error", err)
	}
}

// Check if there's any processes running in background

// Triggers a warn for confirm quiting
func (m *model) warnModalForQuit() {
	m.notifyModel = notify.New(true, "Confirm to quit superfile",
		"You still have files being processed. Are you sure you want to exit?",
		notify.QuitAction)
}

// Implement View function for bubble tea model to handle visualization.
func (m *model) View() string {
	slog.Debug("model.View() called", "mainPanelHeight", m.mainPanelHeight,
		"footerHeight", m.footerHeight, "fullHeight", m.fullHeight,
		"fullWidth", m.fullWidth, "panelCount", m.fileModel.PanelCount(),
		"singlePanelWidth", m.fileModel.SinglePanelWidth,
		"maxPanels", m.fileModel.MaxFilePanel,
		"sideBarWidth", common.Config.SidebarWidth,
		"firstFilePanelWidth", m.fileModel.FilePanels[0].GetWidth())

	if !m.firstLoadingComplete {
		return "Loading..."
	}

	// check is the terminal size enough
	if m.fullHeight < common.MinimumHeight || m.fullWidth < common.MinimumWidth {
		return m.terminalSizeWarnRender()
	}
	if m.fileModel.SinglePanelWidth < filepanel.MinWidth {
		return m.terminalSizeWarnAfterFirstRender()
	}

	// Do validations after min size check above. Validations will fail if user give
	// too less size to the terminal program
	if err := m.validateLayout(); err != nil {
		slog.Error("Invalid layout", "error", err)
	}

	return m.updateRenderForOverlay(m.mainComponentsRender())
}

func (m *model) updateRenderForOverlay(finalRender string) string {
	// check if need pop up modal
	if m.helpMenu.IsOpen() {
		helpMenu := m.helpMenu.Render()
		overlayX := m.fullWidth/common.CenterDivisor - m.helpMenu.GetWidth()/common.CenterDivisor
		overlayY := m.fullHeight/common.CenterDivisor - m.helpMenu.GetHeight()/common.CenterDivisor
		return stringfunction.PlaceOverlay(overlayX, overlayY, helpMenu, finalRender)
	}

	if m.promptModal.IsOpen() {
		promptModal := m.promptModalRender()
		overlayX := m.fullWidth/common.CenterDivisor - m.promptModal.GetWidth()/common.CenterDivisor
		overlayY := m.fullHeight/common.CenterDivisor - m.promptModal.GetMaxHeight()/common.CenterDivisor
		return stringfunction.PlaceOverlay(overlayX, overlayY, promptModal, finalRender)
	}

	if m.zoxideModal.IsOpen() {
		zoxideModal := m.zoxideModalRender()
		overlayX := m.fullWidth/common.CenterDivisor - m.zoxideModal.GetWidth()/common.CenterDivisor
		overlayY := m.fullHeight/common.CenterDivisor - m.zoxideModal.GetMaxHeight()/common.CenterDivisor
		return stringfunction.PlaceOverlay(overlayX, overlayY, zoxideModal, finalRender)
	}

	if m.sortModal.IsOpen() {
		sortOptions := m.sortModal.Render()
		overlayX := m.fullWidth/common.CenterDivisor - m.sortModal.Width/common.CenterDivisor
		overlayY := m.fullHeight/common.CenterDivisor - m.sortModal.Height/common.CenterDivisor
		return stringfunction.PlaceOverlay(overlayX, overlayY, sortOptions, finalRender)
	}

	if m.firstUse {
		introduceModal := m.introduceModalRender()
		overlayX := m.fullWidth/common.CenterDivisor - m.helpMenu.GetWidth()/common.CenterDivisor
		overlayY := m.fullHeight/common.CenterDivisor - m.helpMenu.GetHeight()/common.CenterDivisor
		return stringfunction.PlaceOverlay(overlayX, overlayY, introduceModal, finalRender)
	}

	if m.typingModal.open {
		typingModal := m.typineModalRender()
		overlayX := m.fullWidth/common.CenterDivisor - common.ModalWidth/common.CenterDivisor
		overlayY := m.fullHeight/common.CenterDivisor - common.ModalHeight/common.CenterDivisor
		return stringfunction.PlaceOverlay(overlayX, overlayY, typingModal, finalRender)
	}

	if m.notifyModel.IsOpen() {
		notifyModal := m.notifyModel.Render()
		overlayX := m.fullWidth/common.CenterDivisor - common.ModalWidth/common.CenterDivisor
		overlayY := m.fullHeight/common.CenterDivisor - common.ModalHeight/common.CenterDivisor
		return stringfunction.PlaceOverlay(overlayX, overlayY, notifyModal, finalRender)
	}
	return finalRender
}

func (m *model) mainComponentsRender() string {
	sidebar := m.sidebarRender()
	fileModel := m.fileModel.Render()
	mainPanel := lipgloss.JoinHorizontal(0, sidebar, fileModel)

	if !m.toggleFooter {
		return mainPanel
	}

	processBar := m.processBarRender()
	metaData := m.fileMetaData.Render(m.focusPanel == metadataFocus)
	clipboardBar := m.clipboard.Render()
	footer := lipgloss.JoinHorizontal(0, processBar, metaData, clipboardBar)
	return lipgloss.JoinVertical(0, mainPanel, footer)
}

// Close superfile application. Cd into the current dir if CdOnQuit on and save
// the path in state direcotory
func (m *model) quitSuperfile(cdOnQuit bool) {
	// Resource cleanup
	if common.Config.Metadata && et != nil {
		_ = et.Close()
	}
	m.fileModel.FilePreview.CleanUp()

	// cd on quit
	currentDir := m.getFocusedFilePanel().Location
	variable.SetLastDir(currentDir)

	if cdOnQuit {
		// escape single quote
		currentDir = strings.ReplaceAll(currentDir, "'", "'\\''")
		err := os.WriteFile(variable.LastDirFile, []byte("cd '"+currentDir+"'"), utils.ConfigFilePerm)
		if err != nil {
			slog.Error("Error during writing lastdir file", "error", err)
		}
	}
	m.modelQuitState = quitDone
	slog.Debug("Quitting superfile", "current dir", currentDir)
}
