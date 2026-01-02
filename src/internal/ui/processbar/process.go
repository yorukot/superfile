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

type OperationType int

const (
	OpCopy OperationType = iota
	OpCut
	OpDelete
	OpCompress
	OpExtract
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

// GetIcon returns the appropriate icon for the operation type
func (op OperationType) GetIcon() string {
	switch op {
	case OpCopy:
		return icon.Copy
	case OpCut:
		return icon.Cut
	case OpDelete:
		return icon.Delete
	case OpCompress:
		return icon.CompressFile
	case OpExtract:
		return icon.ExtractFile
	default:
		return icon.InOperation
	}
}

// GetVerb returns the present tense verb for the operation
func (op OperationType) GetVerb() string {
	switch op {
	case OpCopy:
		return "Copying"
	case OpCut:
		return "Moving"
	case OpDelete:
		return "Deleting"
	case OpCompress:
		return "Compressing"
	case OpExtract:
		return "Extracting"
	default:
		return "Processing"
	}
}

// GetPastVerb returns the past tense verb for the operation
func (op OperationType) GetPastVerb() string {
	switch op {
	case OpCopy:
		return "Copied"
	case OpCut:
		return "Moved"
	case OpDelete:
		return "Deleted"
	case OpCompress:
		return "Compressed"
	case OpExtract:
		return "Extracted"
	default:
		return "Processed"
	}
}

// GetDisplayName returns the appropriate display name for the process
func (p *Process) GetDisplayName() string {
	icon := p.Operation.GetIcon()

	if p.State == InOperation {
		return fmt.Sprintf("%s %s %s", icon, p.Operation.GetVerb(), p.CurrentFile)
	}

	// Process completed (successful, failed, or cancelled)
	if p.Total > 1 {
		return fmt.Sprintf("%s %s %d files", icon, p.Operation.GetPastVerb(), p.Total)
	}
	return fmt.Sprintf("%s %s %s", icon, p.Operation.GetPastVerb(), p.CurrentFile)
}
