package search

type Result struct {
	Path      string // path relative to the search root
	Dir       bool   // whether the matched entry is a directory
	Score     int    // fzf match score
	Positions []int  // UTF-8 byte offsets of matched characters in Path
}

type Progress struct {
	Results        []Result
	MatchCount     int64 // total matches, including those beyond the result limit
	UnreadableDirs int
	Done           bool
}
