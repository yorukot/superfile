package utils

import tea "charm.land/bubbletea/v2"

func TeaRuneKeyMsg(msg string) tea.KeyPressMsg {
	runes := []rune(msg)
	if len(runes) == 1 {
		return tea.KeyPressMsg{Code: runes[0], Text: msg}
	}
	return tea.KeyPressMsg{Code: tea.KeyExtended, Text: msg}
}

// HotkeyMatches reports whether msg matches any configured hotkey string.
// Bubbletea returns shifted printable keys via String() as their text form
// (for example "J"), while hotkey config often uses the keystroke form
// (for example "shift+j"). Both representations are checked.
func HotkeyMatches(msg tea.KeyPressMsg, hotkeys []string) bool {
	return HotkeyMatchesStrings(msg.String(), msg.Keystroke(), hotkeys)
}

// HotkeyMatchesStrings reports whether keyStr or keystroke matches any hotkey.
func HotkeyMatchesStrings(keyStr, keystroke string, hotkeys []string) bool {
	for _, hotkey := range hotkeys {
		if hotkey == "" {
			continue
		}
		if hotkey == keyStr || hotkey == keystroke {
			return true
		}
	}
	return false
}
