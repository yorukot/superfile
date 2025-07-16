package utils

import "github.com/reinhrst/fzf-lib"

// Returning a string slice causes inefficiency in current usage
func FzfSearch(query string, source []string) []fzf.MatchResult {
	fzfSearcher := fzf.New(source, fzf.DefaultOptions())
	fzfSearcher.Search(query)
	// TODO : This is a blocking call, which will cause the UI to freeze if the query is slow.
	// Need to put a timeout on this
	fzfResults := <-fzfSearcher.GetResultChannel()
	fzfSearcher.End()
	return fzfResults.Matches
}
