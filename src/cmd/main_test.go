package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
)

func TestBuildChooserRequest(t *testing.T) {
	originalWD, err := os.Getwd()
	require.NoError(t, err)

	tempDir := t.TempDir()
	require.NoError(t, os.Chdir(tempDir))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(originalWD))
	})

	testdata := []struct {
		name           string
		firstPanelPath []string
		chooserFile    string
		saveFile       string
		expected       variable.ChooserRequest
		wantErr        bool
	}{
		{
			name:           "open chooser with relative path",
			firstPanelPath: []string{"relative/file.txt"},
			chooserFile:    "out.txt",
			expected: variable.ChooserRequest{
				Mode:          variable.ChooserModeOpen,
				OutputFile:    "out.txt",
				SuggestedPath: filepath.Join(tempDir, "relative", "file.txt"),
			},
		},
		{
			name:           "save chooser with no startup path",
			firstPanelPath: []string{""},
			saveFile:       "save-out.txt",
			expected: variable.ChooserRequest{
				Mode:       variable.ChooserModeSave,
				OutputFile: "save-out.txt",
			},
		},
		{
			name:           "mutually exclusive flags",
			firstPanelPath: []string{"one"},
			chooserFile:    "open.txt",
			saveFile:       "save.txt",
			wantErr:        true,
		},
		{
			name:           "too many startup paths",
			firstPanelPath: []string{"one", "two"},
			chooserFile:    "open.txt",
			wantErr:        true,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			req, err := buildChooserRequest(tt.firstPanelPath, tt.chooserFile, tt.saveFile)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, req)
		})
	}
}
