package zoxide

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestIsKeyAlphaNum(t *testing.T) {
	assert.True(t, isKeyAlphaNum(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}), "'j' should be alphanumeric")
	assert.True(t, isKeyAlphaNum(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}), "'k' should be alphanumeric")
	assert.True(t, isKeyAlphaNum(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'5'}}), "'5' should be alphanumeric")
	assert.True(t, isKeyAlphaNum(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}}), "'A' should be alphanumeric")
	assert.False(t, isKeyAlphaNum(tea.KeyMsg{Type: tea.KeyUp}), "up arrow should not be alphanumeric")
	assert.False(t, isKeyAlphaNum(tea.KeyMsg{Type: tea.KeyEnter}), "enter should not be alphanumeric")
	assert.False(
		t,
		isKeyAlphaNum(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}),
		"space should not be alphanumeric",
	)
}
