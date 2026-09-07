package search

import (
	"context"

	"github.com/reinhrst/fzf-lib"
)

// matchBatchSize is the number of entries scored in one matcher run. Between
// batches the pipeline checks for cancellation and reports progress.
const matchBatchSize = 2000

// Run performs one search session: it walks the tree under root, matches
// entry paths against query with fzf semantics (fuzzy by default, extended
// syntax available), and reports progress through emit. emit is called from
// the calling goroutine, so Run is safe to run inside a tea.Cmd.
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
	flush()
	emit(snapshot(true))
}

// matchBatch scores one batch of candidate paths against the query and
// returns the matches best first.
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
			Positions: match.Positions,
		})
	}
	return matches
}
