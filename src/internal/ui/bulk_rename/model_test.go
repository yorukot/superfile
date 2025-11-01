package bulkrename

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
)

func init() {
	common.Theme.GradientColor = []string{"#FF0000", "#00FF00"}
}

func setupTestFiles(t *testing.T, dir string, filenames []string) []string {
	t.Helper()
	var paths []string
	for _, name := range filenames {
		path := filepath.Join(dir, name)
		require.NoError(t, os.WriteFile(path, []byte("test content"), 0644))
		paths = append(paths, path)
	}
	return paths
}

func withProcessBar(t *testing.T, fn func(*processbar.Model)) {
	t.Helper()
	pb := processbar.NewModelWithOptions(80, 25)
	pb.ListenForChannelUpdates()
	defer pb.SendStopListeningMsgBlocking()
	fn(&pb)
}

func TestApplyFindReplace(t *testing.T) {
	m := DefaultModel(25, 80)

	tests := []struct {
		find, replace, input, want string
	}{
		{"old", "new", "old_file.txt", "new_file.txt"},
		{"test", "demo", "test_test_file.txt", "demo_demo_file.txt"},
		{"", "new", "file.txt", "file.txt"},
		{"old_", "", "old_file.txt", "file.txt"},
		{"xyz", "abc", "file.txt", "file.txt"},
	}

	for _, tt := range tests {
		m.findInput.SetValue(tt.find)
		m.replaceInput.SetValue(tt.replace)
		assert.Equal(t, tt.want, m.applyFindReplace(tt.input))
	}
}

func TestApplyPrefix(t *testing.T) {
	m := DefaultModel(25, 80)

	tests := []struct {
		prefix, input, want string
	}{
		{"new_", "file.txt", "new_file.txt"},
		{"test_", "document.pdf", "test_document.pdf"},
		{"prefix_", "file", "prefix_file"},
		{"", "file.txt", "file.txt"},
		{"new_", ".gitignore", "new_.gitignore"},
	}

	for _, tt := range tests {
		m.prefixInput.SetValue(tt.prefix)
		assert.Equal(t, tt.want, m.applyPrefix(tt.input))
	}
}

func TestApplySuffix(t *testing.T) {
	m := DefaultModel(25, 80)

	tests := []struct {
		suffix, input, want string
	}{
		{"_copy", "file.txt", "file_copy.txt"},
		{"_backup", "document.pdf", "document_backup.pdf"},
		{"_new", "file", "file_new"},
		{"", "file.txt", "file.txt"},
	}

	for _, tt := range tests {
		m.suffixInput.SetValue(tt.suffix)
		assert.Equal(t, tt.want, m.applySuffix(tt.input))
	}
}

func TestApplyNumbering(t *testing.T) {
	m := DefaultModel(25, 80)
	m.startNumber = 1

	tests := []struct {
		input string
		idx   int
		want  string
	}{
		{"file.txt", 0, "file_1.txt"},
		{"file.txt", 1, "file_2.txt"},
		{"file", 0, "file_1"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, m.applyNumbering(tt.input, tt.idx))
	}
}

func TestApplyCaseConversion(t *testing.T) {
	m := DefaultModel(25, 80)

	tests := []struct {
		caseType CaseType
		input    string
		want     string
	}{
		{CaseLower, "MyFile.TXT", "myfile.TXT"},
		{CaseUpper, "myfile.txt", "MYFILE.txt"},
		{CaseTitle, "my document file.pdf", "My Document File.pdf"},
		{CaseTitle, "file.txt", "File.txt"},
	}

	for _, tt := range tests {
		m.caseType = tt.caseType
		assert.Equal(t, tt.want, m.applyCaseConversion(tt.input))
	}
}

func TestValidateRename(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestFiles(t, tmpDir, []string{"existing.txt", "test.txt"})
	testFile := filepath.Join(tmpDir, "test.txt")

	m := DefaultModel(25, 80)

	tests := []struct {
		old, new, wantErr string
	}{
		{"test.txt", "newname.txt", ""},
		{"test.txt", "", "Empty filename"},
		{"test.txt", "test.txt", "No change"},
		{"test.txt", "existing.txt", "File already exists"},
		{"test.txt", "TEST.txt", ""},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.wantErr, m.validateRename(testFile, tt.old, tt.new))
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

	require.Len(t, previews, 3)
	for i, p := range previews {
		assert.Equal(t, filepath.Base(files[i]), p.OldName)
		assert.Equal(t, "new_"+filepath.Base(files[i]), p.NewName)
		assert.Empty(t, p.Error)
	}
}

func TestOpenAndClose(t *testing.T) {
	tmpDir := t.TempDir()
	files := setupTestFiles(t, tmpDir, []string{"file1.txt", "file2.txt"})

	m := DefaultModel(25, 80)
	assert.False(t, m.IsOpen())

	m.Open(files, tmpDir)
	assert.True(t, m.IsOpen())
	assert.Equal(t, files, m.selectedFiles)
	assert.Equal(t, tmpDir, m.currentDir)
	assert.Equal(t, FindReplace, m.renameType)

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

	types := []RenameType{FindReplace, AddPrefix, AddSuffix, AddNumbering, ChangeCase, EditorMode}
	
	for i := 0; i < len(types); i++ {
		assert.Equal(t, types[i], m.renameType)
		m.nextType()
	}
	assert.Equal(t, FindReplace, m.renameType)

	m.prevType()
	assert.Equal(t, EditorMode, m.renameType)
	m.prevType()
	assert.Equal(t, ChangeCase, m.renameType)
}

func TestAdjustValue(t *testing.T) {
	m := DefaultModel(25, 80)

	m.renameType = AddNumbering
	m.startNumber = 5
	m.adjustValue(1)
	assert.Equal(t, 6, m.startNumber)
	m.adjustValue(-2)
	assert.Equal(t, 4, m.startNumber)

	m.startNumber = 0
	m.adjustValue(-1)
	assert.Equal(t, 0, m.startNumber)

	m.renameType = ChangeCase
	m.caseType = CaseLower
	m.adjustValue(1)
	assert.Equal(t, CaseUpper, m.caseType)
	m.adjustValue(1)
	assert.Equal(t, CaseTitle, m.caseType)
	m.adjustValue(1)
	assert.Equal(t, CaseTitle, m.caseType)
}

func TestBulkRenameOperation(t *testing.T) {
	t.Run("successful rename", func(t *testing.T) {
		tmpDir := t.TempDir()
		files := map[string]string{"old1.txt": "content1", "old2.txt": "content2"}
		
		for name, content := range files {
			require.NoError(t, os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644))
		}

		previews := []RenamePreview{
			{filepath.Join(tmpDir, "old1.txt"), "old1.txt", "new1.txt", ""},
			{filepath.Join(tmpDir, "old2.txt"), "old2.txt", "new2.txt", ""},
		}

		withProcessBar(t, func(pb *processbar.Model) {
			state := bulkRenameOperation(pb, previews)
			assert.Equal(t, processbar.Successful, state)

			for old, content := range files {
				newPath := filepath.Join(tmpDir, "new"+old[3:])
				assert.FileExists(t, newPath)
				data, _ := os.ReadFile(newPath)
				assert.Equal(t, content, string(data))
				assert.NoFileExists(t, filepath.Join(tmpDir, old))
			}
		})
	})

	t.Run("empty previews returns cancelled", func(t *testing.T) {
		withProcessBar(t, func(pb *processbar.Model) {
			assert.Equal(t, processbar.Cancelled, bulkRenameOperation(pb, []RenamePreview{}))
		})
	})

	t.Run("handles rename error", func(t *testing.T) {
		tmpDir := t.TempDir()
		oldFile := filepath.Join(tmpDir, "old1.txt")
		require.NoError(t, os.WriteFile(oldFile, []byte("content1"), 0644))

		previews := []RenamePreview{
			{oldFile, "old1.txt", "new1.txt", ""},
			{filepath.Join(tmpDir, "nonexistent.txt"), "nonexistent.txt", "new2.txt", ""},
		}

		withProcessBar(t, func(pb *processbar.Model) {
			state := bulkRenameOperation(pb, previews)
			assert.Equal(t, processbar.Failed, state)
			assert.FileExists(t, filepath.Join(tmpDir, "new1.txt"))
		})
	})

	t.Run("large file count", func(t *testing.T) {
		tmpDir := t.TempDir()
		fileCount := 100
		var previews []RenamePreview

		for i := 0; i < fileCount; i++ {
			name := "file" + strconv.Itoa(i) + ".txt"
			path := filepath.Join(tmpDir, name)
			require.NoError(t, os.WriteFile(path, []byte("content"), 0644))
			previews = append(previews, RenamePreview{path, name, "renamed" + strconv.Itoa(i) + ".txt", ""})
		}

		withProcessBar(t, func(pb *processbar.Model) {
			state := bulkRenameOperation(pb, previews)
			assert.Equal(t, processbar.Successful, state)

			for i := 0; i < fileCount; i++ {
				assert.FileExists(t, filepath.Join(tmpDir, "renamed"+strconv.Itoa(i)+".txt"))
			}
		})
	})
}
