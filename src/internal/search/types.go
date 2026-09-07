package search

// Result is a single match of a search session.
type Result struct {
	Path      string // path relative to the search root
	Dir       bool   // whether the matched entry is a directory
	Score     int    // fzf match score
	Positions []int  // byte offsets of matched characters in Path
}

// Progress is a snapshot of a running search session.
type Progress struct {
	Results        []Result
	MatchCount     int64 // total matches seen so far, including those beyond the result limit
	UnreadableDirs int
	Done           bool
}
