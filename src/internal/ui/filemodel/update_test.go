package filemodel

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTruncateRemotePreviewPreservesUTF8(t *testing.T) {
	data := bytes.Repeat([]byte{'a'}, remotePreviewLimit-1)
	data = append(data, []byte("€")...)
	require.Greater(t, len(data), remotePreviewLimit)

	truncated := truncateRemotePreviewData(data)
	assert.True(t, utf8.Valid(truncated))
	assert.Contains(t, string(truncated), "[preview truncated at 1 MiB]")
}
