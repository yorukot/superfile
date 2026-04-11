package utils

import tea "charm.land/bubbletea/v2"

func TeaRuneKeyMsg(msg string) tea.KeyMsg {
	return tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(msg),
	}
}
