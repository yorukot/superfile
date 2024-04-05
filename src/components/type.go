package components

import (
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
	metaDataFocus
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
	inOperation processState = iota
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
	fileMetaData        fileMetaData
	firstTextInput      bool
	filePanelFocusIndex int
	mainPanelHeight     int
	fullWidth           int
	fullHeight          int
}

type processBarMessage struct {
	processId       string
	processNewState process
}

type fileMetaData struct {
	metaData    [][2]string
	renderIndex int
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
	renaming   bool
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
	rename       textinput.Model
	renaming     bool
}

type folderRecord struct {
	folderCursor int
	folderRender int
}

type element struct {
	name     string
	location string
	folder   bool
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
	cursor      int
	processList []string
	process     map[string]process
}

type process struct {
	name     string
	progress progress.Model
	state    processState
	speed    string
	total     int
	done     int
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
	Reload [2]string

	Quit     [2]string
	ListUp   [2]string
	ListDown [2]string

	NextFilePanel      [2]string
	PreviousFilePanel  [2]string
	CloseFilePanel     [2]string
	CreateNewFilePanel [2]string

	ChangePanelMode [2]string

	FocusOnSideBar    [2]string
	FocusOnProcessBar [2]string
	FocusOnMetaData   [2]string

	PasteItem [2]string

	FilePanelFolderCreate [2]string
	FilePanelFileCreate   [2]string
	FilePanelItemRename   [2]string

	PinnedFolder [2]string

	Cancel  [2]string
	Confirm [2]string

	DeleteItem     [2]string
	SelectItem     [2]string
	ParentFolder   [2]string
	CopySingleItem [2]string
	CutSingleItem  [2]string

	FilePanelSelectModeItemSingleSelect [2]string
	FilePanelSelectModeItemSelectDown   [2]string
	FilePanelSelectModeItemSelectUp     [2]string
	FilePanelSelectModeItemDelete       [2]string
	FilePanelSelectModeItemCopy         [2]string
	FilePanelSelectModeItemPast         [2]string
	FilePanelSelectModeItemCut          [2]string
	FilePanelSelectAllItem              [2]string
}
