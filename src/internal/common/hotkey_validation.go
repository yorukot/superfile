package common

import (
	"strings"
	"unicode"
)

// Accepted modifier prefixes for search_toggle_hidden.
var searchToggleHiddenModifiers = []string{"ctrl", "alt", "shift", "super", "meta", "hyper"} //nolint:gochecknoglobals // static allowlist

// Single-rune ctrl bases with ASCII control codes. Letters handled separately.
const ctrlBaseExceptions = "@[\\]^_?" //nolint:gochecknoglobals // fixed control-code symbol set

// IsModifierHotkey reports whether s is a modifier combo like "alt+.".
// Shape only. See ValidateSearchToggleHidden for terminal deliverability.
func IsModifierHotkey(s string) bool {
	_, _, ok := parseModifierHotkey(s)
	return ok
}

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

// isCtrlDeliverableBase reports whether base works as ctrl+base without Kitty protocol.
func isCtrlDeliverableBase(base string) bool {
	if len(base) != 1 {
		return true
	}
	r := rune(base[0])
	return unicode.IsLetter(r) || strings.ContainsRune(ctrlBaseExceptions, r)
}

// ValidateSearchToggleHidden returns the first rejected binding, or "" if all valid.
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
