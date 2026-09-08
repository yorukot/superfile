package common

import "strings"

// searchToggleHiddenModifiers are the modifier prefixes accepted for
// search_toggle_hidden bindings, matching bubbletea/ultraviolet keystroke
// prefixes (ctrl+alt+shift+meta+hyper+super+).
var searchToggleHiddenModifiers = []string{"ctrl", "alt", "shift", "super", "meta", "hyper"} //nolint:gochecknoglobals // static allowlist

// IsModifierHotkey reports whether s is a modifier combo (e.g. "ctrl+.").
// Matching is case-insensitive. Empty strings return false; callers skip
// "" padding slots before calling.
func IsModifierHotkey(s string) bool {
	lower := strings.ToLower(strings.TrimSpace(s))
	for _, mod := range searchToggleHiddenModifiers {
		if strings.HasPrefix(lower, mod+"+") && len(lower) > len(mod)+1 {
			return true
		}
	}
	return false
}

// ValidateSearchToggleHidden returns the first non-empty binding that is not
// a modifier combo, or "" when all bindings are valid. Empty "" slots (the
// ['ctrl+.', ''] convention) are skipped.
func ValidateSearchToggleHidden(bindings []string) string {
	for _, binding := range bindings {
		if binding == "" {
			continue
		}
		if !IsModifierHotkey(binding) {
			return binding
		}
	}
	return ""
}
