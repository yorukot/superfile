package pinnedmodal

func (m *Model) navigateUp() {
	if len(m.results) == 0 {
		return
	}
	if m.cursor > 0 {
		m.cursor--
	} else {
		m.cursor = len(m.results) - 1
	}
	m.updateRenderIndex()
}

func (m *Model) navigateDown() {
	if len(m.results) == 0 {
		return
	}
	if m.cursor < len(m.results)-1 {
		m.cursor++
	} else {
		m.cursor = 0
	}
	m.updateRenderIndex()
}

func (m *Model) updateRenderIndex() {
	if len(m.results) == 0 {
		m.renderIndex = 0
		return
	}

	if m.cursor < m.renderIndex {
		m.renderIndex = m.cursor
	}

	if m.cursor >= m.renderIndex+maxVisibleResults {
		m.renderIndex = m.cursor - maxVisibleResults + 1
	}

	if m.renderIndex < 0 {
		m.renderIndex = 0
	}
	maxRenderIndex := len(m.results) - maxVisibleResults
	if maxRenderIndex < 0 {
		maxRenderIndex = 0
	}
	if m.renderIndex > maxRenderIndex {
		m.renderIndex = maxRenderIndex
	}
}
