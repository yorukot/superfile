package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveSaveChooserStartPath(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := filepath.Join(tempDir, "dir")
	existingFile := filepath.Join(existingDir, "file.txt")
	nonExistentFile := filepath.Join(existingDir, "save.txt")

	require.NoError(t, os.MkdirAll(existingDir, 0o755))
	require.NoError(t, os.WriteFile(existingFile, []byte("data"), 0o644))

	testdata := []struct {
		name           string
		suggestedPath  string
		fallbackDir    string
		expectedDir    string
		expectedTarget string
	}{
		{
			name:           "empty suggestion uses fallback",
			suggestedPath:  "",
			fallbackDir:    tempDir,
			expectedDir:    tempDir,
			expectedTarget: "",
		},
		{
			name:           "existing directory",
			suggestedPath:  existingDir,
			fallbackDir:    tempDir,
			expectedDir:    existingDir,
			expectedTarget: "",
		},
		{
			name:           "existing file",
			suggestedPath:  existingFile,
			fallbackDir:    tempDir,
			expectedDir:    existingDir,
			expectedTarget: "file.txt",
		},
		{
			name:           "non existent file with existing parent",
			suggestedPath:  nonExistentFile,
			fallbackDir:    tempDir,
			expectedDir:    existingDir,
			expectedTarget: "save.txt",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			dir, target := resolveSaveChooserStartPath(tt.suggestedPath, tt.fallbackDir)
			assert.Equal(t, tt.expectedDir, dir)
			assert.Equal(t, tt.expectedTarget, target)
		})
	}
}
