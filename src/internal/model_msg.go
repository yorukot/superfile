package internal

import (
	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/yorukot/superfile/src/internal/ui/metadata"
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
	state processState
}

func NewPasteOperationMsg(state processState, reqID int) PasteOperationMsg {
	return PasteOperationMsg{
		state: state,
		BaseMessage: BaseMessage{
			reqID: reqID,
		},
	}
}

func (msg PasteOperationMsg) ApplyToModel(m *model) tea.Cmd {
	if (msg.state == failure || msg.state == successful) && m.copyItems.cut {
		m.copyItems.reset(false)
	}
	return nil
}

type DeleteOperationMsg struct {
	BaseMessage
	state processState
}

func NewDeleteOperationMsg(state processState, reqID int) DeleteOperationMsg {
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

type MetadataMsg struct {
	// Using struct embedding over composition, because behaviour with GetReqID will not change
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
