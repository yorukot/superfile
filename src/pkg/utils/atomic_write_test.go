package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFileAtomicallyRemovesReplacementAfterPublishFailure(t *testing.T) {
	parentDir := t.TempDir()
	destinationDir := filepath.Join(parentDir, "config.toml")
	require.NoError(t, os.Mkdir(destinationDir, 0o700))

	err := writeFileAtomically(destinationDir, []byte("debug = true\n"), ConfigFilePerm)
	require.Error(t, err)

	info, statErr := os.Stat(destinationDir)
	require.NoError(t, statErr)
	assert.True(t, info.IsDir(), "failed publication must leave the destination unchanged")

	replacementFiles, globErr := filepath.Glob(filepath.Join(parentDir, ".config.toml.tmp-*"))
	require.NoError(t, globErr)
	assert.Empty(t, replacementFiles, "failed publication must remove its temporary replacement")
}
