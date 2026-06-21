package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveThemeName(t *testing.T) {
	tests := []struct {
		name       string
		theme      string
		themeLight string
		themeDark  string
		hasDarkBG  bool
		expected   string
	}{
		{
			name:       "Non-auto theme is returned unchanged regardless of background",
			theme:      "catppuccin-mocha",
			themeLight: "catppuccin-latte",
			themeDark:  "catppuccin-mocha",
			hasDarkBG:  false,
			expected:   "catppuccin-mocha",
		},
		{
			name:       "Auto mode with dark background picks theme_dark",
			theme:      "auto",
			themeLight: "catppuccin-latte",
			themeDark:  "catppuccin-mocha",
			hasDarkBG:  true,
			expected:   "catppuccin-mocha",
		},
		{
			name:       "Auto mode with light background picks theme_light",
			theme:      "auto",
			themeLight: "catppuccin-latte",
			themeDark:  "catppuccin-mocha",
			hasDarkBG:  false,
			expected:   "catppuccin-latte",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ResolveThemeName(tt.theme, tt.themeLight, tt.themeDark, tt.hasDarkBG)
			assert.Equal(t, tt.expected, result)
		})
	}
}
