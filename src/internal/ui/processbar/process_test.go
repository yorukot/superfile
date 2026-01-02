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
			expected: icon.Copy + " Copying file.txt",
		},
		{
			name: "Multiple files during operation",
			process: Process{
				CurrentFile: "file2.txt",
				Operation:   OpDelete,
				Total:       10,
				State:       InOperation,
			},
			expected: icon.Delete + " Deleting file2.txt",
		},
		{
			name: "Multiple files after completion",
			process: Process{
				CurrentFile: "file.txt",
				Operation:   OpCopy,
				Total:       5,
				State:       Successful,
			},
			expected: icon.Copy + " Copied 5 files",
		},
		{
			name: "Single file after completion",
			process: Process{
				CurrentFile: "file.txt",
				Operation:   OpDelete,
				Total:       1,
				State:       Successful,
			},
			expected: icon.Delete + " Deleted file.txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.process.GetDisplayName())
		})
	}
}
