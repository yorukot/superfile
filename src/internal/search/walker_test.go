package search

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func buildTree(t *testing.T) string {
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
	mkdir := func(rel string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(root, rel), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	mkfile("visible.txt")
	mkfile(".hidden.txt")
	mkfile("sub/nested.txt")
	mkfile("sub/deeper/inner.txt")
	mkdir("empty_dir")
	if runtime.GOOS != "windows" {
		if err := os.Symlink("sub", filepath.Join(root, "link_to_sub")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(".", filepath.Join(root, "cycle")); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("missing_target", filepath.Join(root, "broken")); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestWalkVisitsAllVisibleEntries(t *testing.T) {
	root := buildTree(t)

	var paths []string
	unreadable := 0
	walk(context.Background(), root, false, &unreadable, func(rel string, _ bool) {
		paths = append(paths, rel)
	})

	if unreadable != 0 {
		t.Errorf("unexpected unreadable dirs: %d", unreadable)
	}
	want := map[string]bool{
		"visible.txt":          true,
		"sub":                  true,
		"sub/nested.txt":       true,
		"sub/deeper":           true,
		"sub/deeper/inner.txt": true,
		"empty_dir":            true,
		"broken":               true,
	}
	if runtime.GOOS != "windows" {
		// The symlink target of cycle is the root itself, an ancestor, so it
		// is listed but not descended into.
		want["cycle"] = true
		// link_to_sub points at a sibling subtree: it is not a cycle, so its
		// contents are walked under the link path as well.
		want["link_to_sub"] = true
		want["link_to_sub/nested.txt"] = true
		want["link_to_sub/deeper"] = true
		want["link_to_sub/deeper/inner.txt"] = true
	}
	if len(paths) != len(want) {
		t.Errorf("walked %d entries, want %d: %v", len(paths), len(want), paths)
	}
	for _, p := range paths {
		if !want[p] {
			t.Errorf("unexpected path: %q", p)
		}
		delete(want, p)
	}
	for p := range want {
		t.Errorf("missing path: %q", p)
	}
}

func TestWalkIncludesHiddenWhenAsked(t *testing.T) {
	root := buildTree(t)

	withHidden := 0
	walk(context.Background(), root, true, new(int), func(rel string, _ bool) {
		if rel == ".hidden.txt" {
			withHidden++
		}
	})
	if withHidden != 1 {
		t.Errorf("hidden file visited %d times, want 1", withHidden)
	}

	withoutHidden := 0
	walk(context.Background(), root, false, new(int), func(rel string, _ bool) {
		if rel == ".hidden.txt" {
			withoutHidden++
		}
	})
	if withoutHidden != 0 {
		t.Errorf("hidden file visited %d times, want 0", withoutHidden)
	}
}

func TestWalkDoesNotRecurseIntoHiddenDirs(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".git", "objects"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git", "objects", "pack"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	visitedGit := 0
	walk(context.Background(), root, false, new(int), func(rel string, _ bool) {
		if rel == ".git" || rel == ".git/objects" || rel == ".git/objects/pack" {
			visitedGit++
		}
	})
	if visitedGit != 0 {
		t.Errorf("hidden dir entries visited %d times, want 0", visitedGit)
	}
}

func TestWalkSymlinkCycleTerminates(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require privileges on windows")
	}
	root := buildTree(t)

	done := make(chan int, 1)
	go func() {
		unreadable := 0
		walk(context.Background(), root, false, &unreadable, func(rel string, isDir bool) {
			if rel == "cycle" && !isDir {
				t.Error("cycle entry visited as file")
			}
			if hasPrefixPath(rel, "cycle/") {
				t.Errorf("walked inside cycle: %q", rel)
			}
		})
		done <- 0
	}()
	select {
	case <-done:
	case <-timeout(t):
		t.Fatal("walk did not terminate; symlink cycle not handled")
	}
}

func TestWalkCountsUnreadableDirs(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("permission-based test requires a non-windows, non-root environment")
	}
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "locked"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "ok.txt"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filepath.Join(root, "locked"), 0o000); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chmod(filepath.Join(root, "locked"), 0o755) }()

	unreadable := 0
	visited := 0
	walk(context.Background(), root, false, &unreadable, func(_ string, _ bool) {
		visited++
	})
	if unreadable != 1 {
		t.Errorf("unreadable dirs = %d, want 1", unreadable)
	}
	// ok.txt and the locked dir itself are both visible entries.
	if visited != 2 {
		t.Errorf("visited %d entries, want 2 (ok.txt and locked)", visited)
	}
}

// hasPrefixPath reports whether path starts with the dir prefix at a path
// segment boundary.
func hasPrefixPath(path, prefix string) bool {
	return len(path) >= len(prefix) && path[:len(prefix)] == prefix
}

func TestWalkStopsOnCancellation(t *testing.T) {
	root := t.TempDir()
	for i := range 200 {
		dir := filepath.Join(root, "d"+string(rune('a'+i%26))+string(rune('0'+i/26)))
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		for j := range 10 {
			if err := os.WriteFile(filepath.Join(dir, "f"+string(rune('0'+j))+".txt"), nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	visited := 0
	done := make(chan int, 1)
	go func() {
		unreadable := 0
		walk(ctx, root, false, &unreadable, func(_ string, _ bool) {
			visited++
			if visited == 50 {
				cancel()
			}
		})
		done <- 0
	}()
	select {
	case <-done:
	case <-timeout(t):
		t.Fatal("walk did not stop after cancellation")
	}
	if visited > 500 {
		t.Errorf("walked %d entries after cancellation, expected to stop early", visited)
	}
}

func timeout(t *testing.T) <-chan struct{} {
	t.Helper()
	done := make(chan struct{})
	go func() {
		waitFor := 5 * time.Second
		<-time.After(waitFor)
		close(done)
	}()
	return done
}
