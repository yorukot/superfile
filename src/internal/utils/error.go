package utils

import (
	"errors"
)

type TomlLoadError struct {
	userMessage   string
	wrappedError  error
	isFatal       bool
	missingFields bool
}

func (t *TomlLoadError) Error() string {
	res := t.userMessage
	if t.wrappedError != nil {
		res += " : " + t.wrappedError.Error()
	}
	return res
}

func (t *TomlLoadError) IsFatal() bool {
	return t.isFatal
}

func (t *TomlLoadError) MissingFields() bool {
	return t.missingFields
}

func (t *TomlLoadError) Unwrap() error {
	return t.wrappedError
}

func (t *TomlLoadError) UpdateMessageAndError(msg string, err error) {
	t.userMessage = msg
	t.wrappedError = err
}

// Include another msg. For now we dont need to have this as wrapped error.
func (t *TomlLoadError) AddMessageAndError(msg string, err error) {
	t.userMessage += " " + msg
	t.wrappedError = errors.Join(t.wrappedError, err)
}
