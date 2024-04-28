package internal

// Toggle help menu
func openHelpMenu(m model) model {
	if m.helpMenu.open {
		m.helpMenu.open = false
		return m
	}

	m.helpMenu.open = true
	return m
}