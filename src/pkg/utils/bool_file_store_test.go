package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadBoolFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		fileContent  string
		defaultValue bool
		createFile   bool
		expected     bool
	}{
		{
			name:         "file contains true",
			fileContent:  TrueString,
			defaultValue: false,
			createFile:   true,
			expected:     true,
		},
		{
			name:         "file contains false",
			fileContent:  FalseString,
			defaultValue: true,
			createFile:   true,
			expected:     false,
		},
		{
			name:         "file contains invalid value",
			fileContent:  "invalid",
			defaultValue: true,
			createFile:   true,
			expected:     true,
		},
		{
			name:         "file contains invalid value with default false",
			fileContent:  "invalid",
			defaultValue: false,
			createFile:   true,
			expected:     false,
		},
		{
			name:         "file contains TRUE (uppercase)",
			fileContent:  "TRUE",
			defaultValue: false,
			createFile:   true,
			expected:     false, // Should not accept uppercase
		},
		{
			name:         "file contains empty string",
			fileContent:  "",
			defaultValue: true,
			createFile:   true,
			expected:     true,
		},
		{
			name:         "file does not exist - default true",
			fileContent:  "",
			defaultValue: true,
			createFile:   false,
			expected:     true,
		},
		{
			name:         "file does not exist - default false",
			fileContent:  "",
			defaultValue: false,
			createFile:   false,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a unique file path for each test
			filePath := filepath.Join(tempDir, tt.name+".txt")

			// Create and write to the file if needed
			if tt.createFile {
				err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
				require.NoError(t, err)
			}

			// Call the function
			result := ReadBoolFile(filePath, tt.defaultValue)

			// Assert the result
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriteBoolFile(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name  string
		value bool
	}{
		{
			name:  "write true value",
			value: true,
		},
		{
			name:  "write false value",
			value: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a unique file path for each test
			filePath := filepath.Join(tempDir, tt.name+".txt")

			// Call the function
			err := WriteBoolFile(filePath, tt.value)
			require.NoError(t, err)

			// Verify file content
			content, err := os.ReadFile(filePath)
			require.NoError(t, err)

			expected := FalseString
			if tt.value {
				expected = TrueString
			}

			assert.Equal(t, expected, string(content))

			// Verify permissions (Unix only)
			if runtime.GOOS != OsWindows {
				info, err := os.Stat(filePath)
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(ConfigFilePerm), info.Mode().Perm())
			}
		})
	}
}

func TestWriteBoolFileError(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentDir := filepath.Join(tempDir, "non_existent_dir", "file.txt")

	err := WriteBoolFile(nonExistentDir, true)
	assert.Error(t, err)
}

func TestReadBoolFilePermissionDenied(t *testing.T) {
	// Skip on Windows as permission handling differs
	if runtime.GOOS == OsWindows {
		t.Skip("Skipping permission test on Windows")
	}

	// Skip when running as root since root can read files with 000 permissions
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tempDir := t.TempDir()

	// Create a file
	filePath := filepath.Join(tempDir, "no_read_perm.txt")
	err := os.WriteFile(filePath, []byte(TrueString), ConfigFilePerm)
	require.NoError(t, err)

	// Remove read permissions
	err = os.Chmod(filePath, 0)
	require.NoError(t, err)
	defer os.Chmod(filePath, ConfigFilePerm) // Reset permissions for cleanup

	// The function should return the default value when it can't read the file
	result := ReadBoolFile(filePath, false)
	assert.False(t, result)

	result = ReadBoolFile(filePath, true)
	assert.True(t, result)
}

func TestWriteBoolFilePermissionDenied(t *testing.T) {
	// Skip on Windows as permission handling differs
	if runtime.GOOS == OsWindows {
		t.Skip("Skipping permission test on Windows")
	}

	// Skip when running as root since root can write to read-only directories
	if os.Geteuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tempDir := t.TempDir()

	// Make the directory read-only
	err := os.Chmod(tempDir, 0500)
	require.NoError(t, err)
	defer os.Chmod(tempDir, 0700) // Reset permissions for cleanup

	filePath := filepath.Join(tempDir, "readonly.txt")
	err = WriteBoolFile(filePath, true)
	assert.Error(t, err)
}

func TestReadBoolFile_CornerCases(t *testing.T) {
	tempDir := t.TempDir()

	// Test with a directory instead of a file
	dirPath := filepath.Join(tempDir, "directory")
	err := os.Mkdir(dirPath, 0755)
	require.NoError(t, err)

	// Should return the default value when path is a directory
	result := ReadBoolFile(dirPath, true)
	assert.True(t, result)
}
