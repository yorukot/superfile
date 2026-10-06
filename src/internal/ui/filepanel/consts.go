package filepanel

import (
	"time"

	"github.com/yorukot/superfile/src/internal/common"
)

const (
	contentPadding = 3 // Title + Searchbar + middle border line

	// elementRowHeight is the number of terminal rows a single file entry
	// occupies. Every column of a row is rendered on the same line, so this
	// stays at one regardless of which extra columns are configured.
	elementRowHeight = 1

	MinHeight = contentPadding + common.BorderPadding + 1
	MinWidth  = 18 // minimal width for rename input to render

	FileSizeColumnWidth       = 15
	ModifyTimeSizeColumnWidth = 18
	PermissionsColumnWidth    = 12
	ColumnHeaderHeight        = 1

	// Delimiter between columns in the file panel.
	ColumnDelimiter      = "  "
	ReRenderChunkDivisor = 100
	ReRenderMaxDelay     = 3

	nonFocussedPanelReRenderTime = 3 * time.Second

	emptyCursor = " "
)
