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
			name: "Error message displayed",
			process: Process{
				CurrentFile: "file.txt",
				ErrorMsg:    "File already exists",
				Operation:   OpCompress,
				Total:       1,
				State:       Cancelled,
			},
			expected: icon.CompressFile + icon.Space + "Compressing cancelled : File already exists",
		},
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
