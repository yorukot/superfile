package filepreview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTmuxPassthrough(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single APC sequence is wrapped with doubled escapes",
			input:    "\x1b_Ga=d\x1b\\",
			expected: "\x1bPtmux;\x1b\x1b_Ga=d\x1b\x1b\\\x1b\\",
		},
		{
			name:  "each sequence is wrapped separately",
			input: "\x1b_Ga=d\x1b\\\x1b_Gm=1;AAAA\x1b\\",
			expected: "\x1bPtmux;\x1b\x1b_Ga=d\x1b\x1b\\\x1b\\" +
				"\x1bPtmux;\x1b\x1b_Gm=1;AAAA\x1b\x1b\\\x1b\\",
		},
		{
			name:     "unterminated trailing data is still wrapped",
			input:    "\x1b_Gm=0",
			expected: "\x1bPtmux;\x1b\x1b_Gm=0\x1b\\",
		},
		{
			name:     "empty input produces no wrapper",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tmuxPassthrough(tt.input))
		})
	}
}

func TestRawForTerminalOutsideTmux(t *testing.T) {
	t.Setenv("TMUX", "")
	raw := "\x1b_Ga=d\x1b\\"
	assert.Equal(t, raw, rawForTerminal(raw))
}

func TestRawForTerminalInsideTmux(t *testing.T) {
	t.Setenv("TMUX", "/tmp/tmux-1000/default,1,0")
	assert.Equal(t, tmuxPassthrough("\x1b_Ga=d\x1b\\"), rawForTerminal("\x1b_Ga=d\x1b\\"))
	assert.Empty(t, rawForTerminal(""))
}
