package bulkrename

import (
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

type RenameType int

const (
	FindReplace RenameType = iota
	AddPrefix
	AddSuffix
	AddNumbering
	ChangeCase
	EditorMode
)

type CaseType int

const (
	CaseLower CaseType = iota
	CaseUpper
	CaseTitle
)

type Model struct {
	open         bool
	renameType   RenameType
	caseType     CaseType
	cursor       int
	startNumber  int
	errorMessage string

	findInput    textinput.Model
	replaceInput textinput.Model
	prefixInput  textinput.Model
	suffixInput  textinput.Model

	preview []RenamePreview

	reqCnt int

	width  int
	height int

	selectedFiles []string
	currentDir    string
}

type RenamePreview struct {
	OldPath string
	OldName string
	NewName string
	Error   string
}

type UpdateMsg struct {
	reqID int
}

func (msg UpdateMsg) GetReqID() int {
	return msg.reqID
}

func DefaultModel(maxHeight int, width int) Model {
	return Model{
		open:        false,
		renameType:  FindReplace,
		caseType:    CaseLower,
		cursor:      0,
		startNumber: 1,
		width:       width,
		height:      maxHeight,
		reqCnt:      0,
	}
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Open(selectedFiles []string, currentDir string) {
	if len(selectedFiles) == 0 {
		return
	}

	m.open = true
	m.selectedFiles = selectedFiles
	m.currentDir = currentDir
	m.renameType = FindReplace
	m.caseType = CaseLower
	m.cursor = 0
	m.startNumber = 1
	m.errorMessage = ""
	m.preview = nil

	m.findInput = common.GenerateBulkRenameTextInput("Find text")
	m.replaceInput = common.GenerateBulkRenameTextInput("Replace with")
	m.prefixInput = common.GenerateBulkRenameTextInput("Add prefix")
	m.suffixInput = common.GenerateBulkRenameTextInput("Add suffix")

	m.focusInput()
}

func (m *Model) Close() {
	m.open = false
	m.selectedFiles = nil
	m.currentDir = ""
	m.errorMessage = ""
	m.preview = nil

	m.findInput.Blur()
	m.replaceInput.Blur()
	m.prefixInput.Blur()
	m.suffixInput.Blur()
}

func (m *Model) HandleUpdate(msg tea.Msg) (common.ModelAction, tea.Cmd) {
	slog.Debug("bulk_rename.Model HandleUpdate()", "msg", msg)
	var action common.ModelAction = common.NoAction{}
	var cmd tea.Cmd

	if !m.IsOpen() {
		slog.Error("HandleUpdate called on closed bulk rename modal")
		return action, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case slices.Contains(common.Hotkeys.CancelTyping, msg.String()):
			m.Close()
		case slices.Contains(common.Hotkeys.ConfirmTyping, msg.String()):
			action = m.handleConfirm()
		case slices.Contains(common.Hotkeys.ListUp, msg.String()):
			m.adjustValue(-1)
		case slices.Contains(common.Hotkeys.ListDown, msg.String()):
			m.adjustValue(1)
		case slices.Contains(common.Hotkeys.NavBulkRename, msg.String()):
			m.nextType()
		case slices.Contains(common.Hotkeys.RevNavBulkRename, msg.String()):
			m.prevType()
		default:
			cmd = m.handleTextInputUpdate(msg)
		}
	default:
		cmd = m.handleTextInputUpdate(msg)
	}

	return action, cmd
}

func (m *Model) handleConfirm() common.ModelAction {
	if m.renameType == EditorMode {
		return m.handleEditorMode()
	}

	previews := m.GeneratePreview()

	validPreviews := make([]RenamePreview, 0, len(previews))
	for _, p := range previews {
		if p.Error == "" {
			validPreviews = append(validPreviews, p)
		}
	}

	if len(validPreviews) == 0 {
		m.errorMessage = "No valid renames to apply"
		return common.NoAction{}
	}

	m.Close()
	return BulkRenameAction{
		Previews: validPreviews,
	}
}

func (m *Model) handleEditorMode() common.ModelAction {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}

	tmpfile, err := os.CreateTemp("", "superfile-bulk-rename-*.txt")
	if err != nil {
		m.errorMessage = "Failed to create temporary file: " + err.Error()
		return common.NoAction{}
	}
	tmpfilePath := tmpfile.Name()

	for _, itemPath := range m.selectedFiles {
		filename := filepath.Base(itemPath)
		_, err := tmpfile.WriteString(filename + "\n")
		if err != nil {
			tmpfile.Close()
			os.Remove(tmpfilePath)
			m.errorMessage = "Failed to write to temporary file: " + err.Error()
			return common.NoAction{}
		}
	}
	tmpfile.Close()

	return EditorModeAction{
		TmpfilePath:   tmpfilePath,
		Editor:        editor,
		SelectedFiles: m.selectedFiles,
		CurrentDir:    m.currentDir,
	}
}

type EditorModeAction struct {
	TmpfilePath   string
	Editor        string
	SelectedFiles []string
	CurrentDir    string
}

func (a EditorModeAction) String() string {
	return "EditorModeAction with editor: " + a.Editor
}

type BulkRenameAction struct {
	Previews []RenamePreview
}

func (a BulkRenameAction) String() string {
	return "BulkRenameAction with " + strconv.Itoa(len(a.Previews)) + " items"
}

func (m *Model) handleTextInputUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch m.renameType {
	case FindReplace:
		if m.cursor == 0 {
			m.findInput, cmd = m.findInput.Update(msg)
		} else {
			m.replaceInput, cmd = m.replaceInput.Update(msg)
		}
	case AddPrefix:
		m.prefixInput, cmd = m.prefixInput.Update(msg)
	case AddSuffix:
		m.suffixInput, cmd = m.suffixInput.Update(msg)
	}

	m.preview = nil

	return cmd
}

func (m *Model) adjustValue(delta int) {
	switch m.renameType {
	case FindReplace:
		m.navigateCursor(delta)
	case AddNumbering:
		newValue := m.startNumber + delta
		if newValue >= 0 {
			m.startNumber = newValue
			m.preview = nil
		}
	case ChangeCase:
		newValue := int(m.caseType) + delta
		if newValue >= 0 && newValue <= 2 {
			m.caseType = CaseType(newValue)
			m.preview = nil
		}
	}
}
func (m *Model) navigateCursor(delta int) {
	if delta > 0 {
		m.navigateDown()
	} else {
		m.navigateUp()
	}
}

func (m *Model) navigateUp() {
	if m.cursor > 0 {
		m.cursor--
		m.focusInput()
	}
}

func (m *Model) navigateDown() {
	if m.cursor < 1 {
		m.cursor++
		m.focusInput()
	}
}

func (m *Model) focusInput() {
	m.findInput.Blur()
	m.replaceInput.Blur()
	m.prefixInput.Blur()
	m.suffixInput.Blur()

	switch m.renameType {
	case FindReplace:
		if m.cursor == 0 {
			m.findInput.Focus()
		} else {
			m.replaceInput.Focus()
		}
	case AddPrefix:
		m.prefixInput.Focus()
	case AddSuffix:
		m.suffixInput.Focus()

	}
}

func (m *Model) nextType() {
	m.renameType = RenameType((int(m.renameType) + 1) % 6)
	m.focusInput()
	m.preview = nil
}

func (m *Model) prevType() {
	newType := int(m.renameType) - 1
	if newType < 0 {
		newType = 5
	}
	m.renameType = RenameType(newType)
	m.focusInput()
	m.preview = nil
}
func (m *Model) GeneratePreview() []RenamePreview {
	previews := make([]RenamePreview, 0, len(m.selectedFiles))

	for i, itemPath := range m.selectedFiles {
		preview := m.createRenamePreview(itemPath, i)
		previews = append(previews, preview)
	}

	m.preview = previews
	return previews
}

func (m *Model) createRenamePreview(itemPath string, index int) RenamePreview {
	oldName := filepath.Base(itemPath)
	newName := m.applyTransformation(oldName, index)

	return RenamePreview{
		OldPath: itemPath,
		OldName: oldName,
		NewName: newName,
		Error:   m.validateRename(itemPath, oldName, newName),
	}
}

func (m *Model) applyTransformation(oldName string, index int) string {
	switch m.renameType {
	case FindReplace:
		return m.applyFindReplace(oldName)
	case AddPrefix:
		return m.applyPrefix(oldName)
	case AddSuffix:
		return m.applySuffix(oldName)
	case AddNumbering:
		return m.applyNumbering(oldName, index)
	case ChangeCase:
		return m.applyCaseConversion(oldName)
	default:
		return oldName
	}
}

func (m *Model) applyFindReplace(filename string) string {
	find := m.findInput.Value()
	replace := m.replaceInput.Value()
	if find == "" {
		return filename
	}
	return strings.ReplaceAll(filename, find, replace)
}

func (m *Model) applyPrefix(filename string) string {
	prefix := m.prefixInput.Value()
	if prefix == "" {
		return filename
	}
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return prefix + nameWithoutExt + ext
}

func (m *Model) applySuffix(filename string) string {
	suffix := m.suffixInput.Value()
	if suffix == "" {
		return filename
	}
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return nameWithoutExt + suffix + ext
}

func (m *Model) applyNumbering(filename string, number int) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return nameWithoutExt + "_" + strconv.Itoa(m.startNumber+number) + ext
}

func (m *Model) applyCaseConversion(filename string) string {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	var converted string
	switch m.caseType {
	case CaseLower:
		converted = strings.ToLower(nameWithoutExt)
	case CaseUpper:
		converted = strings.ToUpper(nameWithoutExt)
	case CaseTitle:
		converted = toTitleCase(nameWithoutExt)
	default:
		converted = nameWithoutExt
	}

	return converted + ext
}

func toTitleCase(text string) string {
	words := strings.Fields(strings.ToLower(text))
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

func (m *Model) validateRename(itemPath, oldName, newName string) string {
	if newName == "" {
		return "Empty filename"
	}
	if newName == oldName {
		return "No change"
	}

	newPath := filepath.Join(filepath.Dir(itemPath), newName)

	if strings.EqualFold(itemPath, newPath) {
		return ""
	}

	if _, statErr := os.Stat(newPath); statErr == nil {
		return "File already exists"
	}
	return ""
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

func ExecuteBulkRename(processBarModel *processbar.Model, previews []RenamePreview) tea.Cmd {
	return func() tea.Msg {
		state := bulkRenameOperation(processBarModel, previews)
		return NewBulkRenameResultMsg(state, len(previews))
	}
}

func bulkRenameOperation(processBarModel *processbar.Model, previews []RenamePreview) processbar.ProcessState {
	if len(previews) == 0 {
		return processbar.Cancelled
	}

	p, err := processBarModel.SendAddProcessMsg(icon.Terminal+icon.Space+"Bulk Rename", len(previews), true)
	if err != nil {
		slog.Error("Cannot spawn bulk rename process", "error", err)
		return processbar.Failed
	}

	for _, preview := range previews {
		newPath := filepath.Join(filepath.Dir(preview.OldPath), preview.NewName)

		err = os.Rename(preview.OldPath, newPath)
		if err != nil {
			p.State = processbar.Failed
			slog.Error("Error in bulk rename operation", "old", preview.OldPath, "new", newPath, "error", err)
			break
		}

		p.Name = icon.Terminal + icon.Space + preview.NewName
		p.Done++
		processBarModel.TrySendingUpdateProcessMsg(p)
	}

	if p.State != processbar.Failed {
		p.State = processbar.Successful
		processBarModel.TrySendingUpdateProcessMsg(p)
	}

	return p.State
}

type BulkRenameResultMsg struct {
	state processbar.ProcessState
	count int
}

func NewBulkRenameResultMsg(state processbar.ProcessState, count int) BulkRenameResultMsg {
	return BulkRenameResultMsg{state: state, count: count}
}

func (msg BulkRenameResultMsg) GetState() processbar.ProcessState {
	return msg.state
}

func (msg BulkRenameResultMsg) GetCount() int {
	return msg.count
}

func GetCursorColor() lipgloss.Color {
	return lipgloss.Color(common.Theme.Cursor)
}
