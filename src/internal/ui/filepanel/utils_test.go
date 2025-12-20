package filepanel

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/internal/utils"
)

func TestReturnDirElement(t *testing.T) {
	curTestDir := t.TempDir()
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
		sortOptions       SortOptionsModelData
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
			sortOptionsModel := SortOptionsModelData{
				Options:  []string{tt.sortOption},
				Selected: 0,
				Reversed: tt.reversed,
			}
			var res []Element
			if tt.searchString == "" {
				res = ReturnDirElement(tt.location, tt.dotFiles, sortOptionsModel)
			} else {
				res = ReturnDirElementBySearchString(tt.location, tt.dotFiles, tt.searchString, sortOptionsModel)
			}

			assert.Len(t, res, len(tt.expectedElemNames))
			actualNames := []string{}
			for i := range res {
				actualNames = append(actualNames, res[i].Name)
			}
			assert.Equal(t, tt.expectedElemNames, actualNames)
		})
	}
}
