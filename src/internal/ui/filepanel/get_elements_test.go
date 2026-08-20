package filepanel

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yorukot/superfile/src/pkg/utils"

	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

func TestReturnDirElement(t *testing.T) {
	curTestDir := t.TempDir()
	dir1 := filepath.Join(curTestDir, "dir1")
	dir2 := filepath.Join(curTestDir, "dir2")
	dirNatural := filepath.Join(curTestDir, "dirNatural")
	hiddenOnlyDir := t.TempDir()
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
		{filepath.Join(hiddenOnlyDir, ".hidden"), []byte("hidden")},
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
			name:              "Search in directory with only hidden files",
			location:          hiddenOnlyDir,
			dotFiles:          false,
			sortKind:          sortmodel.SortByName,
			searchString:      "/",
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
			res := panel.elementsRequest(tt.dotFiles).Read()

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

// Cursor movement must not touch the filesystem. Re-reading on every message
// made each keystroke wait on IO, which desynced the cursor from held keys on
// network mounts.
func TestUpdateElementsIfNeededThrottlesRereads(t *testing.T) {
	curTestDir := t.TempDir()
	utils.SetupFiles(t, filepath.Join(curTestDir, "file1.txt"))

	panel := testModel(0, 0, 12, BrowserMode, nil)
	panel.Location = curTestDir
	panel.IsFocused = true

	// First load has nothing to render without, so it reads right away
	panel.UpdateElementsIfNeeded(false, false)
	require.Equal(t, 1, panel.ElemCount())

	utils.SetupFiles(t, filepath.Join(curTestDir, "file2.txt"))

	// Within the interval, no re-read at all
	panel.UpdateElementsIfNeeded(false, false)
	assert.Equal(t, 1, panel.ElemCount(), "must not re-read the directory on every update")

	// Past the interval, the read is handed back to be run off the event loop
	panel.lastTimeGetElement = time.Now().Add(-2 * focussedPanelReRenderTime)
	req := panel.UpdateElementsIfNeeded(false, false)
	require.NotNil(t, req)
	assert.Equal(t, 1, panel.ElemCount(), "the periodic refresh must not read inline")

	// While one is in flight, no duplicate is dispatched
	panel.lastTimeGetElement = time.Now().Add(-2 * focussedPanelReRenderTime)
	assert.Nil(t, panel.UpdateElementsIfNeeded(false, false))

	panel.ApplyElements(*req, req.Read(), false)
	assert.Equal(t, 2, panel.ElemCount())
}

// A result that arrives after the panel has moved on must not be installed, and
// must not leave the panel unable to refresh again.
func TestApplyElementsDropsStaleResults(t *testing.T) {
	curTestDir := t.TempDir()
	otherDir := filepath.Join(curTestDir, "other")
	utils.SetupDirectories(t, otherDir)
	utils.SetupFiles(t, filepath.Join(curTestDir, "file1.txt"))

	panel := testModel(0, 0, 12, BrowserMode, nil)
	panel.Location = curTestDir
	panel.IsFocused = true
	require.Nil(t, panel.UpdateElementsIfNeeded(false, false))
	require.Equal(t, 2, panel.ElemCount())

	panel.lastTimeGetElement = time.Now().Add(-2 * focussedPanelReRenderTime)
	req := panel.UpdateElementsIfNeeded(false, false)
	require.NotNil(t, req)

	panel.Location = otherDir
	require.Nil(t, panel.UpdateElementsIfNeeded(false, false))
	require.Equal(t, 0, panel.ElemCount())

	panel.ApplyElements(*req, req.Read(), false)
	assert.Equal(t, 0, panel.ElemCount(), "a listing for the previous directory must be dropped")

	panel.lastTimeGetElement = time.Now().Add(-2 * focussedPanelReRenderTime)
	assert.NotNil(t, panel.UpdateElementsIfNeeded(false, false))
}

// A changed request is something the panel cannot render without, so it must be
// read immediately rather than waiting for the interval.
func TestUpdateElementsIfNeededReadsChangedRequestAtOnce(t *testing.T) {
	curTestDir := t.TempDir()
	otherDir := filepath.Join(curTestDir, "other")
	utils.SetupDirectories(t, otherDir)
	utils.SetupFiles(t, filepath.Join(curTestDir, "file1.txt"),
		filepath.Join(otherDir, "a.txt"), filepath.Join(otherDir, "b.txt"))

	panel := testModel(0, 0, 12, BrowserMode, nil)
	panel.Location = curTestDir
	panel.IsFocused = true
	require.Nil(t, panel.UpdateElementsIfNeeded(false, false))
	require.Equal(t, 2, panel.ElemCount())

	// Directory change, well within the refresh interval
	panel.Location = otherDir
	require.Nil(t, panel.UpdateElementsIfNeeded(false, false))
	assert.Equal(t, 2, panel.ElemCount())
	assert.Equal(t, "a.txt", panel.GetElementAtIdx(0).Name)

	// Search text change, also within the interval
	panel.SearchBar.SetValue("b")
	require.Nil(t, panel.UpdateElementsIfNeeded(false, false))
	assert.Equal(t, 1, panel.ElemCount())
	assert.Equal(t, "b.txt", panel.GetElementAtIdx(0).Name)
}

// Changes superfile makes itself have to show up at once, not at the next refresh.
func TestMarkStaleForcesReread(t *testing.T) {
	curTestDir := t.TempDir()
	utils.SetupFiles(t, filepath.Join(curTestDir, "file1.txt"))

	panel := testModel(0, 0, 12, BrowserMode, nil)
	panel.Location = curTestDir
	panel.IsFocused = true
	panel.UpdateElementsIfNeeded(false, false)
	require.Equal(t, 1, panel.ElemCount())

	utils.SetupFiles(t, filepath.Join(curTestDir, "file2.txt"))
	panel.MarkStale()

	assert.Nil(t, panel.UpdateElementsIfNeeded(false, false))
	assert.Equal(t, 2, panel.ElemCount())
}

// A read that started before superfile changed the directory must not land
// afterwards and replace the listing read after the change - otherwise a file
// superfile just created briefly disappears again.
func TestMarkStaleInvalidatesInFlightRead(t *testing.T) {
	curTestDir := t.TempDir()
	utils.SetupFiles(t, filepath.Join(curTestDir, "file1.txt"))

	panel := testModel(0, 0, 12, BrowserMode, nil)
	panel.Location = curTestDir
	panel.IsFocused = true
	require.Nil(t, panel.UpdateElementsIfNeeded(false, false))
	require.Equal(t, 1, panel.ElemCount())

	// A background refresh is dispatched and reads the pre-operation state
	panel.lastTimeGetElement = time.Now().Add(-2 * focussedPanelReRenderTime)
	req := panel.UpdateElementsIfNeeded(false, false)
	require.NotNil(t, req)
	staleElements := req.Read()

	// superfile itself adds a file, e.g. a compress finishing
	utils.SetupFiles(t, filepath.Join(curTestDir, "made.zip"))
	panel.MarkStale()
	require.Nil(t, panel.UpdateElementsIfNeeded(false, false))
	require.Equal(t, 2, panel.ElemCount(), "the read after MarkStale must see both files")

	// The older background read now lands and must be ignored
	panel.ApplyElements(*req, staleElements, false)
	assert.Equal(t, 2, panel.ElemCount(),
		"a read from before the change must not replace the listing read after it")

	// and the panel must still be able to refresh afterwards
	panel.lastTimeGetElement = time.Now().Add(-2 * focussedPanelReRenderTime)
	assert.NotNil(t, panel.UpdateElementsIfNeeded(false, false))
}
