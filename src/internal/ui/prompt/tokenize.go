package prompt

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/yorukot/superfile/src/internal/common/utils"
)

// split into tokens
func tokenizePromptCommand(command string, cwdLocation string) ([]string, error) {
	command, err := resolveShellSubstitution(shellSubTimeout, command, cwdLocation)
	if err != nil {
		return nil, err
	}
	return strings.Fields(command), nil
}

// Replace ${} and $() with values
func resolveShellSubstitution(subCmdTimeout time.Duration, command string, cwdLocation string) (string, error) {
	resCommand := strings.Builder{}
	cmdRunes := []rune(command)
	i := 0
	for i < len(cmdRunes) {
		if i+1 < len(cmdRunes) && cmdRunes[i] == '$' {
			switch cmdRunes[i+1] {
			// ${ spotted
			case '{':
				// Look for Ending '}'
				end := findEndingBracket(cmdRunes, i+1, '{', '}')
				if end == -1 {
					return "", errors.New("unexpected error in tokenization")
				}
				if end == len(cmdRunes) {
					return "", curlyBracketMatchError()
				}

				envVarName := string(cmdRunes[i+2 : end])

				// We can add a layer of abstraction for better unit testing
				if value, ok := os.LookupEnv(envVarName); !ok {
					return "", envVarNotFoundError{varName: envVarName}
				} else {
					// Might Handle values being too big, or having multiple lines
					// But this is based on user input, so it is probably okay for now
					// Same comment for command substitution
					resCommand.WriteString(value)
				}

				i = end + 1
			case '(':
				// Look for ending ')'
				end := findEndingBracket(cmdRunes, i+1, '(', ')')
				if end == -1 {
					return "", errors.New("unexpected error in tokenization")
				}

				if end == len(cmdRunes) {
					return "", roundBracketMatchError()
				}

				subCmd := string(cmdRunes[i+2 : end])
				retCode, output, err := utils.ExecuteCommandInShell(subCmdTimeout, cwdLocation, subCmd)

				if retCode == -1 {
					return "", fmt.Errorf("could not execute shell substitution command : %s : %w", subCmd, err)
				} else {
					// We are allowing commands that exit with non zero status code
					// We still use its output
					if retCode != 0 {
						slog.Debug("substitution command exited with non zero status", "retCode", retCode,
							"command", subCmd)
					}
					resCommand.WriteString(output)
				}

				i = end + 1
			default:
				resCommand.WriteRune(cmdRunes[i])
				i++
			}
		} else {
			resCommand.WriteRune(cmdRunes[i])
			i++
		}
	}

	return resCommand.String(), nil
}

func findEndingBracket(r []rune, openIdx int, openParan rune, closeParan rune) int {
	if openIdx < 0 || openIdx >= len(r) || r[openIdx] != openParan {
		return -1
	}

	openCount := 1
	i := openIdx + 1
	for i < len(r) && openCount != 0 {
		if r[i] == openParan {
			openCount++
		} else if r[i] == closeParan {
			openCount--
		}
		if openCount != 0 {
			i++
		}
	}
	return i
}
