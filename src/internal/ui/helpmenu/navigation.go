package helpmenu

import "github.com/yorukot/superfile/src/internal/common"

// Help menu panel list up
func (m *Model) ListUp() {
	if m.cursor > 1 {
		m.cursor--
		if m.cursor < m.renderIndex {
			m.renderIndex = m.cursor
		}
		if m.filteredData[m.cursor].subTitle != "" {
			m.cursor--
		}
	} else {
		// Set the cursor to the last item in the list.
		// We use max(..., 0) as a safeguard to prevent a negative cursor index
		// in case the filtered list is empty.
		m.cursor = max(len(m.filteredData)-1, 0)

		// Adjust the render index to show the bottom of the list.
		// Similarly, we use max(..., 0) to ensure the renderIndex doesn't become negative,
		// which can happen if the number of items is less than the view height.
		// This prevents a potential out-of-bounds panic during rendering.
		m.renderIndex = max(len(m.filteredData)-(m.height-common.InnerPadding), 0)
	}
}

// Help menu panel list down
func (m *Model) ListDown() {
	if len(m.filteredData) == 0 {
		return
	}

	if m.cursor < len(m.filteredData)-1 {
		// Compute the next selectable row (skip subtitles).
		next := m.cursor + 1
		for next < len(m.filteredData) && m.filteredData[next].subTitle != "" {
			next++
		}
		if next >= len(m.filteredData) {
			// Wrap if no more selectable rows.
			m.cursor = 1
			m.renderIndex = 0
			return
		}
		m.cursor = next

		// Scroll down if cursor moved past the viewport.
		if m.cursor > m.renderIndex+m.height-5 {
			m.renderIndex++
		}
		// Clamp renderIndex to bottom.
		bottom := len(m.filteredData) - (m.height - common.InnerPadding)
		if bottom < 0 {
			bottom = 0
		}
		if m.renderIndex > bottom {
			m.renderIndex = bottom
		}
	} else {
		m.cursor = 1
		m.renderIndex = 0
	}
}
