package internal

import (
	"os"
	"path/filepath"

	"github.com/yorukot/superfile/src/internal/common"
)

func DefaultFilePanel(path string, focused bool) FilePanel {
	targetFile := ""
	panelPath := path
	// If path refers to a file, switch to its parent and remember the filename
	if stat, err := os.Stat(panelPath); err == nil && !stat.IsDir() {
		targetFile = filepath.Base(panelPath)
		panelPath = filepath.Dir(panelPath)
	}

	return FilePanel{
		RenderIndex: 0,
		Cursor:      0,
		Location:    panelPath,
		SortOptions: SortOptionsModel{
			//nolint:mnd // default sort options dimensions
			Width: 20,
			//nolint:mnd // default sort options dimensions
			Height: 4,
			Open:   false,
			Cursor: common.Config.DefaultSortType,
			Data: SortOptionsModelData{
				Options: []string{
					string(sortingName), string(sortingSize),
					string(sortingDateModified), string(sortingFileType),
				},
				Selected: common.Config.DefaultSortType,
				Reversed: common.Config.SortOrderReversed,
			},
		},
		PanelMode:        BrowserMode,
		IsFocused:        focused,
		DirectoryRecords: make(map[string]DirectoryRecord),
		SearchBar:        common.GenerateSearchBar(),
		TargetFile:       targetFile,
	}
}

func (p PanelMode) String() string {
	switch p {
	case SelectMode:
		return "selectMode"
	case BrowserMode:
		return "browserMode"
	default:
		return invalidTypeString
	}
}
