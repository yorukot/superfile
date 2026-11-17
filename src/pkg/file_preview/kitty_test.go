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
