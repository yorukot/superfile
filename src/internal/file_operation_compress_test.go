package internal

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

func TestZipSources(t *testing.T) {
	processBar := processbar.New()
	processBar.ListenForChannelUpdates()
	t.Cleanup(processBar.SendStopListeningMsgBlocking)
	tests := []struct {
		name          string
		setupFunc     func(t *testing.T, tempDir string) ([]string, error)
		expectedFiles map[string]string
		expectError   bool
	}{
		{
			name: "multiple directories with subdirectories",
			setupFunc: func(t *testing.T, tempDir string) ([]string, error) {
				testDir1 := filepath.Join(tempDir, "testdir1")
				testDir2 := filepath.Join(tempDir, "testdir2")
				subDir := filepath.Join(testDir1, "subdir")
				utils.SetupDirectories(t, testDir1, testDir2, subDir)
				utils.SetupFilesWithData(t, []byte("Content of file1"), filepath.Join(testDir1, "file1.txt"))
				utils.SetupFilesWithData(t, []byte("Content of file2"), filepath.Join(subDir, "file2.txt"))
				utils.SetupFilesWithData(t, []byte("Content of file3"), filepath.Join(testDir2, "file3.txt"))

				return []string{testDir1, testDir2}, nil
			},

			// End for directory is always "/" regardless of windows and linux for zipReader library
			expectedFiles: map[string]string{
				"testdir1/":                                      "",
				filepath.Join("testdir1", "file1.txt"):           "Content of file1",
				filepath.Join("testdir1", "subdir") + "/":        "",
				filepath.Join("testdir1", "subdir", "file2.txt"): "Content of file2",
				"testdir2/":                            "",
				filepath.Join("testdir2", "file3.txt"): "Content of file3",
			},
			expectError: false,
		},
		{
			name: "single file",
			setupFunc: func(t *testing.T, tempDir string) ([]string, error) {
				testFile := filepath.Join(tempDir, "single.txt")
				utils.SetupFilesWithData(t, []byte("Single file content"), testFile)
				return []string{testFile}, nil
			},
			expectedFiles: map[string]string{
				"single.txt": "Single file content",
			},
			expectError: false,
		},
		{
			name: "empty list",
			setupFunc: func(_ *testing.T, _ string) ([]string, error) {
				return []string{}, nil
			},
			expectedFiles: map[string]string{},
			expectError:   false,
		},
		{
			name: "non-existent source",
			setupFunc: func(_ *testing.T, _ string) ([]string, error) {
				return []string{"/non/existent/path"}, nil
			},
			expectedFiles: nil,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			sources, err := tt.setupFunc(t, tempDir)
			if err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			targetZip := filepath.Join(tempDir, "test.zip")
			err = zipSources(sources, targetZip, &processBar)

			if tt.expectError {
				require.Error(t, err, "zipSources should return error")
				return
			}

			require.NoError(t, err, "zipSources should not return error")

			zipReader, err := zip.OpenReader(targetZip)
			require.NoError(t, err, "should be able to open ZIP file")
			defer zipReader.Close()
			validateZipExtraction(t, zipReader, tt.expectedFiles)
		})
	}
}

func validateZipExtraction(t *testing.T, zipReader *zip.ReadCloser, expectedFiles map[string]string) {
	require.Len(t, zipReader.File, len(expectedFiles), "ZIP should contain expected number of files")

	foundFiles := make(map[string]string)
	for _, file := range zipReader.File {
		foundFiles[file.Name] = ""
		if !strings.HasSuffix(file.Name, "/") {
			rc, err := file.Open()
			require.NoError(t, err, "should be able to open file %s in ZIP", file.Name)

			content, err := io.ReadAll(rc)
			rc.Close()
			require.NoError(t, err, "should be able to read file %s", file.Name)

			foundFiles[file.Name] = string(content)
		}
	}

	for expectedFile, expectedContent := range expectedFiles {
		foundContent, exists := foundFiles[expectedFile]
		require.True(t, exists, "expected file %s should be found in ZIP", expectedFile)
		if expectedContent != "" {
			require.Equal(t, expectedContent, foundContent, "content should match for file %s", expectedFile)
		}
	}

	for foundFile := range foundFiles {
		_, expected := expectedFiles[foundFile]
		require.True(t, expected, "unexpected file %s found in ZIP", foundFile)
	}
}

func TestZipSourcesInvalidTarget(t *testing.T) {
	processBar := processbar.New()
	processBar.ListenForChannelUpdates()
	t.Cleanup(processBar.SendStopListeningMsgBlocking)
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err, "should be able to create test file")

	invalidTarget := "/invalid/path/test.zip"
	err = zipSources([]string{testFile}, invalidTarget, &processBar)
	require.Error(t, err, "zipSources should return error for invalid target")
}
