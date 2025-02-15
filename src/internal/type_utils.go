package internal

import "log/slog"

// reset the items slice and set the cut value
func (c *copyItems) reset(cut bool) {
	c.cut = cut
	c.items = c.items[:0]
}

// ================ Sidebar related =====================
// Hopefully compiler inlines it
func (d directory) isDivider() bool {
	return d == pinnedDivider || d == diskDivider
}
func (d directory) requiredHeight() int {
	if d.isDivider() {
		return 3
	}
	return 1
}

// True of only divicers are in directories slice, 
// but no actual directories
// This will be pretty quick. But we can replace it with
// len(s.directories) <= 2 - More hacky and hardcoded-like, but faster
func (s *sidebarModel) noActualDir() bool {
	for _,d := range s.directories {
		if !d.isDivider() {
			return false
		}
	}
	return true
}


// Return till what indexes we will render, if we start from startIndex
// if returned value is `startIndex - 1`, that means nothing can be rendered 
// This could be made constant time by keeping Indexes ot special directories saved, 
// but that too much.
func (s *sidebarModel) lastRenderedIndex(mainPanelHeight int, startIndex int) int {
	
	curHeight := sideBarInitialHeight
	endIndex := startIndex - 1
	for i := startIndex; i<len(s.directories); i++ {
		curHeight += s.directories[i].requiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		endIndex = i
	}
	return endIndex
}

// Return what will be the startIndex, if we end at endIndex
// if returned value is `endIndex + 1`, that means nothing can be rendered 
func (s *sidebarModel) firstRenderedIndex(mainPanelHeight int, endIndex int) int {
	// This should ideally never happen. Maybe we should panic ?
	if endIndex >= len(s.directories) {
		return endIndex + 1
	}
	
	curHeight := sideBarInitialHeight
	startIndex := endIndex + 1
	for i := endIndex; i>=0; i-- {
		curHeight += s.directories[i].requiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		startIndex = i
	}
	return startIndex
}

func (s *sidebarModel) updateRenderIndex(mainPanelHeight int) {
	// Case I : New cursor moved above current renderable range
	if s.cursor < s.renderIndex {
		// We will start rendering from there
		s.renderIndex = s.cursor
		return
	}

	curEndIndex := s.lastRenderedIndex(mainPanelHeight, s.renderIndex)

	// Case II : new cursor also comes in range of rendered directores
	// Taking this case later avoid extra lastRenderedIndex() call
	if s.renderIndex <= s.cursor && s.cursor <= curEndIndex {
		// no need to update s.renderIndex
		return 
	}

	// Case III : New cursor is too below
	if curEndIndex < s.cursor {
		s.renderIndex = s.firstRenderedIndex(mainPanelHeight, s.cursor)
		return 
	} 

	// Code should never reach here
	slog.Error("Unexpected situation in updateRenderIndex", "cursor", s.cursor, 
		"renderIndex", s.renderIndex, "directory count", len(s.directories) )
}

// ================ String method for easy logging =====================

func (f focusPanelType) String() string {
	switch f {
	case nonePanelFocus:
		return "nonePanelFocus"
	case processBarFocus:
		return "processBarFocus"
	case sidebarFocus:
		return "sidebarFocus"
	case metadataFocus:
		return "metadataFocus"
	default:
		return "Invalid"
	}
}

func (f filePanelFocusType) String() string {
	switch f {
	case noneFocus:
		return "noneFocus"
	case secondFocus:
		return "secondFocus"
	case focus:
		return "focus"
	default:
		return "Invalid"
	}
}

func (p panelMode) String() string {
	switch p {
	case selectMode:
		return "selectMode"
	case browserMode:
		return "browserMode"
	default:
		return "Invalid"
	}
}
