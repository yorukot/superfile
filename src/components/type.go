package components

type fileState uint
type sideBarStatus uint

const (
	selectDisk sideBarStatus = iota
	selectPinned
)

const (
	selectMultipleFileMode fileState = iota
	normal
)

// main model
type model struct {
	fileModel    fileModel
	sideBarModel sideBarModel
	mainPanelHeight int
}

/* FILE WINDOWS TYPE START*/
type fileModel struct {
	fileWindows    []fileWindows
	width int
}

type fileWindows struct {
	location  string
	fileState fileState
}
/* FILE WINDOWS TYPE END*/

/* SIDE BAR COMPONENTS TYPE START*/
type sideBarModel struct {
	pinnedModel pinnedModel
	choice      string
	state       sideBarStatus
	cursor      int
}

// PINNED MODEL
type pinnedModel struct {
	folder   []folder
	selected string
}

type folder struct {
	location string
	name     string
	endPinned bool
}

/* SIDE BAR COMPONENTS TYPE END*/
