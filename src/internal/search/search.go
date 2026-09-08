package search

import (
	"context"
	"unicode/utf8"

	"github.com/reinhrst/fzf-lib"
)

// Matcher scores entries per run and checks cancellation between batches.
const matchBatchSize = 2000

// Run walks root, matches paths with fzf semantics, and reports via emit.
// Run calls emit on the same goroutine, so it stays safe in a tea.Cmd.
func Run(ctx context.Context, root, query string, includeHidden bool, emit func(Progress)) {
	if query == "" {
		emit(Progress{Done: true})
		return
	}

	top := newTopN(ResultLimit)
	unreadable := 0
	var matchCount int64
	batch := make([]Result, 0, matchBatchSize)
	snapshot := func(done bool) Progress {
		return Progress{
			Results:        top.sorted(),
			MatchCount:     matchCount,
			UnreadableDirs: unreadable,
			Done:           done,
		}
	}
	flush := func() {
		if len(batch) == 0 {
			return
		}
		matches := matchBatch(query, batch)
		matchCount += int64(len(matches))
		for _, match := range matches {
			top.add(match)
		}
		batch = batch[:0]
		emit(snapshot(false))
	}

	walk(ctx, root, includeHidden, &unreadable, func(rel string, isDir bool) {
		batch = append(batch, Result{Path: rel, Dir: isDir})
		if len(batch) == matchBatchSize {
			flush()
		}
	})
	if ctx.Err() != nil {
		// Cancelled. Emits nothing.
		return
	}
	flush()
	emit(snapshot(true))
}

func matchBatch(query string, candidates []Result) []Result {
	if len(candidates) == 0 {
		return nil
	}
	items := make([]string, len(candidates))
	for i, candidate := range candidates {
		items[i] = candidate.Path
	}
	searcher := fzf.New(items, fzf.DefaultOptions())
	searcher.Search(query)
	searchResult := <-searcher.GetResultChannel()
	searcher.End()

	matches := make([]Result, 0, len(searchResult.Matches))
	for _, match := range searchResult.Matches {
		idx := int(match.HayIndex)
		matches = append(matches, Result{
			Path:      candidates[idx].Path,
			Dir:       candidates[idx].Dir,
			Score:     match.Score,
			Positions: runePositionsToByteOffsets(candidates[idx].Path, match.Positions),
		})
	}
	return matches
}

// Converts fzf rune indexes to byte offsets for Result.Positions.
func runePositionsToByteOffsets(path string, positions []int) []int {
	if len(positions) == 0 {
		return positions
	}
	ascii := true
	for i := range len(path) {
		if path[i] >= utf8.RuneSelf {
			ascii = false
			break
		}
	}
	if ascii {
		return positions
	}
	posSet := make(map[int]struct{}, len(positions))
	for _, p := range positions {
		posSet[p] = struct{}{}
	}
	out := make([]int, 0, len(positions))
	runeIdx := 0
	for byteOff := range path {
		if _, ok := posSet[runeIdx]; ok {
			out = append(out, byteOff)
		}
		runeIdx++
	}
	return out
}
