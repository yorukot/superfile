//go:build linux

package systemclipboard

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRunCopyReturnsWhileChildKeepsStderrOpen guards against inherited stderr
// blocking helper completion and verifies that failure diagnostics are retained.
func TestRunCopyReturnsWhileChildKeepsStderrOpen(t *testing.T) {
	for _, exitCode := range []string{"0", "7"} {
		t.Run("exit_"+exitCode, func(t *testing.T) {
			dir := t.TempDir()
			input := filepath.Join(dir, "input")
			release := filepath.Join(dir, "release")
			finished := filepath.Join(dir, "finished")
			t.Cleanup(func() {
				require.NoError(t, os.WriteFile(release, nil, 0o600))
				require.Eventually(t, func() bool {
					_, err := os.Stat(finished)
					return err == nil
				}, 3*time.Second, 10*time.Millisecond)
			})
			tool := linuxTool{
				name: "test-helper",
				copyArgs: func(string) []string {
					return []string{"sh", "-c", `
cat > "$1"
(while [ ! -f "$2" ]; do sleep 0.05; done; touch "$3") < /dev/null &
echo "helper diagnostic" >&2
exit "$4"
`, "test-helper", input, release, finished, exitCode}
				},
			}
			done := make(chan error, 1)
			go func() {
				done <- runCopy(tool, gnomeCopiedFilesMime, []byte("copy\nfile:///tmp/example"))
			}()

			var err error
			select {
			case err = <-done:
			case <-time.After(3 * time.Second):
				t.Error("runCopy waited for the background child to close stderr")
				require.NoError(t, os.WriteFile(release, nil, 0o600))
				err = <-done
			}
			if exitCode == "0" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, "test-helper copy failed: exit status 7: helper diagnostic")
			}
			payload, readErr := os.ReadFile(input)
			require.NoError(t, readErr)
			assert.Equal(t, "copy\nfile:///tmp/example", string(payload))
		})
	}
}

// TestPathURIRoundTrip checks that spaces, URI delimiters, and Unicode survive conversion.
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

// TestPathToURIEncoding checks percent-encoding without altering plain paths.
func TestPathToURIEncoding(t *testing.T) {
	assert.Equal(t, "file:///home/user/a%20b.txt", pathToURI("/home/user/a b.txt"))
	assert.Equal(t, "file:///plain/path", pathToURI("/plain/path"))
}

// TestBuildGnomeCopiedFiles checks copy and cut markers and escaped file URIs.
func TestBuildGnomeCopiedFiles(t *testing.T) {
	payload := buildGnomeCopiedFiles([]string{"/a/b", "/c d"}, true)
	assert.Equal(t, "cut\nfile:///a/b\nfile:///c%20d", payload)

	payload = buildGnomeCopiedFiles([]string{"/a"}, false)
	assert.Equal(t, "copy\nfile:///a", payload)
}

// TestParseGnomeCopiedFiles checks operation flags, trailing NULs, and rejection
// of empty payloads or invalid operation markers.
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

// TestParseURIList checks CRLF handling, comment removal, and bare-path compatibility.
func TestParseURIList(t *testing.T) {
	data := []byte("# comment\r\nfile:///a/b\r\nfile:///c%20d\r\n\r\n")
	assert.Equal(t, []string{"/a/b", "/c d"}, parseURIList(data))

	// Bare paths are tolerated.
	assert.Equal(t, []string{"/plain/path"}, parseURIList([]byte("/plain/path")))
}

// TestParseGnomeCopiedFilesSkipsNonFileURIs prevents remote URLs from becoming file paths.
func TestParseGnomeCopiedFilesSkipsNonFileURIs(t *testing.T) {
	paths, _, ok := parseGnomeCopiedFiles([]byte("copy\nhttp://example.com\nfile:///real"))
	require.True(t, ok)
	assert.Equal(t, []string{"/real"}, paths)
}
