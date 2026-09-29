package spferror

import (
	"slices"
	"strings"

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

func Skip(fles *FileListErrorState, runner processbar.ProcessRunner, reqID int) tea.Msg {
	if len(fles.fileList) <= 1 {
		return Abort(fles, runner, reqID)
	}
	return runner(fles.continuationFun, fles.finalizer, fles.fileList[1:], reqID)
}

func Abort(fles *FileListErrorState, runner processbar.ProcessRunner, reqID int) tea.Msg {
	return runner(fles.continuationFun, fles.finalizer, []string{}, reqID)
}

func KeySkip() []string {
	return common.Hotkeys.ConfirmTyping
}

func KeyAbort() []string {
	return common.Hotkeys.Quit
}

type UserActionType int
type ActionKeyChecker func(key string) bool
type Action func(fles *FileListErrorState, runner processbar.ProcessRunner, reqID int) tea.Msg

type UserAction struct {
	Title      string
	KeyChecker ActionKeyChecker
	Run        Action
}

func SkipAction() *UserAction {
	return &UserAction{
		Title:      " (" + KeySkip()[0] + ") Skip ",
		KeyChecker: func(msg string) bool { return slices.Contains(KeySkip(), msg) },
		Run:        Skip,
	}
}

func AbortAction() *UserAction {
	return &UserAction{
		Title:      " (" + KeyAbort()[0] + ") Abort ",
		KeyChecker: func(msg string) bool { return slices.Contains(KeyAbort(), msg) },
		Run:        Abort,
	}
}

func OkAction() *UserAction {
	return &UserAction{
		Title:      common.ModalOkayInputText,
		KeyChecker: func(msg string) bool { return slices.Contains(common.Hotkeys.ConfirmTyping, msg) },
		Run:        Abort,
	}
}

type Model struct {
	open    bool
	title   string
	content string
	state   *FileListErrorState
	actions []*UserAction
}

func New(open bool, title string, content string, state *FileListErrorState, actions []*UserAction) Model {
	return Model{
		open:    open,
		title:   title,
		content: content,
		state:   state,
		actions: actions,
	}
}

func (m *Model) GetActions() []*UserAction {
	return m.actions
}

func (m *Model) IsOpen() bool {
	return m.open
}

func (m *Model) Open() {
	m.open = true
}

func (m *Model) Close() (*FileListErrorState, []*UserAction) {
	m.open = false
	tmpState := m.state
	m.state = nil
	tmpActions := m.actions
	m.actions = nil
	return tmpState, tmpActions
}

func (m *Model) State() *FileListErrorState {
	return m.state
}

func (m *Model) Render() string {
	// TODO: needs "skip all" and "retry" buttons
	var buttonTitles = make([]string, 0, len(m.actions))
	for _, v := range m.actions {
		buttonTitles = append(buttonTitles, common.ModalConfirm.Render(v.Title))
	}
	tip := strings.Join(buttonTitles, common.ModalInputSpacingText)

	var errHeader = common.ModalErrorStyle.Render("Error")
	return common.ModalBorderStyle(common.ModalHeight, common.ModalWidth).
		Render(errHeader + "\n" + m.content + "\n\n" + tip)
}
