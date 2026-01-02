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
	CurrentFile  string        // Renamed from Name
	Operation    OperationType // NEW
	FileCount    int          // NEW - total files being processed
	Progress     progress.Model
	State        ProcessState
	Total        int
	Done         int
	DoneTime     time.Time
}
```

Update NewProcess function (lines 23-34) to:
```go
func NewProcess(id string, name string, operation OperationType, fileCount int, total int) Process {
	prog := progress.New(common.GenerateGradientColor())
	prog.PercentageStyle = common.FooterStyle
	return Process{
		ID:          id,
		CurrentFile: name,
		Operation:   operation,
		FileCount:   fileCount,
		Progress:    prog,
		State:       InOperation,
		Total:       total,
		Done:        0,
	}
}
```

### Step 3: Add Verb Helper Methods
**File**: `src/internal/ui/processbar/process.go`

Add after the Icon() method (around line 55):
```go
// GetVerb returns the appropriate verb for the operation
// isPast: true for past tense ("Copied"), false for present ("Copying")
func (op OperationType) GetVerb(isPast bool) string {
	switch op {
	case OpCopy:
		if isPast {
			return "Copied"
		}
		return "Copying"
	case OpCut:
		if isPast {
			return "Moved"
		}
		return "Moving"
	case OpDelete:
		if isPast {
			return "Deleted"
		}
		return "Deleting"
	case OpCompress:
		if isPast {
			return "Compressed"
		}
		return "Compressing"
	case OpExtract:
		if isPast {
			return "Extracted"
		}
		return "Extracting"
	default:
		if isPast {
			return "Processed"
		}
		return "Processing"
	}
}

// GetDisplayName returns the appropriate display name for the process
func (p *Process) GetDisplayName() string {
	// During operation, show current file
	if p.State == InOperation {
		return p.CurrentFile
	}
	
	// After completion, show summary for multi-files
	if p.FileCount > 1 {
		verb := p.Operation.GetVerb(true) // past tense
		return fmt.Sprintf("%s %d files", verb, p.FileCount)
	}
	
	// Single file - show filename with verb
	return p.CurrentFile
}
```

Add import at the top of the file:
```go
import (
	"fmt"
	"time"
	// ... other imports
)
```

### Step 4: Update Model's SendAddProcessMsg
**File**: `src/internal/ui/processbar/model.go`

Find SendAddProcessMsg method (around line 60-70) and update signature:
```go
func (m *Model) SendAddProcessMsg(name string, operation OperationType, fileCount int, totalOperations int, blocking bool) (Process, error) {
	// ... existing code ...
	p := NewProcess(id, name, operation, fileCount, totalOperations)
	// ... rest of the method
}
```

### Step 5: Update View Method to Use GetDisplayName
**File**: `src/internal/ui/processbar/model.go`

In the View() method (around line 120-150), find where process name is displayed and update:
```go
// Replace p.Name with p.GetDisplayName()
// The exact line will be something like:
// processName := truncateText(p.Name, nameWidth)
// Change to:
processName := truncateText(p.GetDisplayName(), nameWidth)
```

### Step 6: Update Paste Operation
**File**: `src/internal/handle_file_operations.go`

Update executePasteOperation (line ~276-278):
```go
// Determine operation type
var opType processbar.OperationType
if cut {
	opType = processbar.OpCut
} else {
	opType = processbar.OpCopy
}

// Create initial display name
initialName := icon.GetCopyOrCutIcon(cut) + icon.Space
if len(copyItems) > 1 {
	initialName += fmt.Sprintf("%s %d files", opType.GetVerb(false), len(copyItems))
} else {
	initialName += filepath.Base(copyItems[0])
}

p, err := processBarModel.SendAddProcessMsg(
	initialName,
	opType,
	len(copyItems), // fileCount
	getTotalFilesCnt(copyItems),
	true)
```

Update the loop (line ~296):
```go
// Update current file being processed
p.CurrentFile = icon.GetCopyOrCutIcon(cut) + icon.Space + filepath.Base(filePath)
```

### Step 7: Update Delete Operation
**File**: `src/internal/handle_file_operations.go`

Update deleteOperation function (line ~187):
```go
// Create initial display name
initialName := icon.Delete + icon.Space
if len(items) > 1 {
	initialName += fmt.Sprintf("Deleting %d files", len(items))
} else {
	initialName += filepath.Base(items[0])
}

p, err := processBarModel.SendAddProcessMsg(
	initialName,
	processbar.OpDelete,
	len(items), // fileCount
	len(items),
	true)
```

Update in the loop (line ~200):
```go
p.CurrentFile = icon.Delete + icon.Space + filepath.Base(item)
```

### Step 8: Update Compress Operation
**File**: `src/internal/handle_modal.go`

Find compress operation (search for "NewCompressOperationMsg") and update:
1. Add operation type when creating process
2. Update to use CurrentFile instead of Name
3. Pass fileCount as number of items being compressed

Example:
```go
p, _ := m.processBarModel.SendAddProcessMsg(
	icon.Compress + icon.Space + "Creating " + archiveName,
	processbar.OpCompress,
	len(itemsToCompress), // fileCount
	totalCount,
	true)
```

### Step 9: Update Extract Operation  
**File**: `src/internal/handle_modal.go`

Find extract operation (search for "NewExtractOperationMsg") and update similarly:
```go
p, _ := m.processBarModel.SendAddProcessMsg(
	icon.Extract + icon.Space + "Extracting " + archiveName,
	processbar.OpExtract,
	1, // single archive file
	totalFiles,
	true)
```

### Step 10: Fix All References to p.Name
Search and replace all occurrences of `p.Name` with `p.CurrentFile` in:
- `src/internal/handle_file_operations.go`
- `src/internal/handle_modal.go`
- Any other files that reference Process.Name

Use command:
```bash
grep -r "p\.Name" src/ --include="*.go" | grep -v "_test.go"
```

### Step 11: Update Tests
**File**: `src/internal/model_layout_test.go`

Update test that creates processes (line with `processbar.NewProcess`):
```go
processbar.NewProcess(strconv.Itoa(i), "test", processbar.OpCopy, 1, 1)
```

### Step 12: Add Unit Tests
Create new test file `src/internal/ui/processbar/process_test.go`:
```go
package processbar

import (
	"testing"
	"github.com/stretchr/testify/assert"
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
				CurrentFile: "Copying file.txt",
				Operation:   OpCopy,
				FileCount:   1,
				State:       InOperation,
			},
			expected: "Copying file.txt",
		},
		{
			name: "Multiple files after completion",
			process: Process{
				CurrentFile: "file.txt",
				Operation:   OpCopy,
				FileCount:   5,
				State:       Successful,
			},
			expected: "Copied 5 files",
		},
		{
			name: "Single file after completion",
			process: Process{
				CurrentFile: "Deleted file.txt",
				Operation:   OpDelete,
				FileCount:   1,
				State:       Successful,
			},
			expected: "Deleted file.txt",
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
1. [ ] Copy single file - shows filename during and after
2. [ ] Copy multiple files - shows "Copying N files" then "Copied N files"
3. [ ] Cut/Move single file - shows filename during and after
4. [ ] Cut/Move multiple files - shows "Moving N files" then "Moved N files"
5. [ ] Delete single file - shows filename during and after
6. [ ] Delete multiple files - shows "Deleting N files" then "Deleted N files"
7. [ ] Compress files - shows progress and completion message
8. [ ] Extract archive - shows progress and completion message
9. [ ] Failed operations show correct message
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

## Notes
- Always use `p.CurrentFile` instead of `p.Name` after refactoring
- The GetDisplayName() method handles the logic for what to show based on state and file count
- Import "fmt" package where needed for Sprintf
- Maintain icon prefixes (icon.Delete, icon.Compress, etc.) for visual consistency
