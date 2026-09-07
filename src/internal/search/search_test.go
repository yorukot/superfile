package search

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func buildSearchTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	mkfile := func(rel string) {
		t.Helper()
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mkfile("alpha.txt")
	mkfile(".secret.txt")
	mkfile("docs/readme.md")
	mkfile("src/main.go")
	mkfile("src/main_test.go")
	return root
}

func runCollect(t *testing.T, ctx context.Context, root, query string, includeHidden bool) []Progress {
	t.Helper()
	var snapshots []Progress
	Run(ctx, root, query, includeHidden, func(p Progress) {
		snapshots = append(snapshots, p)
	})
	return snapshots
}

func TestRunEmptyQuery(t *testing.T) {
	root := buildSearchTree(t)
	snapshots := runCollect(t, context.Background(), root, "", false)

	if len(snapshots) != 1 || !snapshots[0].Done {
		t.Fatalf("empty query should emit a single done progress, got %v", snapshots)
	}
	if len(snapshots[0].Results) != 0 {
		t.Errorf("empty query produced %d results, want 0", len(snapshots[0].Results))
	}
}

func TestRunMatchesRelativePaths(t *testing.T) {
	root := buildSearchTree(t)
	snapshots := runCollect(t, context.Background(), root, "main", false)
	last := snapshots[len(snapshots)-1]

	if !last.Done {
		t.Fatal("final progress is not marked done")
	}
	found := map[string]bool{}
	for _, r := range last.Results {
		found[r.Path] = true
	}
	for _, path := range []string{"src/main.go", "src/main_test.go"} {
		if !found[path] {
			t.Errorf("missing match %q in %v", path, last.Results)
		}
	}
	if found["alpha.txt"] || found["docs/readme.md"] {
		t.Errorf("unexpected matches: %v", last.Results)
	}
}

func TestRunMatchPositions(t *testing.T) {
	root := buildSearchTree(t)
	snapshots := runCollect(t, context.Background(), root, "main", false)
	last := snapshots[len(snapshots)-1]

	for _, r := range last.Results {
		if r.Path != "src/main.go" {
			continue
		}
		// "main" is the subsequence matched in "src/main.go".
		if len(r.Positions) != 4 {
			t.Errorf("positions = %v, want 4 byte offsets", r.Positions)
		}
		for _, pos := range r.Positions {
			if pos < 0 || pos >= len(r.Path) {
				t.Errorf("position %d out of range for %q", pos, r.Path)
			}
		}
		return
	}
	t.Error("match for src/main.go not found")
}

func TestRunExactPrefixSyntax(t *testing.T) {
	root := buildSearchTree(t)
	// '^src/' is fzf's prefix operator: only paths starting with "src/".
	snapshots := runCollect(t, context.Background(), root, "^src/", false)
	last := snapshots[len(snapshots)-1]

	if len(last.Results) != 2 {
		t.Fatalf("prefix search matched %d results, want 2: %v", len(last.Results), last.Results)
	}
	for _, r := range last.Results {
		if r.Path != "src/main.go" && r.Path != "src/main_test.go" {
			t.Errorf("unexpected match %q", r.Path)
		}
	}
}

func TestRunHiddenToggle(t *testing.T) {
	root := buildSearchTree(t)
	hidden := runCollect(t, context.Background(), root, "secret", true)
	if len(hidden[len(hidden)-1].Results) != 1 {
		t.Errorf("hidden search matched %d results, want 1", len(hidden[len(hidden)-1].Results))
	}
	visible := runCollect(t, context.Background(), root, "secret", false)
	if len(visible[len(visible)-1].Results) != 0 {
		t.Errorf("non-hidden search matched %d results, want 0", len(visible[len(visible)-1].Results))
	}
}

func TestRunCountsUnreadableDirs(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("permission-based test requires a non-windows, non-root environment")
	}
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "locked"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "plain.txt"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filepath.Join(root, "locked"), 0o000); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chmod(filepath.Join(root, "locked"), 0o755) }()

	snapshots := runCollect(t, context.Background(), root, "txt", false)
	last := snapshots[len(snapshots)-1]
	if last.UnreadableDirs != 1 {
		t.Errorf("unreadable dirs = %d, want 1", last.UnreadableDirs)
	}
	if len(last.Results) != 1 || last.Results[0].Path != "plain.txt" {
		t.Errorf("results = %v, want [plain.txt]", last.Results)
	}
}

func TestRunStopsOnCancellation(t *testing.T) {
	root := buildSearchTree(t)
	for i := range 100 {
		dir := filepath.Join(root, "d"+string(rune('a'+i%26)))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		for j := range 20 {
			if err := os.WriteFile(filepath.Join(dir, "f"+string(rune('0'+j))+".txt"), nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() {
		emits := 0
		Run(ctx, root, "f", false, func(p Progress) {
			emits++
			if emits == 2 {
				cancel()
			}
		})
		done <- emits
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}

func TestRunResultCap(t *testing.T) {
	root := t.TempDir()
	for i := range ResultLimit + 100 {
		name := "m" + string(rune('0'+i%10)) + string(rune('0'+i/10)) + ".txt"
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	snapshots := runCollect(t, context.Background(), root, "m", false)
	last := snapshots[len(snapshots)-1]
	if len(last.Results) > ResultLimit {
		t.Errorf("results = %d, want at most %d", len(last.Results), ResultLimit)
	}
	if last.MatchCount <= ResultLimit {
		t.Errorf("match count = %d, want more than the result limit", last.MatchCount)
	}
}
