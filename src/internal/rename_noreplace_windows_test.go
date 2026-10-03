//go:build windows

package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExtendedLengthPath(t *testing.T) {
	t.Run("local path", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "file.txt")
		got, err := extendedLengthPath(path)
		if err != nil {
			t.Fatal(err)
		}
		want := `\\?\` + strings.ReplaceAll(path, "/", `\`)
		if got != want {
			t.Fatalf("extendedLengthPath(%q) = %q, want %q", path, got, want)
		}
	})

	t.Run("UNC path", func(t *testing.T) {
		const path = `\\server\share\directory\file.txt`
		const want = `\\?\UNC\server\share\directory\file.txt`

		got, err := extendedLengthPath(path)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Fatalf("extendedLengthPath(%q) = %q, want %q", path, got, want)
		}
	})

	for _, path := range []string{`\\?\C:\directory\file.txt`, `\??\C:\directory\file.txt`, `\\.\C:`} {
		t.Run("preserves "+path[:4], func(t *testing.T) {
			got, err := extendedLengthPath(path)
			if err != nil {
				t.Fatal(err)
			}
			if got != path {
				t.Fatalf("extendedLengthPath(%q) = %q, want unchanged", path, got)
			}
		})
	}
}

func TestRenameNoReplaceSupportsLongPaths(t *testing.T) {
	longDir := filepath.Join(
		t.TempDir(),
		strings.Repeat("a", 100),
		strings.Repeat("b", 100),
		strings.Repeat("c", 100),
	)
	if len(longDir) <= 260 {
		t.Fatalf("test path length = %d, want more than MAX_PATH", len(longDir))
	}
	if err := os.MkdirAll(longDir, 0o755); err != nil {
		t.Fatal(err)
	}

	src := filepath.Join(longDir, "source.txt")
	dst := filepath.Join(longDir, "destination.txt")
	if err := os.WriteFile(src, []byte("source"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := renameNoReplace(src, dst); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "source" {
		t.Fatalf("destination content = %q, want %q", content, "source")
	}

	if err := os.WriteFile(src, []byte("new source"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := renameNoReplace(src, dst); err == nil {
		t.Fatal("renameNoReplace replaced an existing destination")
	}
	content, err = os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "new source" {
		t.Fatalf("source content after collision = %q, want %q", content, "new source")
	}
	content, err = os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "source" {
		t.Fatalf("destination content after collision = %q, want %q", content, "source")
	}
}
