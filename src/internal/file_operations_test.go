package internal

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoveElementDoesNotReplaceConcurrentDestination(t *testing.T) {
	tempDir := t.TempDir()
	destination := filepath.Join(tempDir, "destination.txt")
	sources := []string{
		filepath.Join(tempDir, "first.txt"),
		filepath.Join(tempDir, "second.txt"),
	}
	contents := []string{"first", "second"}
	contentBySource := make(map[string]string, len(sources))
	for index, source := range sources {
		contentBySource[source] = contents[index]
		require.NoError(t, os.WriteFile(source, []byte(contents[index]), 0o644))
	}

	type result struct {
		source string
		err    error
	}
	start := make(chan struct{})
	results := make(chan result, len(sources))
	for _, source := range sources {
		go func() {
			<-start
			results <- result{source: source, err: moveElement(source, destination)}
		}()
	}
	close(start)

	var succeeded, rejected result
	for range sources {
		moveResult := <-results
		if moveResult.err == nil {
			succeeded = moveResult
		} else {
			rejected = moveResult
		}
	}

	require.NotEmpty(t, succeeded.source)
	require.NotEmpty(t, rejected.source)
	assert.True(t, errors.Is(rejected.err, os.ErrExist), rejected.err)
	_, err := os.Stat(succeeded.source)
	assert.True(t, os.IsNotExist(err))
	assert.FileExists(t, rejected.source)
	destinationContents, err := os.ReadFile(destination)
	require.NoError(t, err)
	assert.Equal(t, contentBySource[succeeded.source], string(destinationContents))
}
