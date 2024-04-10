package components

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
)

// Type representing the mode of the panel
type panelMode uint

// Type representing the focus type of the file panel
type filePanelFocusType uint

// Type representing the state of a process
type processState int

// Type representing the type of focused panel
type focusPanelType int

// Type representing the type of item
type itemType int

type warnType int

const (
	confirmDeleteItem warnType = iota
)

// Constants for new file or new directory
const (
	newFile itemType = iota
	newDirectory
)

// Constants for panel with no focus
const (
	nonePanelFocus focusPanelType = iota
	processBarFocus
	sideBarFocus
	metaDataFocus
)

// Constants for file panel with no focus
const (
	noneFocus filePanelFocusType = iota
	secondFocus
	focus
)

// Constants for select mode or browser mode
const (
	selectMode panelMode = iota
	browserMode
)

// Constants for operation, success, cancel, failure
const (
	inOperation processState = iota
	successful
	cancel
	failure
)

// Main model
type model struct {
	fileModel           fileModel
	sideBarModel        sideBarModel
	processBarModel     processBarModel
	focusPanel          focusPanelType
	copyItems           copyItems
	typingModal         typingModal
	warnModal           warnModal
	fileMetaData        fileMetaData
	firstTextInput      bool
	toggleDotFile       bool
	filePanelFocusIndex int
	mainPanelHeight     int
	fullWidth           int
	fullHeight          int
}

// Modal

type warnModal struct {
	open     bool
	warnType warnType
	title    string
	content  string
}

type typingModal struct {
	location  string
	open      bool
	itemType  itemType
	textInput textinput.Model
}

// File metadata
type fileMetaData struct {
	metaData    [][2]string
	renderIndex int
}

// Copied items
type copyItems struct {
	items         []string
	originalPanel originalPanel
	cut           bool
}

// Original panel
type originalPanel struct {
	index    int
	location string
}

/* FILE WINDOWS TYPE START*/
// Model for file windows
type fileModel struct {
	filePanels []filePanel
	width      int
	renaming   bool
}

// Panel representing a file
type filePanel struct {
	cursor          int
	render          int
	focusType       filePanelFocusType
	location        string
	panelMode       panelMode
	selected        []string
	element         []element
	directoryRecord map[string]directoryRecord
	rename          textinput.Model
	renaming        bool
}

// Record for directory navigation
type directoryRecord struct {
	directoryCursor int
	directoryRender int
}

// Element within a file panel
type element struct {
	name      string
	location  string
	directory bool
	metaData  [][2]string
}

/* FILE WINDOWS TYPE END*/

/* SIDE BAR COMPONENTS TYPE START*/
// Model for sidebar components
type sideBarModel struct {
	pinnedModel pinnedModel
	cursor      int
}

// Model for pinned items in sidebar
type pinnedModel struct {
	directory []directory
}

// Directory within pinned items
type directory struct {
	location  string
	name      string
	endPinned bool
}

/* SIDE BAR COMPONENTS TYPE END*/

/*PROCESS BAR COMPONENTS TYPE START*/

// Model for process bar components
type processBarModel struct {
	render      int
	cursor      int
	processList []string
	process     map[string]process
}

// Model for an individual process
type process struct {
	name     string
	progress progress.Model
	state    processState
	total    int
	done     int
	doneTime time.Time
}

// Message for process bar
type channelMessage struct {
	processId       string
	processNewState process
	returnWarnModal bool
	warnModal       warnModal
}

/*PROCESS BAR COMPONENTS TYPE END*/

// Style for icons
type iconStyle struct {
	icon  string
	color string
}

// Theme configuration
type ThemeType struct {
	Border string
	Cursor string

	TerminalTooSmallError string
	TerminalSizeCorrect   string

	BrowserMode string
	SelectMode  string

	SideBarTitle    string
	SideBarItem     string
	SideBarSelected string
	SideBarFocus    string

	FilePanelFocus            string
	FilePanelTopDirectoryIcon string
	FilePanelTopPath          string
	FilePanelItem             string
	FilePanelItemSelected     string

	BottomBarFocus string

	ProcessBarSideLine string
	ProcessBarGradient [2]string
	InOperation        string
	Done               string
	Fail               string
	Cancel             string
	Warn               string

	ModalForeground string
	ModalCancel     string
	ModalConfirm    string
}

// Configuration settings
type ConfigType struct {
	Theme           string
	Terminal        string
	TerminalWorkDir string

	Reload [2]string
	Quit   [2]string

	ListUp   [2]string
	ListDown [2]string

	OpenTerminal [2]string

	PinnedDirectory [2]string

	CloseFilePanel     [2]string
	CreateNewFilePanel [2]string

	NextFilePanel     [2]string
	PreviousFilePanel [2]string
	FocusOnProcessBar [2]string
	FocusOnSideBar    [2]string
	FocusOnMetaData   [2]string

	ChangePanelMode [2]string

	FilePanelDirectoryCreate [2]string
	FilePanelFileCreate      [2]string
	FilePanelItemRename      [2]string
	PasteItem                [2]string
	ToggleDotFile            [2]string

	Cancel  [2]string
	Confirm [2]string

	DeleteItem      [2]string
	SelectItem      [2]string
	ParentDirectory [2]string
	CopySingleItem  [2]string
	CutSingleItem   [2]string

	FilePanelSelectModeItemSingleSelect [2]string
	FilePanelSelectModeItemSelectDown   [2]string
	FilePanelSelectModeItemSelectUp     [2]string
	FilePanelSelectModeItemDelete       [2]string
	FilePanelSelectModeItemCopy         [2]string
	FilePanelSelectModeItemCut          [2]string
	FilePanelSelectAllItem              [2]string
}
