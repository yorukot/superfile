package rendering

import (
	"bytes"
	"strings"
)

func LineWiseAnsiTree(s string) []string {
	var res []string
	for line := range strings.SplitSeq(s, "\n") {
		res = append(res, AnsiTree(line))
	}
	return res
}

// Majorly used for debugging and logging purposes
func AnsiTree(s string) string {
	res := bytes.Buffer{}

	startEsc := false
	b := []byte(s)

	for i := range b {
		switch {
		case startEsc:
			if b[i] == 'm' {
				if b[i-1] == '4' && b[i-2] == '3' {
					// \x1b34m
					res.WriteByte('<')
					startEsc = false
				} else if b[i-1] == '0' && b[i-2] == '[' {
					// \x1b[0m
					res.WriteByte('>')
					startEsc = false
				}
			}
		case b[i] == '\x1b':
			startEsc = true
		default:
			res.WriteByte(b[i])
		}
	}

	return res.String()
}
