package filemodel

import (
	"errors"

	"github.com/yorukot/superfile/src/internal/ui/filepanel"
)

// Now they are same doesn't means that they will be forever.
// Explicitly stating here tells that they are derived from same
// source, but have inherently different meaning
const (
	FileModelMinHeight      = filepanel.MinHeight
	FileModelMinWidth       = filepanel.MinWidth
	FilePreviewResizingText = "Resizing..."
	FilePreviewLoadingText  = "Loading..."
)

var ErrMaximumPanelCount = errors.New("maximum panel count reached")

var ErrMinimumPanelCount = errors.New("minimum panel count reached")
