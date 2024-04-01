package components

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
)

type panelMode uint
type filePanelFocusType uint
type processState int
type focusPanelType int
type itemType int

const (
	newFile itemType = iota
	newFolder
)
const (
	nonePanelFocus focusPanelType = iota
	processBarFocus
	sideBarFocus
)

const (
	noneFocus filePanelFocusType = iota
	secondFocus
	focus
)

const (
	selectMode panelMode = iota
	browserMode
)

const (
	copying processState = iota
	deleting
	moving
	successful
	cancel
	failure
)

// main model
type model struct {
	fileModel           fileModel
	sideBarModel        sideBarModel
	processBarModel     processBarModel
	focusPanel          focusPanelType
	copyItems           copyItems
	createNewItem       createNewItemModal
	filePanelFocusIndex int
	mainPanelHeight     int
	fullWidth           int
	fullHeight          int
}

type createNewItemModal struct {
	location  string
	open      bool
	itemType  itemType
	textInput textinput.Model
}

type copyItems struct {
	items         []string
	oringnalPanel orignalPanel
	cut           bool
}

type orignalPanel struct {
	index    int
	location string
}

/* FILE WINDOWS TYPE START*/
type fileModel struct {
	filePanels []filePanel
	width      int
}

type filePanel struct {
	cursor       int
	render       int
	focusType    filePanelFocusType
	location     string
	panelMode    panelMode
	selected     []string
	element      []element
	folderRecord map[string]folderRecord
}

type folderRecord struct {
	folderCursor int
	folderRender int
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

type processBarModel struct {
	cursor  int
	process []process
}

type process struct {
	name     string
	progress progress.Model
	state    processState
}

/*PROCESS BAR COMPONENTS TYPE END*/

type iconStyle struct {
	icon  string
	color string
}

type ThemeType struct {
	Border string
	Cursor string

	TerminalTooSmallError string
	TerminalSizeCurrect   string

	BrowserMode string
	SelectMode  string

	SideBarTitle    string
	SideBarItem     string
	SideBarSelected string
	SideBarFocus    string

	FilePanelFocus         string
	FilePanelTopFolderIcon string
	FilePanelTopPath       string
	FilePanelItem          string
	FilePanelItemSelected  string

	BottomBarFocus string

	ProcessBarSideLine string
	ProcessBarGradient [2]string
	Done               string
	Fail               string
	Cancel             string

	ModalForeground string
	ModalCancel     string
	ModalConfirm    string
}

type ConfigType struct {
	TrashCanPath string

	// HotKey setting
	Quit     [2]string
	ListUp   [2]string
	ListDown [2]string

	NextFilePanel      [2]string
	PreviousFilePanel  [2]string
	CloseFilePanel     [2]string
	CreateNewFilePanel [2]string
	FocusOnSideBar     [2]string
	FocusOnProcessBar  [2]string

	ChangePanelMode [2]string

	PasteItem [2]string

	Cancel  [2]string
	Confirm [2]string

	DeleteItem     [2]string
	SelectItem     [2]string
	ParentFolder   [2]string
	CopySingleItem [2]string
	CutSingleItem  [2]string

	FilePanelFolderCreate [2]string
	FilePanelFileCreate   [2]string

	FilePanelSelectModeItemSingleSelect [2]string
	FilePanelSelectModeItemSelectDown   [2]string
	FilePanelSelectModeItemSelectUp     [2]string
	FilePanelSelectModeItemDelete       [2]string
	FilePanelSelectModeItemCopy         [2]string
	FilePanelSelectModeItemPast         [2]string
	FilePanelSelectModeItemCut          [2]string
	FilePanelSelectAllItem              [2]string
}
