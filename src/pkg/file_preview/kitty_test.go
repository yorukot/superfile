package filepreview

import "testing"

func TestMatchesTerminal(t *testing.T) {
	tests := []struct {
		name           string
		termProgram    string
		term           string
		knownTerminals []string
		wantMatch      bool
	}{
		{
			name:           "matches TERM_PROGRAM case insensitively",
			termProgram:    "ghostty",
			knownTerminals: []string{"ghostty"},
			wantMatch:      true,
		},
		{
			name:           "matches TERM",
			term:           "xterm-kitty",
			knownTerminals: []string{"xterm-kitty"},
			wantMatch:      true,
		},
		{
			name:           "matches with different case",
			termProgram:    "WEZTERM",
			knownTerminals: []string{"WezTerm"},
			wantMatch:      true,
		},
		{
			name:           "does not match an unrelated terminal",
			termProgram:    "Apple_Terminal",
			term:           "xterm-256color",
			knownTerminals: []string{"ghostty", "WezTerm"},
			wantMatch:      false,
		},
		{
			name:           "does not match on empty environment",
			knownTerminals: []string{"ghostty"},
			wantMatch:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesTerminal(tt.termProgram, tt.term, tt.knownTerminals); got != tt.wantMatch {
				t.Fatalf("matchesTerminal(%q, %q) = %v, want %v",
					tt.termProgram, tt.term, got, tt.wantMatch)
			}
		})
	}
}

// WezTerm accepts Kitty graphics transmissions but has no Unicode placeholder
// cell support, so it must not take the Kitty render path or the preview fills
// with U+10EEEE placeholder glyphs.
func TestIsKittyCapableExcludesTerminalsWithoutPlaceholders(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "WezTerm")
	if isKittyCapable() {
		t.Fatal("isKittyCapable() = true for WezTerm, want false: " +
			"WezTerm cannot render Kitty Unicode placeholder cells")
	}
}

func TestIsKittyCapableAllowsPlaceholderCapableTerminals(t *testing.T) {
	for _, terminal := range []string{"ghostty", "kitty", "xterm-kitty"} {
		t.Run(terminal, func(t *testing.T) {
			t.Setenv("TERM_PROGRAM", terminal)
			t.Setenv("TERM", "")
			if !isKittyCapable() {
				t.Fatalf("isKittyCapable() = false for %s, want true", terminal)
			}
		})
	}
}

func TestIsKittyCapableUnknownTerminal(t *testing.T) {
	t.Setenv("TERM_PROGRAM", "Apple_Terminal")
	t.Setenv("TERM", "xterm-256color")
	if isKittyCapable() {
		t.Fatal("isKittyCapable() = true for an unknown terminal, want false")
	}
}
