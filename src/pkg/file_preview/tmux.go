package filepreview

import (
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const (
	apcStart             = "\x1b_"
	apcEnd               = "\x1b\\"
	tmuxPassthroughStart = "\x1bPtmux;\x1b\x1b"
	tmuxPassthroughEnd   = "\x1b\\"
	tmuxMinMajorVersion  = 3
	tmuxMinMinorVersion  = 3
)

// inTmux reports whether the process is running inside a tmux pane.
func inTmux() bool {
	return os.Getenv("TMUX") != "" ||
		strings.EqualFold(os.Getenv("TERM_PROGRAM"), "tmux") ||
		strings.HasPrefix(os.Getenv("TERM"), "tmux") ||
		strings.HasPrefix(os.Getenv("TERM"), "screen")
}

var (
	tmuxPassthroughOnce sync.Once
	tmuxPassthroughOK   bool
)

func ensureTmuxPassthrough() bool {
	tmuxPassthroughOnce.Do(func() {
		if !inTmux() || !tmuxVersionSupportsPassthrough() {
			return
		}
		cmd := exec.Command("tmux", "set", "-p", "allow-passthrough", "all",
			";", "set", "-s", "input-buffer-size", "104857600")
		if err := cmd.Run(); err != nil {
			slog.Error("Failed to enable tmux passthrough options", "error", err)
			return
		}
		tmuxPassthroughOK = true
		slog.Debug("tmux passthrough enabled", "allow-passthrough", "all",
			"input-buffer-size", 104857600)
	})
	return tmuxPassthroughOK
}

func tmuxVersionSupportsPassthrough() bool {
	out, err := exec.Command("tmux", "-V").Output()
	if err != nil {
		return false
	}
	major, minor, ok := parseTmuxVersion(string(out))
	if !ok {
		return false
	}
	return major > tmuxMinMajorVersion ||
		(major == tmuxMinMajorVersion && minor >= tmuxMinMinorVersion)
}

func parseTmuxVersion(output string) (int, int, bool) {
	for _, field := range strings.Fields(output) {
		if idx := strings.LastIndex(field, "-"); idx >= 0 && idx+1 < len(field) {
			field = field[idx+1:]
		}
		if len(field) == 0 || field[0] < '0' || field[0] > '9' {
			continue
		}
		majorStr, minorRest, _ := strings.Cut(field, ".")
		major, err := strconv.Atoi(majorStr)
		if err != nil {
			continue
		}
		digits := minorRest
		for len(digits) > 0 && digits[0] >= '0' && digits[0] <= '9' {
			digits = digits[1:]
		}
		minor, err := strconv.Atoi(strings.TrimSuffix(minorRest, digits))
		if err != nil {
			continue
		}
		return major, minor, true
	}
	return 0, 0, false
}

func tmuxPassthroughWrap(raw string) string {
	if raw == "" || !ensureTmuxPassthrough() {
		return raw
	}
	return tmuxWrap(raw)
}

// tmuxWrap performs the actual passthrough wrapping.
// Every ESC inside the wrapped sequence must be doubled.
// The algorithm mirrors yazi's TmuxPassthrough implementation.
func tmuxWrap(raw string) string {
	var out strings.Builder
	for raw != "" {
		start := strings.Index(raw, apcStart)
		if start < 0 {
			out.WriteString(raw)
			break
		}
		out.WriteString(raw[:start])
		raw = raw[start:]

		end := strings.Index(raw, apcEnd)
		if end < 0 {
			out.WriteString(raw)
			break
		}
		apc := raw[:end+len(apcEnd)]
		raw = raw[end+len(apcEnd):]

		// Drop only the leading ESC; the APC intro '_' must remain, tmux
		// un-doubles the ESCs and forwards the original sequence as-is.
		out.WriteString(tmuxPassthroughStart)
		out.WriteString(strings.ReplaceAll(apc[1:], "\x1b", "\x1b\x1b"))
		out.WriteString(tmuxPassthroughEnd)
	}
	return out.String()
}
