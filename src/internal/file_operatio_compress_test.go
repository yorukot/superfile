package internal

import (
	"archive/zip"

	"io"

	"os"

	"path/filepath"

	"strings"

	"testing"
)

func TestZipSources(t *testing.T) {

	tempDir := t.TempDir()

	testDir1 := filepath.Join(tempDir, "testdir1")

	testDir2 := filepath.Join(tempDir, "testdir2")

	subDir := filepath.Join(testDir1, "subdir")

	err := os.MkdirAll(subDir, 0755)

	if err != nil {

		t.Fatalf("Error creating test directory: %v", err)

	}

	err = os.MkdirAll(testDir2, 0755)

	if err != nil {

		t.Fatalf("Error creating test directory: %v", err)

	}

	testFiles := map[string]string{

		filepath.Join(testDir1, "file1.txt"): "Content of file1",

		filepath.Join(subDir, "file2.txt"): "Content of file2",

		filepath.Join(testDir2, "file3.txt"): "Content of file3",
	}

	for filePath, content := range testFiles {

		err := os.WriteFile(filePath, []byte(content), 0644)

		if err != nil {

			t.Fatalf("Error creating test file %s: %v", filePath, err)

		}

	}

	targetZip := filepath.Join(tempDir, "test.zip")

	sources := []string{testDir1, testDir2}

	err = zipSources(sources, targetZip)

	if err != nil {

		t.Fatalf("zipSources returned error: %v", err)

	}

	zipReader, err := zip.OpenReader(targetZip)

	if err != nil {

		t.Fatalf("Error opening ZIP file: %v", err)

	}

	defer zipReader.Close()

	expectedFiles := map[string]string{

		"testdir1/": "",

		"testdir1/file1.txt": "Content of file1",

		"testdir1/subdir/": "",

		"testdir1/subdir/file2.txt": "Content of file2",

		"testdir2/": "",

		"testdir2/file3.txt": "Content of file3",
	}

	foundFiles := make(map[string]string)

	for _, file := range zipReader.File {

		foundFiles[file.Name] = ""

		if !strings.HasSuffix(file.Name, "/") {

			rc, err := file.Open()

			if err != nil {

				t.Errorf("Error opening file %s in ZIP: %v", file.Name, err)

				continue

			}

			content, err := io.ReadAll(rc)

			rc.Close()

			if err != nil {

				t.Errorf("Error reading file %s: %v", file.Name, err)

				continue

			}

			foundFiles[file.Name] = string(content)

		}

	}

	for expectedFile, expectedContent := range expectedFiles {

		if foundContent, exists := foundFiles[expectedFile]; !exists {

			t.Errorf("Expected file %s not found in ZIP", expectedFile)

		} else if expectedContent != "" && foundContent != expectedContent {

			t.Errorf("Content mismatch in file %s. Expected: %s, Found: %s",

				expectedFile, expectedContent, foundContent)

		}

	}

	for foundFile := range foundFiles {

		if _, expected := expectedFiles[foundFile]; !expected {

			t.Errorf("Unexpected file %s found in ZIP", foundFile)

		}

	}

}

func TestZipSourcesWithSingleFile(t *testing.T) {

	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "single.txt")

	testContent := "Single file content"

	err := os.WriteFile(testFile, []byte(testContent), 0644)

	if err != nil {

		t.Fatalf("Error creating test file: %v", err)

	}

	targetZip := filepath.Join(tempDir, "single.zip")

	err = zipSources([]string{testFile}, targetZip)

	if err != nil {

		t.Fatalf("zipSources returned error: %v", err)

	}

	zipReader, err := zip.OpenReader(targetZip)

	if err != nil {

		t.Fatalf("Error opening ZIP file: %v", err)

	}

	defer zipReader.Close()

	if len(zipReader.File) != 1 {

		t.Fatalf("Expected 1 file in ZIP, found: %d", len(zipReader.File))

	}

	file := zipReader.File[0]

	if file.Name != "single.txt" {

		t.Errorf("Expected filename: single.txt, found: %s", file.Name)

	}

	rc, err := file.Open()

	if err != nil {

		t.Fatalf("Error opening file in ZIP: %v", err)

	}

	defer rc.Close()

	content, err := io.ReadAll(rc)

	if err != nil {

		t.Fatalf("Error reading file: %v", err)

	}

	if string(content) != testContent {

		t.Errorf("Content mismatch. Expected: %s, Found: %s",

			testContent, string(content))

	}

}

func TestZipSourcesEmptyList(t *testing.T) {

	tempDir := t.TempDir()

	targetZip := filepath.Join(tempDir, "empty.zip")

	err := zipSources([]string{}, targetZip)

	if err != nil {

		t.Fatalf("zipSources should work with empty list but got error: %v", err)

	}

	zipReader, err := zip.OpenReader(targetZip)

	if err != nil {

		t.Fatalf("Error opening ZIP file: %v", err)

	}

	defer zipReader.Close()

	if len(zipReader.File) != 0 {

		t.Errorf("Expected 0 files in ZIP, found: %d", len(zipReader.File))

	}

}

func TestZipSourcesNonExistentSource(t *testing.T) {

	tempDir := t.TempDir()

	targetZip := filepath.Join(tempDir, "test.zip")

	err := zipSources([]string{"/non/existent/path"}, targetZip)

	if err == nil {

		t.Error("zipSources should return error for non-existent source")

	}

}

func TestZipSourcesInvalidTarget(t *testing.T) {

	tempDir := t.TempDir()

	testFile := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(testFile, []byte("test"), 0644)

	if err != nil {

		t.Fatalf("Error creating test file: %v", err)

	}

	invalidTarget := "/invalid/path/test.zip"

	err = zipSources([]string{testFile}, invalidTarget)

	if err == nil {

		t.Error("zipSources should return error for invalid target")

	}

}
