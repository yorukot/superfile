package internal

import "log/slog"

// True ff only dividers are in directories slice,
// but no actual directories
// This will be pretty quick. But we can replace it with
// len(s.directories) <= 2 - More hacky and hardcoded-like, but faster
func (s *SidebarModel) NoActualDir() bool {
	for _, d := range s.Directories {
		if !d.isDivider() {
			return false
		}
	}
	return true
}

func (s *SidebarModel) IsCursorInvalid() bool {
	return s.Cursor < 0 || s.Cursor >= len(s.Directories) || s.Directories[s.Cursor].isDivider()
}

func (s *SidebarModel) ResetCursor() {
	s.Cursor = 0
	// Move to first non Divider dir
	for i, d := range s.Directories {
		if !d.isDivider() {
			s.Cursor = i
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
func (s *SidebarModel) LastRenderedIndex(mainPanelHeight int, startIndex int) int {
	curHeight := SideBarInitialHeight
	endIndex := startIndex - 1
	for i := startIndex; i < len(s.Directories); i++ {
		curHeight += s.Directories[i].requiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		endIndex = i
	}
	return endIndex
}

// Return what will be the startIndex, if we end at endIndex
// if returned value is `endIndex + 1`, that means nothing can be rendered
func (s *SidebarModel) FirstRenderedIndex(mainPanelHeight int, endIndex int) int {
	// This should ideally never happen. Maybe we should panic ?
	if endIndex >= len(s.Directories) {
		return endIndex + 1
	}

	curHeight := SideBarInitialHeight
	startIndex := endIndex + 1
	for i := endIndex; i >= 0; i-- {
		curHeight += s.Directories[i].requiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		startIndex = i
	}
	return startIndex
}

func (s *SidebarModel) UpdateRenderIndex(mainPanelHeight int) {
	// Case I : New cursor moved above current renderable range
	if s.Cursor < s.RenderIndex {
		// We will start rendering from there
		s.RenderIndex = s.Cursor
		return
	}

	curEndIndex := s.LastRenderedIndex(mainPanelHeight, s.RenderIndex)

	// Case II : new cursor also comes in range of rendered directories
	// Taking this case later avoid extra lastRenderedIndex() call
	if s.RenderIndex <= s.Cursor && s.Cursor <= curEndIndex {
		// no need to update s.renderIndex
		return
	}

	// Case III : New cursor is too below
	if curEndIndex < s.Cursor {
		s.RenderIndex = s.FirstRenderedIndex(mainPanelHeight, s.Cursor)
		return
	}

	// Code should never reach here
	slog.Error("Unexpected situation in updateRenderIndex", "cursor", s.Cursor,
		"renderIndex", s.RenderIndex, "directory count", len(s.Directories))
}

// ======================================== Sidebar controller ========================================

func (s *SidebarModel) ListUp(mainPanelHeight int) {
	slog.Debug("controlListUp called", "cursor", s.Cursor,
		"renderIndex", s.RenderIndex, "directory count", len(s.Directories))
	if s.NoActualDir() {
		return
	}
	if s.Cursor > 0 {
		// Not at the top, can safely decrease
		s.Cursor--
	} else {
		// We are at the top. Move to the bottom
		s.Cursor = len(s.Directories) - 1
	}
	// We should update even if cursor is at divider for now
	// Otherwise dividers are sometimes skipped in render in case of
	// large pinned directories
	s.UpdateRenderIndex(mainPanelHeight)
	if s.Directories[s.Cursor].isDivider() {
		// cause another listUp trigger to move up.
		s.ListUp(mainPanelHeight)
	}
}

func (s *SidebarModel) ListDown(mainPanelHeight int) {
	slog.Debug("controlListDown called", "cursor", s.Cursor,
		"renderIndex", s.RenderIndex, "directory count", len(s.Directories))
	if s.NoActualDir() {
		return
	}
	if s.Cursor < len(s.Directories)-1 {
		// Not at the bottom, can safely increase
		s.Cursor++
	} else {
		// We are at the bottom. Move to the top
		s.Cursor = 0
	}

	// We should update even if cursor is at divider for now
	// Otherwise dividers are sometimes skipped in render in case of
	// large pinned directories
	s.UpdateRenderIndex(mainPanelHeight)

	// Move below special divider directories
	if s.Directories[s.Cursor].isDivider() {
		// cause another listDown trigger to move down.
		s.ListDown(mainPanelHeight)
	}
}
