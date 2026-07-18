package internal

import (
	"log/slog"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/filesystem"
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/spferror"
)

type ModelUpdateMessage interface {
	ApplyToModel(m *model) tea.Cmd
	GetReqID() int
}

type BaseMessage struct {
	reqID int
}

type FileMutationKind string

const (
	FileMutationCreate FileMutationKind = "create"
	FileMutationRename FileMutationKind = "rename"
)

type FileMutationMsg struct {
	BaseMessage

	kind       FileMutationKind
	location   filesystem.Location
	panelIndex int
	conflict   bool
	err        error
}

func NewFileMutationMsg(
	kind FileMutationKind,
	location filesystem.Location,
	panelIndex int,
	conflict bool,
	err error,
	reqID int,
) FileMutationMsg {
	return FileMutationMsg{
		BaseMessage: BaseMessage{reqID: reqID},
		kind:        kind,
		location:    location,
		panelIndex:  panelIndex,
		conflict:    conflict,
		err:         err,
	}
}

func (msg FileMutationMsg) ApplyToModel(m *model) tea.Cmd {
	if msg.err != nil {
		m.handleRemoteSessionError(msg.location, msg.err)
	}

	switch msg.kind {
	case FileMutationCreate:
		if msg.err != nil {
			m.notifyModel = notify.New(true, "Create failed", msg.err.Error(), notify.NoAction)
		}
	case FileMutationRename:
		m.renameOperationPending = false
		if msg.conflict && msg.err == nil {
			m.notifyModel = notify.New(
				true,
				common.SameRenameWarnTitle,
				common.SameRenameWarnContent,
				notify.RenameAction,
			)
			return nil
		}
		m.resetRenameState(msg.panelIndex)
		if msg.err != nil {
			m.notifyModel = notify.New(true, "Rename failed", msg.err.Error(), notify.NoAction)
		}
	}

	if msg.location.Provider == filesystem.ProviderLocal {
		m.fileModel.UpdateLocalFilePanelsIfNeeded(true)
		return nil
	}
	return m.fileModel.GetRemoteFilePanelUpdateCmd(true)
}

type RemoteNavigationMsg struct {
	BaseMessage

	panelIndex int
	source     filesystem.Location
	target     filesystem.Path
	generation uint64
	elements   []filepanel.Element
	loadedAt   time.Time
	err        error
}

func NewRemoteNavigationMsg(
	panelIndex int,
	source filesystem.Location,
	target filesystem.Path,
	generation uint64,
	elements []filepanel.Element,
	loadedAt time.Time,
	err error,
	reqID int,
) RemoteNavigationMsg {
	return RemoteNavigationMsg{
		BaseMessage: BaseMessage{reqID: reqID},
		panelIndex:  panelIndex,
		source:      source,
		target:      target,
		generation:  generation,
		elements:    elements,
		loadedAt:    loadedAt,
		err:         err,
	}
}

func (msg RemoteNavigationMsg) ApplyToModel(m *model) tea.Cmd {
	if msg.panelIndex < 0 || msg.panelIndex >= len(m.fileModel.FilePanels) {
		return nil
	}
	panel := &m.fileModel.FilePanels[msg.panelIndex]
	if panel.CurrentLocation() != msg.source {
		return nil
	}
	if msg.err != nil {
		m.handleRemoteSessionErrorIfCurrent(msg.source, msg.generation, msg.err)
		m.notifyModel = notify.New(true, "Remote navigation failed", msg.err.Error(), notify.NoAction)
		return nil
	}
	if err := panel.ApplyCurrentFilePanelDir(msg.target.String()); err != nil {
		m.notifyModel = notify.New(true, "Remote navigation failed", err.Error(), notify.NoAction)
		return nil
	}
	panel.ApplyLoadedElements(msg.elements, msg.loadedAt)
	m.fileModel.SyncPaneSessionLocations()
	return nil
}

func (msg BaseMessage) GetReqID() int {
	return msg.reqID
}

type PasteOperationMsg struct {
	BaseMessage

	state            processbar.ProcessState
	failureLocation  filesystem.Location
	refreshLocations []filesystem.Location
	remainingSources []filesystem.Location
	err              error
}

func NewProviderPasteOperationMsg(
	state processbar.ProcessState,
	failureLocation filesystem.Location,
	refreshLocations []filesystem.Location,
	remainingSources []filesystem.Location,
	err error,
	reqID int,
) PasteOperationMsg {
	return PasteOperationMsg{
		state:            state,
		failureLocation:  failureLocation,
		refreshLocations: append([]filesystem.Location(nil), refreshLocations...),
		remainingSources: append([]filesystem.Location(nil), remainingSources...),
		err:              err,
		BaseMessage:      BaseMessage{reqID: reqID},
	}
}

func NewPasteOperationMsg(state processbar.ProcessState, reqID int) PasteOperationMsg {
	return PasteOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg PasteOperationMsg) ApplyToModel(m *model) tea.Cmd {
	refreshRemote := false
	if msg.err != nil && msg.failureLocation.Provider != "" {
		m.handleRemoteSessionError(msg.failureLocation, msg.err)
	}
	for _, location := range msg.refreshLocations {
		refreshRemote = refreshRemote || location.Provider != filesystem.ProviderLocal
	}
	if m.clipboard.IsCut() && len(msg.remainingSources) > 0 {
		m.clipboard.SetLocations(msg.remainingSources)
	} else if msg.state == processbar.Successful && m.clipboard.IsCut() {
		m.clipboard.Reset(false)
	}
	if refreshRemote {
		return m.fileModel.GetRemoteFilePanelUpdateCmd(true)
	}
	return nil
}

type CreateOperationMsg struct {
	BaseMessage

	state    processbar.ProcessState
	location filesystem.Location
}

func NewProviderCreateOperationMsg(
	state processbar.ProcessState,
	location filesystem.Location,
	reqID int,
) CreateOperationMsg {
	return CreateOperationMsg{
		state:       state,
		location:    location,
		BaseMessage: BaseMessage{reqID: reqID},
	}
}

func NewCreateOperationMsg(state processbar.ProcessState, reqID int) CreateOperationMsg {
	return CreateOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg CreateOperationMsg) ApplyToModel(m *model) tea.Cmd {
	m.typingModal.submitting = false
	if msg.location.Provider == filesystem.ProviderLocal {
		m.fileModel.UpdateLocalFilePanelsIfNeeded(true)
		return nil
	}
	if msg.location.Provider != "" {
		return m.fileModel.GetRemoteFilePanelUpdateCmd(true)
	}
	return nil
}

type DeleteOperationMsg struct {
	BaseMessage

	state    processbar.ProcessState
	location filesystem.Location
	err      error
}

func NewProviderDeleteOperationMsg(
	state processbar.ProcessState,
	location filesystem.Location,
	err error,
	reqID int,
) DeleteOperationMsg {
	return DeleteOperationMsg{
		state:       state,
		location:    location,
		err:         err,
		BaseMessage: BaseMessage{reqID: reqID},
	}
}

func NewDeleteOperationMsg(state processbar.ProcessState, reqID int) DeleteOperationMsg {
	return DeleteOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg DeleteOperationMsg) ApplyToModel(m *model) tea.Cmd {
	m.handleRemoteSessionError(msg.location, msg.err)
	// Remove selection
	m.getFocusedFilePanel().ResetSelected()
	if msg.location.Provider != "" && msg.location.Provider != filesystem.ProviderLocal {
		return m.fileModel.GetRemoteFilePanelUpdateCmd(true)
	}
	return nil
}

type ProcessBarUpdateMsg struct {
	BaseMessage

	pMsg processbar.UpdateMsg
}

func (msg ProcessBarUpdateMsg) ApplyToModel(m *model) tea.Cmd {
	cmd, err := msg.pMsg.Apply(&m.processBarModel)
	if err != nil {
		slog.Error("Error applying processbar update", "error", err)
	}
	return processCmdToTeaCmd(cmd)
}

type CompressOperationMsg struct {
	BaseMessage

	state processbar.ProcessState
}

func NewCompressOperationMsg(state processbar.ProcessState, reqID int) CompressOperationMsg {
	return CompressOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg CompressOperationMsg) ApplyToModel(_ *model) tea.Cmd {
	return nil
}

type ExtractOperationMsg struct {
	BaseMessage

	state processbar.ProcessState
}

func NewExtractOperationMsg(state processbar.ProcessState, reqID int) ExtractOperationMsg {
	return ExtractOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg ExtractOperationMsg) ApplyToModel(_ *model) tea.Cmd {
	return nil
}

type MetadataMsg struct {
	BaseMessage

	meta            metadata.Metadata
	metadataFocused bool
}

func NewMetadataMsg(meta metadata.Metadata, metadataFocused bool, reqID int) MetadataMsg {
	return MetadataMsg{
		meta:            meta,
		metadataFocused: metadataFocused,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg MetadataMsg) ApplyToModel(m *model) tea.Cmd {
	m.fileMetaData.SetMetadataCache(msg.meta, msg.metadataFocused)
	selectedItem := m.getFocusedFilePanel().GetFocusedItemPtr()
	if selectedItem == nil {
		slog.Debug("Panel empty or cursor invalid. Ignoring MetadataMsg")
		return nil
	}
	if selectedItem.Location != msg.meta.GetPath() {
		slog.Debug("MetadataMsg for older files. Ignoring",
			"currentItem", selectedItem.Location, "msgItem", msg.meta.GetPath())
		return nil
	}
	if (m.focusPanel == metadataFocus) != msg.metadataFocused {
		slog.Debug("MetadataMsg for older state. Ignoring",
			"actualFocus", m.focusPanel, "msgFocus", msg.metadataFocused)
		return nil
	}
	m.fileMetaData.SetMetadata(msg.meta, msg.metadataFocused)
	return nil
}

type NotifyModalUpdateMsg struct {
	BaseMessage

	m notify.Model
}

func NewNotifyModalMsg(m notify.Model, reqID int) NotifyModalUpdateMsg {
	return NotifyModalUpdateMsg{
		m: m,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg NotifyModalUpdateMsg) ApplyToModel(m *model) tea.Cmd {
	m.notifyModel = msg.m
	return nil
}

type SpfErrorModalUpdateMsg struct {
	BaseMessage

	m spferror.Model
}

func NewSpfErrorModalMsg(m spferror.Model, reqID int) SpfErrorModalUpdateMsg {
	return SpfErrorModalUpdateMsg{
		m: m,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg SpfErrorModalUpdateMsg) ApplyToModel(m *model) tea.Cmd {
	m.typingModal.submitting = false
	m.spfError = msg.m
	return nil
}
