package search

import "sort"

// Caps kept and rendered matches.
const ResultLimit = 5000

type rankedResult struct {
	result Result
	order  int64
}

// Keeps best results. Ties break by arrival. Use from one goroutine.
type topN struct {
	limit int
	items []rankedResult // min-heap, worst kept result at front
	order int64
}

// Makes an empty ranking that keeps at most limit results.
func newTopN(limit int) *topN {
	return &topN{limit: limit}
}

// Keeps the result when it beats the worst one held so far.
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

// Reports whether a sorts ahead of b: higher score first, earlier arrival on ties.
func outranks(a, b rankedResult) bool {
	if a.result.Score != b.result.Score {
		return a.result.Score > b.result.Score
	}
	return a.order < b.order
}

// Returns the kept results ordered best first.
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

// Restores the heap order after appending at index i.
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

// Restores the heap order after replacing the front at index i.
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
