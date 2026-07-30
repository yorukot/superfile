package contentpanel

import (
	"log/slog"
	"os"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
)

func (m *Model) GetContent() string {
	return m.content
}

func (m *Model) GetContentWidth() int {
	return m.contentWidth
}

func (m *Model) GetContentHeight() int {
	return m.contentHeight
}

func (m *Model) GetLocation() string {
	return m.location
}

func (m *Model) SetOpen(open bool) {
	m.open = open
}

func (m *Model) SetLocation(location string) {
	m.location = location
}

func (m *Model) SetLoading() {
	m.loading = true
}

// All content change happen via this only, to ensure the sync between
// content and width x height, and the loading variable reset
func (m *Model) setContent(content string, width int, height int, location string) {
	m.content = content
	m.contentWidth = width
	m.contentHeight = height
	m.location = location
	m.loading = false
}

func (m *Model) SetEmptyWithDimensions(width int, height int) {
	m.setContent(m.RenderTextWithDimension("", height, width), width, height, "")
}

func (m *Model) IsLoading() bool {
	return m.loading
}

func (m *Model) ToggleOpen() {
	m.open = !m.open
}

func (m *Model) CleanUp() {
	if m.thumbnailGenerator != nil {
		err := m.thumbnailGenerator.CleanUp()
		if err != nil {
			slog.Error("Error While cleaning up TempDirectory", "error", err)
		}
	}
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Open() {
	m.open = true
}

func (m *Model) Close() {
	m.open = false
}

func (m *Model) SetFocused(focused bool) {
	if m.focused != focused {
		m.focused = focused
		// Invalidate cached content so border is re-rendered
		m.location = ""
	}
}

func (m *Model) IsFocused() bool {
	return m.focused
}

func (m *Model) ScrollUp() {
	if m.scrollOffset > 0 {
		m.scrollOffset--
	}
}

func (m *Model) ScrollDown() {
	m.scrollOffset++
}

func (m *Model) ScrollPageUp(pageSize int) {
	m.scrollOffset = max(0, m.scrollOffset-pageSize)
}

func (m *Model) ScrollPageDown(pageSize int) {
	m.scrollOffset += pageSize
}

func (m *Model) ResetScroll() {
	m.scrollOffset = 0
}

// IsEditing returns true when the panel is in edit mode.
func (m *Model) IsEditing() bool {
	return m.editMode
}

// IsAutoRefreshPaused returns true when auto-refresh is paused (during edit).
func (m *Model) IsAutoRefreshPaused() bool {
	return m.autoRefreshPaused
}

// EnterEditMode attempts to enter edit mode for the file at path.
// Returns false with a reason if the file cannot be edited.
func (m *Model) EnterEditMode(path string) (bool, string) {
	info, err := os.Stat(path)
	if err != nil {
		return false, "Cannot access file"
	}

	editable, reason := checkEditable(path, info, common.Config.ContentPanelMaxEditSize)
	if !editable {
		m.editableReason = reason
		return false, reason
	}

	data, err := readFileContent(path, common.Config.ContentPanelMaxEditSize)
	if err != nil {
		return false, "Cannot read file content"
	}

	m.editMode = true
	m.autoRefreshPaused = true
	m.originalPath = path
	m.originalPerm = info.Mode().Perm()
	m.editableReason = ""

	// Seed textarea with file content
	m.textarea = textarea.New()
	m.textarea.SetValue(data)
	m.textarea.Focus()
	m.textarea.SetWidth(m.contentWidth)
	m.textarea.SetHeight(m.contentHeight)

	return true, ""
}

// SaveEdit writes the textarea content back to the file, preserving
// the original file mode.
func (m *Model) SaveEdit() error {
	if !m.editMode || m.originalPath == "" {
		return nil
	}
	content := m.textarea.Value()
	err := os.WriteFile(m.originalPath, []byte(content), m.originalPerm)
	if err != nil {
		return err
	}
	return nil
}

// ExitEditMode leaves edit mode and resumes auto-refresh.
func (m *Model) ExitEditMode() {
	m.editMode = false
	m.autoRefreshPaused = false
	m.originalPath = ""
	m.editableReason = ""
}

// GetEditableReason returns why a file cannot be edited (empty if editable).
func (m *Model) GetEditableReason() string {
	return m.editableReason
}

// HandleEditKey forwards a key message to the textarea for character input.
func (m *Model) HandleEditKey(msg tea.KeyPressMsg) tea.Cmd {
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return cmd
}

// SetShellOutput stores shell command results for display.
func (m *Model) SetShellOutput(command string, exitCode int, output string) {
	m.shellCommand = command
	m.shellExit = exitCode
	m.shellOutput = output
	m.autoRefreshPaused = true // keep output visible
}

// ClearShellOutput resets the shell output display so the preview can resume.
func (m *Model) ClearShellOutput() {
	m.shellOutput = ""
	m.shellCommand = ""
	m.autoRefreshPaused = false
}

// HasShellOutput returns true when shell results are being displayed.
func (m *Model) HasShellOutput() bool {
	return m.shellOutput != "" || m.shellCommand != ""
}

// EnterNvimMode starts embedded nvim in a PTY for the given file path.
func (m *Model) EnterNvimMode(path string) error {
	ns, err := StartNvim(path, max(m.contentWidth, 40), max(m.contentHeight, 5))
	if err != nil {
		return err
	}
	m.nvim = ns
	m.editMode = true
	m.autoRefreshPaused = true
	m.originalPath = path
	return nil
}

// IsNvimRunning returns true when embedded nvim is active.
func (m *Model) IsNvimRunning() bool {
	return m.nvim != nil && m.nvim.IsRunning()
}

// StopNvim terminates the embedded nvim session and marks for content refresh.
func (m *Model) StopNvim() {
	if m.nvim != nil {
		m.nvim.Stop()
		m.nvim = nil
	}
	m.editMode = false
	m.autoRefreshPaused = false
	m.originalPath = ""
	// Clear location so GetFilePreviewCmd re-renders the updated file
	m.location = ""
	m.content = ""
	m.loading = false
}

// ResizeNvim updates the PTY dimensions.
func (m *Model) ResizeNvim(width, height int) {
	if m.nvim != nil && m.nvim.IsRunning() {
		if err := m.nvim.Resize(width, height); err != nil {
			slog.Debug("Failed to resize nvim PTY", "error", err)
		}
	}
}

// WriteNvimKey sends a key sequence to the nvim PTY.
func (m *Model) WriteNvimKey(data []byte) {
	if m.nvim != nil && m.nvim.IsRunning() {
		m.nvim.WriteKey(data)
	}
}

// GetNvimOutput returns the rendered nvim output (ANSI stripped for now).
func (m *Model) GetNvimOutput() string {
	if m.nvim == nil {
		return ""
	}
	// Poll nvim status
	if !m.nvim.IsRunning() {
		// Nvim exited — clean up and refresh
		m.StopNvim()
		return ""
	}
	return m.nvim.GetOutput()
}
