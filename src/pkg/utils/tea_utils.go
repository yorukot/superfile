package utils

import tea "charm.land/bubbletea/v2"

func TeaRuneKeyMsg(msg string) tea.KeyPressMsg {
	runes := []rune(msg)
	if len(runes) == 1 {
		return tea.KeyPressMsg{Code: runes[0], Text: msg}
	}
	return tea.KeyPressMsg{Code: tea.KeyExtended, Text: msg}
}
