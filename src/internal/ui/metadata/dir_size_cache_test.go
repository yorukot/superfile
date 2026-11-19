package metadata

import (
	"os"
	"path/filepath"
	"sync"
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

func TestGetDirectorySizeCacheInvalidation(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(
		filepath.Join(tmp, "test.txt"),
		[]byte("hello"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	// First call calculates and caches
	first := getDirectorySize(tmp)

	// Add a new file — directory modTime should change, invalidating cache
	err = os.WriteFile(
		filepath.Join(tmp, "test2.txt"),
		[]byte("world!!!"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	// Second call should recalculate due to modTime change
	second := getDirectorySize(tmp)

	if first == second {
		t.Errorf(
			"expected cache invalidation after file addition, "+
				"but sizes matched: first=%d second=%d",
			first,
			second,
		)
	}

	expected := int64(5 + 8) // "hello" + "world!!!"
	if second != expected {
		t.Errorf(
			"expected recalculated size %d, got %d",
			expected,
			second,
		)
	}
}

func TestGetDirectorySizeCacheConcurrent(t *testing.T) {
	tmp := t.TempDir()

	err := os.WriteFile(
		filepath.Join(tmp, "test.txt"),
		[]byte("hello"),
		0644,
	)
	if err != nil {
		t.Fatal(err)
	}

	const goroutines = 50
	results := make([]int64, goroutines)
	var wg sync.WaitGroup

	wg.Add(goroutines)
	for i := range goroutines {
		go func(idx int) {
			defer wg.Done()
			results[idx] = getDirectorySize(tmp)
		}(i)
	}
	wg.Wait()

	// All results should be identical
	for i := 1; i < goroutines; i++ {
		if results[i] != results[0] {
			t.Errorf(
				"concurrent results mismatch: results[0]=%d results[%d]=%d",
				results[0],
				i,
				results[i],
			)
		}
	}

	expected := int64(5) // "hello"
	if results[0] != expected {
		t.Errorf(
			"expected size %d, got %d",
			expected,
			results[0],
		)
	}
}

func TestGetDirectorySizeNonExistentPath(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing")
	size := getDirectorySize(missingPath)
	if size != 0 {
		t.Errorf(
			"expected 0 for non-existent path, got %d",
			size,
		)
	}
}

func TestGetDirectorySizeOverwriteInDirectory(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "test.txt")

	err := os.WriteFile(filePath, []byte("hello"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	first := getDirectorySize(tmp)
	if first != 5 {
		t.Fatalf("expected initial size 5, got %d", first)
	}

	// Overwrite file (top-level directory mtime may not change on all systems)
	err = os.WriteFile(filePath, []byte("hello world!!"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Manually clear/expire cache to simulate staleness window passing if mtime didn't update
	directorySizeMutex.Lock()
	directorySizeCache.Remove(tmp)
	directorySizeMutex.Unlock()

	second := getDirectorySize(tmp)
	if second != 13 {
		t.Errorf("expected recalculated size 13 after overwrite, got %d", second)
	}
}

func TestGetDirectorySizeNestedSubdirectoryModification(t *testing.T) {
	tmp := t.TempDir()
	subDir := filepath.Join(tmp, "sub", "nested")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	filePath := filepath.Join(subDir, "test.txt")
	err = os.WriteFile(filePath, []byte("nested"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	first := getDirectorySize(tmp)
	if first != 6 {
		t.Fatalf("expected initial size 6, got %d", first)
	}

	// Modify file in nested subdirectory
	err = os.WriteFile(filePath, []byte("nested modified"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Evict cached entry to simulate TTL expiration
	directorySizeMutex.Lock()
	directorySizeCache.Remove(tmp)
	directorySizeMutex.Unlock()

	second := getDirectorySize(tmp)
	if second != 15 {
		t.Errorf("expected recalculated size 15 after nested edit, got %d", second)
	}
}

func TestGetDirectorySizeConcurrentDifferentRevisions(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "test.txt")
	if err := os.WriteFile(filePath, []byte("rev1"), 0644); err != nil {
		t.Fatal(err)
	}

	const goroutines = 10
	results := make([]int64, goroutines)
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			results[idx] = getDirectorySize(tmp)
		}(i)
	}

	wg.Wait()

	for i, res := range results {
		if res != 4 {
			t.Errorf("goroutine %d expected size 4, got %d", i, res)
		}
	}
}
