package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTTYForRedirectedStdout(t *testing.T) {
	openFile := func(t *testing.T, name string) *os.File {
		t.Helper()
		f, err := os.Create(filepath.Join(t.TempDir(), name))
		require.NoError(t, err)
		t.Cleanup(func() { f.Close() })
		return f
	}

	t.Run("stdout is a terminal", func(t *testing.T) {
		called := false
		tty := ttyForRedirectedStdout(true, func() (*os.File, *os.File, error) {
			called = true
			return nil, nil, nil
		})
		assert.Nil(t, tty)
		assert.False(t, called, "the terminal should not be opened")
	})

	t.Run("stdout is redirected", func(t *testing.T) {
		f := openFile(t, "tty")
		tty := ttyForRedirectedStdout(false, func() (*os.File, *os.File, error) {
			return f, f, nil
		})
		assert.Same(t, f, tty)
		_, err := f.Stat()
		assert.NoError(t, err, "the returned terminal should stay open")
	})

	t.Run("separate input handle is closed", func(t *testing.T) {
		in := openFile(t, "in")
		out := openFile(t, "out")
		tty := ttyForRedirectedStdout(false, func() (*os.File, *os.File, error) {
			return in, out, nil
		})
		assert.Same(t, out, tty)
		_, err := in.Stat()
		require.ErrorIs(t, err, os.ErrClosed)
		_, err = out.Stat()
		assert.NoError(t, err)
	})

	t.Run("terminal cannot be opened", func(t *testing.T) {
		tty := ttyForRedirectedStdout(false, func() (*os.File, *os.File, error) {
			return nil, nil, errors.New("no tty")
		})
		assert.Nil(t, tty)
	})
}
