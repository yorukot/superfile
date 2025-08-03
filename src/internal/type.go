package internal

import (
	"time"

	"github.com/yorukot/superfile/src/internal/ui/metadata"
	"github.com/yorukot/superfile/src/internal/ui/processbar"
	"github.com/yorukot/superfile/src/internal/ui/sidebar"
	filepreview "github.com/yorukot/superfile/src/pkg/file_preview"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/yorukot/superfile/src/internal/ui/prompt"
)

// Type representing the mode of the panel
type panelMode uint

// Type representing the focus type of the file panel
type filePanelFocusType uint

// Type representing the type of focused panel
type focusPanelType int

type warnType int

type hotkeyType int

type channelMessageType int

type modelQuitStateType int

const (
	globalType hotkeyType = iota
	normalType
	selectType
)

const (
	confirmDeleteItem warnType = iota
	confirmRenameItem
)

// Constants for panel with no focus
const (
	nonePanelFocus focusPanelType = iota
	processBarFocus
	sidebarFocus
	metadataFocus
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

const (
	sendWarnModal channelMessageType = iota
	sendNotifyModal
)

const (
	notQuitting modelQuitStateType = iota
	quitInitiated
	confirmToQuit
	quitDone
)

// Main model
// TODO : We could consider using *model as tea.Model, instead of model.
// for reducing re-allocations. The struct is 20K bytes. But this could lead to
// issues like race conditions and whatnot, which are hidden since we are creating
// new model in each tea update.
type model struct {
	// Main Panels
	fileModel       fileModel
	sidebarModel    sidebar.Model
	processBarModel processbar.Model
	focusPanel      focusPanelType
	copyItems       copyItems

	// Modals
	notifyModal notifyModal
	typingModal typingModal
	warnModal   warnModal
	helpMenu    helpMenuModal
	promptModal prompt.Model

	fileMetaData         metadata.Model
	ioReqCnt             int
	imagePreviewer       *filepreview.ImagePreviewer
	modelQuitState       modelQuitStateType
	firstTextInput       bool
	toggleDotFile        bool
	updatedToggleDotFile bool
	toggleFooter         bool
	firstLoadingComplete bool
	firstUse             bool
	filePanelFocusIndex  int

	// Height in number of lines of actual viewport of
	// main panel and sidebar excluding border
	mainPanelHeight int

	// Height in number of lines of actual viewport of
	// footer panels - process/metadata/clipboard - excluding border
	footerHeight int
	fullWidth    int
	fullHeight   int
}

// Modal
type helpMenuModal struct {
	height      int
	width       int
	open        bool
	renderIndex int
	cursor      int
	data        []helpMenuModalData
}

type helpMenuModalData struct {
	hotkey         []string
	description    string
	hotkeyWorkType hotkeyType
	subTitle       string
}

type warnModal struct {
	open     bool
	warnType warnType
	title    string
	content  string
}

type typingModal struct {
	location      string
	open          bool
	textInput     textinput.Model
	errorMesssage string
}

type notifyModal struct {
	open    bool
	title   string
	content string
}

// Copied items
type copyItems struct {
	items []string
	cut   bool
}

/* FILE WINDOWS TYPE START*/
// Model for file windows
type fileModel struct {
	filePanels   []filePanel
	width        int
	renaming     bool
	maxFilePanel int
	filePreview  filePreviewPanel
}

type filePreviewPanel struct {
	open  bool
	width int
}

// Panel representing a file
type filePanel struct {
	cursor             int
	render             int
	focusType          filePanelFocusType
	location           string
	sortOptions        sortOptionsModel
	panelMode          panelMode
	selected           []string
	element            []element
	directoryRecords   map[string]directoryRecord
	rename             textinput.Model
	renaming           bool
	searchBar          textinput.Model
	lastTimeGetElement time.Time
}

// Sort options
type sortOptionsModel struct {
	width  int
	height int
	open   bool
	cursor int
	data   sortOptionsModelData
}

type sortOptionsModelData struct {
	options  []string
	selected int
	reversed bool
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

/* SIDE BAR internal TYPE END*/

/*PROCESS BAR internal TYPE START*/

// Message for process bar
type channelMessage struct {
	messageID   string
	messageType channelMessageType
	// TODO : We will stop using channel for this, and use tea.Cmd
	warnModal   warnModal
	notifyModal notifyModal
}

/*PROCESS BAR internal TYPE END*/

type editorFinishedMsg struct{ err error }
