package search

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func benchmarkPaths(count int) []Result {
	paths := make([]Result, 0, count)
	for i := range count {
		dir := fmt.Sprintf("src/pkg%02d", i%40)
		name := fmt.Sprintf("file%d_%s.txt", i, string(rune('a'+i%26)))
		if i%5 == 0 {
			paths = append(paths, Result{Path: dir + "/" + name + ".go", Dir: false})
		} else {
			paths = append(paths, Result{Path: dir + "/" + name, Dir: false})
		}
	}
	return paths
}

func BenchmarkMatchBatch100k(b *testing.B) {
	paths := benchmarkPaths(100_000)
	b.ResetTimer()
	for range b.N {
		_ = matchBatch("f1_a", paths)
	}
}

func BenchmarkRunPipeline100k(b *testing.B) {
	// 10k files cover walk+match and save 100k disk writes.
	root := b.TempDir()
	for i := range 10_000 {
		dir := filepath.Join(root, "d"+string(rune('a'+i%26)))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			b.Fatal(err)
		}
		path := filepath.Join(dir, fmt.Sprintf("file%d.txt", i))
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	for range b.N {
		Run(context.Background(), root, "file1", false, func(Progress) {})
	}
}
