package components

import (
	"bytes"
	"strings"
	
	charmansi "github.com/charmbracelet/x/exp/term/ansi"
	"github.com/mattn/go-runewidth"
	ansi "github.com/muesli/reflow/ansi"
	"github.com/muesli/reflow/truncate"
	"github.com/muesli/termenv"
)

// whitespace is a whitespace renderer.
type whitespace struct {
	style termenv.Style
	chars string
}

type WhitespaceOption func(*whitespace)

// Render whitespaces.
func (w whitespace) render(width int) string {
	if w.chars == "" {
		w.chars = " "
	}

	r := []rune(w.chars)
	j := 0
	b := strings.Builder{}

	// Cycle through runes and print them into the whitespace.
	for i := 0; i < width; {
		b.WriteRune(r[j])
		j++
		if j >= len(r) {
			j = 0
		}
		i += charmansi.StringWidth(string(r[j]))
	}

	// Fill any extra gaps white spaces. This might be necessary if any runes
	// are more than one cell wide, which could leave a one-rune gap.
	short := width - charmansi.StringWidth(b.String())
	if short > 0 {
		b.WriteString(strings.Repeat(" ", short))
	}

	return w.style.Styled(b.String())
}

// PlaceOverlay places fg on top of bg.
func PlaceOverlay(x, y int, fg, bg string, opts ...WhitespaceOption) string {
	fgLines, fgWidth := getLines(fg)
	bgLines, bgWidth := getLines(bg)
	bgHeight := len(bgLines)
	fgHeight := len(fgLines)

	if fgWidth >= bgWidth && fgHeight >= bgHeight {
		// FIXME: return fg or bg?
		return fg
	}
	// TODO: allow placement outside of the bg box?
	x = clamp(x, 0, bgWidth-fgWidth)
	y = clamp(y, 0, bgHeight-fgHeight)

	ws := &whitespace{}
	for _, opt := range opts {
		opt(ws)
	}

	var b strings.Builder
	for i, bgLine := range bgLines {
		if i > 0 {
			b.WriteByte('\n')
		}
		if i < y || i >= y+fgHeight {
			b.WriteString(bgLine)
			continue
		}

		pos := 0
		if x > 0 {
			left := truncate.String(bgLine, uint(x))
			pos = ansi.PrintableRuneWidth(left)
			b.WriteString(left)
			if pos < x {
				b.WriteString(ws.render(x - pos))
				pos = x
			}
		}

		fgLine := fgLines[i-y]
		b.WriteString(fgLine)
		pos += ansi.PrintableRuneWidth(fgLine)

		right := cutLeft(bgLine, pos)
		bgWidth := ansi.PrintableRuneWidth(bgLine)
		rightWidth := ansi.PrintableRuneWidth(right)
		if rightWidth <= bgWidth-pos {
			b.WriteString(ws.render(bgWidth - rightWidth - pos))
		}

		b.WriteString(right)
	}

	return b.String()
}

// cutLeft cuts printable characters from the left.
// This function is heavily based on muesli's ansi and truncate packages.
func cutLeft(s string, cutWidth int) string {
	var (
		pos    int
		isAnsi bool
		ab     bytes.Buffer
		b      bytes.Buffer
	)
	for _, c := range s {
		var w int
		if c == ansi.Marker || isAnsi {
			isAnsi = true
			ab.WriteRune(c)
			if ansi.IsTerminator(c) {
				isAnsi = false
				if bytes.HasSuffix(ab.Bytes(), []byte("[0m")) {
					ab.Reset()
				}
			}
		} else {
			w = runewidth.RuneWidth(c)
		}

		if pos >= cutWidth {
			if b.Len() == 0 {
				if ab.Len() > 0 {
					b.Write(ab.Bytes())
				}
				if pos-cutWidth > 1 {
					b.WriteByte(' ')
					continue
				}
			}
			b.WriteRune(c)
		}
		pos += w
	}
	return b.String()
}

func clamp(v, lower, upper int) int {
	return min(max(v, lower), upper)
}

// Split a string into lines, additionally returning the size of the widest
// line.
func getLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")

	for _, l := range lines {
		w := charmansi.StringWidth(l)
		if widest < w {
			widest = w
		}
	}

	return lines, widest
}


// import (
// 	"bytes"
// 	"regexp"
// 	"strings"
// 	"unicode/utf8"

// 	"github.com/mattn/go-runewidth"
// 	"github.com/muesli/reflow/ansi"
// 	"github.com/muesli/reflow/truncate"
// 	"github.com/muesli/termenv"
// 	"github.com/rivo/uniseg"
// )

// var re = regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)

// func placeOverlay(x, y int, bg, fg string) string {
// 	BGLines := strings.Split(bg, "\n")
// 	FGLines := strings.Split(fg, "\n")

// 	FGNoStyleText := noStyleString(fg)
// 	FGNoStyleTextLines := strings.Split(FGNoStyleText, "\n")

// 	BGNoStyleText := noStyleString(bg)
// 	BGNoStyleTextLines := strings.Split(BGNoStyleText, "\n")
	
// 	for Y := y; Y < y+len(FGNoStyleTextLines); Y++ {
// 		end := x+utf8.RuneCountInString(FGNoStyleTextLines[Y - y])

// 		// BGAsciiIndexes := re.FindAllStringIndex(BGLines[Y], -1)
// 		// BGAsciiStringList := mapCoordsList(BGLines[Y], BGAsciiIndexes)
		
// 		BGLines[Y] = replaceString(x,end, BGNoStyleTextLines[Y], FGLines[Y-y])
// 	}

// 	newBGText := strings.Join(BGLines, "\n")
// 	return newBGText
// }

// func substring(start, end int, str string) string {
// 	runes := []rune(str)

// 	substr := string(runes[start:end])

// 	return substr
// }

// func insertString(index int, str, insertStr string) string {
// 	runes := []rune(str)

// 	beforeIndex := string(runes[:index])

// 	afterIndex := string(runes[index:])

// 	return beforeIndex + insertStr + afterIndex
// }

// func replaceString(start, end int, str, replaceStr string) string {
// 	runes := []rune(str)

// 	beforeX := string(runes[:start])

// 	afterY := string(runes[end:])

// 	return beforeX + replaceStr + afterY
// }

// func noStyleString(styledText string) string {
// 	re := regexp.MustCompile(`\x1b\[[0-9;]*[mK]`)

// 	plainText := re.ReplaceAllString(styledText, "")

// 	return plainText
// }

// func checkIsAsciiString(locationX int, asciiStringList [][]int) bool {
// 	for _, loc := range asciiStringList {
// 		if locationX < loc[0] || locationX > loc[1] {
// 			return false
// 		}
// 	}
// 	return true
// }

// func mapCoordsList(s string, indexes [][]int) (graphemeCoordsList [][]int) {
// 	for _, loc := range indexes {
// 		loc = mapCoords(s, loc)
// 		graphemeCoordsList = append(graphemeCoordsList, loc)
// 	}
// 	return
// }

// func mapCoords(s string, byteCoords []int) (graphemeCoords []int) {
// 	graphemeCoords = make([]int, 2)
// 	gr := uniseg.NewGraphemes(s)
// 	graphemeIndex := -1
// 	for gr.Next() {
// 		graphemeIndex++
// 		a, b := gr.Positions()
// 		if a == byteCoords[0] {
// 			graphemeCoords[0] = graphemeIndex
// 		}
// 		if b == byteCoords[1] {
// 			graphemeCoords[1] = graphemeIndex + 1
// 			break
// 		}
// 	}
// 	return
// }

// type Renderer struct {
// 	output            *termenv.Output
// 	colorProfile      termenv.Profile
// 	hasDarkBackground bool
// }

// type whitespace struct {
// 	re    *Renderer
// 	style termenv.Style
// 	chars string
// }
// type WhitespaceOption func(*whitespace)

// func PlaceOverlay(x, y int, fg, bg string, opts ...WhitespaceOption) string {
// 	fgLines, fgWidth := getLines(fg)
// 	bgLines, bgWidth := getLines(bg)
// 	bgHeight := len(bgLines)
// 	fgHeight := len(fgLines)

// 	if fgWidth >= bgWidth && fgHeight >= bgHeight {
// 		// FIXME: return fg or bg?
// 		return fg
// 	}
// 	// TODO: allow placement outside of the bg box?
// 	x = clamp(x, 0, bgWidth-fgWidth)
// 	y = clamp(y, 0, bgHeight-fgHeight)

// 	ws := &whitespace{}
// 	for _, opt := range opts {
// 		opt(ws)
// 	}

// 	var b strings.Builder
// 	for i, bgLine := range bgLines {
// 		if i > 0 {
// 			b.WriteByte('\n')
// 		}
// 		if i < y || i >= y+fgHeight {
// 			b.WriteString(bgLine)
// 			continue
// 		}

// 		pos := 0
// 		if x > 0 {
// 			left := truncate.String(bgLine, uint(x))
// 			pos = ansi.PrintableRuneWidth(left)
// 			b.WriteString(left)
// 			if pos < x {
// 				b.WriteString(ws.render(x - pos))
// 				pos = x
// 			}
// 		}

// 		fgLine := fgLines[i-y]
// 		b.WriteString(fgLine)
// 		pos += ansi.PrintableRuneWidth(fgLine)

// 		right := cutLeft(bgLine, pos)
// 		bgWidth := ansi.PrintableRuneWidth(bgLine)
// 		rightWidth := ansi.PrintableRuneWidth(right)
// 		if rightWidth <= bgWidth-pos {
// 			b.WriteString(ws.render(bgWidth - rightWidth - pos))
// 		}

// 		b.WriteString(right)
// 	}

// 	return b.String()
// }

// // cutLeft cuts printable characters from the left.
// // This function is heavily based on muesli's ansi and truncate packages.
// func cutLeft(s string, cutWidth int) string {
// 	var (
// 		pos    int
// 		isAnsi bool
// 		ab     bytes.Buffer
// 		b      bytes.Buffer
// 	)
// 	for _, c := range s {
// 		var w int
// 		if c == ansi.Marker || isAnsi {
// 			isAnsi = true
// 			ab.WriteRune(c)
// 			if ansi.IsTerminator(c) {
// 				isAnsi = false
// 				if bytes.HasSuffix(ab.Bytes(), []byte("[0m")) {
// 					ab.Reset()
// 				}
// 			}
// 		} else {
// 			w = runewidth.RuneWidth(c)
// 		}

// 		if pos >= cutWidth {
// 			if b.Len() == 0 {
// 				if ab.Len() > 0 {
// 					b.Write(ab.Bytes())
// 				}
// 				if pos-cutWidth > 1 {
// 					b.WriteByte(' ')
// 					continue
// 				}
// 			}
// 			b.WriteRune(c)
// 		}
// 		pos += w
// 	}
// 	return b.String()
// }

// func clamp(v, lower, upper int) int {
// 	return min(max(v, lower), upper)
// }

// func getLines(s string) (lines []string, widest int) {
// 	lines = strings.Split(s, "\n")

// 	for _, l := range lines {
// 		w := ansi.PrintableRuneWidth(l)
// 		if widest < w {
// 			widest = w
// 		}
// 	}

// 	return lines, widest
// }