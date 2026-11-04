package filepreview

import (
	"strings"
	"testing"
)

func TestKittyCapabilityResolution(t *testing.T) {
	reset := func() { kittyCapability.Store(kittyCapUnknown) }
	p := &ImagePreviewer{}

	// A graphics reply marks support, and only reports a change the first time.
	reset()
	if !MarkKittySupported() {
		t.Fatal("first graphics reply should report a change")
	}
	if MarkKittySupported() {
		t.Fatal("repeated graphics reply should not report a change")
	}
	if !p.IsKittyCapable() {
		t.Fatal("expected capable after graphics reply")
	}

	// DA1 must not downgrade a terminal that already answered positively.
	if MarkKittyUnsupportedIfUnknown() {
		t.Fatal("DA1 must not override a positive reply")
	}
	if !p.IsKittyCapable() {
		t.Fatal("capability was downgraded by DA1")
	}

	// DA1 alone (the fence arriving with no graphics reply) means unsupported,
	// and must win over the allowlist even on a terminal that is on it.
	reset()
	t.Setenv("TERM", "xterm-kitty")
	if !MarkKittyUnsupportedIfUnknown() {
		t.Fatal("DA1 from unknown state should report a change")
	}
	if p.IsKittyCapable() {
		t.Fatal("terminal answered DA1 only; must not be treated as capable")
	}

	// While unknown, fall back to the allowlist.
	reset()
	if !p.IsKittyCapable() {
		t.Fatal("expected allowlist fallback to report capable for xterm-kitty")
	}
	t.Setenv("TERM", "dumb")
	t.Setenv("TERM_PROGRAM", "dumb")
	if p.IsKittyCapable() {
		t.Fatal("expected allowlist fallback to report incapable for dumb")
	}
	reset()
}

func TestKittyGraphicsQuery(t *testing.T) {
	t.Setenv("TMUX", "")
	q := KittyGraphicsQuery()
	if !strings.Contains(q, "a=q") {
		t.Fatalf("query must ask for a report (a=q), got %q", q)
	}
	if !strings.HasSuffix(q, "\x1b[c") {
		t.Fatalf("DA1 fence must come last so it arrives after any graphics reply, got %q", q)
	}

	// Inside tmux the fence must be wrapped alongside the query, so the host
	// terminal answers both rather than tmux answering DA1 itself.
	t.Setenv("TMUX", "/tmp/tmux-501/default,123,0")
	wrapped := KittyGraphicsQuery()
	if strings.HasSuffix(wrapped, "\x1b[c") {
		t.Fatal("DA1 must be inside the tmux passthrough, not trailing it")
	}
	if !strings.HasPrefix(wrapped, "\x1bPtmux;") {
		t.Fatalf("query must be wrapped for tmux, got %q", wrapped)
	}
}

func TestTmuxPassthrough(t *testing.T) {
	t.Setenv("TMUX", "")
	if got := tmuxPassthrough("\x1b_Ga=d\x1b\\"); got != "\x1b_Ga=d\x1b\\" {
		t.Fatalf("outside tmux input must be unchanged, got %q", got)
	}

	t.Setenv("TMUX", "/tmp/tmux-501/default,123,0")

	if got := tmuxPassthrough(""); got != "" {
		t.Fatalf("empty input must stay empty, got %q", got)
	}

	// Single sequence: wrapped in DCS, inner ESCs doubled.
	got := tmuxPassthrough("\x1b_Ga=d\x1b\\")
	want := "\x1bPtmux;\x1b\x1b_Ga=d\x1b\x1b\\\x1b\\"
	if got != want {
		t.Fatalf("single sequence\n got %q\nwant %q", got, want)
	}

	// Chunked data must produce one DCS per sequence, not one giant DCS.
	two := tmuxPassthrough("\x1b_Gm=1;AAA\x1b\\\x1b_Gm=0;BBB\x1b\\")
	if n := strings.Count(two, "\x1bPtmux;"); n != 2 {
		t.Fatalf("expected 2 passthrough wrappers, got %d in %q", n, two)
	}

	// A trailing fragment with no ST must not be dropped.
	if frag := tmuxPassthrough("\x1b_Gnoterm"); !strings.Contains(frag, "_Gnoterm") {
		t.Fatalf("unterminated fragment was dropped: %q", frag)
	}
}
