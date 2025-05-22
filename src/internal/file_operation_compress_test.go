package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZipSources(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func(tempDir string) ([]string, error)
		expectedFiles map[string]string
		expectError   bool
	}{
		{
			name: "multiple directories with subdirectories",
			setupFunc: func(tempDir string) ([]string, error) {
				testDir1 := filepath.Join(tempDir, "testdir1")
				testDir2 := filepath.Join(tempDir, "testdir2")
				subDir := filepath.Join(testDir1, "subdir")

				if err := os.MkdirAll(subDir, 0755); err != nil {
					return nil, err
				}
				if err := os.MkdirAll(testDir2, 0755); err != nil {
					return nil, err
				}

				testFiles := map[string]string{
					filepath.Join(testDir1, "file1.txt"): "Content of file1",
					filepath.Join(subDir, "file2.txt"):   "Content of file2",
					filepath.Join(testDir2, "file3.txt"): "Content of file3",
				}

				for filePath, content := range testFiles {
					if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
						return nil, err
					}
				}

				return []string{testDir1, testDir2}, nil
			},
			expectedFiles: map[string]string{
				"testdir1/":                 "",
				"testdir1/file1.txt":        "Content of file1",
				"testdir1/subdir/":          "",
				"testdir1/subdir/file2.txt": "Content of file2",
				"testdir2/":                 "",
				"testdir2/file3.txt":        "Content of file3",
			},
			expectError: false,
		},
		{
			name: "single file",
			setupFunc: func(tempDir string) ([]string, error) {
				testFile := filepath.Join(tempDir, "single.txt")
				testContent := "Single file content"
				if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
					return nil, err
				}
				return []string{testFile}, nil
			},
			expectedFiles: map[string]string{
				"single.txt": "Single file content",
			},
			expectError: false,
		},
		{
			name: "empty list",
			setupFunc: func(tempDir string) ([]string, error) {
				return []string{}, nil
			},
			expectedFiles: map[string]string{},
			expectError:   false,
		},
		{
			name: "non-existent source",
			setupFunc: func(tempDir string) ([]string, error) {
				return []string{"/non/existent/path"}, nil
			},
			expectedFiles: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			sources, err := tt.setupFunc(tempDir)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			targetZip := filepath.Join(tempDir, "test.zip")
			err = zipSources(sources, targetZip)

			if tt.expectError {
				assert.Error(t, err, "zipSources should return error")
				return
			}

			assert.NoError(t, err, "zipSources should not return error")

			zipReader, err := zip.OpenReader(targetZip)
			assert.NoError(t, err, "should be able to open ZIP file")
			defer zipReader.Close()

			assert.Equal(t, len(tt.expectedFiles), len(zipReader.File), "ZIP should contain expected number of files")

			foundFiles := make(map[string]string)
			for _, file := range zipReader.File {
				foundFiles[file.Name] = ""
				if !strings.HasSuffix(file.Name, "/") {
					rc, err := file.Open()
					assert.NoError(t, err, "should be able to open file %s in ZIP", file.Name)

					content, err := io.ReadAll(rc)
					rc.Close()
					assert.NoError(t, err, "should be able to read file %s", file.Name)

					foundFiles[file.Name] = string(content)
				}
			}

			for expectedFile, expectedContent := range tt.expectedFiles {
				foundContent, exists := foundFiles[expectedFile]
				assert.True(t, exists, "expected file %s should be found in ZIP", expectedFile)
				if expectedContent != "" {
					assert.Equal(t, expectedContent, foundContent, "content should match for file %s", expectedFile)
				}
			}

			for foundFile := range foundFiles {
				_, expected := tt.expectedFiles[foundFile]
				assert.True(t, expected, "unexpected file %s found in ZIP", foundFile)
			}
		})
	}
}

func TestZipSourcesInvalidTarget(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	assert.NoError(t, err, "should be able to create test file")

	invalidTarget := "/invalid/path/test.zip"
	err = zipSources([]string{testFile}, invalidTarget)
	assert.Error(t, err, "zipSources should return error for invalid target")
}
