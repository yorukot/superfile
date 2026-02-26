package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"
)

func TestCheckFileNameValidity(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		{name: "Invalid - single dot",
			input:   ".",
			wantErr: true,
			errMsg:  "file name cannot be '.' or '..'",
		}, {
			name:    "invalid - double dot",
			input:   "..",
			wantErr: true,
			errMsg:  "file name cannot be '.' or '..'",
		}, {
			name:    "invalid - ends with /.. (platform separator)",
			input:   fmt.Sprintf("testDir%c..", filepath.Separator),
			wantErr: true,
			errMsg:  fmt.Sprintf("file name cannot end with '%c.' or '%c..'", filepath.Separator, filepath.Separator),
		}, {
			name:    "invalid - ends with /. (platform separator)",
			input:   fmt.Sprintf("testDir%c.", filepath.Separator),
			wantErr: true,
			errMsg:  fmt.Sprintf("file name cannot end with '%c.' or '%c..'", filepath.Separator, filepath.Separator),
		}, {
			name:    "valid - normal file name",
			input:   "valid_file.txt",
			wantErr: false,
		},
		{
			name:    "valid - contains dot inside",
			input:   "some.folder.name/file.txt",
			wantErr: false,
		},
		{
			name:    "valid - ends with dot not after separator",
			input:   "somefile.",
			wantErr: false,
		},
		{
			name:    "valid - ends with .. not after separator",
			input:   "somefile..",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkFileNameValidity(tt.input)

			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func Test_renameIfDuplicate(t *testing.T) {
	curTestDir := t.TempDir()
	f1NonExistent := filepath.Join(curTestDir, "file.txt")
	f2 := filepath.Join(curTestDir, "file2.txt")
	f3 := filepath.Join(curTestDir, "file3(3).txt")
	d1 := filepath.Join(curTestDir, "dir1")

	utils.SetupFiles(t, f2, f3)
	utils.SetupDirectories(t, d1)

	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{
			name:     "file does not exist",
			fileName: f1NonExistent,
			want:     filepath.Base(f1NonExistent),
		},
		{
			name:     "file exists without suffix",
			fileName: f2,
			want:     "file2(1).txt",
		},
		{
			name:     "file exists with suffix",
			fileName: f3,
			want:     "file3(4).txt",
		},
		{
			name:     "directory exists",
			fileName: d1,
			want:     "dir1(1)", // without extension
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := renameIfDuplicate(tt.fileName)
			require.NoError(t, err)
			assert.Equal(t, filepath.Base(tt.want), filepath.Base(results))
		})
	}
}

func Benchmark_renameIfDuplicate(b *testing.B) {
	dir := b.TempDir()

	existingFile := filepath.Join(dir, "file.txt")
	err := os.WriteFile(existingFile, utils.SampleDataBytes, 0644)
	require.NoError(b, err)

	existingDir := filepath.Join(dir, "docs")
	err = os.Mkdir(existingDir, 0o755)
	require.NoError(b, err)

	b.Run("file_exists", func(b *testing.B) {
		for range b.N {
			_, err := renameIfDuplicate(existingFile)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("dir_exists", func(b *testing.B) {
		for range b.N {
			_, err := renameIfDuplicate(existingDir)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("file_not_exists", func(b *testing.B) {
		nonExistent := filepath.Join(dir, "nofile.txt")
		for range b.N {
			_, err := renameIfDuplicate(nonExistent)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
