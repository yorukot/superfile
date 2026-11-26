package filepreview

import (
	"strings"
	"testing"
)

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
