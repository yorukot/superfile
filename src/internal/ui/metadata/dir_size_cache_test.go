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
	size := getDirectorySize("/nonexistent/path/that/should/not/exist")
	if size != 0 {
		t.Errorf(
			"expected 0 for non-existent path, got %d",
			size,
		)
	}
}
