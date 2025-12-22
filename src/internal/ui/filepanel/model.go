package filepanel

import (
	"os"
	"path/filepath"

	"github.com/yorukot/superfile/src/internal/common"
)

// FilePanelSlice creates a slice of FilePanels from the given paths
func FilePanelSlice(paths []string) []Model {
	res := make([]Model, len(paths))
	for i := range paths {
		// Making the first panel as the focussed
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
	sortOptions := sortOptionsModel{
		Width:  SortOptionsDefaultWidth,
		Height: SortOptionsDefaultHeight,
		Open:   false,
		Cursor: common.Config.DefaultSortType,
		Data: sortOptionsModelData{
			Options: []string{
				string(sortingName), string(sortingSize),
				string(sortingDateModified), string(sortingFileType),
			},
			Selected: common.Config.DefaultSortType,
			Reversed: common.Config.SortOrderReversed,
		},
	}
	return New(panelPath, sortOptions, focused, targetFile)
}

func New(location string, sortOptions sortOptionsModel, focused bool, targetFile string) Model {
	return Model{
		Cursor:           0,
		RenderIndex:      0,
		Location:         location,
		SortOptions:      sortOptions,
		PanelMode:        BrowserMode,
		IsFocused:        focused,
		DirectoryRecords: make(map[string]directoryRecord),
		SearchBar:        common.GenerateSearchBar(),
		TargetFile:       targetFile,
		width:            MinWidth,
		height:           MinHeight,
	}
}
