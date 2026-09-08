package common

import "strings"

// searchToggleHiddenModifiers are the modifier prefixes accepted for
// search_toggle_hidden bindings, matching bubbletea/ultraviolet keystroke
// prefixes (ctrl+alt+shift+meta+hyper+super+).
var searchToggleHiddenModifiers = []string{"ctrl", "alt", "shift", "super", "meta", "hyper"} //nolint:gochecknoglobals // static allowlist

// IsModifierHotkey reports whether s is a modifier combo (e.g. "ctrl+.").
// Matching is case-insensitive. Empty strings return false; callers skip
// "" padding slots before calling. Strings ending with "+" (e.g. "ctrl+",
// "ctrl+alt+") and modifier-only chains (e.g. "ctrl+alt") are rejected:
// every "+"-separated segment except the last must be a known modifier and
// the last segment must be a non-empty, non-modifier key.
func IsModifierHotkey(s string) bool {
	lower := strings.ToLower(strings.TrimSpace(s))
	parts := strings.Split(lower, "+")
	if len(parts) < 2 {
		return false
	}
	for _, part := range parts {
		if part == "" {
			return false
		}
	}
	for _, part := range parts[:len(parts)-1] {
		if !isHotkeyModifier(part) {
			return false
		}
	}
	return !isHotkeyModifier(parts[len(parts)-1])
}

func isHotkeyModifier(s string) bool {
	for _, mod := range searchToggleHiddenModifiers {
		if s == mod {
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
