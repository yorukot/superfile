package processbar

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/yorukot/superfile/src/config/icon"
	"github.com/yorukot/superfile/src/internal/common"
)

// Model for an individual process
// Note : Its size is ~ 800 bytes
type Process struct {
	ID       string
	Name     string
	Progress progress.Model
	State    ProcessState
	Total    int
	Done     int
	DoneTime time.Time
}

func NewProcess(id string, name string, total int) Process {
	prog := progress.New(common.GenerateGradientColor())
	prog.PercentageStyle = common.FooterStyle
	return Process{
		ID:       id,
		Name:     name,
		Progress: prog,
		State:    InOperation,
		Total:    total,
		Done:     0,
	}
}

// Type representing the state of a process
type ProcessState int

// Constants for operation, success, Cancelled, Failed
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
