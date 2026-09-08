package common

import (
	"strings"
	"unicode"
)

// searchToggleHiddenModifiers are the modifier prefixes accepted for
// search_toggle_hidden bindings, matching bubbletea/ultraviolet keystroke
// prefixes (ctrl+alt+shift+meta+hyper+super+).
var searchToggleHiddenModifiers = []string{"ctrl", "alt", "shift", "super", "meta", "hyper"} //nolint:gochecknoglobals // static allowlist

// ctrlBaseExceptions are single-rune base keys that do have legacy ASCII
// control-code equivalents (and so are deliverable as ctrl combos without
// Kitty/ModifyOtherKeys). Letters are handled via unicode.IsLetter.
const ctrlBaseExceptions = "@[\\]^_?" //nolint:gochecknoglobals // fixed control-code symbol set

// IsModifierHotkey reports whether s is a modifier combo (e.g. "alt+.").
// Matching is case-insensitive. Empty strings return false; callers skip
// "" padding slots before calling. Strings ending with "+" (e.g. "ctrl+",
// "ctrl+alt+") and modifier-only chains (e.g. "ctrl+alt") are rejected:
// every "+"-separated segment except the last must be a known modifier and
// the last segment must be a non-empty, non-modifier key.
//
// Note: this is a pure shape check. Whether the combo is actually
// deliverable by most terminals (see ValidateSearchToggleHidden) is a
// separate policy.
func IsModifierHotkey(s string) bool {
	_, _, ok := parseModifierHotkey(s)
	return ok
}

// parseModifierHotkey splits s into its modifier segments and base key.
// ok is false unless every segment is non-empty, every leading segment is a
// known modifier, and the last segment is a non-modifier key.
func parseModifierHotkey(s string) (mods []string, base string, ok bool) {
	lower := strings.ToLower(strings.TrimSpace(s))
	parts := strings.Split(lower, "+")
	if len(parts) < 2 {
		return nil, "", false
	}
	for _, part := range parts {
		if part == "" {
			return nil, "", false
		}
	}
	for _, part := range parts[:len(parts)-1] {
		if !isHotkeyModifier(part) {
			return nil, "", false
		}
	}
	base = parts[len(parts)-1]
	if isHotkeyModifier(base) {
		return nil, "", false
	}
	return parts[:len(parts)-1], base, true
}

func isHotkeyModifier(s string) bool {
	for _, mod := range searchToggleHiddenModifiers {
		if s == mod {
			return true
		}
	}
	return false
}

// isCtrlDeliverableBase reports whether base works as a ctrl-combo key on
// terminals without Kitty/ModifyOtherKeys: letters and the few symbols with
// ASCII control-code equivalents, plus multi-rune named keys (f1, enter,
// ...) which travel as escape sequences. Single-rune punctuation and digits
// (e.g. ".", "/") have no control-code equivalent: most terminals deliver
// them as the bare key, so the binding would type into the search query
// instead of toggling.
func isCtrlDeliverableBase(base string) bool {
	if len(base) != 1 {
		return true
	}
	r := rune(base[0])
	return unicode.IsLetter(r) || strings.ContainsRune(ctrlBaseExceptions, r)
}

// ValidateSearchToggleHidden returns the first non-empty binding that must
// be rejected, or "" when all bindings are valid. Empty "" slots (the
// ['alt+.', ''] convention) are skipped. A binding is rejected when it is
// not a modifier combo, or when it needs ctrl with a base key most
// terminals cannot deliver (e.g. "ctrl+.").
func ValidateSearchToggleHidden(bindings []string) string {
	for _, binding := range bindings {
		if binding == "" {
			continue
		}
		mods, base, ok := parseModifierHotkey(binding)
		if !ok {
			return binding
		}
		if slicesContain(mods, "ctrl") && !isCtrlDeliverableBase(base) {
			return binding
		}
	}
	return ""
}

func slicesContain(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}

// SearchToggleHiddenErrorMessage explains why offender was rejected: either
// it is not a modifier combo (it would steal query characters), or it is a
// ctrl combo most terminals cannot deliver (it would type into the query).
func SearchToggleHiddenErrorMessage(offender string) string {
	if mods, base, ok := parseModifierHotkey(offender); ok &&
		slicesContain(mods, "ctrl") && !isCtrlDeliverableBase(base) {
		return "\"" + offender + "\" is invalid : " +
			"ctrl+" + base + " is not delivered by most terminals without Kitty keyboard " +
			"protocol (it types \"" + base + "\" into the search query instead of toggling); " +
			"use \"alt+" + base + "\" or a ctrl+letter combo instead."
	}
	return "\"" + offender + "\" is invalid : " +
		"search_toggle_hidden must use a modifier combo (e.g. alt+.) " +
		"so it does not steal characters from the recursive search query."
}
