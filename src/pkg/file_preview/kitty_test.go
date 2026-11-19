package filepreview

import "testing"

func TestNormalizeImagePreviewRenderer(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", rendererAuto},
		{"auto", rendererAuto},
		{"kitty", rendererKitty},
		{"ansi", rendererANSI},
		{"KITTY", rendererKitty},
		{"  ansi  ", rendererANSI},
		{"invalid", rendererAuto},
	}

	for _, tt := range tests {
		if got := normalizeImagePreviewRenderer(tt.input); got != tt.want {
			t.Errorf("normalizeImagePreviewRenderer(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// clearKittyEnv unsets every env var consulted by detectKittySupport so each
// test case starts from a clean slate (tests may run in any order and the
// real environment may be inside a terminal).
func clearKittyEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"TERM", "TERM_PROGRAM", "TMUX", "ORIG_TERM_PROGRAM", "LC_TERMINAL",
	} {
		t.Setenv(key, "")
	}
	for _, key := range kittyEnvHints {
		t.Setenv(key, "")
	}
}

func TestIsKittyCapable(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		want bool
	}{
		{
			name: "kitty alone with TERM overridden by user config",
			env:  map[string]string{"KITTY_WINDOW_ID": "1", "TERM": "xterm-256color"},
			want: true,
		},
		{
			name: "kitty alone with default TERM",
			env:  map[string]string{"TERM": "xterm-kitty"},
			want: true,
		},
		{
			name: "no evidence of a kitty-capable terminal",
			env:  map[string]string{"TERM": "xterm-256color"},
			want: false,
		},
		{
			name: "ghostty alone via env hint",
			env:  map[string]string{"GHOSTTY_RESOURCES_DIR": "/opt/ghostty", "TERM": "xterm-256color"},
			want: true,
		},
		{
			name: "TERM_PROGRAM from allowlist",
			env:  map[string]string{"TERM_PROGRAM": "Ghostty", "TERM": "xterm-256color"},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearKittyEnv(t)
			for key, val := range tt.env {
				t.Setenv(key, val)
			}
			if got := detectKittySupport(); got != tt.want {
				t.Errorf("detectKittySupport() = %v, want %v (env: %v)", got, tt.want, tt.env)
			}
		})
	}
}
