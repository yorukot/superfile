package common

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
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
		{TruncateText, "TruncateText", "Hello world", 4, "...", "H..."},
		{TruncateText, "TruncateText", "Hello world", 6, "...", "Hel..."},
		{TruncateText, "TruncateText", "Hello", 100, "...", "Hello"},
		{TruncateTextBeginning, "TruncateTextBeginning", "Hello world", 4, "...", "...d"},
		{TruncateTextBeginning, "TruncateTextBeginning", "Hello world", 6, "...", "...rld"},
		{TruncateTextBeginning, "TruncateTextBeginning", "Hello", 100, "...", "Hello"},
		{TruncateMiddleText, "TruncateMiddleText", "Hello world", 5, "...", "H...d"},
		{TruncateMiddleText, "TruncateMiddleText", "Hello world", 7, "...", "He...ld"},
		{TruncateMiddleText, "TruncateMiddleText", "Hello", 100, "...", "Hello"},
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
			result := FileNameWithoutExtension(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestHelpHotkeyString(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			name:     "Single key",
			input:    []string{"a"},
			expected: "a",
		},
		{
			name:     "Multiple keys",
			input:    []string{"a", "b", "c"},
			expected: "a | b | c",
		},
		{
			name:     "Empty key",
			input:    []string{"a", "", "b"},
			expected: "a | b",
		},
		{
			name:     "Trailing empty",
			input:    []string{"a", ""},
			expected: "a",
		},
		{
			name:     "Trailing empty with multiple keys",
			input:    []string{"a", "b", ""},
			expected: "a | b",
		},
		{
			name:     "Space key",
			input:    []string{" "},
			expected: "space",
		},

		// Starting with an empty key ("", "a") is not allowed by the file parser,
		// so a test is not needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetHelpMenuHotkeyString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsBufferPrintable(t *testing.T) {
	var inputs = []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"hello", true},
		{"abcdABCD0123~!@#$%^&*()_+-={}|:\"<>?,./;'[]", true},
		{"Horizontal Tab and NewLine\t\t\n\n", true},
		{"\xa0(NBSP)", true},
		{"\x0b(Vertical Tab)", true},
		{"\x0d(CR)", true},
		{"ASCII control characters : \x00(NULL)", false},
		{"\x05(ENQ)", false},
		{"\x0f(SI)", false},
		{"\x1b(ESC)", false},
		{"\x7f(DEL)", false},
	}
	for _, tt := range inputs {
		t.Run(fmt.Sprintf("Testing if buffer %q is printable", tt.input), func(t *testing.T) {
			result := IsBufferPrintable([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsExtensionExtractable(t *testing.T) {
	inputs := []struct {
		ext      string
		expected bool
	}{
		{".zip", true},
		{".rar", true},
		{".7z", true},
		{".tar.gz", true},
		{".tar.bz2", true},
		{".exe", false},
		{".txt", false},
		{".tar", true},
		{"", false},    // Empty string case
		{".ZIP", true}, // Case sensitivity check
		{".Zip", true}, // Case sensitivity check
		{".bz", true},
		{".gz", true},
		{".iso", true},
	}

	for _, tt := range inputs {
		t.Run(tt.ext, func(t *testing.T) {
			result := IsExtensionExtractable(tt.ext)
			if result != tt.expected {
				t.Errorf("IsExensionExtractable (%q) = %v; want %v", tt.ext, result, tt.expected)
			}
		})
	}
}

func TestMakePrintable(t *testing.T) {
	var inputs = []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "hello"},
		{"abcdABCD0123~!@#$%^&*()_+-={}|:\"<>?,./;'[]", "abcdABCD0123~!@#$%^&*()_+-={}|:\"<>?,./;'[]"},
		{"Horizontal Tab and NewLine\t\t\n\n", "Horizontal Tab and NewLine      \n\n"},
		// Tab expansion tests - tabs should expand to reach next tab stop (TabWidth=4)
		{"a\tb", "a   b"},        // Position 1: need 3 spaces to reach position 4
		{"ab\tc", "ab  c"},       // Position 2: need 2 spaces to reach position 4
		{"abc\td", "abc d"},      // Position 3: need 1 space to reach position 4
		{"abcd\te", "abcd    e"}, // Position 4: need 4 spaces to reach position 8
		{"a\tb\tc", "a   b   c"}, // Position 1: 3 spaces, then position 4: 4 spaces to reach 8
		{"(NBSP)\u00a0\u00a0\u00a0\u00a0;", "(NBSP)\u00a0\u00a0\u00a0\u00a0;"},
		{"\x0b(Vertical Tab)", "(Vertical Tab)"},
		{"\x0d(CR)", "(CR)"},
		{"ASCII control characters : \x00(NULL)", "ASCII control characters : (NULL)"},
		{"\x05(ENQ)", "(ENQ)"},
		{"\x0f(SI)", "(SI)"},
		{"\x1b(ESC)", "\x1b(ESC)"},
		{"\x7f(DEL)", "(DEL)"},
		{"\x7f(DEL)", "(DEL)"},
		{"Valid unicodes like nerdfont \uf410 \U000f0868", "Valid unicodes like nerdfont \uf410 \U000f0868"},
		{"Invalid Unicodes\ufffd", "Invalid Unicodes"},
		{"Invalid Unicodes\xa0", "Invalid Unicodes"},
		{"Ascii color sequence\x1b[38;2;230;219;116;48;2;39;40;34m\ue68f \x1b[0m",
			"Ascii color sequence\x1b[38;2;230;219;116;48;2;39;40;34m\ue68f \x1b[0m"},
		{"Unicodes spaces\u202f\u205f\u2029", "Unicodes spaces   "},
		{"IDEOGRAPHIC SPACE\u3000", "IDEOGRAPHIC SPACE "},
	}
	for _, tt := range inputs {
		t.Run(fmt.Sprintf("Make %q printable", tt.input), func(t *testing.T) {
			result := MakePrintable(tt.input)
			if result != tt.expected {
				t.Errorf("Expected '%v', got '%v' (input : '%v')", tt.expected, result, tt.input)
			}
		})
	}
	t.Run("ESC is skipped", func(t *testing.T) {
		assert.Equal(t, "(ESC)", MakePrintableWithEscCheck("\x1b(ESC)", false))
	})
	t.Run("ESC is not skipped", func(t *testing.T) {
		assert.Equal(t, "\x1b(ESC)", MakePrintableWithEscCheck("\x1b(ESC)", true))
	})
}

func TestFormatSizeInternal(t *testing.T) {
	t.Run("max int size", func(t *testing.T) {
		actual := formatSizeInternal(math.MaxInt64, KilobyteSize, unitsDec())
		assert.Equal(t, "9.22 EB", actual)
	})
	t.Run("zero size", func(t *testing.T) {
		actual := formatSizeInternal(0, KilobyteSize, unitsDec())
		assert.Equal(t, "0 B", actual)
	})
	t.Run("100 bytes size", func(t *testing.T) {
		actual := formatSizeInternal(100, KilobyteSize, unitsDec())
		assert.Equal(t, "100 B", actual)
	})
	t.Run("1005 bytes size", func(t *testing.T) {
		actual := formatSizeInternal(1005, KilobyteSize, unitsDec())
		assert.Equal(t, "1.00 kB", actual)
	})
	t.Run("1005 bytes size kibi", func(t *testing.T) {
		actual := formatSizeInternal(1005, KibibyteSize, unitsBin())
		assert.Equal(t, "1005 B", actual)
	})
	t.Run("1025 bytes size kibi", func(t *testing.T) {
		actual := formatSizeInternal(1025, KibibyteSize, unitsBin())
		assert.Equal(t, "1.00 KiB", actual)
	})
	t.Run("1035 bytes size kibi", func(t *testing.T) {
		actual := formatSizeInternal(1035, KibibyteSize, unitsBin())
		assert.Equal(t, "1.01 KiB", actual)
	})
}
