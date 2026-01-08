package filepanel

import "github.com/yorukot/superfile/src/internal/common"

const (
	contentPadding = 3 // Title + Searchbar + middle border line
	MinHeight      = contentPadding + common.BorderPadding + 1
	MinWidth       = 18 // minimal width for rename input to render

	sortOptionsDefaultWidth  = 20
	sortOptionsDefaultHeight = 4

	FileSizeColumnWidth       = 15
	ModifyTimeSizeColumnWidth = 18
	PermissionsColumnWidth    = 12
	ColumnHeaderHeight        = 1

	// If the percentage column is smaller than this number, the additional columns will be hidden.
	FileNameRatioDefault = 25

	FileNameRatioMax = 100

	// Delimiter between columns in the file panel.
	ColumnDelimiter = "  "
)
