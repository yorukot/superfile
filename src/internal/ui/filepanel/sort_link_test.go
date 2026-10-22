package filepanel

import (
	"os"
	"path/filepath"
	"testing"
)

type irregularFileInfo struct{ os.FileInfo }

func (i irregularFileInfo) Mode() os.FileMode { return i.FileInfo.Mode() | os.ModeIrregular }

func TestIsLinkedDirectoryWithJunctionMode(t *testing.T) {
	root := t.TempDir()
	directory := filepath.Join(root, "directory")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name string
		want bool
	}{
		{"directory", true},
		{"file", false},
		{"missing", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			info, err := os.Stat(file)
			if err != nil {
				t.Fatal(err)
			}
			if got := isLinkedDirectory(root, irregularFileInfo{info}, tt.name); got != tt.want {
				t.Errorf("isLinkedDirectory(%q) = %v; want %v", tt.name, got, tt.want)
			}
		})
	}
}
