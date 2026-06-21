package common

import (
	"errors"
	"image/color"
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

func TestShouldWarnAutoDetectFailed(t *testing.T) {
	someColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	tests := []struct {
		name     string
		theme    string
		bg       color.Color
		err      error
		expected bool
	}{
		{
			name:     "Non-auto theme never warns, even if detection failed",
			theme:    "catppuccin-mocha",
			bg:       nil,
			err:      errors.New("input/output is not a terminal"),
			expected: false,
		},
		{
			name:     "Auto theme with detection error warns",
			theme:    "auto",
			bg:       nil,
			err:      errors.New("input/output is not a terminal"),
			expected: true,
		},
		{
			name:     "Auto theme with nil color and no error warns (inconclusive query)",
			theme:    "auto",
			bg:       nil,
			err:      nil,
			expected: true,
		},
		{
			name:     "Auto theme with a resolved color does not warn",
			theme:    "auto",
			bg:       someColor,
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldWarnAutoDetectFailed(tt.theme, tt.bg, tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func validConfigType() ConfigType {
	return ConfigType{
		Theme:                "catppuccin-mocha",
		DefaultSortType:      0,
		FilePanelNamePercent: 50,
		BorderTop:            "─",
		BorderBottom:         "─",
		BorderLeft:           "│",
		BorderRight:          "│",
		BorderTopLeft:        "╭",
		BorderTopRight:       "╮",
		BorderBottomLeft:     "╰",
		BorderBottomRight:    "╯",
		BorderMiddleLeft:     "├",
		BorderMiddleRight:    "┤",
	}
}

func TestValidateConfig_AutoTheme(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(c *ConfigType)
		expectErr bool
	}{
		{
			name:      "Non-auto theme is unaffected by empty theme_light/theme_dark",
			mutate:    func(_ *ConfigType) {},
			expectErr: false,
		},
		{
			name: "Auto theme with both theme_light and theme_dark set passes",
			mutate: func(c *ConfigType) {
				c.Theme = "auto"
				c.ThemeLight = "catppuccin-latte"
				c.ThemeDark = "catppuccin-mocha"
			},
			expectErr: false,
		},
		{
			name: "Auto theme with theme_light unset fails",
			mutate: func(c *ConfigType) {
				c.Theme = "auto"
				c.ThemeDark = "catppuccin-mocha"
			},
			expectErr: true,
		},
		{
			name: "Auto theme with theme_dark unset fails",
			mutate: func(c *ConfigType) {
				c.Theme = "auto"
				c.ThemeLight = "catppuccin-latte"
			},
			expectErr: true,
		},
		{
			name: "Auto theme with both theme_light and theme_dark unset fails",
			mutate: func(c *ConfigType) {
				c.Theme = "auto"
			},
			expectErr: true,
		},
		{
			name: "Auto theme with theme_light set to auto fails",
			mutate: func(c *ConfigType) {
				c.Theme = "auto"
				c.ThemeLight = "auto"
				c.ThemeDark = "catppuccin-mocha"
			},
			expectErr: true,
		},
		{
			name: "Auto theme with theme_dark set to auto fails",
			mutate: func(c *ConfigType) {
				c.Theme = "auto"
				c.ThemeLight = "catppuccin-latte"
				c.ThemeDark = "auto"
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfigType()
			tt.mutate(&cfg)
			err := ValidateConfig(&cfg)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
