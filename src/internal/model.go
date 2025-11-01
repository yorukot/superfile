package internal

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	bulkrename "github.com/yorukot/superfile/src/internal/ui/bulk_rename"
	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/utils"

	"github.com/barasher/go-exiftool"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

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
func InitialModel(firstFilePanelDirs []string, firstUseCheck bool) tea.Model {
	toggleDotFile, toggleFooter, zClient := initialConfig(firstFilePanelDirs)
	return defaultModelConfig(toggleDotFile, toggleFooter, firstUseCheck, firstFilePanelDirs, zClient)
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
	// TODO : We could check for m.modelQuitState and skip doing anything
	// If its quitDone. But if we are at this state, its already bad, so we need
	// to first figure out if its possible in testing, and fix it.
	slog.Debug("model.Update() called", "msgType", reflect.TypeOf(msg))
	var sidebarCmd, inputCmd, updateCmd, panelCmd, metadataCmd, filePreviewCmd tea.Cmd
	gotModelUpdateMsg := false

	sidebarCmd = m.sidebarModel.UpdateState(msg)

	// this is similar to m.sidebarModel.UpdateState(msg) but since helpMenu is not a Model
	// we call .Update() manually here
	var helpMenuCmd tea.Cmd
	if m.helpMenu.searchBar.Focused() {
		m.helpMenu.searchBar, helpMenuCmd = m.helpMenu.searchBar.Update(msg)
	}

	forcePreviewRender := false

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Must re-render file preview on resize
		m.handleWindowResize(msg)
		forcePreviewRender = true
	case tea.MouseMsg:
		m.handleMouseMsg(msg)
	case tea.KeyMsg:
		inputCmd = m.handleKeyInput(msg)

	case editorFinishedForBulkRenameMsg:
		updateCmd = m.handleEditorFinishedForBulkRename(msg.err)

	// Has to handle zoxide messages separately as they could be generated via
	// zoxide update commands, or batched commands from textinput
	// Cannot do it like processbar messages
	case zoxideui.UpdateMsg:
		slog.Debug("Got ModelUpdate message", "id", msg.GetReqID())
		gotModelUpdateMsg = true
		updateCmd = msg.Apply(&m.zoxideModal)
	case ModelUpdateMessage:
		// TODO: Some of these updates messages should trigger filePanel state update
		// For example a success message for delete operation
		// But, we cant do that as of now, because if every opertion including metadata operation
		// keeps triggering model update below, which will trigger another metadata fetch,
		// we will be stuck in a loop.
		slog.Debug("Got ModelUpdate message", "id", msg.GetReqID())
		gotModelUpdateMsg = true
		updateCmd = msg.ApplyToModel(m)

	default:
		slog.Debug("Message of type that is not handled")
	}

	// This is needed for blink, etc to work
	panelCmd = m.updateFilePanelsState(msg)

	m.updateModelStateAfterMsg()

	// Temp fix till we add metadata cache, to prevent multiple metadata fetch spawns
	// Ideally we might want to fetch only if the current file selected in filepanel changes
	if !gotModelUpdateMsg {
		metadataCmd = m.getMetadataCmd()
		filePreviewCmd = m.getFilePreviewCmd(forcePreviewRender)
	}

	return m, tea.Batch(sidebarCmd, helpMenuCmd, inputCmd, updateCmd, panelCmd, metadataCmd, filePreviewCmd)
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
	m.getFilePanelItems()
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

func (m *model) getFilePreviewCmd(forcePreviewRender bool) tea.Cmd {
	if !m.fileModel.filePreview.IsOpen() {
		return nil
	}
	if len(m.getFocusedFilePanel().element) == 0 {
		// Sync call because this will be fast
		m.fileModel.filePreview.SetContentWithRenderText("")
		return nil
	}
	selectedItem := m.getFocusedFilePanel().getSelectedItem()
	if m.fileModel.filePreview.GetLocation() == selectedItem.location && !forcePreviewRender {
		return nil
	}

	m.fileModel.filePreview.SetLocation(selectedItem.location)
	m.fileModel.filePreview.SetContentWithRenderText("Loading...")
	reqCnt := m.ioReqCnt
	m.ioReqCnt++
	slog.Debug("Submitting file preview render request", "id", reqCnt, "path", selectedItem.location)

	// Copy to a local variable to be used in below closure.
	fullModalWidth := m.fullWidth

	return func() tea.Msg {
		return NewFilePreviewUpdateMsg(selectedItem.location,
			m.fileModel.filePreview.RenderWithPath(selectedItem.location, fullModalWidth), reqCnt)
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
	if len(m.getFocusedFilePanel().element) == 0 {
		m.fileMetaData.SetBlank()
		return nil
	}
	selectedItem := m.getFocusedFilePanel().getSelectedItem()

	// Note : This will cause metadata not being refreshed when you are not scrolling,
	// or filepanel is not getting updated. Its not a big problem as we repeatedly refresh filepanel
	// In case this is a significant issue, we would implement metadata caching.
	// But need to implement it carefully if we do. Make sure cache is not unbounded
	// Remove metadata from filepanel.elemets[] and have cache as source of truth.
	// Have a TTL for expiry, or lister for file update events.
	if len(selectedItem.metaData) > 0 {
		m.fileMetaData.SetMetadata(metadata.NewMetadata(selectedItem.metaData,
			selectedItem.location, ""))
		return nil
	}
	if m.fileMetaData.IsBlank() {
		m.fileMetaData.SetInfoMsg(icon.InOperation + icon.Space + "Loading metadata...")
	}
	metadataFocussed := m.focusPanel == metadataFocus
	reqCnt := m.ioReqCnt
	m.ioReqCnt++
	// If there are too many metadata fetches, we need to have a cache with path as a key
	// and timeout based eviction
	slog.Debug("Submitting metadata fetch request", "id", reqCnt, "path", selectedItem.location)
	return func() tea.Msg {
		return NewMetadataMsg(
			metadata.GetMetadata(selectedItem.location, metadataFocussed, et), reqCnt)
	}
}

// Adjust window size based on msg information
func (m *model) handleWindowResize(msg tea.WindowSizeMsg) {
	m.fullHeight = msg.Height
	m.fullWidth = msg.Width

	m.setHeightValues(msg.Height)

	if m.fileModel.filePreview.IsOpen() {
		// File preview panel width same as file panel
		m.setFilePreviewPanelSize()
	}

	m.setFilePanelsSize(msg.Width)
	m.setHelpMenuSize()
	m.setMetadataModelSize()
	m.setProcessBarModelSize()
	m.setPromptModelSize()
	m.setZoxideModelSize()

	if m.fileModel.maxFilePanel >= 10 {
		m.fileModel.maxFilePanel = 10
	}
}

func (m *model) setFilePreviewPanelSize() {
	m.fileModel.filePreview.SetWidth(m.getFilePreviewWidth())
	m.fileModel.filePreview.SetHeight(m.mainPanelHeight + 2)
}

// Set file preview panel Widht to width. Assure that
func (m *model) getFilePreviewWidth() int {
	if common.Config.FilePreviewWidth == 0 {
		return (m.fullWidth - common.Config.SidebarWidth -
			(4 + (len(m.fileModel.filePanels))*2)) / (len(m.fileModel.filePanels) + 1)
	}
	return (m.fullWidth - common.Config.SidebarWidth) / common.Config.FilePreviewWidth
}

// Proper set panels size. Assure that panels do not overlap
func (m *model) setFilePanelsSize(width int) {
	// set each file panel size and max file panel amount
	m.fileModel.width = (width - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth() -
		(4 + (len(m.fileModel.filePanels)-1)*2)) / len(m.fileModel.filePanels)
	m.fileModel.maxFilePanel = (width - common.Config.SidebarWidth - m.fileModel.filePreview.GetWidth()) / 20
	for i := range m.fileModel.filePanels {
		m.fileModel.filePanels[i].searchBar.Width = m.fileModel.width - 4
	}
}

func (m *model) setHeightValues(height int) {
	//nolint: gocritic // This is to be separated out to a function, and made better later. No need to refactor here
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
	// TODO : Make it grow even more for bigger screen sizes.
	// TODO : Calculate the value , instead of manually hard coding it.

	// Main panel height = Total terminal height- 2(file panel border) - footer height
	m.mainPanelHeight = height - 2 - utils.FullFooterHeight(m.footerHeight, m.toggleFooter)

	for index := range m.fileModel.filePanels {
		m.fileModel.filePanels[index].handleResize(m.mainPanelHeight)
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
	// 2 for border, 1 for left padding, 2 for placeholder icon of searchbar
	// 1 for additional character that View() of search bar function mysteriously adds.
	m.helpMenu.searchBar.Width = m.helpMenu.width - 6
}

func (m *model) setPromptModelSize() {
	// Scale prompt model's maxHeight - 33% of total height
	m.promptModal.SetMaxHeight(m.fullHeight / 3)

	// Scale prompt model's maxHeight - 50% of total height
	m.promptModal.SetWidth(m.fullWidth / 2)
}

func (m *model) setZoxideModelSize() {
	// Scale zoxide model's maxHeight - 50% of total height to accommodate scroll indicators
	m.zoxideModal.SetMaxHeight(m.fullHeight / 2)

	// Scale zoxide model's width - 50% of total width
	m.zoxideModal.SetWidth(m.fullWidth / 2)
}

func (m *model) setMetadataModelSize() {
	m.fileMetaData.SetDimensions(utils.FooterWidth(m.fullWidth)+2, m.footerHeight+2)
}

// TODO: Remove this code duplication with footer models
func (m *model) setProcessBarModelSize() {
	m.processBarModel.SetDimensions(utils.FooterWidth(m.fullWidth)+2, m.footerHeight+2)
}

// Identify the current state of the application m and properly handle the
// msg keybind pressed
func (m *model) handleKeyInput(msg tea.KeyMsg) tea.Cmd {
	slog.Debug("model.handleKeyInput", "msg", msg, "typestr", msg.Type.String(),
		"runes", msg.Runes, "type", int(msg.Type), "paste", msg.Paste,
		"alt", msg.Alt)
	slog.Debug("model.handleKeyInput. model info. ",
		"filePanelFocusIndex", m.filePanelFocusIndex,
		"filePanel.isFocused", m.fileModel.filePanels[m.filePanelFocusIndex].isFocused,
		"filePanel.panelMode", m.fileModel.filePanels[m.filePanelFocusIndex].panelMode,
		"typingModal.open", m.typingModal.open,
		"notifyModel.open", m.notifyModel.IsOpen(),
		"promptModal.open", m.promptModal.IsOpen(),
		"fileModel.renaming", m.fileModel.renaming,
		"searchBar.focussed", m.fileModel.filePanels[m.filePanelFocusIndex].searchBar.Focused(),
		"helpMenu.open", m.helpMenu.open,
		"firstTextInput", m.firstTextInput,
		"focusPanel", m.focusPanel,
	)
	if m.firstUse {
		m.firstUse = false
		return nil
	}

	cmd := m.handleModalOrDefaultKey(msg)
	return m.handleQuitState(cmd, common.Config.CdOnQuit)
}

func (m *model) handleModalOrDefaultKey(msg tea.KeyMsg) tea.Cmd {
	if cmd := m.routeModalInput(msg.String()); cmd != nil {
		return cmd
	}
	if m.isAnyModalActive() {
		return nil
	}
	return m.handleDefaultKey(msg.String())
}

func (m *model) routeModalInput(msg string) tea.Cmd {
	if m.typingModal.open {
		m.typingModalOpenKey(msg)
		return nil
	}
	if m.promptModal.IsOpen() || m.zoxideModal.IsOpen() || m.bulkRenameModel.IsOpen() {
		return nil
	}
	if m.notifyModel.IsOpen() {
		return m.notifyModelOpenKey(msg)
	}
	if m.fileModel.renaming {
		return m.renamingKey(msg)
	}
	if m.sidebarModel.IsRenaming() {
		m.sidebarRenamingKey(msg)
		return nil
	}

	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	if panel.searchBar.Focused() {
		m.focusOnSearchbarKey(msg)
		return nil
	}
	if m.sidebarModel.SearchBarFocused() {
		m.sidebarModel.HandleSearchBarKey(msg)
		return nil
	}
	if panel.sortOptions.open {
		m.sortOptionsKey(msg)
		return nil
	}
	if m.helpMenu.open {
		m.helpMenuKey(msg)
		return nil
	}
	return nil
}

func (m *model) isAnyModalActive() bool {
	panel := m.fileModel.filePanels[m.filePanelFocusIndex]
	return m.typingModal.open || m.promptModal.IsOpen() || m.zoxideModal.IsOpen() ||
		m.notifyModel.IsOpen() || m.bulkRenameModel.IsOpen() || m.fileModel.renaming ||
		m.sidebarModel.IsRenaming() || panel.searchBar.Focused() ||
		m.sidebarModel.SearchBarFocused() || panel.sortOptions.open || m.helpMenu.open
}

func (m *model) handleDefaultKey(msg string) tea.Cmd {
	switch {
	case slices.Contains(common.Hotkeys.Quit, msg):
		m.modelQuitState = quitInitiated
		return nil
	case slices.Contains(common.Hotkeys.CdQuit, msg):
		m.modelQuitState = quitInitiated
		common.Config.CdOnQuit = true
		return nil
	default:
		return m.mainKey(msg)
	}
}

func (m *model) handleQuitState(cmd tea.Cmd, cdOnQuit bool) tea.Cmd {
	if m.modelQuitState == quitInitiated {
		if m.processBarModel.HasRunningProcesses() {
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

func (m *model) updateFilePanelsState(msg tea.Msg) tea.Cmd {
	focusPanel := &m.fileModel.filePanels[m.filePanelFocusIndex]
	var cmd tea.Cmd

	if m.firstTextInput {
		m.firstTextInput = false
		return nil
	}

	cmd = m.handleInputUpdates(msg, focusPanel)

	// The code should never reach this state.
	if focusPanel.cursor < 0 {
		focusPanel.cursor = 0
	}

	return cmd
}

func (m *model) handleInputUpdates(msg tea.Msg, focusPanel *filePanel) tea.Cmd {
	var cmd tea.Cmd
	var action common.ModelAction

	switch {
	case m.bulkRenameModel.IsOpen():
		action, cmd = m.bulkRenameModel.HandleUpdate(msg)
		editorCmd := m.applyBulkRenameAction(action)
		if editorCmd != nil {
			cmd = tea.Batch(cmd, editorCmd)
		}
	case m.fileModel.renaming:
		focusPanel.rename, cmd = focusPanel.rename.Update(msg)
	case focusPanel.searchBar.Focused():
		focusPanel.searchBar, cmd = focusPanel.searchBar.Update(msg)
	case m.typingModal.open:
		m.typingModal.textInput, cmd = m.typingModal.textInput.Update(msg)
	case m.promptModal.IsOpen():
		cwdLocation := m.fileModel.filePanels[m.filePanelFocusIndex].location
		action, cmd = m.promptModal.HandleUpdate(msg, cwdLocation)
		m.applyPromptModalAction(action)
	case m.zoxideModal.IsOpen():
		action, cmd = m.zoxideModal.HandleUpdate(msg)
		m.applyZoxideModalAction(action)
	}

	return cmd
}

// Apply the Action and notify the promptModal
func (m *model) applyPromptModalAction(action common.ModelAction) {
	successMsg, actionErr := m.logAndExecuteAction(action)
	if actionErr != nil {
		m.promptModal.HandleSPFActionResults(false, actionErr.Error())
	} else if successMsg != "" {
		m.promptModal.HandleSPFActionResults(true, successMsg)
	}
}

// Utility function to log and execute actions, reducing duplication
func (m *model) logAndExecuteAction(action common.ModelAction) (string, error) {
	// Only log actions that aren't NoAction to reduce debug noise
	if _, ok := action.(common.NoAction); !ok {
		slog.Debug("Applying model action", "action", action)
	}

	switch action := action.(type) {
	case common.NoAction:
		return "", nil
	case common.ShellCommandAction:
		// Shell commands are handled separately and don't return here
		m.applyShellCommandAction(action.Command)
		return "", nil
	case common.SplitPanelAction:
		return "Panel successfully split", m.splitPanel()
	case common.CDCurrentPanelAction:
		return "Panel directory changed", m.updateCurrentFilePanelDir(action.Location)
	case common.OpenPanelAction:
		return "New panel opened", m.createNewFilePanelRelativeToCurrent(action.Location)
	default:
		return "", errors.New("unhandled action type")
	}
}

// Apply the Action for zoxide modal (no result notifications needed)
func (m *model) applyZoxideModalAction(action common.ModelAction) {
	_, _ = m.logAndExecuteAction(action)
}

func (m *model) applyBulkRenameAction(action common.ModelAction) tea.Cmd {
	if editorAction, ok := action.(bulkrename.EditorModeAction); ok {
		return m.handleEditorModeAction(editorAction)
	}
	if brAction, ok := action.(bulkrename.BulkRenameAction); ok {
		return bulkrename.ExecuteBulkRename(&m.processBarModel, brAction.Previews)
	}
	_, _ = m.logAndExecuteAction(action)
	return nil
}

func (m *model) handleEditorModeAction(action bulkrename.EditorModeAction) tea.Cmd {
	m.pendingEditorAction = &action

	parts := strings.Fields(action.Editor)
	cmd := parts[0]
	args := append(parts[1:], action.TmpfilePath)

	c := exec.Command(cmd, args...)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedForBulkRenameMsg{err}
	})
}

func (m *model) handleEditorFinishedForBulkRename(err error) tea.Cmd {
	if m.pendingEditorAction == nil {
		slog.Error("No pending editor action found")
		return nil
	}

	action := m.pendingEditorAction
	defer func() {
		os.Remove(action.TmpfilePath)
		m.pendingEditorAction = nil
	}()

	if err != nil {
		slog.Error("Editor finished with error", "error", err)
		return nil
	}

	lines, readErr := m.readEditorOutputLines(action.TmpfilePath)
	if readErr != nil {
		return nil
	}

	if len(lines) != len(action.SelectedFiles) {
		slog.Error("Number of lines in tmpfile doesn't match number of selected files",
			"expected", len(action.SelectedFiles), "got", len(lines))
		return nil
	}

	previews := m.buildRenamePreviews(action.SelectedFiles, lines)
	validPreviews := m.filterValidPreviews(previews)

	if len(validPreviews) == 0 {
		slog.Info("No valid renames to apply from tmpfile")
		m.bulkRenameModel.Close()
		return nil
	}

	m.bulkRenameModel.Close()
	return bulkrename.ExecuteBulkRename(&m.processBarModel, validPreviews)
}

func (m *model) readEditorOutputLines(tmpfilePath string) ([]string, error) {
	content, err := os.ReadFile(tmpfilePath)
	if err != nil {
		slog.Error("Failed to read tmpfile", "error", err)
		return nil, err
	}
	return strings.Split(strings.TrimSpace(string(content)), "\n"), nil
}

func (m *model) buildRenamePreviews(selectedFiles []string, newNames []string) []bulkrename.RenamePreview {
	previews := make([]bulkrename.RenamePreview, 0, len(newNames))
	for i, itemPath := range selectedFiles {
		oldName := filepath.Base(itemPath)
		newName := strings.TrimSpace(newNames[i])

		if newName == "" {
			slog.Warn("Empty filename in tmpfile, skipping", "line", i+1)
			continue
		}

		preview := bulkrename.RenamePreview{
			OldPath: itemPath,
			OldName: oldName,
			NewName: newName,
		}

		preview.Error = m.validateRename(itemPath, oldName, newName)
		previews = append(previews, preview)
	}
	return previews
}

func (m *model) validateRename(itemPath, oldName, newName string) string {
	if newName == oldName {
		return "No change"
	}

	newPath := filepath.Join(filepath.Dir(itemPath), newName)
	if strings.EqualFold(itemPath, newPath) {
		return ""
	}

	if _, err := os.Stat(newPath); err == nil {
		return "File already exists"
	}
	return ""
}

func (m *model) filterValidPreviews(previews []bulkrename.RenamePreview) []bulkrename.RenamePreview {
	validPreviews := make([]bulkrename.RenamePreview, 0, len(previews))
	for _, p := range previews {
		if p.Error == "" {
			validPreviews = append(validPreviews, p)
		}
	}
	return validPreviews
}

// TODO : Move them around to appropriate places
func (m *model) applyShellCommandAction(shellCommand string) {
	focusPanelDir := m.fileModel.filePanels[m.filePanelFocusIndex].location

	retCode, output, err := utils.ExecuteCommandInShell(common.DefaultCommandTimeout, focusPanelDir, shellCommand)

	m.promptModal.HandleShellCommandResults(retCode, output)

	if err != nil {
		slog.Error("Command execution failed", "retCode", retCode,
			"error", err, "output", output)
		return
	}
}

func (m *model) splitPanel() error {
	return m.createNewFilePanel(m.fileModel.filePanels[m.filePanelFocusIndex].location)
}

func (m *model) createNewFilePanelRelativeToCurrent(path string) error {
	currentDir := m.fileModel.filePanels[m.filePanelFocusIndex].location
	return m.createNewFilePanel(utils.ResolveAbsPath(currentDir, path))
}

// simulates a 'cd' action
func (m *model) updateCurrentFilePanelDir(path string) error {
	panel := m.getFocusedFilePanel()
	err := panel.updateCurrentFilePanelDir(path)
	if err == nil {
		// Track the directory change with zoxide
		m.trackDirectoryWithZoxide(panel.location)
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
		"fullWidth", m.fullWidth)

	if !m.firstLoadingComplete {
		return "Loading..."
	}

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

	if common.Config.Debug {
		showRenderDebugStatsMain(sidebar, filePanel, filePreview)
	}

	var footer string
	if m.toggleFooter {
		processBar := m.processBarRender()

		metaData := m.fileMetaData.Render(m.focusPanel == metadataFocus)

		clipboardBar := m.clipboardRender()

		footer = lipgloss.JoinHorizontal(0, processBar, metaData, clipboardBar)
		if common.Config.Debug {
			showRenderDebugStatsFooter(processBar, metaData, clipboardBar)
		}
	}

	var finalRender string

	if m.toggleFooter {
		finalRender = lipgloss.JoinVertical(0, mainPanel, footer)
	} else {
		finalRender = mainPanel
	}

	finalRender = m.updateRenderForOverlay(finalRender)

	return finalRender
}

func (m *model) updateRenderForOverlay(finalRender string) string {
	// check if need pop up modal
	if m.helpMenu.open {
		helpMenu := m.helpMenuRender()
		overlayX := m.fullWidth/2 - m.helpMenu.width/2
		overlayY := m.fullHeight/2 - m.helpMenu.height/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, helpMenu, finalRender)
	}

	if m.promptModal.IsOpen() {
		promptModal := m.promptModalRender()
		overlayX := m.fullWidth/2 - m.promptModal.GetWidth()/2
		overlayY := m.fullHeight/2 - m.promptModal.GetMaxHeight()/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, promptModal, finalRender)
	}

	if m.zoxideModal.IsOpen() {
		zoxideModal := m.zoxideModalRender()
		overlayX := m.fullWidth/2 - m.zoxideModal.GetWidth()/2
		overlayY := m.fullHeight/2 - m.zoxideModal.GetMaxHeight()/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, zoxideModal, finalRender)
	}

	panel := m.fileModel.filePanels[m.filePanelFocusIndex]

	if panel.sortOptions.open {
		sortOptions := m.sortOptionsRender()
		overlayX := m.fullWidth/2 - panel.sortOptions.width/29
		overlayY := m.fullHeight/2 - panel.sortOptions.height/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, sortOptions, finalRender)
	}

	if m.firstUse {
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

	if m.bulkRenameModel.IsOpen() {
		bulkRenameModal := m.bulkRenameModel.Render()
		overlayX := m.fullWidth/2 - m.bulkRenameModel.GetWidth()/2
		overlayY := m.fullHeight/2 - m.bulkRenameModel.GetHeight()/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, bulkRenameModal, finalRender)
	}

	if m.notifyModel.IsOpen() {
		notifyModal := m.notifyModel.Render()
		overlayX := m.fullWidth/2 - common.ModalWidth/2
		overlayY := m.fullHeight/2 - common.ModalHeight/2
		return stringfunction.PlaceOverlay(overlayX, overlayY, notifyModal, finalRender)
	}
	return finalRender
}

func showRenderDebugStatsMain(sidebar, filePanel, filePreview string) {
	slog.Debug("Render stats for main panel",
		"sidebarLineCnt", getLineCnt(sidebar), "sidebarMaxW", getMaxW(sidebar),
		"filePanelLineCnt", getLineCnt(filePanel), "filePanelMaxW", getMaxW(filePanel),
		"filePreviewLineCnt", getLineCnt(filePreview), "filePreviewMaxW", getMaxW(filePreview),
	)
}

func showRenderDebugStatsFooter(processBar, metaData, clipboardBar string) {
	slog.Debug("Render stats for footer",
		"processBarLineCnt", getLineCnt(processBar), "processBarMaxW", getMaxW(processBar),
		"metaDataLineCnt", getLineCnt(metaData), "metaDataMaxW", getMaxW(metaData),
		"clipboardBarLineCnt", getLineCnt(clipboardBar), "clipboardBarMaxW", getMaxW(clipboardBar),
	)
}

func getLineCnt(s string) int {
	return strings.Count(s, "\n") + 1
}

func getMaxW(s string) int {
	maxW := 0
	for line := range strings.Lines(s) {
		maxW = max(maxW, ansi.StringWidth(line))
	}
	return maxW
}

// Render and update file panel items. Check for changes and updates in files and
// folders in the current directory.
func (m *model) getFilePanelItems() {
	focusPanel := m.fileModel.filePanels[m.filePanelFocusIndex]
	for i, filePanel := range m.fileModel.filePanels {
		var fileElement []element
		nowTime := time.Now()
		// Check last time each element was updated, if less then 3 seconds ignore
		if !filePanel.isFocused && nowTime.Sub(filePanel.lastTimeGetElement) < 3*time.Second {
			// TODO : revisit this. This feels like a duct tape solution of an actual
			// deep rooted problem. This feels very hacky.
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

		if filePanel.isFocused && !focusPanelReRender &&
			nowTime.Sub(filePanel.lastTimeGetElement) < time.Duration(reRenderTime)*time.Second {
			continue
		}

		// Get file names based on search bar filter
		if filePanel.searchBar.Value() != "" {
			fileElement = returnDirElementBySearchString(filePanel.location, m.toggleDotFile,
				filePanel.searchBar.Value(), filePanel.sortOptions.data)
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

// Close superfile application. Cd into the current dir if CdOnQuit on and save
// the path in state direcotory
func (m *model) quitSuperfile(cdOnQuit bool) {
	// close exiftool session
	if common.Config.Metadata && et != nil {
		et.Close()
	}
	// cd on quit
	currentDir := m.fileModel.filePanels[m.filePanelFocusIndex].location
	variable.SetLastDir(currentDir)

	if cdOnQuit {
		// escape single quote
		currentDir = strings.ReplaceAll(currentDir, "'", "'\\''")
		err := os.WriteFile(variable.LastDirFile, []byte("cd '"+currentDir+"'"), 0o755)
		if err != nil {
			slog.Error("Error during writing lastdir file", "error", err)
		}
	}
	m.modelQuitState = quitDone
	slog.Debug("Quitting superfile", "current dir", currentDir)
}
