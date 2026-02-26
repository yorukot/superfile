package metadata

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/barasher/go-exiftool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
