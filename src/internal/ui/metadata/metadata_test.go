package metadata

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/barasher/go-exiftool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestGetMetadata(t *testing.T) {
	if runtime.GOOS != utils.OsLinux {
		t.Skip("Skipping metatada fetch test in windows and macOS")
	}
	et, err := exiftool.NewExiftool()
	require.NoError(t, err)
	_, curFilename, _, ok := runtime.Caller(0)
	testdataDir := filepath.Join(filepath.Dir(curFilename), "testdata")

	defaultKeys := []string{keyName, keySize, keyDataModified, keyPermissions}

	require.True(t, ok)
	testdata := []struct {
		name            string
		filepath        string
		metadataFocused bool
	}{
		{
			name:            "Basic Metadata fetching",
			filepath:        filepath.Join(testdataDir, "file1.txt"),
			metadataFocused: true,
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			meta := GetMetadata(tt.filepath, tt.metadataFocused, et)
			assert.Empty(t, meta.infoMsg)
			assert.Equal(t, tt.filepath, meta.filepath)
			for _, key := range defaultKeys {
				_, err := meta.GetValue(key)
				require.NoError(t, err)
			}
		})
	}
}

func TestDirectorySize(t *testing.T) {
	dirPath := t.TempDir()
	// not 4096, or it would match the directory inode's own stat size and the
	// focused assertion would pass without DirSize being called at all
	fileContent := make([]byte, 9000)
	filePath := filepath.Join(dirPath, "file1.txt")
	require.NoError(t, os.WriteFile(filePath, fileContent, 0644))

	// a file one level down, so the expected total below is only reachable by
	// walking the tree, not by stat-ing the directory or its direct children
	nestedContent := make([]byte, 3000)
	nestedDir := filepath.Join(dirPath, "sub")
	require.NoError(t, os.Mkdir(nestedDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(nestedDir, "file2.txt"), nestedContent, 0644))

	dirInfo, err := os.Lstat(dirPath)
	require.NoError(t, err)
	statSize := common.FormatFileSize(dirInfo.Size())
	// computed from the known file sizes rather than from DirSize, which is the
	// function under test
	recursiveSize := common.FormatFileSize(int64(len(fileContent) + len(nestedContent)))
	fileSize := common.FormatFileSize(int64(len(fileContent)))

	unfocusedDir, err := GetMetadata(dirPath, false, nil).GetValue(keySize)
	require.NoError(t, err)
	assert.NotEqual(t, statSize, unfocusedDir)
	assert.Equal(t, dirSizeUnfocusedMsg, unfocusedDir)

	focusedDir, err := GetMetadata(dirPath, true, nil).GetValue(keySize)
	require.NoError(t, err)
	assert.Equal(t, recursiveSize, focusedDir)

	for _, focused := range []bool{false, true} {
		fileVal, err := GetMetadata(filePath, focused, nil).GetValue(keySize)
		require.NoError(t, err)
		assert.Equal(t, fileSize, fileVal)
	}
}
