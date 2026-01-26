package filepanel

import (
	"os"
	"path/filepath"

	"github.com/yorukot/superfile/src/internal/common"
	"github.com/yorukot/superfile/src/internal/ui/sortmodel"
)

// FilePanelSlice creates a slice of FilePanels from the given paths
func FilePanelSlice(paths []string) []Model {
	res := make([]Model, len(paths))
	for i := range paths {
		// Making the first panel as the focused
		isFocus := i == 0
		res[i] = defaultFilePanel(paths[i], isFocus)
	}
	return res
}

// defaultFilePanel creates a new FilePanel with default settings
func defaultFilePanel(path string, focused bool) Model {
	targetFile := ""
	panelPath := path
	// If path refers to a file, switch to its parent and remember the filename
	if stat, err := os.Stat(panelPath); err == nil && !stat.IsDir() {
		targetFile = filepath.Base(panelPath)
		panelPath = filepath.Dir(panelPath)
	}
	return New(panelPath, focused, targetFile, sortmodel.SortKind(common.Config.DefaultSortType),
		common.Config.SortOrderReversed)
}

func New(location string, focused bool, targetFile string, sortKind sortmodel.SortKind, sortReversed bool) Model {
	return Model{
		cursor:           0,
		renderIndex:      0,
		Location:         location,
		SortKind:         sortKind,
		SortReversed:     sortReversed,
		PanelMode:        BrowserMode,
		IsFocused:        focused,
		DirectoryRecords: make(map[string]directoryRecord),
		SearchBar:        common.GenerateSearchBar(),
		TargetFile:       targetFile,
		width:            MinWidth,
		height:           MinHeight,
		selected:         make(map[string]int),
	}
}
