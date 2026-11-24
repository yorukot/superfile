package filepreview

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/yorukot/superfile/src/pkg/utils"
)

const (
	apcStart                    = "\x1b_"
	apcEnd                      = "\x1b\\"
	tmuxPassthroughStart        = "\x1bPtmux;\x1b\x1b" //nolint:gosec // tmux passthrough prefix, not a credential.
	tmuxPassthroughEnd          = "\x1b\\"
	tmuxMinMajorVersion         = 3
	tmuxMinMinorVersion         = 3
	tmuxAllMajorVersion         = 3
	tmuxAllMinorVersion         = 4
	tmuxInputBufferMajorVersion = 3
	tmuxInputBufferMinorVersion = 6
	tmuxInputBufferSize         = 104857600
)

// tmuxPassthroughOnce/tmuxPassthroughOK cache the result of the one-time tmux setup.
var (
	tmuxPassthroughOnce sync.Once //nolint:gochecknoglobals // Process-wide one-time tmux setup.
	tmuxPassthroughOK   bool      //nolint:gochecknoglobals // Process-wide one-time tmux setup.
)

func ensureTmuxPassthrough() bool {
	tmuxPassthroughOnce.Do(func() {
		if !inTmux() {
			return
		}
		major, minor, ok := tmuxVersion()
		if !ok || !tmuxVersionAtLeast(major, minor, tmuxMinMajorVersion, tmuxMinMinorVersion) {
			return
		}

		passthroughValue := tmuxAllowPassthroughValue(major, minor)
		args := []string{"set", "-p", "allow-passthrough", passthroughValue}
		if tmuxSupportsInputBufferSize(major, minor) {
			args = append(args, ";", "set", "-s", "input-buffer-size", strconv.Itoa(tmuxInputBufferSize))
		}

		cmd := exec.Command("tmux", args...)
		if err := cmd.Run(); err != nil {
			slog.Error("Failed to enable tmux passthrough options", "error", err)
			return
		}
		tmuxPassthroughOK = true
		slog.Debug("tmux passthrough enabled", "allow-passthrough", passthroughValue)
	})
	return tmuxPassthroughOK
}

func inTmux() bool {
	return os.Getenv("TMUX") != "" ||
		strings.EqualFold(os.Getenv("TERM_PROGRAM"), "tmux") ||
		strings.HasPrefix(os.Getenv("TERM"), "tmux") ||
		strings.HasPrefix(os.Getenv("TERM"), "screen")
}

func tmuxVersion() (int, int, bool) {
	out, err := exec.Command("tmux", "-V").Output()
	if err != nil {
		return 0, 0, false
	}
	return parseTmuxVersion(string(out))
}

func tmuxVersionAtLeast(major, minor, wantMajor, wantMinor int) bool {
	return major > wantMajor || (major == wantMajor && minor >= wantMinor)
}

// tmuxSupportsInputBufferSize reports whether tmux supports input-buffer-size (3.6+).
func tmuxSupportsInputBufferSize(major, minor int) bool {
	return tmuxVersionAtLeast(major, minor, tmuxInputBufferMajorVersion, tmuxInputBufferMinorVersion)
}

// tmuxAllowPassthroughValue returns "all" on tmux 3.4+, otherwise "on".
func tmuxAllowPassthroughValue(major, minor int) string {
	if tmuxVersionAtLeast(major, minor, tmuxAllMajorVersion, tmuxAllMinorVersion) {
		return "all"
	}
	return "on"
}

func parseTmuxVersion(output string) (int, int, bool) {
	for _, field := range strings.Fields(output) {
		if idx := strings.LastIndex(field, "-"); idx >= 0 && idx+1 < len(field) {
			field = field[idx+1:]
		}
		if field == "" || !utils.IsASCIIDigit(field[0]) {
			continue
		}

		majorStr, minorRest, _ := strings.Cut(field, ".")
		major, err := strconv.Atoi(majorStr)
		if err != nil {
			continue
		}

		minorDigits := utils.LeadingASCIIDigits(minorRest)
		if minorDigits == "" {
			continue
		}
		minor, err := strconv.Atoi(minorDigits)
		if err != nil {
			continue
		}

		return major, minor, true
	}
	return 0, 0, false
}

// tmuxPassthroughWrap wraps APC data. Outside tmux it returns raw unchanged;
// inside tmux a failed setup returns an error so callers can fall back.
func tmuxPassthroughWrap(raw string) (string, error) {
	if raw == "" {
		return "", nil
	}
	if !inTmux() {
		return raw, nil
	}
	if !ensureTmuxPassthrough() {
		return "", errors.New("tmux passthrough setup failed")
	}
	return tmuxWrap(raw), nil
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
