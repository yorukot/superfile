package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/internal/utils"
)

func TestReturnDirElement(t *testing.T) {
	// using 'testDir' set up by testMain of our internal package.
	require.DirExists(t, testDir, "Main test directory should be pre created")
	curTestDir := filepath.Join(testDir, "TestRDE")
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	utils.SetupDirectories(t, curTestDir, dir1, dir2)

	creationDelay := time.Millisecond * 5
	// Cleanup is handled by TestMain

	// Setup files
	// All files with 10 bytes of text

	// dir1
	// - file1.txt
	// dir2 (Empty)
	// .xyz
	// 1.json
	// abc - Add 15 bytes of text
	// aBcD
	// file1.txt
	// file2.txt - Add 20 bytes of text
	// xyz.json

	fileSetup := []struct {
		path string
		data []byte
	}{
		{filepath.Join(curTestDir, ".xyz"), []byte("0123456789")},
		{filepath.Join(dir1, "file1.txt"), []byte("0123456789")},
		{filepath.Join(curTestDir, "aBcD"), []byte("0123456789")},
		{filepath.Join(curTestDir, "file1.txt"), []byte("0123456789")},
		{filepath.Join(curTestDir, "xyz.json"), []byte("0123456789")},
		{filepath.Join(curTestDir, "abc"), []byte("012345678901234")},
		{filepath.Join(curTestDir, "file2.txt"), []byte("01234567890123456789")},
		{filepath.Join(curTestDir, "1.json"), []byte("0123456789")},
	}

	for _, f := range fileSetup {
		utils.SetupFilesWithData(t, f.data, f.path)
		time.Sleep(creationDelay)
	}

	testdata := []struct {
		name              string
		location          string
		dotFiles          bool
		sortOption        string
		reversed          bool
		sortOptions       sortOptionsModelData
		searchString      string
		expectedElemNames []string
	}{
		{
			name:              "Empty Directory",
			location:          dir2,
			dotFiles:          false,
			sortOption:        "Name",
			reversed:          false,
			expectedElemNames: []string{},
		},
		{
			name:       "Sort by Name",
			location:   curTestDir,
			dotFiles:   false,
			sortOption: "Name",
			reversed:   false,
			expectedElemNames: []string{"dir1", "dir2", "1.json", "abc", "aBcD", "file1.txt",
				"file2.txt", "xyz.json"},
		},
		{
			name:       "Sort by Name, with dotfiles",
			location:   curTestDir,
			dotFiles:   true,
			sortOption: "Name",
			reversed:   false,
			expectedElemNames: []string{"dir1", "dir2", ".xyz", "1.json", "abc", "aBcD",
				"file1.txt", "file2.txt", "xyz.json"},
		},
		{
			name:       "Sort by Name Reversed",
			location:   curTestDir,
			dotFiles:   false,
			sortOption: "Name",
			reversed:   true,
			expectedElemNames: []string{"dir2", "dir1", "xyz.json", "file2.txt",
				"file1.txt", "aBcD", "abc", "1.json"},
		},
		{
			name:       "Sort by Size",
			location:   curTestDir,
			dotFiles:   false,
			sortOption: "Size",
			reversed:   false,
			expectedElemNames: []string{"dir2", "dir1", "1.json", "aBcD",
				"file1.txt", "xyz.json", "abc", "file2.txt"},
		},
		{
			name:       "Sort by Size Reversed",
			location:   curTestDir,
			dotFiles:   false,
			sortOption: "Size",
			reversed:   true,
			expectedElemNames: []string{"dir1", "dir2", "file2.txt", "abc", "xyz.json",
				"file1.txt", "aBcD", "1.json"},
		},
		// This one could be flakey if files are created to quickly, or maybe created in
		// parallel
		{
			name:       "Sort by Date",
			location:   curTestDir,
			dotFiles:   false,
			sortOption: "Date Modified",
			reversed:   false,
			expectedElemNames: []string{"1.json", "file2.txt", "abc",
				"xyz.json", "file1.txt", "aBcD", "dir1", "dir2"},
		},
		{
			name:       "Sort by Type",
			location:   curTestDir,
			dotFiles:   false,
			sortOption: "Type",
			reversed:   false,
			expectedElemNames: []string{"dir1", "dir2", "abc", "aBcD", "1.json", "xyz.json",
				"file1.txt", "file2.txt"},
		},
		{
			name:       "Sort by Type Reversed and dotfiles",
			location:   curTestDir,
			dotFiles:   true,
			sortOption: "Type",
			reversed:   true,
			expectedElemNames: []string{"dir2", "dir1", ".xyz", "file2.txt", "file1.txt",
				"xyz.json", "1.json", "aBcD", "abc"},
		},
		{
			name:              "Sort by Type Reversed and dotfiles with search",
			location:          curTestDir,
			dotFiles:          true,
			sortOption:        "Type",
			reversed:          true,
			searchString:      "x",
			expectedElemNames: []string{".xyz", "file2.txt", "file1.txt", "xyz.json"},
		},
		{
			name:              "Sort by Size Reversed with search ftt",
			location:          curTestDir,
			dotFiles:          false,
			sortOption:        "Size",
			reversed:          true,
			searchString:      "ftt",
			expectedElemNames: []string{"file2.txt", "file1.txt"},
		},
		{
			name:              "Sort by Size Reversed with search d",
			location:          curTestDir,
			dotFiles:          false,
			sortOption:        "Size",
			reversed:          true,
			searchString:      "d",
			expectedElemNames: []string{"dir1", "dir2", "aBcD"},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			sortOptionsModel := sortOptionsModelData{
				options:  []string{tt.sortOption},
				selected: 0,
				reversed: tt.reversed,
			}
			var res []element
			if tt.searchString == "" {
				res = returnDirElement(tt.location, tt.dotFiles, sortOptionsModel)
			} else {
				res = returnDirElementBySearchString(tt.location, tt.dotFiles, tt.searchString, sortOptionsModel)
			}

			assert.Len(t, res, len(tt.expectedElemNames))
			actualNames := []string{}
			for i := range res {
				actualNames = append(actualNames, res[i].name)
			}
			assert.Equal(t, tt.expectedElemNames, actualNames)
		})
	}
}

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
