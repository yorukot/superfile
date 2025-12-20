package internal

import (
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

// Panel representing a file
type FilePanel struct {
	Cursor      int
	RenderIndex int
	IsFocused   bool
	Location    string
	// TODO: Every file panel doesn't needs sort options model
	// They just need to store their current sort config.
	SortOptions        SortOptionsModel
	PanelMode          PanelMode
	Selected           []string
	Element            []Element
	DirectoryRecords   map[string]DirectoryRecord
	Rename             textinput.Model
	Renaming           bool
	SearchBar          textinput.Model
	LastTimeGetElement time.Time
	TargetFile         string // filename to position cursor on after load
}

// Sort options
type SortOptionsModel struct {
	Width  int
	Height int
	Open   bool
	Cursor int
	Data   SortOptionsModelData
}

type SortOptionsModelData struct {
	Options  []string
	Selected int
	Reversed bool
}

// Record for directory navigation
type DirectoryRecord struct {
	DirectoryCursor int
	DirectoryRender int
}

// Element within a file panel
type Element struct {
	Name      string
	Location  string
	Directory bool
	MetaData  [][2]string
	Info      os.FileInfo
}

// Type representing the mode of the panel
type PanelMode uint

// Constants for select mode or browser mode
const (
	selectMode PanelMode = iota
	browserMode
)
