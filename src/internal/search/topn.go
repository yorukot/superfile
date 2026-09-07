package search

import "sort"

// ResultLimit bounds how many matches a search session keeps and renders.
const ResultLimit = 5000

type rankedResult struct {
	result Result
	order  int64
}

// topN keeps the highest-scoring results seen so far, breaking score ties by
// arrival (walk) order. Not safe for concurrent use.
type topN struct {
	limit int
	items []rankedResult // min-heap: the worst kept result is at the front
	order int64
}

func newTopN(limit int) *topN {
	return &topN{limit: limit}
}

func (t *topN) add(r Result) {
	item := rankedResult{result: r, order: t.order}
	t.order++
	if len(t.items) < t.limit {
		t.items = append(t.items, item)
		t.siftUp(len(t.items) - 1)
		return
	}
	if !outranks(item, t.items[0]) {
		return
	}
	t.items[0] = item
	t.siftDown(0)
}

// outranks reports whether a sorts ahead of b: a higher score wins, and an
// earlier arrival breaks ties.
func outranks(a, b rankedResult) bool {
	if a.result.Score != b.result.Score {
		return a.result.Score > b.result.Score
	}
	return a.order < b.order
}

// sorted returns the top results, best first.
func (t *topN) sorted() []Result {
	items := make([]rankedResult, len(t.items))
	copy(items, t.items)
	sort.Slice(items, func(i, j int) bool {
		return outranks(items[i], items[j])
	})
	out := make([]Result, len(items))
	for i, item := range items {
		out[i] = item.result
	}
	return out
}

func (t *topN) siftUp(i int) {
	for i > 0 {
		parent := (i - 1) / 2 //nolint:mnd // binary heap parent index
		if outranks(t.items[i], t.items[parent]) {
			return
		}
		t.items[i], t.items[parent] = t.items[parent], t.items[i]
		i = parent
	}
}

func (t *topN) siftDown(i int) {
	n := len(t.items)
	for {
		worst := i
		left, right := 2*i+1, 2*i+2 //nolint:mnd // binary heap child indices
		if left < n && outranks(t.items[worst], t.items[left]) {
			worst = left
		}
		if right < n && outranks(t.items[worst], t.items[right]) {
			worst = right
		}
		if worst == i {
			return
		}
		t.items[i], t.items[worst] = t.items[worst], t.items[i]
		i = worst
	}
}
