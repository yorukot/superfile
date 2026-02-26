package prompt

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/yorukot/superfile/src/pkg/utils"
)

// split into tokens
func tokenizePromptCommand(command string, cwdLocation string) ([]string, error) {
	command, err := resolveShellSubstitution(shellSubTimeout, command, cwdLocation)
	if err != nil {
		return nil, err
	}
	return tokenizeWithQuotes(command)
}

// Replace ${} and $() with values
func resolveShellSubstitution(subCmdTimeout time.Duration, command string, cwdLocation string) (string, error) {
	resCommand := strings.Builder{}
	cmdRunes := []rune(command)
	i := 0
	for i < len(cmdRunes) {
		if i+1 < len(cmdRunes) && cmdRunes[i] == '$' {
			if !isOpenBracket(cmdRunes[i+1]) {
				resCommand.WriteRune(cmdRunes[i])
				i++
				continue
			}
			openChar := cmdRunes[i+1]
			closeChar := getClosingBracket(openChar)
			end := findEndingBracket(cmdRunes, i+1, openChar, closeChar)
			if end == -1 {
				return "", errors.New("unexpected error in tokenization")
			}
			if end == len(cmdRunes) {
				return "", bracketMatchError{openChar: openChar, closeChar: closeChar}
			}

			err := updateResCommand(&resCommand, openChar, string(cmdRunes[i+2:end]),
				subCmdTimeout, cwdLocation)
			if err != nil {
				return "", err
			}
			i = end + 1
		} else {
			resCommand.WriteRune(cmdRunes[i])
			i++
		}
	}

	return resCommand.String(), nil
}

func updateResCommand(resCommand *strings.Builder, openChar rune, token string,
	subCmdTimeout time.Duration, cwdLocation string) error {
	switch openChar {
	case '{':
		value, ok := os.LookupEnv(token)
		if !ok {
			return envVarNotFoundError{varName: token}
		}
		// Might Handle values being too big, or having multiple lines
		// But this is based on user input, so it is probably okay for now
		// Same comment for command substitution
		resCommand.WriteString(value)
	case '(':
		retCode, output, err := utils.ExecuteCommandInShell(subCmdTimeout, cwdLocation, token)

		if retCode == -1 {
			return fmt.Errorf("could not execute shell substitution command : %s : %w", token, err)
		}
		// We are allowing commands that exit with non zero status code
		// We still use its output
		if retCode != 0 {
			slog.Debug("substitution command exited with non zero status", "retCode", retCode,
				"command", token)
		}
		resCommand.WriteString(output)
	default:
		return fmt.Errorf("unexpected openChar %v in tokenization", openChar)
	}
	return nil
}

func findEndingBracket(r []rune, openIdx int, openParan rune, closeParan rune) int {
	if openIdx < 0 || openIdx >= len(r) || r[openIdx] != openParan {
		return -1
	}

	openCount := 1
	i := openIdx + 1
	for i < len(r) && openCount != 0 {
		switch r[i] {
		case openParan:
			openCount++
		case closeParan:
			openCount--
		}
		if openCount != 0 {
			i++
		}
	}
	return i
}

func isOpenBracket(r rune) bool {
	switch r {
	case '(', '{':
		return true
	default:
		return false
	}
}

func getClosingBracket(r rune) rune {
	switch r {
	case '(':
		return ')'
	case '{':
		return '}'
	default:
		return ' '
	}
}

// splits command into tokens while respecting quotes and escapes
func tokenizeWithQuotes(command string) ([]string, error) {
	var (
		tokens    []string
		buffer    strings.Builder
		quoteOpen rune // 0:none, '\'' or '"'
		escaped   bool
	)

	// Initialize tokens as empty slice instead of nil
	tokens = []string{}

	// Helper function to flush the current buffer into tokens
	flush := func() {
		tokens = append(tokens, buffer.String())
		buffer.Reset()
	}

	for _, r := range command {
		switch {
		case escaped:
			// Only allow escaping of specific characters that have special meaning
			switch r {
			case '"', '\'', '\\', ' ':
				// These are valid escape sequences
				buffer.WriteRune(r)
			default:
				// Invalid escape sequence - treat backslash as literal
				buffer.WriteRune('\\')
				buffer.WriteRune(r)
			}
			escaped = false
		case r == '\\':
			escaped = true
		case quoteOpen == 0 && (r == '"' || r == '\''):
			quoteOpen = r
		case quoteOpen == r:
			// End of quoted section - always flush (even if empty)
			flush()
			quoteOpen = 0
		case unicode.IsSpace(r) && quoteOpen == 0:
			// Only flush if we have content
			if buffer.Len() > 0 {
				flush()
			}
		default:
			buffer.WriteRune(r)
		}
	}

	if escaped || quoteOpen != 0 {
		return nil, errors.New("unmatched quotes or escape characters in command")
	}

	// Flush any remaining content
	if buffer.Len() > 0 {
		flush()
	}

	return tokens, nil
}
