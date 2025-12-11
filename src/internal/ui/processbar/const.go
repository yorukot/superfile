package processbar

const (
	// Min width and height for borders
	minHeight = 2
	minWidth  = 2

	// This should allow smooth tracking of 5-10 active processes
	// In case we have issues in future, we could attempt to change this
	msgChannelSize = 50

	// UI dimension constants for process bar rendering
	// borderSize is the border width for the process bar panel
	borderSize = 2

	// progressBarRightPadding is padding after progress bar
	progressBarRightPadding = 3

	// processNameTruncatePadding is the space reserved for ellipsis and icon in process name
	processNameTruncatePadding = 7

	// linesPerProcess is the number of lines needed to render one process
	linesPerProcess = 3
)
