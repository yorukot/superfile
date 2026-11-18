package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsASCIIDigit(t *testing.T) {
	for b := byte('0'); b <= '9'; b++ {
		assert.True(t, IsASCIIDigit(b), "expected %q to be a digit", b)
	}

	for _, b := range []byte{'a', 'Z', '/', ':', ' ', '-', '.'} {
		assert.False(t, IsASCIIDigit(b), "expected %q to not be a digit", b)
	}
}

func TestLeadingASCIIDigits(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty", input: "", expected: ""},
		{name: "digits only", input: "123", expected: "123"},
		{name: "digits with suffix", input: "3a", expected: "3"},
		{name: "digits with complex suffix", input: "42rc1", expected: "42"},
		{name: "no leading digits", input: "abc", expected: ""},
		{name: "starts with non-digit", input: ".1", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, LeadingASCIIDigits(tt.input))
		})
	}
}
