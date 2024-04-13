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
	sidebarFocus
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
	sidebarModel        sidebarModel
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
type sidebarModel struct {
	directories []directory
	// wellKnownModel []directory
	// pinnedModel    []directory
	// disksModel     []directory
	cursor int
}

type directory struct {
	location string
	name     string
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
	messageId       string
	processNewState process
	returnWarnModal bool
	warnModal       warnModal
	loadMetadata    bool
	metadata        [][2]string
}

/*PROCESS BAR COMPONENTS TYPE END*/

// Style for icons
type iconStyle struct {
	icon  string
	color string
}

// Theme configuration
type ThemeType struct {
	Border         string `toml:"border"`
	Cursor         string `toml:"cursor"`
	MainBackground string `toml:"main_background"`

	TerminalTooSmallError string `toml:"terminal_too_small_error"`
	TerminalSizeCorrect   string `toml:"terminal_size_correct"`

	SidebarTitle    string `toml:"sidebar_title"`
	SidebarItem     string `toml:"sidebar_item"`
	SidebarSelected string `toml:"sidebar_selected"`
	SidebarFocus    string `toml:"sidebar_focus"`

	FilePanelFocus            string `toml:"file_panel_focus"`
	FilePanelTopDirectoryIcon string `toml:"file_panel_top_directory_icon"`
	FilePanelTopPath          string `toml:"file_panel_top_path"`
	FilePanelItem             string `toml:"file_panel_item"`

	FilePanelItemSelected string    `toml:"file_panel_item_selected"`
	FooterFocus           string    `toml:"footer_focus"`
	ProcessBarGradient    [2]string `toml:"process_bar_gradient"`
	InOperation           string    `toml:"in_operation"`
	Done                  string    `toml:"done"`
	Fail                  string    `toml:"fail"`
	Cancel                string    `toml:"cancel"`
	ModalForeground       string    `toml:"modal_foreground"`
	ModalCancel           string    `toml:"modal_cancel"`
	ModalConfirm          string    `toml:"modal_confirm"`
}

// Configuration settings
type ConfigType struct {
	Theme    string
	Terminal string
}

type HotkeysType struct {
	CommitHotkey string `toml:"_COMMIT_HOTKEY"`

	CommitGlobalHotkey string   `toml:"_COMMIT_global_hotkey"`
	Quit               []string `toml:"quit"`

	ListUp   []string `toml:"list_up"`
	ListDown []string `toml:"list_down"`

	PinnedDirectory []string `toml:"pinned_directory"`

	CloseFilePanel           []string `toml:"close_file_panel"`
	CreateNewFilePanel       []string `toml:"create_new_file_panel"`
	NextFilePanel            []string `toml:"next_file_panel"`
	PreviousFilePanel        []string `toml:"previous_file_panel"`
	FocusOnProcessBar        []string `toml:"focus_on_process_bar"`
	FocusOnSideBar           []string `toml:"focus_on_side_bar"`
	FocusOnMetaData          []string `toml:"focus_on_meta_data"`
	ChangePanelMode          []string `toml:"change_panel_mode"`
	FilePanelDirectoryCreate []string `toml:"file_panel_directory_create"`
	FilePanelFileCreate      []string `toml:"file_panel_file_create"`
	FilePanelItemRename      []string `toml:"file_panel_item_rename"`
	PasteItem                []string `toml:"paste_item"`
	ExtractFile              []string `toml:"extract_file"`
	CompressFile             []string `toml:"compress_file"`
	ToggleDotFile            []string `toml:"toggle_dot_file"`

	CommitSpecialHotkey string   `toml:"_COMMIT_special_hotkey"`
	Cancel              []string `toml:"cancel"`
	Confirm             []string `toml:"confirm"`

	CommitNormalModeHotkey string   `toml:"_COMMIT_normal_mode_hotkey"`
	DeleteItem             []string `toml:"delete_item"`
	SelectItem             []string `toml:"select_item"`
	ParentDirectory        []string `toml:"parent_directory"`
	CopySingleItem         []string `toml:"copy_single_item"`
	CutSingleItem          []string `toml:"cut_single_item"`

	CommitSelectModeHotkey              string   `toml:"_COMMIT_select_mode_hotkey"`
	FilePanelSelectModeItemSingleSelect []string `toml:"file_panel_select_mode_item_single_select"`
	FilePanelSelectModeItemSelectDown   []string `toml:"file_panel_select_mode_item_select_down"`
	FilePanelSelectModeItemSelectUp     []string `toml:"file_panel_select_mode_item_select_up"`
	FilePanelSelectModeItemDelete       []string `toml:"file_panel_select_mode_item_delete"`
	FilePanelSelectModeItemCopy         []string `toml:"file_panel_select_mode_item_copy"`
	FilePanelSelectModeItemCut          []string `toml:"file_panel_select_mode_item_cut"`
	FilePanelSelectAllItem              []string `toml:"file_panel_select_all_item"`

	CommitProcessBarHotkey string `toml:"_COMMIT_process_bar_hotkey"`
}
