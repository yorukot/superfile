# Process Bar Multi-File Operations Improvement Plan

## Issue #883: Fix process name display for multi-file operations

### Problem Statement
When processing multiple files (copy/cut/delete/compress/extract), only the last selected file name is shown in the process bar. Users should see:
- During execution: Current file being processed (e.g., "Copying file123.txt")
- After completion: Summary for multi-files (e.g., "Copied 5 files" or "Deleted 10 items")

## Implementation Steps

### Step 1: Add OperationType Enum
**File**: `src/internal/ui/processbar/process.go`

Add after line 35 (after `type ProcessState int`):
```go
type OperationType int

const (
	OpCopy OperationType = iota
	OpCut
	OpDelete
	OpCompress
	OpExtract
)
```

### Step 2: Update Process Struct
**File**: `src/internal/ui/processbar/process.go`

Update Process struct (lines 13-21) to:
```go
type Process struct {
	ID           string
	CurrentFile  string        // Renamed from Name - contains only filename, no icons
	Operation    OperationType // NEW - type of operation being performed
	Progress     progress.Model
	State        ProcessState
	Total        int          // Already exists - represents total files/items
	Done         int
	DoneTime     time.Time
}
```

Update NewProcess function signature (lines 23-34) to:
```go
func NewProcess(id string, currentFile string, operation OperationType, total int) Process {
	prog := progress.New(common.GenerateGradientColor())
	prog.PercentageStyle = common.FooterStyle
	return Process{
		ID:          id,
		CurrentFile: currentFile,  // Just the filename, no icons
		Operation:   operation,
		Progress:    prog,
		State:       InOperation,
		Total:       total,
		Done:        0,
	}
}
```

### Step 3: Add Helper Methods for Operations
**File**: `src/internal/ui/processbar/process.go`

Add after the Icon() method (around line 55):
```go
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
		return icon.Compress
	case OpExtract:
		return icon.Extract
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
		return fmt.Sprintf("%s%s%s %s", icon, icon.Space, p.Operation.GetVerb(), p.CurrentFile)
	}
	
	if p.Total > 1 {
		verb := p.Operation.GetPastVerb()
		return fmt.Sprintf("%s%s%s %d files", icon, icon.Space, verb, p.Total)
	}
	
	verb := p.Operation.GetPastVerb()
	return fmt.Sprintf("%s%s%s %s", icon, icon.Space, verb, p.CurrentFile)
}
```

Add import at the top of the file:
```go
import (
	"fmt"
	"time"
	"github.com/yorukot/superfile/src/config/icon"
	// ... other imports
)
```

### Step 4: Update Model's SendAddProcessMsg
**File**: `src/internal/ui/processbar/model.go`

Find SendAddProcessMsg method (around line 60-70) and update:
```go
func (m *Model) SendAddProcessMsg(currentFile string, operation OperationType, total int, blocking bool) (Process, error) {
	// ... existing code ...
	p := NewProcess(id, currentFile, operation, total)
	// ... rest of the method
}
```

### Step 5: Update View Method to Use GetDisplayName
**File**: `src/internal/ui/processbar/model.go`

In the View() method (around line 120-150), find where process name is displayed and update:
```go
// Find line with truncateText(p.Name, nameWidth) and change to:
processName := truncateText(p.GetDisplayName(), nameWidth)
```

### Step 6: Update Paste Operation
**File**: `src/internal/handle_file_operations.go`

Update executePasteOperation (line ~276-278):
```go
var opType processbar.OperationType
if cut {
	opType = processbar.OpCut
} else {
	opType = processbar.OpCopy
}

currentFileName := filepath.Base(copyItems[0])
if len(copyItems) > 1 {
	currentFileName = fmt.Sprintf("%s (%d files)", currentFileName, len(copyItems))
}

p, err := processBarModel.SendAddProcessMsg(
	currentFileName,
	opType,
	getTotalFilesCnt(copyItems),
	true)
```

Update the loop (line ~296):
```go
p.CurrentFile = filepath.Base(filePath)
```

### Step 7: Update Delete Operation
**File**: `src/internal/handle_file_operations.go`

Update deleteOperation function (line ~187):
```go
currentFileName := filepath.Base(items[0])
if len(items) > 1 {
	currentFileName = fmt.Sprintf("%s (%d files)", currentFileName, len(items))
}

p, err := processBarModel.SendAddProcessMsg(
	currentFileName,
	processbar.OpDelete,
	len(items),
	true)
```

Update in the loop (line ~200):
```go
p.CurrentFile = filepath.Base(item)
```

### Step 8: Update Compress Operation
**File**: `src/internal/handle_modal.go`

Find compress operation (search for "NewCompressOperationMsg") and update:
```go
p, _ := m.processBarModel.SendAddProcessMsg(
	archiveName,
	processbar.OpCompress,
	len(itemsToCompress),
	true)
```

When updating progress in the loop:
```go
p.CurrentFile = filepath.Base(currentItem)
```

### Step 9: Update Extract Operation  
**File**: `src/internal/handle_modal.go`

Find extract operation (search for "NewExtractOperationMsg") and update similarly:
```go
p, _ := m.processBarModel.SendAddProcessMsg(
	archiveName,
	processbar.OpExtract,
	totalFiles,
	true)
```

When updating progress:
```go
p.CurrentFile = filepath.Base(currentExtractedFile)
```

### Step 10: Fix All References to p.Name or Process.Name
Search and replace all occurrences of `p.Name` with `p.CurrentFile` in:
- `src/internal/handle_file_operations.go`
- `src/internal/handle_modal.go`
- Any other files that reference Process.Name

Use command:
```bash
grep -r "p\.Name" src/ --include="*.go" | grep -v "_test.go"
grep -r "Process.*Name" src/ --include="*.go" | grep -v "_test.go"
```

Also search for any direct assignments that include icons:
```bash
grep -r "icon\.[A-Z].*+.*icon\.Space" src/ --include="*.go" | grep "p\."
```
These should be changed to only assign the filename without icons.

### Step 11: Update Tests
**File**: `src/internal/model_layout_test.go`

Update test that creates processes (line with `processbar.NewProcess`):
```go
processbar.NewProcess(strconv.Itoa(i), "test.txt", processbar.OpCopy, 1)
```

### Step 12: Add Unit Tests
Create new test file `src/internal/ui/processbar/process_test.go`:
```go
package processbar

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/yorukot/superfile/src/config/icon"
)

func TestGetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		process  Process
		expected string
	}{
		{
			name: "Single file during operation",
			process: Process{
				CurrentFile: "file.txt",
				Operation:   OpCopy,
				Total:       1,
				State:       InOperation,
			},
			expected: icon.Copy + icon.Space + "Copying file.txt",
		},
		{
			name: "Multiple files during operation",
			process: Process{
				CurrentFile: "file2.txt",
				Operation:   OpDelete,
				Total:       10,
				State:       InOperation,
			},
			expected: icon.Delete + icon.Space + "Deleting file2.txt",
		},
		{
			name: "Multiple files after completion",
			process: Process{
				CurrentFile: "file.txt",
				Operation:   OpCopy,
				Total:       5,
				State:       Successful,
			},
			expected: icon.Copy + icon.Space + "Copied 5 files",
		},
		{
			name: "Single file after completion",
			process: Process{
				CurrentFile: "file.txt",
				Operation:   OpDelete,
				Total:       1,
				State:       Successful,
			},
			expected: icon.Delete + icon.Space + "Deleted file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.process.GetDisplayName())
		})
	}
}
```

## Testing Checklist

After implementation, test these scenarios:
1. [ ] Copy single file - shows "Copying file.txt" during, "Copied file.txt" after
2. [ ] Copy multiple files - shows "Copying file.txt" during, "Copied N files" after
3. [ ] Cut/Move single file - shows "Moving file.txt" during, "Moved file.txt" after
4. [ ] Cut/Move multiple files - shows "Moving file.txt" during, "Moved N files" after
5. [ ] Delete single file - shows "Deleting file.txt" during, "Deleted file.txt" after
6. [ ] Delete multiple files - shows "Deleting file.txt" during, "Deleted N files" after
7. [ ] Compress files - shows "Compressing archive.zip" during, "Compressed archive.zip" after
8. [ ] Extract archive - shows "Extracting archive.zip" during, "Extracted archive.zip" after
9. [ ] Failed operations show correct message with icon
10. [ ] Verify process bar updates in real-time during operations

## Build and Lint Commands
```bash
# Build
go build -o bin/spf

# Run tests
go test ./src/internal/ui/processbar/...
go test ./src/internal/...

# Lint
golangci-lint run

# Format
go fmt ./...
```

## Key Principles
- **CurrentFile field**: Always contains just the filename, no icons or verbs
- **GetDisplayName() method**: Handles all formatting logic including icons and verbs
- **Icons**: Automatically fetched based on OperationType via GetIcon()
- **Verbs**: Use GetVerb() for present tense, GetPastVerb() for past tense
- **Total field**: Reuse existing field instead of adding new FileCount
- **No manual concatenation**: Never manually build display strings with icons - let GetDisplayName() handle it
