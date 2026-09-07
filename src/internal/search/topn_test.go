package search

import (
	"testing"
)

func TestTopNKeepsBestScores(t *testing.T) {
	top := newTopN(3)
	top.add(Result{Path: "low", Score: 10})
	top.add(Result{Path: "high", Score: 100})
	top.add(Result{Path: "mid", Score: 50})
	top.add(Result{Path: "top", Score: 200})

	got := top.sorted()
	if len(got) != 3 {
		t.Fatalf("got %d results, want 3", len(got))
	}
	want := []string{"top", "high", "mid"}
	for i, path := range want {
		if got[i].Path != path {
			t.Errorf("result %d = %q, want %q", i, got[i].Path, path)
		}
	}
}

func TestTopNTiebreakByArrivalOrder(t *testing.T) {
	top := newTopN(2)
	top.add(Result{Path: "first", Score: 50})
	top.add(Result{Path: "second", Score: 50})
	top.add(Result{Path: "third", Score: 60})

	got := top.sorted()
	want := []string{"third", "first"}
	for i, path := range want {
		if got[i].Path != path {
			t.Errorf("result %d = %q, want %q", i, got[i].Path, path)
		}
	}
}

func TestTopNRespectsLimit(t *testing.T) {
	top := newTopN(5)
	for i := range 100 {
		top.add(Result{Path: "p", Score: i})
	}
	got := top.sorted()
	if len(got) != 5 {
		t.Fatalf("got %d results, want 5", len(got))
	}
	// Best scores win.
	for i, item := range got {
		if want := 99 - i; item.Score != want {
			t.Errorf("result %d score = %d, want %d", i, item.Score, want)
		}
	}
}
