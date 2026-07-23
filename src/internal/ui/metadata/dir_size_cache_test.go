package metadata

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDirectorySizeCache(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(
		filepath.Join(tmp, "test.txt"),
		[]byte("hello"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	first := getDirectorySize(tmp)
	second := getDirectorySize(tmp)

	if first != second {
		t.Errorf(
			"cached size mismatch: first=%d second=%d",
			first,
			second,
		)
	}
}