//go:build windows

package common

import (
	"fmt"
	"testing"
)

func TestFilenameWithouTextWindows(t *testing.T) {
	var inputs = []struct {
		input    string
		expected string
	}{
		{"c:\\home\\user\\temp\\.tmp\\.dockerignore.zip", "c:\\home\\user\\temp\\.tmp\\.dockerignore"},
		{"c:\\home\\user\\temp\\.tmp\\.dockerignore", "c:\\home\\user\\temp\\.tmp\\.dockerignore"},
		{"c:\\tmp\\aaa.bbb\\file", "c:\\tmp\\aaa.bbb\\file"},
		{"c:\\tmp\\aaa.bbb\\.file", "c:\\tmp\\aaa.bbb\\.file"},
		{"c:\\tmp\\aaa.bbb\\.file.txt", "c:\\tmp\\aaa.bbb\\.file"},
	}

	for _, tt := range inputs {
		t.Run(fmt.Sprintf("Remove extension from %s", tt.input), func(t *testing.T) {
			result := FileNameWithoutExtension(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
