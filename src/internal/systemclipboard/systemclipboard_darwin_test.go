//go:build darwin && cgo

package systemclipboard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSplitNULTerminatedIsLossless checks that names containing newlines or
// surrounding whitespace are not split or trimmed.
func TestSplitNULTerminatedIsLossless(t *testing.T) {
	raw := []byte("/tmp/a\x00/tmp/line\nbreak\x00/tmp/trailing \x00/tmp/ lead\x00")
	assert.Equal(t,
		[]string{"/tmp/a", "/tmp/line\nbreak", "/tmp/trailing ", "/tmp/ lead"},
		splitNULTerminated(raw))
}

// TestSplitNULTerminatedSkipsEmptyRecords checks empty input and stray separators.
func TestSplitNULTerminatedSkipsEmptyRecords(t *testing.T) {
	assert.Empty(t, splitNULTerminated(nil))
	assert.Equal(t, []string{"/x"}, splitNULTerminated([]byte("\x00/x\x00\x00")))
}
