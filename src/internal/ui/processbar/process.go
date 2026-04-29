package processbar

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/progress"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// Model for an individual process
// Note : Its size is ~ 800 bytes
type Process struct {
	ID          string
	CurrentFile string
	// TODO : We always want ErrorMsg to be set when State is
	// moved to Cancelled or Failed. To ensure it, we need to only allow state
	// change via helper functions and ask for the  errorMsg
	ErrorMsg  string
	Operation OperationType
	Progress  progress.Model
	State     ProcessState
	Total     int
	Done      int
	DoneTime  time.Time
}

type FileListProcessor func(items []string) (Process, []string)
type ProcessFinalizer func(state ProcessState, reqID int) tea.Msg
type ProcessRunner func(processor FileListProcessor, finalizer ProcessFinalizer, items []string, reqID int) tea.Msg

func NewProcess(id string, currentFile string, operation OperationType, total int) Process {
	prog := progress.New(
		progress.WithColors(
			lipgloss.Color(common.Theme.GradientColor[0]),
			lipgloss.Color(common.Theme.GradientColor[1])),
		progress.WithScaled(true))
	prog.PercentageStyle = common.FooterStyle
	return Process{
		ID:          id,
		CurrentFile: currentFile,
		Operation:   operation,
		Progress:    prog,
		State:       InOperation,
		Total:       total,
		Done:        0,
	}
}

type ProcessState int

const (
	InOperation ProcessState = iota
	Successful
	Cancelled
	Failed
)

// TODO : Should we store in a global map for efficiency ? At least need to prerender
// Yes, this is a Render() call, which is expensive
func (p ProcessState) Icon() string {
	switch p {
	case Failed:
		return common.ProcessErrorStyle.Render(icon.Warn)
	case Successful:
		return common.ProcessSuccessfulStyle.Render(icon.Done)
	case InOperation:
		return common.ProcessInOperationStyle.Render(icon.InOperation)
	case Cancelled:
		fallthrough
	default:
		return common.ProcessCancelStyle.Render(icon.Error)
	}
}

// GetDisplayName returns the appropriate display name for the process
func (p *Process) GetDisplayName() string {
	return p.Operation.GetIcon() + icon.Space + p.displayNameWithoutIcon()
}

func (p *Process) displayNameWithoutIcon() string {
	if p.State == Cancelled {
		return p.Operation.GetVerb() + " cancelled : " + p.ErrorMsg
	}
	if p.State == Failed {
		return p.Operation.GetVerb() + " failed : " + p.ErrorMsg
	}

	if p.State == InOperation {
		return p.Operation.GetVerb() + " " + p.CurrentFile
	}

	if p.Total > 1 {
		return fmt.Sprintf("%s %d files", p.Operation.GetPastVerb(), p.Total)
	}
	return p.Operation.GetPastVerb() + " " + p.CurrentFile
}
