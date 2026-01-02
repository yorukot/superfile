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
	// If there's an error message, display it
	if p.ErrorMsg != "" {
		return p.ErrorMsg
	}

	ic := p.Operation.GetIcon()

	if p.State == InOperation {
		return fmt.Sprintf("%s%s%s %s", ic, icon.Space, p.Operation.GetVerb(), p.CurrentFile)
	}

	// Process completed (successful, failed, or cancelled)
	if p.Total > 1 {
		return fmt.Sprintf("%s%s%s %d files", ic, icon.Space, p.Operation.GetPastVerb(), p.Total)
	}
	return fmt.Sprintf("%s%s%s %s", ic, icon.Space, p.Operation.GetPastVerb(), p.CurrentFile)
}
