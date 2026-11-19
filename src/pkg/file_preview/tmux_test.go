package filepreview

import (
	"strings"
	"testing"
)

func TestParseTmuxVersion(t *testing.T) {
	tests := []struct {
		output string
		major  int
		minor  int
		ok     bool
	}{
		{"tmux 3.7c", 3, 7, true},
		{"tmux 3.3", 3, 3, true},
		{"tmux 3.3a", 3, 3, true},
		{"tmux 2.9a", 2, 9, true},
		{"tmux next-3.4", 3, 4, true},
		{"openbsd-7.6", 7, 6, true},
		{"no version here", 0, 0, false},
		{"", 0, 0, false},
	}

	for _, tt := range tests {
		major, minor, ok := parseTmuxVersion(tt.output)
		if ok != tt.ok || major != tt.major || minor != tt.minor {
			t.Errorf("parseTmuxVersion(%q) = (%d, %d, %v), want (%d, %d, %v)",
				tt.output, major, minor, ok, tt.major, tt.minor, tt.ok)
		}
	}
}

func TestParseTmuxVersionPassthroughSupport(t *testing.T) {
	supportsList := []string{"tmux 3.3", "tmux 3.3a", "tmux 3.7c", "tmux 4.0"}
	unsupported := []string{"tmux 2.9a", "tmux 3.2", "tmux 3.2a"}

	supports := func(major, minor int) bool {
		return major > tmuxMinMajorVersion ||
			(major == tmuxMinMajorVersion && minor >= tmuxMinMinorVersion)
	}

	for _, out := range supportsList {
		major, minor, ok := parseTmuxVersion(out)
		if !ok || !supports(major, minor) {
			t.Errorf("%q should support passthrough, got (%d, %d, %v)", out, major, minor, ok)
		}
	}

	for _, out := range unsupported {
		major, minor, ok := parseTmuxVersion(out)
		if !ok || supports(major, minor) {
			t.Errorf("%q should NOT support passthrough, got (%d, %d, %v)", out, major, minor, ok)
		}
	}
}

func TestTmuxAllowPassthroughValue(t *testing.T) {
	tests := []struct {
		version string
		want    string
	}{
		{"tmux 3.3", "on"},
		{"tmux 3.3a", "on"},
		{"tmux 3.4", "all"},
		{"tmux 3.6", "all"},
		{"tmux 4.0", "all"},
	}

	for _, tt := range tests {
		major, minor, ok := parseTmuxVersion(tt.version)
		if !ok {
			t.Errorf("parseTmuxVersion(%q) failed", tt.version)
			continue
		}
		if got := tmuxAllowPassthroughValue(major, minor); got != tt.want {
			t.Errorf("tmuxAllowPassthroughValue(%q) = %q, want %q", tt.version, got, tt.want)
		}
	}
}

func TestTmuxSupportsInputBufferSize(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{"tmux 3.3", false},
		{"tmux 3.5a", false},
		{"tmux 3.6", true},
		{"tmux 4.0", true},
	}

	for _, tt := range tests {
		major, minor, ok := parseTmuxVersion(tt.version)
		if !ok {
			t.Errorf("parseTmuxVersion(%q) failed", tt.version)
			continue
		}
		if got := tmuxSupportsInputBufferSize(major, minor); got != tt.want {
			t.Errorf("tmuxSupportsInputBufferSize(%q) = %v, want %v", tt.version, got, tt.want)
		}
	}
}

func TestTmuxPassthroughWrap(t *testing.T) {
	wrap := func(raw string) string {
		return strings.ReplaceAll(raw[1:], "\x1b", "\x1b\x1b")
	}

	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "single APC sequence",
			raw:  "\x1b_Ga=d\x1b\\",
			want: tmuxPassthroughStart + wrap("\x1b_Ga=d\x1b\\") + tmuxPassthroughEnd,
		},
		{
			name: "chunked APC sequences with leading delete",
			raw:  "\x1b_Ga=d,d=i,i=1\x1b\\\x1b_Gf=32,s=4,v=2,m=1;q=2;AAAA\x1b\\\x1b_Gm=0;BBBB\x1b\\",
			want: tmuxPassthroughStart + wrap("\x1b_Ga=d,d=i,i=1\x1b\\") + tmuxPassthroughEnd +
				tmuxPassthroughStart + wrap("\x1b_Gf=32,s=4,v=2,m=1;q=2;AAAA\x1b\\") + tmuxPassthroughEnd +
				tmuxPassthroughStart + wrap("\x1b_Gm=0;BBBB\x1b\\") + tmuxPassthroughEnd,
		},
		{
			name: "ESC inside payload is doubled",
			raw:  "\x1b_Ga=d,d=i,i=\x1b5\x1b\\",
			want: tmuxPassthroughStart + "_Ga=d,d=i,i=\x1b\x1b5\x1b\x1b\\" + tmuxPassthroughEnd,
		},
		{
			name: "unterminated APC is passed through unchanged",
			raw:  "\x1b_Gf=32;base64",
			want: "\x1b_Gf=32;base64",
		},
		{
			name: "non APC content preserved around sequences",
			raw:  "\x1b[38;5;1m\x1b_Ga=d\x1b\\\x1b[39m",
			want: "\x1b[38;5;1m" +
				tmuxPassthroughStart + wrap("\x1b_Ga=d\x1b\\") + tmuxPassthroughEnd +
				"\x1b[39m",
		},
	}

	for _, tt := range tests {
		if got := tmuxWrap(tt.raw); got != tt.want {
			t.Errorf("%s: tmuxWrap() = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestTmuxPassthroughWrapEmpty(t *testing.T) {
	if got := tmuxWrap(""); got != "" {
		t.Errorf("tmuxWrap(\"\") = %q, want \"\"", got)
	}

	got, err := tmuxPassthroughWrap("")
	if err != nil || got != "" {
		t.Errorf("tmuxPassthroughWrap(\"\") = (%q, %v), want (\"\", nil)", got, err)
	}
}

func TestTmuxPassthroughWrapOutsideTmux(t *testing.T) {
	t.Setenv("TMUX", "")
	t.Setenv("TERM_PROGRAM", "")
	t.Setenv("TERM", "")

	raw := "\x1b_Ga=d\x1b\\"
	got, err := tmuxPassthroughWrap(raw)
	if err != nil {
		t.Fatalf("tmuxPassthroughWrap() error = %v", err)
	}
	if got != raw {
		t.Errorf("tmuxPassthroughWrap() = %q, want %q", got, raw)
	}
}
