//go:build linux || darwin

package common

import (
	"fmt"
	"testing"
)

func TestFileNameWithoutExtensionUnix(t *testing.T) {
	var inputs = []struct {
		input    string
		expected string
	}{
		{"/home/user/temp/.tmp/.dockerignore.zip", "/home/user/temp/.tmp/.dockerignore"},
		{"/home/user/temp/.tmp/.dockerignore", "/home/user/temp/.tmp/.dockerignore"},
		{"/tmp/aaa.bbb/file", "/tmp/aaa.bbb/file"},
		{"/tmp/aaa.bbb/.file", "/tmp/aaa.bbb/.file"},
		{"/tmp/aaa.bbb/.file.txt", "/tmp/aaa.bbb/.file"},
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
