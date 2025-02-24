package internal

import "log/slog"

// True of only divicers are in directories slice,
// but no actual directories
// This will be pretty quick. But we can replace it with
// len(s.directories) <= 2 - More hacky and hardcoded-like, but faster
func (s *sidebarModel) noActualDir() bool {
	for _, d := range s.directories {
		if !d.isDivider() {
			return false
		}
	}
	return true
}

func (s *sidebarModel) isCursorInvalid() bool {
	return s.cursor < 0 || s.cursor >= len(s.directories) || s.directories[s.cursor].isDivider()
}

func (s *sidebarModel) resetCursor() {
	s.cursor = 0
	// Move to first non Divider dir
	for i, d := range s.directories {
		if !d.isDivider() {
			s.cursor = i
			return
		}
	}
	// If all directories are divider, code will reach here. and s.cursor will stay 0
	// Or s.directories is empty
}

// Return till what indexes we will render, if we start from startIndex
// if returned value is `startIndex - 1`, that means nothing can be rendered
// This could be made constant time by keeping Indexes ot special directories saved,
// but that too much.
func (s *sidebarModel) lastRenderedIndex(mainPanelHeight int, startIndex int) int {

	curHeight := sideBarInitialHeight
	endIndex := startIndex - 1
	for i := startIndex; i < len(s.directories); i++ {
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
	for i := endIndex; i >= 0; i-- {
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
		"renderIndex", s.renderIndex, "directory count", len(s.directories))
}

// ======================================== Sidebar controller ========================================

func (s *sidebarModel) controlListUp(wheel bool, mainPanelHeight int) {

	// Todo : This snippet is duplicated everywhere. It can be better refractored outside
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}
	for i := 0; i < runTime; i++ {
		s.listUp(mainPanelHeight)
	}
}

func (s *sidebarModel) listUp(mainPanelHeight int) {
	slog.Debug("controlListUp called", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "directory count", len(s.directories))
	if s.noActualDir() {
		return
	}
	if s.cursor > 0 {
		// Not at the top, can safely decrease
		s.cursor--
	} else {
		// We are at the top. Move to the bottom
		s.cursor = len(s.directories) - 1
	}
	// We should update even if cursor is at divider for now
	// Otherwise dividers are sometimes skipped in render in case of
	// large pinned directories
	s.updateRenderIndex(mainPanelHeight)
	if s.directories[s.cursor].isDivider() {
		// cause another listUp trigger to move up.
		s.listUp(mainPanelHeight)
	}

}

func (s *sidebarModel) controlListDown(wheel bool, mainPanelHeight int) {

	// Todo : This snippet is duplicated everywhere. It can be better refractored outside
	runTime := 1
	if wheel {
		runTime = wheelRunTime
	}
	for i := 0; i < runTime; i++ {
		s.listDown(mainPanelHeight)
	}
}

func (s *sidebarModel) listDown(mainPanelHeight int) {
	slog.Debug("controlListDown called", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "directory count", len(s.directories))
	if s.noActualDir() {
		return
	}
	if s.cursor < len(s.directories)-1 {
		// Not at the bottom, can safely increase
		s.cursor++
	} else {
		// We are at the bottom. Move to the top
		s.cursor = 0
	}

	// We should update even if cursor is at divider for now
	// Otherwise dividers are sometimes skipped in render in case of
	// large pinned directories
	s.updateRenderIndex(mainPanelHeight)

	// Move below special divider directories
	if s.directories[s.cursor].isDivider() {
		// cause another listDown trigger to move down.
		s.listDown(mainPanelHeight)
	}
}
