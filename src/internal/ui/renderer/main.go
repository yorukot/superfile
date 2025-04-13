package renderer

type Renderer struct {
	contentLines []string
	// Final rendered string should have exactly this many lines, including borders
	totalHeight int
	// Every line should have at most this many characters, including borders
	totalWidth int

	// Usually it would be the height of the terminal - 2 (if border is required)
	maxLines int
	// Usually it would be totalWidth - 2 (if border is required)
	maxLineWidth   int
	borderRequired bool
	truncateStyle  TruncateStyle
}



type TruncateStyle int

const (
	TruncateStyleLeft TruncateStyle = iota
	TruncateStyleMiddle
	TruncateStyleRight
)

func (r *Renderer) AddLine(line string) {
	r.contentLines = append(r.contentLines, line)
}

func New(totalHeight int, totalWidth int, borderRequired bool, truncateStyle TruncateStyle) *Renderer {
	res := Renderer{
		contentLines:   make([]string, 0),
		totalHeight:    totalHeight,
		totalWidth:     totalWidth,
		borderRequired: borderRequired,
		truncateStyle:  truncateStyle,
	}

	res.maxLines = res.totalHeight
	if res.borderRequired {
		res.maxLines -= 2
	}

	res.maxLineWidth = res.totalWidth
	if res.borderRequired {
		res.maxLineWidth -= 2
	}

	return &res
}
