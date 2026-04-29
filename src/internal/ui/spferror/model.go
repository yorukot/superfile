package spferror

import (
	"github.com/yorukot/superfile/src/internal/common"
	processbar "github.com/yorukot/superfile/src/internal/ui/processbar"

	tea "charm.land/bubbletea/v2"
)

type FileListErrorState struct {
	fileList        []string
	continuationFun processbar.FileListProcessor
	finalizer       processbar.ProcessFinalizer
}

func NewFileListError(fileList []string,
	continuationFun processbar.FileListProcessor,
	finalizer processbar.ProcessFinalizer) *FileListErrorState {
	return &FileListErrorState{fileList: fileList, continuationFun: continuationFun, finalizer: finalizer}
}

func (fles *FileListErrorState) Skip(runner processbar.ProcessRunner, reqID int) tea.Msg {
	if len(fles.fileList) <= 1 {
		return fles.Abort(runner, reqID)
	}
	return runner(fles.continuationFun, fles.finalizer, fles.fileList[1:], reqID)
}

func (fles *FileListErrorState) Abort(runner processbar.ProcessRunner, reqID int) tea.Msg {
	return runner(fles.continuationFun, fles.finalizer, []string{}, reqID)
}

type Model struct {
	open    bool
	title   string
	content string
	state   *FileListErrorState
}

func New(open bool, title string, content string, state *FileListErrorState) Model {
	return Model{
		open:    open,
		title:   title,
		content: content,
		state:   state,
	}
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Open() {
	m.open = true
}

func (m *Model) Close() *FileListErrorState {
	m.open = false
	tmpState := m.state
	m.state = nil
	return tmpState
}

func (m *Model) State() *FileListErrorState {
	return m.state
}

func KeySkip() []string {
	return common.Hotkeys.ConfirmTyping
}

func KeyAbort() []string {
	return common.Hotkeys.Quit
}

func (m *Model) Render() string {
	// TODO: needs "skip all" and "retry" buttons
	skip := common.ModalConfirm.Render(" (" + KeySkip()[0] + ") Skip ")
	abort := common.ModalCancel.Render(" (" + KeyAbort()[0] + ") Abort ")

	tip := skip + common.ModalInputSpacingText + abort

	var errHeader = common.ModalErrorStyle.Render("Error")
	return common.ModalBorderStyle(common.ModalHeight, common.ModalWidth).
		Render(errHeader + "\n" + m.content + "\n\n" + tip)
}
