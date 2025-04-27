package rendering

import (
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/stretchr/testify/assert"
)

func TestTruncate(t *testing.T) {
	testStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
	testdata := []struct {
		name     string
		line     string
		maxWidth int
		style    TruncateStyle
		expected string
	}{
		{
			name:     "No truncate",
			line:     "abc",
			maxWidth: 10,
			style:    PlainTruncateRight,
			expected: "abc",
		},
		{
			name:     "Plain truncate",
			line:     "abcdefgh",
			maxWidth: 5,
			style:    PlainTruncateRight,
			expected: "abcde",
		},
		{
			name:     "Tails truncate",
			line:     "abcdefgh",
			maxWidth: 5,
			style:    TailsTruncateRight,
			expected: "ab...",
		},
		{
			name:     "Tails truncate with too less width",
			line:     "abcdefgh",
			maxWidth: 2,
			style:    TailsTruncateRight,
			expected: "",
		},
		{
			name:     "Wide characters",
			line:     "✅1✅2✅3",
			maxWidth: 3,
			style:    PlainTruncateRight,
			expected: "✅1",
		},
		{
			name:     "Wide characters 2",
			line:     "✅1✅2✅3",
			maxWidth: 4,
			style:    PlainTruncateRight,
			expected: "✅1",
		},
		{
			name:     "Wide characters 3",
			line:     "✅1✅2✅3",
			maxWidth: 4,
			style:    TailsTruncateRight,
			expected: "...",
		},
		{
			name:     "Ansi color sequence",
			line:     testStyle.Render("1234"),
			maxWidth: 4,
			style:    TailsTruncateRight,
			expected: testStyle.Render("1..."),
		},
	}
	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TruncateBasedOnStyle(tt.line, tt.maxWidth, tt.style))
		})
	}
}
