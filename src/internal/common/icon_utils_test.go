package common

import (
	"testing"

	"github.com/yorukot/superfile/src/config/icon"
)

func TestGetElementIcon(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		isDir    bool
		isLink   bool
		nerdFont bool
		expected icon.Style
	}{
		{
			name:     "Non-nerdfont returns empty icon",
			file:     "test.txt",
			isDir:    false,
			isLink:   false,
			nerdFont: false,
			expected: icon.Style{
				Icon:  "",
				Color: Theme.FilePanelFG,
			},
		},
		{
			name:     "Directory with nerd font",
			file:     "folder",
			isDir:    true,
			isLink:   false,
			nerdFont: true,
			expected: icon.Folders["folder"],
		},
		{
			name:     "File with known extension",
			file:     "test.js",
			isDir:    false,
			isLink:   false,
			nerdFont: true,
			expected: icon.Icons["js"],
		},
		{
			name:     "Full name takes priority over extension",
			file:     "gulpfile.js",
			isDir:    false,
			isLink:   false,
			nerdFont: true,
			expected: icon.Icons["gulpfile.js"],
		},
		{
			name:     ".git directory",
			file:     ".git",
			isDir:    true,
			isLink:   false,
			nerdFont: true,
			expected: icon.Folders[".git"],
		},
		{
			name:     "superfile directory",
			file:     "superfile",
			isDir:    true,
			isLink:   false,
			nerdFont: true,
			expected: icon.Folders["superfile"],
		},
		{
			name:     "package.json file",
			file:     "package.json",
			isDir:    false,
			isLink:   false,
			nerdFont: true,
			expected: icon.Icons["package"],
		},
		{
			name:     "File with unknown extension",
			file:     "test.xyz",
			isDir:    false,
			isLink:   false,
			nerdFont: true,
			expected: icon.Style{
				Icon: icon.Icons["file"].Icon,
				// Theme is not defined here, so this will be blank
				Color: Theme.FilePanelFG,
			},
		},
		{
			name:     "File with aliased name",
			file:     "dockerfile",
			isDir:    false,
			isLink:   false,
			nerdFont: true,
			expected: icon.Icons["dockerfile"],
		},
		{
			name:     "Link to Directory with nerd font",
			file:     "folder",
			isDir:    true,
			isLink:   true,
			nerdFont: true,
			expected: icon.Folders["link_folder"],
		},
		{
			name:     "Link to File",
			file:     "test.js",
			isDir:    false,
			isLink:   true,
			nerdFont: true,
			expected: icon.Icons["link_file"],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetElementIcon(tt.file, tt.isDir, tt.isLink, tt.nerdFont)
			if result.Icon != tt.expected.Icon || result.Color != tt.expected.Color {
				t.Errorf("GetElementIcon() = %v, want %v", result, tt.expected)
			}
		})
	}
}
