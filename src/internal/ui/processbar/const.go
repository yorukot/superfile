package processbar

const (
	// Min width and height for borders
	minHeight = 2
	minWidth  = 2

	// This should allow smooth tracking of 5-10 active processes
	// In case we have issues in future, we could attempt to change this
	msgChannelSize = 50
)
