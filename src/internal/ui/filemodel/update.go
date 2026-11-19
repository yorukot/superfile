package filemodel

import (
	"fmt"
	"log/slog"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/preview"
)

func (m *Model) CreateNewFilePanel(location string) (tea.Cmd, error) {
	if m.PanelCount() >= m.MaxFilePanel {
		return nil, ErrMaximumPanelCount
	}

	if _, err := os.Stat(location); err != nil {
		return nil, fmt.Errorf("cannot access location : %s", location)
	}

	m.FilePanels = append(m.FilePanels, filepanel.New(
		location, false, "", m.GetFocusedFilePanel().SortKind,
		m.GetFocusedFilePanel().SortReversed))

	newPanelIndex := m.PanelCount() - 1

	m.FilePanels[m.FocusedPanelIndex].IsFocused = false
	m.FilePanels[newPanelIndex].IsFocused = true
	m.FilePanels[newPanelIndex].SetHeight(m.Height)
	m.FocusedPanelIndex = newPanelIndex

	m.updateChildComponentWidth()
	return m.ensurePreviewDimensionsSync(), nil
}

func (m *Model) CloseFilePanel() (tea.Cmd, error) {
	if m.PanelCount() <= 1 {
		return nil, ErrMinimumPanelCount
	}

	m.FilePanels = append(m.FilePanels[:m.FocusedPanelIndex],
		m.FilePanels[m.FocusedPanelIndex+1:]...)

	if m.FocusedPanelIndex != 0 {
		m.FocusedPanelIndex--
	}
	m.FilePanels[m.FocusedPanelIndex].IsFocused = true
	m.updateChildComponentWidth()

	return m.ensurePreviewDimensionsSync(), nil
}

func (m *Model) ToggleFilePreviewPanel() tea.Cmd {
	m.FilePreview.ToggleOpen()
	m.updateChildComponentWidth()
	return m.ensurePreviewDimensionsSync()
}

func (m *Model) UpdatePreviewPanel(msg preview.UpdateMsg) tea.Cmd {
	var target string
	if m.previewOverride != "" {
		target = m.previewOverride
	} else {
		selectedItem := m.GetFocusedFilePanel().GetFocusedItemPtr()
		if selectedItem == nil {
			slog.Debug("Panel empty or cursor invalid. Ignoring FilePreviewUpdateMsg")
			return nil
		}
		target = selectedItem.Location
	}
	if target != msg.GetLocation() {
		slog.Debug("FilePreviewUpdateMsg for older files. Ignoring",
			"curLocation", target, "msgLocation", msg.GetLocation())
		return nil
	}

	if m.ExpectedPreviewWidth != msg.GetContentWidth() ||
		m.Height != msg.GetContentHeight() {
		slog.Debug("FilePreviewUpdateMsg for older dimensions. Ignoring",
			"curW", m.ExpectedPreviewWidth, "curH", m.Height,
			"msgW", msg.GetContentWidth(), "msgH", msg.GetContentHeight())
		return nil
	}
	m.FilePreview.Apply(msg)

	// For Kitty images, transmit image data directly to the terminal
	if raw := msg.GetRawTransmit(); raw != "" {
		return tea.Raw(raw)
	}
	return nil
}

func (m *Model) GetFilePreviewCmd(forcePreviewRender bool) tea.Cmd {
	if !m.FilePreview.IsOpen() {
		return nil
	}
	if m.previewOverride != "" {
		// While an override is active, the pane belongs to the caller
		// (e.g. the find modal). Only an explicit re-render (force)
		// touches the pane.
		if !forcePreviewRender {
			return nil
		}
		return m.submitPreviewRender(m.previewOverride)
	}
	panel := m.GetFocusedFilePanel()
	if panel.EmptyOrInvalid() {
		// Sync call because this will be fast
		m.FilePreview.SetEmptyWithDimensions(m.ExpectedPreviewWidth, m.Height)
		return nil
	}
	selectedItem := panel.GetFocusedItem()
	if m.FilePreview.GetLocation() == selectedItem.Location && !forcePreviewRender {
		return nil
	}
	return m.submitPreviewRender(selectedItem.Location)
}

// submitPreviewRender marks the pane as loading path and returns a
// command that renders it asynchronously.
func (m *Model) submitPreviewRender(path string) tea.Cmd {
	m.FilePreview.SetLocation(path)
	m.FilePreview.SetLoading()

	// HACK!!!. fileModel must not be aware of other dimensions. but...
	// Unfortunately, previewPanel isn't completely 'under' fileModel
	// Note: Must save the dimensions for the closure of the Cmd to avoid
	// problems
	fullModalWidth := m.Width + common.Config.SidebarWidth
	if common.Config.SidebarWidth != 0 {
		fullModalWidth += common.BorderPadding
	}
	width := m.ExpectedPreviewWidth
	height := m.Height

	reqCnt := m.ioReqCnt
	m.ioReqCnt++
	slog.Debug("Submitting file preview render request", "id", reqCnt,
		"path", path, "w", width, "h", height)

	return func() tea.Msg {
		content, rawTransmit := m.FilePreview.RenderWithPath(path, width, height, fullModalWidth)
		return preview.NewUpdateMsg(path, content, rawTransmit,
			width, height, reqCnt)
	}
}

// SetPreviewPathCmd overrides the focused panel's selection and makes
// the preview pane render path. The override is remembered even while
// the pane is closed.
func (m *Model) SetPreviewPathCmd(path string) tea.Cmd {
	if path == "" || m.previewOverride == path {
		return nil
	}
	m.previewOverride = path
	if !m.FilePreview.IsOpen() {
		return nil
	}
	if m.FilePreview.GetLocation() == path {
		return nil
	}
	return m.submitPreviewRender(path)
}

// ClearPreviewOverride removes the preview override and restores the
// pane to the focused panel's selection, unless the pane already shows
// it.
func (m *Model) ClearPreviewOverride() tea.Cmd {
	if m.previewOverride == "" {
		return nil
	}
	m.previewOverride = ""
	if !m.FilePreview.IsOpen() {
		return nil
	}
	panel := m.GetFocusedFilePanel()
	if panel.EmptyOrInvalid() {
		// Sync call because this will be fast
		m.FilePreview.SetEmptyWithDimensions(m.ExpectedPreviewWidth, m.Height)
		return nil
	}
	if m.FilePreview.GetLocation() == panel.GetFocusedItem().Location {
		return nil
	}
	return m.GetFilePreviewCmd(true)
}

// HasPreviewOverride reports whether a preview override is active.
func (m *Model) HasPreviewOverride() bool {
	return m.previewOverride != ""
}

func (m *Model) ToggleDotFile() {
	m.DisplayDotFiles = !m.DisplayDotFiles
	m.UpdateFilePanelsIfNeeded(true)
}

func (m *Model) UpdateFilePanelsIfNeeded(force bool) {
	for i := range m.FilePanels {
		m.FilePanels[i].UpdateElementsIfNeeded(force, m.DisplayDotFiles)
	}
}
