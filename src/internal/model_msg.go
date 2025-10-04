package internal

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/notify"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

type ModelUpdateMessage interface {
	ApplyToModel(m *model) tea.Cmd
	GetReqID() int
}

type BaseMessage struct {
	reqID int
}

func (msg BaseMessage) GetReqID() int {
	return msg.reqID
}

type PasteOperationMsg struct {
	BaseMessage

	state processbar.ProcessState
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
	if (msg.state == processbar.Failed || msg.state == processbar.Successful) && m.copyItems.cut {
		m.copyItems.reset(false)
	}
	return nil
}

type DeleteOperationMsg struct {
	BaseMessage

	state processbar.ProcessState
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
	// Remove selection
	m.getFocusedFilePanel().resetSelected()
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

	meta metadata.Metadata
}

func NewMetadataMsg(meta metadata.Metadata, reqID int) MetadataMsg {
	return MetadataMsg{
		meta: meta,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg MetadataMsg) ApplyToModel(m *model) tea.Cmd {
	selectedItem := m.getFocusedFilePanel().getSelectedItemPtr()
	if selectedItem == nil {
		slog.Debug("Panel empty or cursor invalid. Ignoring MetadataMsg")
		return nil
	}
	if selectedItem.location != msg.meta.GetPath() {
		slog.Debug("MetadataMsg for older files. Ignoring")
		return nil
	}
	m.fileMetaData.SetMetadata(msg.meta)
	selectedItem.metaData = msg.meta.GetData()
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

type FilePreviewUpdateMsg struct {
	BaseMessage

	location string
	content  string
}

func NewFilePreviewUpdateMsg(location string, content string, reqID int) FilePreviewUpdateMsg {
	return FilePreviewUpdateMsg{
		location: location,
		content:  content,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg FilePreviewUpdateMsg) ApplyToModel(m *model) tea.Cmd {
	selectedItem := m.getFocusedFilePanel().getSelectedItemPtr()
	if selectedItem == nil {
		slog.Debug("Panel empty or cursor invalid. Ignoring FilePreviewUpdateMsg")
		return nil
	}
	if selectedItem.location != msg.location {
		slog.Debug("FilePreviewUpdateMsg for older files. Ignoring")
		return nil
	}
	m.fileModel.filePreview.SetContent(msg.content)
	return nil
}
