package utils

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
)

func TestHotkeyMatchesShiftedPrintableKeys(t *testing.T) {
	shiftJ := tea.KeyPressMsg{Code: 'j', Text: "J", Mod: tea.ModShift}
	shiftK := tea.KeyPressMsg{Code: 'k', Text: "K", Mod: tea.ModShift}

	assert.True(t, HotkeyMatches(shiftJ, []string{"shift+j"}))
	assert.True(t, HotkeyMatches(shiftK, []string{"shift+k"}))
	assert.True(t, HotkeyMatches(shiftJ, []string{"J"}))
	assert.True(t, HotkeyMatches(shiftK, []string{"K"}))
	assert.False(t, HotkeyMatches(shiftJ, []string{"j"}))
}

func TestHotkeyMatchesSpecialKeys(t *testing.T) {
	shiftUp := tea.KeyPressMsg{Code: tea.KeyUp, Mod: tea.ModShift}
	ctrlJ := tea.KeyPressMsg{Code: 'j', Mod: tea.ModCtrl}

	assert.True(t, HotkeyMatches(shiftUp, []string{"shift+up"}))
	assert.True(t, HotkeyMatches(ctrlJ, []string{"ctrl+j"}))
}

func TestHotkeyMatchesStrings(t *testing.T) {
	hotkeys := []string{"shift+j", "shift+down"}

	assert.True(t, HotkeyMatchesStrings("J", "shift+j", hotkeys))
	assert.True(t, HotkeyMatchesStrings("shift+down", "shift+down", hotkeys))
	assert.False(t, HotkeyMatchesStrings("j", "j", hotkeys))
}
