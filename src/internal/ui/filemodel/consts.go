package filemodel

import (
	"errors"

	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

const (
	FileModelMinWidth  = filepanel.FilePanelMinWidth
	FileModelMinHeight = filepanel.FilePanelMinHeight
)

var ErrMaximumPanelCount = errors.New("maximum panel count reached")
