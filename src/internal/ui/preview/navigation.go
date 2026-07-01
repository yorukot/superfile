package preview

import (
	"github.com/yorukot/superfile/src/internal/common"
)

func bulkScrollSize(viewportHeight int) int {
	pages := common.Config.PreviewScrollBulk
	if pages <= 0 {
		pages = 2
	}
	chunk := viewportHeight / pages
	if chunk <= 0 {
		return 1
	}
	return chunk
}

func (m *Model) ScrollLineUp() bool {
	if m.scrollOffset == 0 {
		return false
	}
	m.scrollOffset--
	return true
}

func (m *Model) ScrollLineDown() bool {
	if !m.canScrollDown {
		return false
	}
	m.scrollOffset++
	return true
}

func (m *Model) ScrollBulkUp(viewportHeight int) bool {
	if m.scrollOffset == 0 {
		return false
	}
	m.scrollOffset -= bulkScrollSize(viewportHeight)
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
	return true
}

func (m *Model) ScrollBulkDown(viewportHeight int) bool {
	if !m.canScrollDown {
		return false
	}
	m.scrollOffset += bulkScrollSize(viewportHeight)
	return true
}

func (m *Model) ScrollTop() bool {
	if m.scrollOffset == 0 {
		return false
	}
	m.scrollOffset = 0
	return true
}

func (m *Model) CanScrollDown() bool {
	return m.canScrollDown
}

func (m *Model) resetScroll() {
	m.scrollOffset = 0
	m.canScrollDown = false
}

func (m *Model) setScrollState(canScrollDown bool) {
	m.canScrollDown = canScrollDown
}

func (m *Model) clampScrollOffset(maxOffset int) {
	if maxOffset < 0 {
		maxOffset = 0
	}
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}
}
