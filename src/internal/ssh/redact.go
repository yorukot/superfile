package ssh

import (
	"strings"
)

const redacted = "[REDACTED]"

var sensitiveMarkers = []string{ //nolint:gochecknoglobals // Immutable redaction marker catalog.
	"secret-password",
	"secret-passphrase",
	"-----BEGIN OPENSSH PRIVATE KEY-----",
	"-----BEGIN RSA PRIVATE KEY-----",
	"-----BEGIN EC PRIVATE KEY-----",
	"-----BEGIN DSA PRIVATE KEY-----",
}

func RedactString(value string) string {
	redactedValue := value
	for _, marker := range sensitiveMarkers {
		redactedValue = strings.ReplaceAll(redactedValue, marker, redacted)
	}
	return redactedValue
}

func RedactError(err error) error {
	if err == nil {
		return nil
	}
	return redactedError{err: err, message: RedactString(err.Error())}
}

type redactedError struct {
	err     error
	message string
}

func (e redactedError) Error() string {
	return e.message
}

func (e redactedError) Unwrap() error {
	return e.err
}

var _ error = redactedError{}
