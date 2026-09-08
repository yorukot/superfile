package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsModifierHotkey(t *testing.T) {
	testdata := []struct {
		name     string
		binding  string
		expected bool
	}{
		{"ctrl dot", "ctrl+.", true},
		{"ctrl lowercase combo", "ctrl+n", true},
		{"alt combo", "alt+.", true},
		{"shift combo", "shift+f1", true},
		{"super combo", "super+.", true},
		{"meta combo", "meta+x", true},
		{"hyper combo", "hyper+x", true},
		{"uppercase prefix", "CTRL+.", true},
		{"dual modifier combo", "ctrl+alt+.", true},
		{"bare letter", "a", false},
		{"bare dot", ".", false},
		{"bare word", "enter", false},
		{"bare f-key", "f1", false},
		{"bare esc", "esc", false},
		{"bare arrow", "up", false},
		{"empty", "", false},
		{"modifier only no key", "ctrl+", false},
		{"modifier chain trailing plus", "ctrl+alt+", false},
		{"modifier-only chain", "ctrl+alt", false},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsModifierHotkey(tt.binding))
		})
	}
}

func TestValidateSearchToggleHidden(t *testing.T) {
	assert.Empty(t, ValidateSearchToggleHidden([]string{"ctrl+.", ""}))
	assert.Empty(t, ValidateSearchToggleHidden([]string{"ctrl+.", "alt+."}))
	assert.Empty(t, ValidateSearchToggleHidden([]string{"ctrl+alt+.", ""}))
	assert.Equal(t, "a", ValidateSearchToggleHidden([]string{"a", ""}))
	assert.Equal(t, ".", ValidateSearchToggleHidden([]string{"ctrl+.", "."}))
	assert.Equal(t, "f1", ValidateSearchToggleHidden([]string{"f1", ""}))
	assert.Equal(t, "ctrl+alt+", ValidateSearchToggleHidden([]string{"ctrl+alt+", ""}))
	assert.Empty(t, ValidateSearchToggleHidden([]string{"", ""}))
}
