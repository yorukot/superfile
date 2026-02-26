package filepanel

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

func TestReturnDirElement(t *testing.T) {
	curTestDir := t.TempDir()
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	dirNatural := filepath.Join(curTestDir, "dirNatural")
	utils.SetupDirectories(t, curTestDir, dir1, dir2, dirNatural)

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
		{filepath.Join(dirNatural, "file1.txt"), []byte("a")},
		{filepath.Join(dirNatural, "file2.txt"), []byte("b")},
		{filepath.Join(dirNatural, "file10.txt"), []byte("c")},
		{filepath.Join(dirNatural, "file20.txt"), []byte("d")},
	}

	for _, f := range fileSetup {
		utils.SetupFilesWithData(t, f.data, f.path)
		time.Sleep(creationDelay)
	}

	testdata := []struct {
		name              string
		location          string
		dotFiles          bool
		sortKind          sortmodel.SortKind
		reversed          bool
		searchString      string
		expectedElemNames []string
	}{
		{
			name:              "Empty Directory",
			location:          dir2,
			dotFiles:          false,
			sortKind:          sortmodel.SortByName,
			reversed:          false,
			expectedElemNames: []string{},
		},
		{
			name:     "Sort by Name",
			location: curTestDir,
			dotFiles: false,
			sortKind: sortmodel.SortByName,
			reversed: false,
			expectedElemNames: []string{"dir1", "dir2", "dirNatural", "1.json", "abc", "aBcD", "file1.txt",
				"file2.txt", "xyz.json"},
		},
		{
			name:     "Sort by Name, with dotfiles",
			location: curTestDir,
			dotFiles: true,
			sortKind: sortmodel.SortByName,
			reversed: false,
			expectedElemNames: []string{"dir1", "dir2", "dirNatural", ".xyz", "1.json", "abc", "aBcD",
				"file1.txt", "file2.txt", "xyz.json"},
		},
		{
			name:     "Sort by Name Reversed",
			location: curTestDir,
			dotFiles: false,
			sortKind: sortmodel.SortByName,
			reversed: true,
			expectedElemNames: []string{"dirNatural", "dir2", "dir1", "xyz.json", "file2.txt",
				"file1.txt", "aBcD", "abc", "1.json"},
		},
		{
			name:     "Sort by Size",
			location: curTestDir,
			dotFiles: false,
			sortKind: sortmodel.SortBySize,
			reversed: false,
			expectedElemNames: []string{"dir2", "dir1", "dirNatural", "1.json", "aBcD",
				"file1.txt", "xyz.json", "abc", "file2.txt"},
		},
		{
			name:     "Sort by Size Reversed",
			location: curTestDir,
			dotFiles: false,
			sortKind: sortmodel.SortBySize,
			reversed: true,
			expectedElemNames: []string{"dirNatural", "dir1", "dir2", "file2.txt", "abc", "xyz.json",
				"file1.txt", "aBcD", "1.json"},
		},
		// This one could be flakey if files are created to quickly, or maybe created in
		// parallel
		{
			name:     "Sort by Date",
			location: curTestDir,
			dotFiles: false,
			sortKind: sortmodel.SortByDate,
			reversed: false,
			expectedElemNames: []string{"dirNatural", "1.json", "file2.txt", "abc",
				"xyz.json", "file1.txt", "aBcD", "dir1", "dir2"},
		},
		{
			name:     "Sort by Type",
			location: curTestDir,
			dotFiles: false,
			sortKind: sortmodel.SortByType,
			reversed: false,
			expectedElemNames: []string{"dir1", "dir2", "dirNatural", "abc", "aBcD", "1.json", "xyz.json",
				"file1.txt", "file2.txt"},
		},
		{
			name:     "Sort by Type Reversed and dotfiles",
			location: curTestDir,
			dotFiles: true,
			sortKind: sortmodel.SortByType,
			reversed: true,
			expectedElemNames: []string{"dirNatural", "dir2", "dir1", ".xyz", "file2.txt", "file1.txt",
				"xyz.json", "1.json", "aBcD", "abc"},
		},
		{
			name:              "Sort by Type Reversed and dotfiles with search",
			location:          curTestDir,
			dotFiles:          true,
			sortKind:          sortmodel.SortByType,
			reversed:          true,
			searchString:      "x",
			expectedElemNames: []string{".xyz", "file2.txt", "file1.txt", "xyz.json"},
		},
		{
			name:              "Sort by Size Reversed with search ftt",
			location:          curTestDir,
			dotFiles:          false,
			sortKind:          sortmodel.SortBySize,
			reversed:          true,
			searchString:      "ftt",
			expectedElemNames: []string{"file2.txt", "file1.txt"},
		},
		{
			name:              "Sort by Size Reversed with search d",
			location:          curTestDir,
			dotFiles:          false,
			sortKind:          sortmodel.SortBySize,
			reversed:          true,
			searchString:      "d",
			expectedElemNames: []string{"dirNatural", "dir1", "dir2", "aBcD"},
		},
		{
			name:              "Sort by Natural",
			location:          dirNatural,
			dotFiles:          false,
			sortKind:          sortmodel.SortByNatural,
			reversed:          false,
			expectedElemNames: []string{"file1.txt", "file2.txt", "file10.txt", "file20.txt"},
		},
		{
			name:              "Sort by Natural Reversed",
			location:          dirNatural,
			dotFiles:          false,
			sortKind:          sortmodel.SortByNatural,
			reversed:          true,
			expectedElemNames: []string{"file20.txt", "file10.txt", "file2.txt", "file1.txt"},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			panel := testModel(0, 0, 0, BrowserMode, nil)
			panel.Location = tt.location
			panel.SortKind = tt.sortKind
			panel.SortReversed = tt.reversed
			panel.SearchBar.SetValue(tt.searchString)
			var res []Element
			if tt.searchString == "" {
				res = panel.getDirectoryElements(tt.dotFiles)
			} else {
				res = panel.getDirectoryElementsBySearch(tt.dotFiles)
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

func TestSingleItemSelect(t *testing.T) {
	testdata := []struct {
		name             string
		panel            Model
		panelToSelect    []string
		expectedSelected map[string]int
	}{
		{
			name: "Select unselected item",
			panel: testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
			}),
			panelToSelect:    []string{},
			expectedSelected: map[string]int{"/tmp/file1.txt": 1},
		},
		{
			name: "Deselect selected item",
			panel: testModel(0, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
				{Name: "file2.txt", Location: "/tmp/file2.txt"},
			}),
			panelToSelect:    []string{"/tmp/file1.txt"},
			expectedSelected: map[string]int{},
		},
		{
			name: "Out of bounds cursor negative",
			panel: testModel(-1, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
			}),
			panelToSelect:    []string{},
			expectedSelected: map[string]int{},
		},
		{
			name: "Out of bounds cursor beyond count",
			panel: testModel(5, 0, 12, SelectMode, []Element{
				{Name: "file1.txt", Location: "/tmp/file1.txt"},
			}),
			panelToSelect:    []string{},
			expectedSelected: map[string]int{},
		},
		{
			name:             "Empty element list",
			panel:            testModel(0, 0, 12, SelectMode, []Element{}),
			panelToSelect:    []string{},
			expectedSelected: map[string]int{},
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			tt.panel.SetSelectedAll(tt.panelToSelect)
			tt.panel.SingleItemSelect()
			assert.Equal(t, tt.expectedSelected, tt.panel.selected)
		})
	}
}
