package filepanel

import (
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

// TODO: Convert to integer enum
type sortingKind string

// Make sure to use New() to ensure that maps are initialized
// zero value `Model{}`, or direct initialization should be avoided
// or used very carefully if needed
type Model struct {
	Cursor      int
	RenderIndex int
	IsFocused   bool
	Location    string
	// Dimension fields
	width  int // Total width including borders
	height int // Total height including borders
	// TODO: Every file panel doesn't needs sort options model
	// They just need to store their current sort config.
	SortOptions sortOptionsModel
	PanelMode   PanelMode
	// key is file location, value order of selection
	selected           map[string]int
	selectOrderCounter int
	Element            []Element
	DirectoryRecords   map[string]directoryRecord
	Rename             textinput.Model
	Renaming           bool
	SearchBar          textinput.Model
	LastTimeGetElement time.Time
	TargetFile         string             // filename to position cursor on after load
	columns            []columnDefinition // columns for rendering
}

// Sort options
type sortOptionsModel struct {
	Width  int
	Height int
	Open   bool
	Cursor int
	Data   sortOptionsModelData
}

type sortOptionsModelData struct {
	Options  []string
	Selected int
	Reversed bool
}

// Record for directory navigation
type directoryRecord struct {
	directoryCursor int
	directoryRender int
}

// Element within a file panel
type Element struct {
	Name      string
	Location  string
	Directory bool
	Info      os.FileInfo
}

// Type representing the mode of the panel
type PanelMode uint

// Constants for select mode or browser mode
const (
	SelectMode PanelMode = iota
	BrowserMode
)

type sliceOrderFunc func(i, j int) bool

// Note: There are here, instead of consts.go as they are definitions of the enum `sortingKind`
const (
	sortingName         sortingKind = "Name"
	sortingSize         sortingKind = "Size"
	sortingDateModified sortingKind = "Date Modified"
	sortingFileType     sortingKind = "Type"
)

type columnRenderer func(indexElement int, columnWidth int) string

type columnDefinition struct {
	Name         string
	Size         int
	columnRender columnRenderer
}
