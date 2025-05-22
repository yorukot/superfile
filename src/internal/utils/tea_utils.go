package utils

import tea "github.com/charmbracelet/bubbletea"

func TeaRuneKeyMsg(msg string) tea.KeyMsg {
	return tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(msg),
	}
}
