package processbar

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"

	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// Model for an individual process
// Note : Its size is ~ 800 bytes
type Process struct {
	ID          string
	CurrentFile string
	ErrorMsg    string
	Operation   OperationType
	Progress    progress.Model
	State       ProcessState
	Total       int
	Done        int
	DoneTime    time.Time
}

func NewProcess(id string, currentFile string, operation OperationType, total int) Process {
	prog := progress.New(common.GenerateGradientColor())
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
	if p.ErrorMsg != "" {
		return "Unexpected failure: " + p.ErrorMsg
	}

	if p.State == InOperation {
		return p.Operation.GetVerb() + " " + p.CurrentFile
	}

	// Process completed (successful, failed, or cancelled)
	if p.Total > 1 {
		return fmt.Sprintf("%s %d files", p.Operation.GetPastVerb(), p.Total)
	}
	return p.Operation.GetPastVerb() + " " + p.CurrentFile
}
