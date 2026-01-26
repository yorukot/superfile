package bulkrename

func (m *Model) navigateCursor(delta int) {
	if delta > 0 {
		m.navigateDown()
	} else {
		m.navigateUp()
	}
}

func (m *Model) navigateUp() {
	if m.cursor > 0 {
		m.cursor--
		m.focusInput()
	}
}

func (m *Model) navigateDown() {
	if m.cursor < 1 {
		m.cursor++
		m.focusInput()
	}
}

func (m *Model) focusInput() {
	m.findInput.Blur()
	m.replaceInput.Blur()
	m.prefixInput.Blur()
	m.suffixInput.Blur()

	switch m.renameType {
	case FindReplace:
		if m.cursor == 0 {
			m.findInput.Focus()
		} else {
			m.replaceInput.Focus()
		}
	case AddPrefix:
		m.prefixInput.Focus()
	case AddSuffix:
		m.suffixInput.Focus()
	}
}

func (m *Model) nextType() {
	m.renameType = RenameType((int(m.renameType) + 1) % 6)
	m.focusInput()
	m.preview = nil
}

func (m *Model) prevType() {
	newType := int(m.renameType) - 1
	if newType < 0 {
		newType = 5
	}
	m.renameType = RenameType(newType)
	m.focusInput()
	m.preview = nil
}