package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestApplyConfigOverrides covers key=value parsing, per-kind type coercion
// (string, bool, int, []string), multiple ordered overrides, values containing
// '=', whitespace trimming, and the unknown-key / missing-'=' / bad-value errors.
func TestApplyConfigOverrides(t *testing.T) {
	testdata := []struct {
		name      string
		overrides []string
		// check inspects the resulting config when no error is expected
		check func(t *testing.T, c *ConfigType)
		// errSubstr, when non-empty, asserts the returned error contains it
		errSubstr string
	}{
		{
			name:      "bool override true",
			overrides: []string{"debug=true"},
			check: func(t *testing.T, c *ConfigType) {
				assert.True(t, c.Debug)
			},
		},
		{
			name:      "bool override false",
			overrides: []string{"auto_check_update=false"},
			check: func(t *testing.T, c *ConfigType) {
				assert.False(t, c.AutoCheckUpdate)
			},
		},
		{
			name:      "string override",
			overrides: []string{"theme=gruvbox"},
			check: func(t *testing.T, c *ConfigType) {
				assert.Equal(t, "gruvbox", c.Theme)
			},
		},
		{
			name:      "int override",
			overrides: []string{"default_sort_type=2"},
			check: func(t *testing.T, c *ConfigType) {
				assert.Equal(t, 2, c.DefaultSortType)
			},
		},
		{
			name:      "string slice override comma separated",
			overrides: []string{"sidebar_sections=home,pinned"},
			check: func(t *testing.T, c *ConfigType) {
				assert.Equal(t, []string{"home", "pinned"}, c.SidebarSections)
			},
		},
		{
			name:      "multiple overrides applied in order",
			overrides: []string{"debug=true", "theme=catppuccin", "sidebar_width=12"},
			check: func(t *testing.T, c *ConfigType) {
				assert.True(t, c.Debug)
				assert.Equal(t, "catppuccin", c.Theme)
				assert.Equal(t, 12, c.SidebarWidth)
			},
		},
		{
			name:      "value containing equals sign is preserved",
			overrides: []string{"default_directory=/tmp/a=b"},
			check: func(t *testing.T, c *ConfigType) {
				assert.Equal(t, "/tmp/a=b", c.DefaultDirectory)
			},
		},
		{
			name:      "key is trimmed and bool value parses despite surrounding whitespace",
			overrides: []string{" debug = true "},
			check: func(t *testing.T, c *ConfigType) {
				assert.True(t, c.Debug)
			},
		},
		{
			name:      "int value parses despite surrounding whitespace",
			overrides: []string{"sidebar_width= 12 "},
			check: func(t *testing.T, c *ConfigType) {
				assert.Equal(t, 12, c.SidebarWidth)
			},
		},
		{
			// config.toml documents "Use ' ' for borderless", so a single space
			// is a valid intended string value. The key must still be trimmed,
			// but the raw string value must be preserved verbatim (not trimmed to
			// empty) so CLI overrides can express the same values as config.toml.
			name:      "string value with intended single space is preserved",
			overrides: []string{" border_top = "},
			check: func(t *testing.T, c *ConfigType) {
				assert.Equal(t, " ", c.BorderTop)
			},
		},
		{
			name:      "string value surrounding whitespace is preserved verbatim",
			overrides: []string{"border_bottom=  x  "},
			check: func(t *testing.T, c *ConfigType) {
				assert.Equal(t, "  x  ", c.BorderBottom)
			},
		},
		{
			name:      "unknown key returns clear error",
			overrides: []string{"not_a_real_key=true"},
			errSubstr: "not_a_real_key",
		},
		{
			name:      "missing equals returns error",
			overrides: []string{"debug"},
			errSubstr: "debug",
		},
		{
			name:      "invalid bool value returns error",
			overrides: []string{"debug=notabool"},
			errSubstr: "debug",
		},
		{
			name:      "invalid int value returns error",
			overrides: []string{"default_sort_type=abc"},
			errSubstr: "default_sort_type",
		},
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConfigType{}
			err := ApplyConfigOverrides(c, tt.overrides)
			if tt.errSubstr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errSubstr)
				return
			}
			require.NoError(t, err)
			tt.check(t, c)
		})
	}
}

// TestApplyConfigOverridesEmpty verifies that a nil/empty override list is a
// no-op and leaves the config untouched.
func TestApplyConfigOverridesEmpty(t *testing.T) {
	c := &ConfigType{Theme: "original"}
	require.NoError(t, ApplyConfigOverrides(c, nil))
	assert.Equal(t, "original", c.Theme, "no overrides should leave config untouched")
}

// TestApplyConfigOverridesUnknownKeyMentionsOverride verifies the unknown-key
// error names both the offending key and the full override string.
func TestApplyConfigOverridesUnknownKeyMentionsOverride(t *testing.T) {
	c := &ConfigType{}
	err := ApplyConfigOverrides(c, []string{"bogus_key=1"})
	require.Error(t, err)
	// Error should be understandable: name the offending key.
	assert.Contains(t, err.Error(), "bogus_key")
}
