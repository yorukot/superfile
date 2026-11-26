package filepreview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Only the non-tmux paths are asserted here, since the tmux paths depend on a
// running tmux server being reachable from the test environment.
func TestDetectKittyCapableOutsideTmux(t *testing.T) {
	testcases := []struct {
		name        string
		termProgram string
		term        string
		expected    bool
	}{
		{
			name:        "known terminal via TERM_PROGRAM",
			termProgram: "ghostty",
			term:        "xterm-256color",
			expected:    true,
		},
		{
			name:     "known terminal via TERM",
			term:     "xterm-kitty",
			expected: true,
		},
		{
			name:     "terminfo name containing a known terminal",
			term:     "xterm-ghostty",
			expected: true,
		},
		{
			name:     "unknown terminal",
			term:     "xterm-256color",
			expected: false,
		},
		{
			name:     "empty environment",
			expected: false,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TMUX", "")
			t.Setenv("TERM_PROGRAM", tt.termProgram)
			t.Setenv("TERM", tt.term)
			assert.Equal(t, tt.expected, detectKittyCapable())
		})
	}
}
