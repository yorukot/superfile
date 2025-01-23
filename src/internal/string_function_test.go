package internal

import (
	"fmt"
	"testing"
)

func TestStringTruncate(t *testing.T) {
	var inputs = []struct {
		function func(string, int, string) string
		funcName string
		input    string
		maxSize  int
		talis    string
		expected string
	}{
		{truncateText, "truncateText", "Hello world", 4, "...", "H..."},
		{truncateText, "truncateText", "Hello world", 6, "...", "Hel..."},
		{truncateText, "truncateText", "Hello", 100, "...", "Hello"},
		{truncateTextBeginning, "truncateTextBeginning", "Hello world", 4, "...", "...d"},
		{truncateTextBeginning, "truncateTextBeginning", "Hello world", 6, "...", "...rld"},
		{truncateTextBeginning, "truncateTextBeginning", "Hello", 100, "...", "Hello"},
		{truncateMiddleText, "truncateMiddleText", "Hello world", 5, "...", "H...d"},
		{truncateMiddleText, "truncateMiddleText", "Hello world", 7, "...", "He...ld"},
		{truncateMiddleText, "truncateMiddleText", "Hello", 100, "...", "Hello"},
	}

	for _, tt := range inputs {
		t.Run(fmt.Sprintf("Run %s on string %s to %d chars", tt.funcName, tt.input, tt.maxSize), func(t *testing.T) {
			result := tt.function(tt.input, tt.maxSize, tt.talis)
			expected := tt.expected
			if result != expected {
				t.Errorf("got \"%s\", expected \"%s\"", result, expected)
			}
		})
	}
}

func TestFilenameWithouText(t *testing.T) {
	var inputs = []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello.zip", "hello"},
		{"hello.tar.gz", "hello"},
		{".gitignore", ".gitignore"},
		{"", ""},
	}

	for _, tt := range inputs {
		t.Run(fmt.Sprintf("Remove extension from %s", tt.input), func(t *testing.T) {
			result := fileNameWithoutExtension(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}


func TestIsBufferPrintable(t *testing.T) {
	var inputs = []struct {
		input    string
		expected bool
	} {
		{"hello", true},
		{"abcdABCD0123~!@#$%^&*()_+-={}|:\"<>?,./;'[]", true},
		{"Horizontal Tab and NewLine\t\t\n\n", true},
		{"ASCII control characters : \x00(NULL)", false},
		{"\x05(ENQ)", false},
		{"\x0f(SI)", false},
		{"\x1b(ESC)", false},
		{"\x7f(DEL)", false},
	}
	for _, tt := range inputs {

		t.Run(fmt.Sprintf("Testing if buffer %q is printable", tt.input), func(t* testing.T){
			result := isBufferPrintable([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}

}