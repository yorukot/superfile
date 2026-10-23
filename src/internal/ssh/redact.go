package ssh

import (
	"slices"
	"strings"
)

const redacted = "[REDACTED]"

const (
	privateKeyBeginPrefix = "-----BEGIN "
	pemMarkerSuffix       = "-----"
	baseRedactionSecrets  = 2
)

type RedactionSecrets struct {
	Password           string
	IdentityPassphrase string
	AdditionalSecrets  []string
}

func RedactString(value string, secrets RedactionSecrets) string {
	redactedValue := redactPrivateKeyBlocks(value)
	for _, secret := range redactionValues(secrets) {
		redactedValue = strings.ReplaceAll(redactedValue, secret, redacted)
	}
	return redactedValue
}

func RedactError(err error, secrets RedactionSecrets) error {
	if err == nil {
		return nil
	}
	return redactedError{err: err, message: RedactString(err.Error(), secrets)}
}

func redactionValues(secrets RedactionSecrets) []string {
	values := make([]string, 0, len(secrets.AdditionalSecrets)+baseRedactionSecrets)
	seen := make(map[string]struct{}, cap(values))
	for _, secret := range append(
		[]string{secrets.Password, secrets.IdentityPassphrase},
		secrets.AdditionalSecrets...,
	) {
		if secret == "" {
			continue
		}
		if _, ok := seen[secret]; ok {
			continue
		}
		seen[secret] = struct{}{}
		values = append(values, secret)
	}
	slices.SortFunc(values, func(a string, b string) int {
		if len(a) != len(b) {
			return len(b) - len(a)
		}
		return strings.Compare(a, b)
	})
	return values
}

func redactPrivateKeyBlocks(value string) string {
	var result strings.Builder
	remaining := value
	for {
		begin := strings.Index(remaining, privateKeyBeginPrefix)
		if begin < 0 {
			result.WriteString(remaining)
			return result.String()
		}

		labelStart := begin + len(privateKeyBeginPrefix)
		headerRemainder := remaining[labelStart:]
		labelEndOffset := strings.Index(headerRemainder, pemMarkerSuffix)
		nextBeginOffset := strings.Index(headerRemainder, privateKeyBeginPrefix)
		if nextBeginOffset >= 0 && (labelEndOffset < 0 || nextBeginOffset <= labelEndOffset) {
			nextBegin := labelStart + nextBeginOffset
			result.WriteString(remaining[:nextBegin])
			remaining = remaining[nextBegin:]
			continue
		}
		lineEndOffset := strings.IndexAny(headerRemainder, "\r\n")
		if labelEndOffset < 0 || (lineEndOffset >= 0 && labelEndOffset > lineEndOffset) {
			result.WriteString(remaining[:labelStart])
			remaining = headerRemainder
			continue
		}
		labelEnd := labelStart + labelEndOffset
		beginMarkerEnd := labelEnd + len(pemMarkerSuffix)
		label := remaining[labelStart:labelEnd]
		if !strings.Contains(label, "PRIVATE KEY") {
			result.WriteString(remaining[:beginMarkerEnd])
			remaining = remaining[beginMarkerEnd:]
			continue
		}

		result.WriteString(remaining[:begin])
		result.WriteString(redacted)
		endMarker := "-----END " + label + pemMarkerSuffix
		endOffset := strings.Index(remaining[beginMarkerEnd:], endMarker)
		if endOffset < 0 {
			return result.String()
		}
		remaining = remaining[beginMarkerEnd+endOffset+len(endMarker):]
	}
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
