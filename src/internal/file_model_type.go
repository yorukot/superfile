package internal

import (
	"github.com/yorukot/superfile/src/internal/ui/filepanel"
	"github.com/yorukot/superfile/src/internal/ui/preview"
)

/* FILE WINDOWS TYPE START*/
// Model for file windows
type FileModel struct {
	FilePanels        []filepanel.Model
	SinglePanelWidth  int
	Width             int
	Height            int
	Renaming          bool
	MaxFilePanel      int
	FilePreview       preview.Model
	FocusedPanelIndex int
}
