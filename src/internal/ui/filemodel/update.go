package filemodel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/preview"
)

const (
	remotePreviewTimeout = 10 * time.Second
	remotePreviewLimit   = 1024 * 1024
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
	closedSessionID := m.FilePanels[m.FocusedPanelIndex].CurrentLocation().SessionID

	m.FilePanels = append(m.FilePanels[:m.FocusedPanelIndex],
		m.FilePanels[m.FocusedPanelIndex+1:]...)

	if m.FocusedPanelIndex != 0 {
		m.FocusedPanelIndex--
	}
	m.FilePanels[m.FocusedPanelIndex].IsFocused = true
	m.updateChildComponentWidth()
	if err := m.closeSessionIfUnused(closedSessionID); err != nil {
		return m.ensurePreviewDimensionsSync(), fmt.Errorf("close unused panel session: %w", err)
	}

	return m.ensurePreviewDimensionsSync(), nil
}

func (m *Model) ToggleFilePreviewPanel() tea.Cmd {
	m.FilePreview.ToggleOpen()
	m.updateChildComponentWidth()
	return m.ensurePreviewDimensionsSync()
}

func (m *Model) UpdatePreviewPanel(msg preview.UpdateMsg) tea.Cmd {
	if errors.Is(msg.GetError(), filesystem.ErrDisconnected) && msg.GetSessionID() != "" {
		if err := m.MarkSessionDisconnectedIfCurrent(
			msg.GetSessionID(),
			msg.GetSessionGeneration(),
			msg.GetError(),
		); err != nil {
			slog.Error("failed to mark remote preview session disconnected", "error", err)
		}
	}
	selectedItem := m.GetFocusedFilePanel().GetFocusedItemPtr()
	if selectedItem == nil {
		slog.Debug("Panel empty or cursor invalid. Ignoring FilePreviewUpdateMsg")
		return nil
	}
	if selectedItem.Location != msg.GetLocation() {
		slog.Debug("FilePreviewUpdateMsg for older files. Ignoring",
			"curLocation", selectedItem.Location, "msgLocation", msg.GetLocation())
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

	m.FilePreview.SetLocation(selectedItem.Location)
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
		"path", selectedItem.Location, "w", width, "h", height)

	if panel.CurrentLocation().Provider != filesystem.ProviderLocal {
		sessionState, ok := m.Sessions.Get(panel.CurrentLocation().SessionID)
		if !ok || sessionState.Browser == nil {
			content, rawTransmit := m.FilePreview.RenderTextPreviewWithDimension(
				"Remote preview unavailable: session is disconnected.",
				height,
				width,
			)
			return func() tea.Msg {
				return preview.NewUpdateMsg(selectedItem.Location, content, rawTransmit, width, height, reqCnt)
			}
		}
		browser := sessionState.Browser
		sessionID := panel.CurrentLocation().SessionID
		return func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), remotePreviewTimeout)
			defer cancel()
			text, previewErr := remotePreviewText(ctx, browser, selectedItem.Path, selectedItem.Directory)
			content, _ := m.FilePreview.RenderTextPreviewWithDimension(text, height, width)
			return preview.NewRemoteUpdateMsg(
				selectedItem.Location,
				content,
				width,
				height,
				reqCnt,
				sessionID,
				sessionState.Generation,
				previewErr,
			)
		}
	}

	return func() tea.Msg {
		content, rawTransmit := m.FilePreview.RenderWithPath(selectedItem.Location, width, height, fullModalWidth)
		return preview.NewUpdateMsg(selectedItem.Location, content, rawTransmit,
			width, height, reqCnt)
	}
}

func remotePreviewText(
	ctx context.Context,
	session filesystem.Session,
	path filesystem.Path,
	directory bool,
) (string, error) {
	if directory {
		entries, err := session.List(ctx, path)
		if err != nil {
			return "Remote preview error: " + err.Error(), err
		}
		lines := make([]string, 0, len(entries))
		for _, entry := range entries {
			name := entry.Name
			if entry.Stat.IsDir {
				name += "/"
			}
			lines = append(lines, name)
		}
		if len(lines) == 0 {
			return "Empty remote directory", nil
		}
		return strings.Join(lines, "\n"), nil
	}

	reader, err := session.Read(ctx, path)
	if err != nil {
		return "Remote preview error: " + err.Error(), err
	}
	defer reader.Close()
	data, err := io.ReadAll(io.LimitReader(reader, remotePreviewLimit+1))
	if err != nil {
		return "Remote preview error: " + err.Error(), err
	}
	if len(data) > remotePreviewLimit {
		data = data[:remotePreviewLimit]
		data = append(data, []byte("\n\n[preview truncated at 1 MiB]")...)
	}
	if !utf8.Valid(data) || bytes.IndexByte(data, 0) >= 0 {
		return "Binary remote file preview is not supported.", nil
	}
	return string(data), nil
}

func (m *Model) ToggleDotFile() tea.Cmd {
	m.DisplayDotFiles = !m.DisplayDotFiles
	m.UpdateLocalFilePanelsIfNeeded(true)
	return m.GetRemoteFilePanelUpdateCmd(true)
}

func (m *Model) UpdateFilePanelsIfNeeded(force bool) {
	for i := range m.FilePanels {
		m.FilePanels[i].UpdateElementsIfNeeded(force, m.DisplayDotFiles)
	}
}

type PanelUpdateMsg struct {
	panelIndex int
	requestID  uint64
	location   filesystem.Location
	elements   []filepanel.Element
	loadedAt   time.Time
	err        error
}

func (m *Model) UpdateLocalFilePanelsIfNeeded(force bool) {
	for i := range m.FilePanels {
		if m.FilePanels[i].CurrentLocation().Provider == filesystem.ProviderLocal {
			m.FilePanels[i].UpdateElementsIfNeeded(force, m.DisplayDotFiles)
		}
	}
}

func (m *Model) GetRemoteFilePanelUpdateCmd(force bool) tea.Cmd {
	now := time.Now()
	commands := make([]tea.Cmd, 0, len(m.FilePanels))
	for i := range m.FilePanels {
		panel := &m.FilePanels[i]
		location := panel.CurrentLocation()
		if location.Provider == filesystem.ProviderLocal {
			continue
		}
		session, ok := m.Sessions.Get(location.SessionID)
		if !ok || session.Status != SessionConnected || session.Browser == nil {
			continue
		}
		requestID, started := panel.BeginElementsLoading(force, now)
		if !started {
			continue
		}
		panelCopy := *panel
		panelIndex := i
		displayDotFiles := m.DisplayDotFiles
		commands = append(commands, func() tea.Msg {
			elements, err := panelCopy.LoadElements(displayDotFiles)
			return PanelUpdateMsg{
				panelIndex: panelIndex,
				requestID:  requestID,
				location:   location,
				elements:   elements,
				loadedAt:   time.Now(),
				err:        err,
			}
		})
	}
	return tea.Batch(commands...)
}

func (m *Model) ApplyPanelUpdate(msg PanelUpdateMsg) tea.Cmd {
	if msg.panelIndex < 0 || msg.panelIndex >= len(m.FilePanels) {
		return nil
	}
	panel := &m.FilePanels[msg.panelIndex]
	accepted, refreshPending := panel.FinishElementsLoading(msg.requestID)
	if !accepted {
		return nil
	}
	if panel.CurrentLocation() != msg.location {
		return nil
	}
	if msg.err != nil {
		slog.Error("Error while loading remote folder elements", "error", msg.err, "location", panel.DisplayLocation())
		if errors.Is(msg.err, filesystem.ErrDisconnected) {
			if err := m.MarkSessionDisconnected(msg.location.SessionID, msg.err); err != nil {
				slog.Error("failed to mark remote panel session disconnected", "error", err)
			}
			return nil
		}
		if refreshPending {
			return m.GetRemoteFilePanelUpdateCmd(true)
		}
		return nil
	}
	panel.ApplyLoadedElements(msg.elements, msg.loadedAt)
	if refreshPending {
		return m.GetRemoteFilePanelUpdateCmd(true)
	}
	return nil
}
