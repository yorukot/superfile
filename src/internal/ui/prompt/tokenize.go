package prompt

import (
	"fmt"
	"github.com/yorukot/superfile/src/internal/common/utils"
	"log/slog"
	"os"
	"strings"
	"time"
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
			// ${ spotted
			if cmdRunes[i+1] == '{' {
				// Look for Ending '}'
				end := findEndingParenthesis(cmdRunes, i+1, '{', '}')
				if end == -1 {
					return "", fmt.Errorf("unexpected error in tokenization")
				}
				if end == len(cmdRunes) {
					return "", curlyBracketParMatchError()
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

			} else if cmdRunes[i+1] == '(' {
				// Look for ending ')'
				end := findEndingParenthesis(cmdRunes, i+1, '(', ')')
				if end == -1 {
					return "", fmt.Errorf("unexpected error in tokenization")
				}

				if end == len(cmdRunes) {
					return "", bracketParMatchError()
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
			} else {
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

func findEndingParenthesis(r []rune, openIdx int, openParan rune, closeParan rune) int {
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
