package common

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// ReadPassword reads a password from stdin without echoing
func ReadPassword() (string, error) {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		var password string
		_, err := fmt.Scanln(&password)
		return password, err
	}

	bytePassword, err := term.ReadPassword(fd)
	if err != nil {
		return "", err
	}

	return string(bytePassword), nil
}
