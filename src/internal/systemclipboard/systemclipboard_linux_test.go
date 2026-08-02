//go:build linux

package systemclipboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathURIRoundTrip(t *testing.T) {
	cases := []string{
		"/home/user/file.txt",
		"/home/user/with space/a b.txt",
		"/tmp/weird#name?.txt",
		"/tmp/ünïcode/файл",
	}
	for _, p := range cases {
		got := uriToPath(pathToURI(p))
		assert.Equal(t, p, got, "round trip for %q", p)
	}
}

func TestPathToURIEncoding(t *testing.T) {
	assert.Equal(t, "file:///home/user/a%20b.txt", pathToURI("/home/user/a b.txt"))
	assert.Equal(t, "file:///plain/path", pathToURI("/plain/path"))
}

func TestBuildGnomeCopiedFiles(t *testing.T) {
	payload := buildGnomeCopiedFiles([]string{"/a/b", "/c d"}, true)
	assert.Equal(t, "cut\nfile:///a/b\nfile:///c%20d", payload)

	payload = buildGnomeCopiedFiles([]string{"/a"}, false)
	assert.Equal(t, "copy\nfile:///a", payload)
}

func TestParseGnomeCopiedFiles(t *testing.T) {
	paths, cut, ok := parseGnomeCopiedFiles([]byte("cut\nfile:///a/b\nfile:///c%20d"))
	require.True(t, ok)
	assert.True(t, cut)
	assert.Equal(t, []string{"/a/b", "/c d"}, paths)

	paths, cut, ok = parseGnomeCopiedFiles([]byte("copy\nfile:///x\x00"))
	require.True(t, ok)
	assert.False(t, cut)
	assert.Equal(t, []string{"/x"}, paths)

	// Missing/invalid operation header.
	_, _, ok = parseGnomeCopiedFiles([]byte("file:///a\nfile:///b"))
	assert.False(t, ok)

	// Empty payload.
	_, _, ok = parseGnomeCopiedFiles([]byte(""))
	assert.False(t, ok)
}

func TestParseURIList(t *testing.T) {
	data := []byte("# comment\r\nfile:///a/b\r\nfile:///c%20d\r\n\r\n")
	assert.Equal(t, []string{"/a/b", "/c d"}, parseURIList(data))

	// Bare paths are tolerated.
	assert.Equal(t, []string{"/plain/path"}, parseURIList([]byte("/plain/path")))
}

func TestParseGnomeCopiedFilesSkipsNonFileURIs(t *testing.T) {
	paths, _, ok := parseGnomeCopiedFiles([]byte("copy\nhttp://example.com\nfile:///real"))
	require.True(t, ok)
	assert.Equal(t, []string{"/real"}, paths)
}
