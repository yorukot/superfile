package bulkrename
import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

type RenameType int
type CaseType int

type Model struct {
	open         bool
	renameType   RenameType
	caseType     CaseType
	cursor       int
	startNumber  int
	errorMessage string

	findInput    textinput.Model
	replaceInput textinput.Model
	prefixInput  textinput.Model
	suffixInput  textinput.Model

	preview []RenamePreview

	reqCnt int

	width  int
	height int

	selectedFiles []string
	currentDir    string
}

type RenamePreview struct {
	OldPath string
	OldName string
	NewName string
	Error   string
}

type BulkRenameAction struct {
	Previews []RenamePreview
}

type UpdateMsg struct {
	reqID int
}

type EditorModeAction struct {
	TmpfilePath   string
	Editor        string
	SelectedFiles []string
	CurrentDir    string
}

type BulkRenameResultMsg struct {
	state processbar.ProcessState
	count int
}