package components

import "time"

type fileState uint
type sideBarStatus uint
type filePanelFocusType uint

const (
	selectDisk sideBarStatus = iota
	selectPinned
)

const (
	noneFocus filePanelFocusType = iota
	secondFocus
	focus
)

const (
	selectMultipleFileMode fileState = iota
	normal
)

// main model
type model struct {
	fileModel           fileModel
	sideBarModel        sideBarModel
	filePanelFocusIndex int
	sideBarFocus        bool
	mainPanelHeight     int
	test                string
	fullWidth           int
	fullHeight          int
}

/* FILE WINDOWS TYPE START*/
type fileModel struct {
	filePanels []filePanel
	width      int
}

type filePanel struct {
	cursor    int
	render    int
	focusType filePanelFocusType
	location  string
	fileState fileState
	selected  []selectedElement
	element   []element
}

type selectedElement struct {
	location string
}

type element struct {
	name       string
	location   string
	folder     bool
	size       int64
	updateTime time.Time
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
	location  string
	name      string
	endPinned bool
}

/* SIDE BAR COMPONENTS TYPE END*/

/*PROCESS BAR COMPONENTS TYPE START*/

type processBar struct {
	process []process
}

type process struct {
	name        string
	process     int
	description string
	command     string
}

/*PROCESS BAR COMPONENTS TYPE END*/

type iconStyle struct {
	icon  string
	color string
}
