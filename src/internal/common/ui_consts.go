package common

import "time"

// Shared UI/layout constants to replace magic numbers flagged by mnd.
const (
	HelpKeyColumnWidth       = 55              // width of help key column in CLI help
	DefaultCLIContextTimeout = 5 * time.Second // default CLI context timeout for CLI ops

	PanelPadding    = 3 // rows reserved around file list (borders/header/footer)
	BorderPadding   = 2 // rows/cols for outer border frame
	InnerPadding    = 4 // cols for inner content padding (truncate widths)
	FooterGroupCols = 3 // columns per group in footer layout math

	DefaultFilePanelWidth    = 10 // default width for file panels
	FilePanelMax             = 10 // max number of file panels supported
	MinWidthForRename        = 18 // minimal width for rename input to render
	ResponsiveWidthThreshold = 95 // width breakpoint for layout behavior

	HeightBreakA = 30 // responsive height tiers
	HeightBreakB = 35
	HeightBreakC = 40
	HeightBreakD = 45

	ReRenderChunkDivisor = 100 // divisor for re-render throttling

	FilePanelWidthUnit    = 20                     // width unit used to calculate max file panels
	DefaultPreviewTimeout = 500 * time.Millisecond // preview operation timeout

	// File permissions

	// UI positioning
	CenterDivisor = 2 // divisor for centering UI elements
)
