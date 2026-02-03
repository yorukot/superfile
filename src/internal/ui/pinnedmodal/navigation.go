package pinnedmodal

func (m *Model) navigateUp() {
	if len(m.results) == 0 {
		return
	}
	if m.cursor > 0 {
		m.cursor--
		if m.cursor < m.renderIndex {
			m.renderIndex = m.cursor
		}
	} else {
		m.cursor = len(m.results) - 1
		if m.cursor >= m.renderIndex+maxVisibleResults {
			m.renderIndex = m.cursor - maxVisibleResults + 1
		}
	}
}

func (m *Model) navigateDown() {
	if len(m.results) == 0 {
		return
	}
	if m.cursor < len(m.results)-1 {
		m.cursor++
		if m.cursor >= m.renderIndex+maxVisibleResults {
			m.renderIndex = m.cursor - maxVisibleResults + 1
		}
	} else {
		m.cursor = 0
		m.renderIndex = 0
	}
}

func (m *Model) navigatePageUp() {
	if len(m.results) == 0 {
		return
	}
	scrollAmount := maxVisibleResults - 1
	if scrollAmount < 1 {
		scrollAmount = 1
	}

	if m.cursor-scrollAmount >= 0 {
		m.cursor -= scrollAmount
		if m.cursor < m.renderIndex {
			m.renderIndex = m.cursor
		}
	} else {
		m.cursor = 0
		m.renderIndex = 0
	}
}

func (m *Model) navigatePageDown() {
	if len(m.results) == 0 {
		return
	}
	scrollAmount := maxVisibleResults - 1
	if scrollAmount < 1 {
		scrollAmount = 1
	}

	if m.cursor+scrollAmount < len(m.results) {
		m.cursor += scrollAmount
		if m.cursor >= m.renderIndex+maxVisibleResults {
			m.renderIndex = m.cursor - maxVisibleResults + 1
		}
	} else {
		m.cursor = len(m.results) - 1
		if m.cursor >= m.renderIndex+maxVisibleResults {
			m.renderIndex = m.cursor - maxVisibleResults + 1
		}
	}
}
