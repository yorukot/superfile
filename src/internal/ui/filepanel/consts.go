package filepanel

import "github.com/yorukot/superfile/src/internal/common"

const (
	FilePanelContentPadding = 3                        // Title + Searchbar + middle border line
	FilePanelMinWidth       = common.BorderPadding + 3 // Must fit the searchbar
	FilePanelMinHeight      = FilePanelContentPadding + common.BorderPadding + 1

	SortOptionsDefaultWidth  = 20
	SortOptionsDefaultHeight = 4
)
