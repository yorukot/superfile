package common

import "time"

// Shared UI/layout constants to replace magic numbers flagged by mnd.
const (
	HelpKeyColumnWidth       = 55              // width of help key column in CLI help
	DefaultCLIContextTimeout = 5 * time.Second // default CLI context timeout for CLI ops

	PanelPadding  = 3 // rows reserved around file list (borders/header/footer)
	BorderPadding = 2 // rows/cols for outer border frame
	// BorderTopHeight and BorderBottomHeight split BorderPadding per edge.
	// A bordered panel reserves one row above its content and one below it, so
	// the top and bottom offsets are not interchangeable: deriving a content
	// offset from BorderPadding instead of BorderTopHeight lands one row too
	// low, which is exactly the kind of error that silently shifts hit-testing
	// away from the rendered output.
	BorderTopHeight    = 1
	BorderBottomHeight = 1
	InnerPadding       = 4 // cols for inner content padding (truncate widths)
	FooterGroupCols    = 3 // columns per group in footer layout math

	DefaultFilePanelWidth    = 10 // default width for file panels
	FilePanelMax             = 10 // max number of file panels supported
	ResponsiveWidthThreshold = 95 // width breakpoint for layout behavior

	HeightBreakA = 30 // responsive height tiers
	HeightBreakB = 35
	HeightBreakC = 40
	HeightBreakD = 45

	FilePanelWidthUnit    = 20                     // width unit used to calculate max file panels
	DefaultPreviewTimeout = 500 * time.Millisecond // preview operation timeout

	FileNameRatioMin = 25
	FileNameRatioMax = 100

	RequiredGradientColorCount = 2

	// UI positioning
	CenterDivisor = 2 // divisor for centering UI elements
)
