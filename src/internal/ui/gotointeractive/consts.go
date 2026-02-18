package gotointeractive

import (
	"github.com/yorukot/superfile/src/config/icon"
)

const (
	GotoMinWidth      = 80
	GotoMinHeight     = 10
	maxVisibleResults = 15
	modalInputPadding = 5
	iconColumnWidth   = 4
	gotoHeadlineText  = "Goto:"
)

func getGotoPrompt() string {
	return icon.Search + icon.Space + "Search: "
}
