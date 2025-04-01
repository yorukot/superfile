package prompt

import (
	"fmt"
	"github.com/yorukot/superfile/src/internal/common/utils"
	"os"
	"strings"
)

// split into tokens
func tokenizePromptCommand(command string, cwdLocation string) ([]string, error) {
	command, err := resolveShellSubstitution(command, cwdLocation)
	if err != nil {
		return nil, err
	}
	return strings.Fields(command), nil
}

// Replace ${} and $() with values
func resolveShellSubstitution(command string, cwdLocation string) (string, error) {
	resCommand := strings.Builder{}
	cmdRunes := []rune(command)
	i := 0
	for i < len(cmdRunes) {

		if i+1 < len(cmdRunes) && cmdRunes[i] == '$' {
			// ${ spotted
			if cmdRunes[i+1] == '{' {
				// Look for Ending '}'
				st := i + 2
				for st < len(cmdRunes) && cmdRunes[st] != '}' {
					st++
				}

				if st == len(cmdRunes) {
					return "", fmt.Errorf("could not find matching for '}' for '${'")
				}

				envVarName := string(cmdRunes[i+2 : st])

				// Todo : add a layer of abstraction for unit testing
				if value, ok := os.LookupEnv(envVarName); !ok {
					return "", fmt.Errorf("env %s not found", envVarName)
				} else {
					// Todo : Handle value being too big ? or having newlines ?
					resCommand.WriteString(value)
				}

				i = st + 1

			} else if cmdRunes[i+1] == '(' {
				// Look for ending ')'
				st := i + 2
				for st < len(cmdRunes) && cmdRunes[st] != ')' {
					st++
				}

				if st == len(cmdRunes) {
					return "", fmt.Errorf("could not find matching for ')' for '$('")
				}

				subCmd := string(cmdRunes[i+2 : st])
				retCode, output, err := utils.ExecuteCommandInShell(shellSubTimeoutMsec, cwdLocation, subCmd)

				if retCode == -1 {
					return "", fmt.Errorf("could not execute shell substitution command : %s : %w", subCmd, err)
				} else {
					// Todo : Handle value being too big ? or having newlines ?
					resCommand.WriteString(output)
				}

				i = st + 1
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
