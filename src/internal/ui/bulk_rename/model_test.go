package bulkrename

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create test files
func setupTestFiles(t *testing.T, dir string, filenames []string) []string {
	t.Helper()
	var paths []string
	for _, name := range filenames {
		path := filepath.Join(dir, name)
		err := os.WriteFile(path, []byte("test content"), 0644)
		require.NoError(t, err)
		paths = append(paths, path)
	}
	return paths
}

func TestApplyFindReplace(t *testing.T) {
	m := DefaultModel(25, 80)

	testCases := []struct {
		name       string
		findVal    string
		replaceVal string
		filename   string
		expected   string
	}{
		{
			name:       "Basic find and replace",
			findVal:    "old",
			replaceVal: "new",
			filename:   "old_file.txt",
			expected:   "new_file.txt",
		},
		{
			name:       "Multiple occurrences",
			findVal:    "test",
			replaceVal: "demo",
			filename:   "test_test_file.txt",
			expected:   "demo_demo_file.txt",
		},
		{
			name:       "Empty find returns original",
			findVal:    "",
			replaceVal: "new",
			filename:   "file.txt",
			expected:   "file.txt",
		},
		{
			name:       "Replace with empty string",
			findVal:    "old_",
			replaceVal: "",
			filename:   "old_file.txt",
			expected:   "file.txt",
		},
		{
			name:       "No match",
			findVal:    "xyz",
			replaceVal: "abc",
			filename:   "file.txt",
			expected:   "file.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m.findInput.SetValue(tc.findVal)
			m.replaceInput.SetValue(tc.replaceVal)
			result := m.applyFindReplace(tc.filename)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestApplyPrefix(t *testing.T) {
	m := DefaultModel(25, 80)

	testCases := []struct {
		name     string
		prefix   string
		filename string
		expected string
	}{
		{
			name:     "Add prefix to simple file",
			prefix:   "new_",
			filename: "file.txt",
			expected: "new_file.txt",
		},
		{
			name:     "Add prefix preserves extension",
			prefix:   "test_",
			filename: "document.pdf",
			expected: "test_document.pdf",
		},
		{
			name:     "Add prefix to file without extension",
			prefix:   "prefix_",
			filename: "file",
			expected: "prefix_file",
		},
		{
			name:     "Empty prefix returns original",
			prefix:   "",
			filename: "file.txt",
			expected: "file.txt",
		},
		{
			name:     "Add prefix to hidden file",
			prefix:   "new_",
			filename: ".gitignore",
			expected: "new_.gitignore",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m.prefixInput.SetValue(tc.prefix)
			result := m.applyPrefix(tc.filename)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestApplySuffix(t *testing.T) {
	m := DefaultModel(25, 80)

	testCases := []struct {
		name     string
		suffix   string
		filename string
		expected string
	}{
		{
			name:     "Add suffix to simple file",
			suffix:   "_copy",
			filename: "file.txt",
			expected: "file_copy.txt",
		},
		{
			name:     "Add suffix preserves extension",
			suffix:   "_backup",
			filename: "document.pdf",
			expected: "document_backup.pdf",
		},
		{
			name:     "Add suffix to file without extension",
			suffix:   "_new",
			filename: "file",
			expected: "file_new",
		},
		{
			name:     "Empty suffix returns original",
			suffix:   "",
			filename: "file.txt",
			expected: "file.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m.suffixInput.SetValue(tc.suffix)
			result := m.applySuffix(tc.filename)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestApplyNumbering(t *testing.T) {
	m := DefaultModel(25, 80)
	m.startNumber = 1

	testCases := []struct {
		name     string
		filename string
		index    int
		expected string
	}{
		{
			name:     "First file",
			filename: "file.txt",
			index:    0,
			expected: "file_1.txt",
		},
		{
			name:     "Second file",
			filename: "file.txt",
			index:    1,
			expected: "file_2.txt",
		},
		{
			name:     "File without extension",
			filename: "file",
			index:    0,
			expected: "file_1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := m.applyNumbering(tc.filename, tc.index)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestApplyCaseConversion(t *testing.T) {
	m := DefaultModel(25, 80)

	testCases := []struct {
		name     string
		caseType CaseType
		filename string
		expected string
	}{
		{
			name:     "Lowercase conversion",
			caseType: CaseLower,
			filename: "MyFile.TXT",
			expected: "myfile.TXT",
		},
		{
			name:     "Uppercase conversion",
			caseType: CaseUpper,
			filename: "myfile.txt",
			expected: "MYFILE.txt",
		},
		{
			name:     "Title case conversion",
			caseType: CaseTitle,
			filename: "my document file.pdf",
			expected: "My Document File.pdf",
		},
		{
			name:     "Title case single word",
			caseType: CaseTitle,
			filename: "file.txt",
			expected: "File.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m.caseType = tc.caseType
			result := m.applyCaseConversion(tc.filename)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateRename(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	existingFile := filepath.Join(tmpDir, "existing.txt")
	err := os.WriteFile(existingFile, []byte("test"), 0644)
	require.NoError(t, err)

	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	require.NoError(t, err)

	m := DefaultModel(25, 80)

	testCases := []struct {
		name        string
		itemPath    string
		oldName     string
		newName     string
		expectedErr string
	}{
		{
			name:        "Valid rename",
			itemPath:    testFile,
			oldName:     "test.txt",
			newName:     "newname.txt",
			expectedErr: "",
		},
		{
			name:        "Empty filename",
			itemPath:    testFile,
			oldName:     "test.txt",
			newName:     "",
			expectedErr: "Empty filename",
		},
		{
			name:        "No change",
			itemPath:    testFile,
			oldName:     "test.txt",
			newName:     "test.txt",
			expectedErr: "No change",
		},
		{
			name:        "File already exists",
			itemPath:    testFile,
			oldName:     "test.txt",
			newName:     "existing.txt",
			expectedErr: "File already exists",
		},
		{
			name:        "Case change only (allowed on case-insensitive fs)",
			itemPath:    testFile,
			oldName:     "test.txt",
			newName:     "TEST.txt",
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := m.validateRename(tc.itemPath, tc.oldName, tc.newName)
			assert.Equal(t, tc.expectedErr, result)
		})
	}
}

func TestGeneratePreview(t *testing.T) {
	tmpDir := t.TempDir()
	files := setupTestFiles(t, tmpDir, []string{"file1.txt", "file2.txt", "file3.txt"})

	m := DefaultModel(25, 80)
	m.selectedFiles = files
	m.currentDir = tmpDir
	m.renameType = AddPrefix
	m.prefixInput.SetValue("new_")

	previews := m.GeneratePreview()

	assert.Len(t, previews, 3)
	assert.Equal(t, "file1.txt", previews[0].OldName)
	assert.Equal(t, "new_file1.txt", previews[0].NewName)
	assert.Empty(t, previews[0].Error)

	assert.Equal(t, "file2.txt", previews[1].OldName)
	assert.Equal(t, "new_file2.txt", previews[1].NewName)
	assert.Empty(t, previews[1].Error)
}

func TestOpenAndClose(t *testing.T) {
	tmpDir := t.TempDir()
	files := setupTestFiles(t, tmpDir, []string{"file1.txt", "file2.txt"})

	m := DefaultModel(25, 80)

	// Initially closed
	assert.False(t, m.IsOpen())

	// Open modal
	m.Open(files, tmpDir)
	assert.True(t, m.IsOpen())
	assert.Equal(t, files, m.selectedFiles)
	assert.Equal(t, tmpDir, m.currentDir)
	assert.Equal(t, FindReplace, m.renameType)

	// Close modal
	m.Close()
	assert.False(t, m.IsOpen())
	assert.Nil(t, m.selectedFiles)
	assert.Empty(t, m.currentDir)
}

func TestNavigateTypes(t *testing.T) {
	tmpDir := t.TempDir()
	files := setupTestFiles(t, tmpDir, []string{"file1.txt"})

	m := DefaultModel(25, 80)
	m.Open(files, tmpDir)
	m.renameType = FindReplace

	// Navigate forward
	m.nextType()
	assert.Equal(t, AddPrefix, m.renameType)

	m.nextType()
	assert.Equal(t, AddSuffix, m.renameType)

	m.nextType()
	assert.Equal(t, AddNumbering, m.renameType)

	m.nextType()
	assert.Equal(t, ChangeCase, m.renameType)

	m.nextType()
	assert.Equal(t, EditorMode, m.renameType)

	// Should wrap around
	m.nextType()
	assert.Equal(t, FindReplace, m.renameType)

	// Navigate backward
	m.prevType()
	assert.Equal(t, EditorMode, m.renameType)

	m.prevType()
	assert.Equal(t, ChangeCase, m.renameType)
}

func TestAdjustValue(t *testing.T) {
	m := DefaultModel(25, 80)

	// Test numbering adjustment
	m.renameType = AddNumbering
	m.startNumber = 5

	m.adjustValue(1)
	assert.Equal(t, 6, m.startNumber)

	m.adjustValue(-2)
	assert.Equal(t, 4, m.startNumber)

	// Should not go below 0
	m.startNumber = 0
	m.adjustValue(-1)
	assert.Equal(t, 0, m.startNumber)

	// Test case type adjustment
	m.renameType = ChangeCase
	m.caseType = CaseLower

	m.adjustValue(1)
	assert.Equal(t, CaseUpper, m.caseType)

	m.adjustValue(1)
	assert.Equal(t, CaseTitle, m.caseType)

	// Should not go beyond bounds
	m.adjustValue(1)
	assert.Equal(t, CaseTitle, m.caseType)
}

func TestBulkRenameOperation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	oldFile1 := filepath.Join(tmpDir, "old1.txt")
	oldFile2 := filepath.Join(tmpDir, "old2.txt")
	err := os.WriteFile(oldFile1, []byte("content1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(oldFile2, []byte("content2"), 0644)
	require.NoError(t, err)

	// Create previews for rename
	previews := []RenamePreview{
		{
			OldPath: oldFile1,
			OldName: "old1.txt",
			NewName: "new1.txt",
			Error:   "",
		},
		{
			OldPath: oldFile2,
			OldName: "old2.txt",
			NewName: "new2.txt",
			Error:   "",
		},
	}

	// Execute renames manually (without processbar for testing)
	for _, preview := range previews {
		newPath := filepath.Join(filepath.Dir(preview.OldPath), preview.NewName)
		err := os.Rename(preview.OldPath, newPath)
		require.NoError(t, err)
	}

	// Verify files were renamed
	newFile1 := filepath.Join(tmpDir, "new1.txt")
	newFile2 := filepath.Join(tmpDir, "new2.txt")

	assert.FileExists(t, newFile1)
	assert.FileExists(t, newFile2)
	assert.NoFileExists(t, oldFile1)
	assert.NoFileExists(t, oldFile2)

	// Verify content preserved
	content1, err := os.ReadFile(newFile1)
	require.NoError(t, err)
	assert.Equal(t, "content1", string(content1))

	content2, err := os.ReadFile(newFile2)
	require.NoError(t, err)
	assert.Equal(t, "content2", string(content2))
}
