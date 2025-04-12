package sidebar

import (
	"log/slog"
)

func (s *Model) ListUp(mainPanelHeight int) {
	slog.Debug("controlListUp called", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "directory count", len(s.directories))
	if s.NoActualDir() {
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
	if s.directories[s.cursor].IsDivider() {
		// cause another listUp trigger to move up.
		s.ListUp(mainPanelHeight)
	}
}

func (s *Model) ListDown(mainPanelHeight int) {
	slog.Debug("controlListDown called", "cursor", s.cursor,
		"renderIndex", s.renderIndex, "directory count", len(s.directories))
	if s.NoActualDir() {
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
	if s.directories[s.cursor].IsDivider() {
		// cause another listDown trigger to move down.
		s.ListDown(mainPanelHeight)
	}
}

// Return till what indexes we will render, if we start from startIndex
// if returned value is `startIndex - 1`, that means nothing can be rendered
// This could be made constant time by keeping Indexes ot special directories saved,
// but that too much.
func (s *Model) lastRenderedIndex(mainPanelHeight int, startIndex int) int {
	curHeight := SideBarInitialHeight
	endIndex := startIndex - 1
	for i := startIndex; i < len(s.directories); i++ {
		curHeight += s.directories[i].RequiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		endIndex = i
	}
	return endIndex
}

// Return what will be the startIndex, if we end at endIndex
// if returned value is `endIndex + 1`, that means nothing can be rendered
func (s *Model) firstRenderedIndex(mainPanelHeight int, endIndex int) int {
	// This should ideally never happen. Maybe we should panic ?
	if endIndex >= len(s.directories) {
		return endIndex + 1
	}

	curHeight := SideBarInitialHeight
	startIndex := endIndex + 1
	for i := endIndex; i >= 0; i-- {
		curHeight += s.directories[i].RequiredHeight()
		if curHeight > mainPanelHeight {
			break
		}
		startIndex = i
	}
	return startIndex
}

func (s *Model) updateRenderIndex(mainPanelHeight int) {
	// Case I : New cursor moved above current renderable range
	if s.cursor < s.renderIndex {
		// We will start rendering from there
		s.renderIndex = s.cursor
		return
	}

	curEndIndex := s.lastRenderedIndex(mainPanelHeight, s.renderIndex)

	// Case II : new cursor also comes in range of rendered directories
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
