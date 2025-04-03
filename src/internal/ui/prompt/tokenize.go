package prompt

import (
	"fmt"
	"github.com/yorukot/superfile/src/internal/common/utils"
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
	// We want to test
	// Empty string
	// Strings without $ - Normal, trailing and leading whitespace, special char
	// Strings with $ but no ${} or $() - $HOME , _$$_xyz
	// Ill formatted substitutions - missing '$', missing ')' or '}' or '{'
	// Multiple correct ${} and $()
	// Empty ${} and $()
	// ${ $() } -> Should not work , $(echo ${} $(echo $())) -> Should work
	// cd $(echo $(echo hi))
	// env var - not found
	// substitution command times out
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
					return "", fmt.Errorf("could not find matching for '}' for '${'")
				}

				envVarName := string(cmdRunes[i+2 : end])

				// Todo : add a layer of abstraction for unit testing
				if value, ok := os.LookupEnv(envVarName); !ok {
					return "", fmt.Errorf("env %s not found", envVarName)
				} else {
					// Todo : Handle value being too big ? or having newlines ?
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
					return "", fmt.Errorf("could not find matching for ')' for '$('")
				}

				subCmd := string(cmdRunes[i+2 : end])
				retCode, output, err := utils.ExecuteCommandInShell(subCmdTimeout, cwdLocation, subCmd)

				if retCode == -1 {
					return "", fmt.Errorf("could not execute shell substitution command : %s : %w", subCmd, err)
				} else {
					// Todo : Handle value being too big ? or having newlines ?
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

func findEndingParenthesis(r []rune, openIdx int, open rune, close rune) int {
	if openIdx < 0 || openIdx >= len(r) || r[openIdx] != open {
		return -1
	}

	openCount := 1
	i := openIdx
	for i < len(r) && openCount != 0 {
		if r[i] == open {
			openCount++
		} else if r[i] == close {
			openCount--
		}
		if openCount != 0 {
			i++
		}
	}
	return i
}
