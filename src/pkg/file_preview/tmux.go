package filepreview

import (
	"os"
	"os/exec"
	"strings"
)

// tmux hides the outer terminal behind its own TERM/TERM_PROGRAM and consumes
// bare APC sequences, so Kitty graphics only work inside tmux when the payload
// is tunnelled through tmux's DCS passthrough and the terminal tmux is drawing
// to is itself Kitty capable. These helpers ask the running server for both.

const (
	escByte           = "\x1b"
	stringTerminator  = "\x1b\\"
	tmuxPassthroughIn = "\x1bPtmux;"
)

// insideTmux reports whether superfile is running inside a tmux pane.
func insideTmux() bool {
	return os.Getenv("TMUX") != ""
}

// tmuxQuery runs a tmux command and returns its trimmed output, or "" on error.
func tmuxQuery(args ...string) string {
	out, err := exec.Command("tmux", args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

// tmuxOuterTerminals returns the terminal name and type reported by the tmux
// client, which identify the terminal emulator tmux is actually drawing to.
func tmuxOuterTerminals() []string {
	return strings.Split(tmuxQuery("display-message", "-p",
		"#{client_termname}\n#{client_termtype}"), "\n")
}

// tmuxAllowsPassthrough reports whether the server forwards DCS passthrough
// sequences to the outer terminal. Without it the image data is discarded and
// only the Unicode placeholder cells would be drawn.
func tmuxAllowsPassthrough() bool {
	switch tmuxQuery("show-options", "-gv", "allow-passthrough") {
	case "on", "all":
		return true
	}
	return false
}

// tmuxPassthrough wraps every escape sequence in s in tmux's DCS passthrough,
// doubling inner ESC bytes as tmux requires. Sequences are wrapped one at a
// time so that a large chunked image transmission never forms a single DCS
// string big enough to hit tmux's input buffer limit.
func tmuxPassthrough(s string) string {
	var b strings.Builder
	for s != "" {
		seq := s
		if i := strings.Index(s, stringTerminator); i >= 0 {
			seq, s = s[:i+len(stringTerminator)], s[i+len(stringTerminator):]
		} else {
			s = ""
		}
		b.WriteString(tmuxPassthroughIn)
		b.WriteString(strings.ReplaceAll(seq, escByte, escByte+escByte))
		b.WriteString(stringTerminator)
	}
	return b.String()
}

// rawForTerminal prepares raw terminal output for the transport in use. Inside
// tmux the sequences must be tunnelled so they reach the outer terminal.
func rawForTerminal(s string) string {
	if s == "" || !insideTmux() {
		return s
	}
	return tmuxPassthrough(s)
}
