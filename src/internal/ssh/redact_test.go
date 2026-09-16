package ssh

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedactStringRemovesRuntimeSecrets(t *testing.T) {
	secrets := RedactionSecrets{
		Password:           "arbitrary-runtime-password",
		IdentityPassphrase: "arbitrary-runtime-passphrase",
		AdditionalSecrets: []string{
			"arbitrary-keyboard-answer-one",
			"arbitrary-keyboard-answer-two",
		},
	}
	message := "visible arbitrary-runtime-password arbitrary-runtime-passphrase " +
		"arbitrary-keyboard-answer-one arbitrary-keyboard-answer-two"

	redactedMessage := RedactString(message, secrets)

	assert.Equal(t, "visible [REDACTED] [REDACTED] [REDACTED] [REDACTED]", redactedMessage)
}

func TestRedactStringSkipsEmptyAndRedactsOverlappingValuesLongestFirst(t *testing.T) {
	secrets := RedactionSecrets{
		Password:           "token",
		IdentityPassphrase: "",
		AdditionalSecrets:  []string{"token-long", "", "token", "token-long"},
	}

	redactedMessage := RedactString("token-long token token-long", secrets)

	assert.Equal(t, "[REDACTED] [REDACTED] [REDACTED]", redactedMessage)
}

func TestRedactStringDoesNotUseFixtureOnlySecrets(t *testing.T) {
	message := "secret-password secret-passphrase ordinary text"

	assert.Equal(t, message, RedactString(message, RedactionSecrets{}))
}

func TestRedactStringRemovesPrivateKeyBlocks(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name: "complete block preserves surrounding LF text",
			value: "before\n-----BEGIN OPENSSH PRIVATE KEY-----\npayload\n" +
				"-----END OPENSSH PRIVATE KEY-----\nafter",
			expected: "before\n[REDACTED]\nafter",
		},
		{
			name: "complete block preserves surrounding CRLF text",
			value: "before\r\n-----BEGIN ENCRYPTED PRIVATE KEY-----\r\npayload\r\n" +
				"-----END ENCRYPTED PRIVATE KEY-----\r\nafter",
			expected: "before\r\n[REDACTED]\r\nafter",
		},
		{
			name: "multiple blocks",
			value: "prefix -----BEGIN RSA PRIVATE KEY-----\none\n-----END RSA PRIVATE KEY-----" +
				" middle -----BEGIN PRIVATE KEY-----\ntwo\n-----END PRIVATE KEY----- suffix",
			expected: "prefix [REDACTED] middle [REDACTED] suffix",
		},
		{
			name:     "unterminated block redacts through end of input",
			value:    "before -----BEGIN EC PRIVATE KEY-----\npayload\nstill-sensitive",
			expected: "before [REDACTED]",
		},
		{
			name: "matching end label controls block boundary",
			value: "before -----BEGIN DSA PRIVATE KEY-----\npayload\n-----END EC PRIVATE KEY-----\n" +
				"more\n-----END DSA PRIVATE KEY----- after",
			expected: "before [REDACTED] after",
		},
		{
			name:     "non-private PEM remains unchanged",
			value:    "before -----BEGIN PUBLIC KEY-----\npayload\n-----END PUBLIC KEY----- after",
			expected: "before -----BEGIN PUBLIC KEY-----\npayload\n-----END PUBLIC KEY----- after",
		},
		{
			name: "malformed header before private key does not bypass redaction",
			value: "before -----BEGIN MALFORMED\n-----BEGIN OPENSSH PRIVATE KEY-----\npayload\n" +
				"-----END OPENSSH PRIVATE KEY----- after",
			expected: "before -----BEGIN MALFORMED\n[REDACTED] after",
		},
		{
			name: "same-line malformed header before private key does not bypass redaction",
			value: "before -----BEGIN MALFORMED -----BEGIN OPENSSH PRIVATE KEY-----\npayload\n" +
				"-----END OPENSSH PRIVATE KEY----- after",
			expected: "before -----BEGIN MALFORMED [REDACTED] after",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, RedactString(tt.value, RedactionSecrets{}))
		})
	}
}

type typedRedactionTestError struct {
	message string
}

func (e *typedRedactionTestError) Error() string {
	return e.message
}

func TestRedactErrorPreservesErrorChain(t *testing.T) {
	target := errors.New("arbitrary-target-secret")
	typed := &typedRedactionTestError{message: "arbitrary-typed-secret"}
	source := fmt.Errorf("outer failure: %w", errors.Join(target, typed))

	redactedErr := RedactError(source, RedactionSecrets{AdditionalSecrets: []string{
		"arbitrary-target-secret",
		"arbitrary-typed-secret",
	}})

	require.Error(t, redactedErr)
	assert.NotContains(t, redactedErr.Error(), "arbitrary-target-secret")
	assert.NotContains(t, redactedErr.Error(), "arbitrary-typed-secret")
	require.ErrorIs(t, redactedErr, target)
	var foundTyped *typedRedactionTestError
	require.ErrorAs(t, redactedErr, &foundTyped)
	assert.Same(t, typed, foundTyped)
}

func TestRedactErrorAcceptsNil(t *testing.T) {
	assert.NoError(t, RedactError(nil, RedactionSecrets{}))
}
