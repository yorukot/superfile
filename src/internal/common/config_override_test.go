package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyConfigOverrides(t *testing.T) {
	c := ConfigType{
		Theme:           "original",
		Editor:          "vim",
		AutoCheckUpdate: true,
		SidebarWidth:    8,
		SidebarSections: []string{"disks"},
	}
	overrides := []string{
		"debug=true",
		"auto_check_update=false",
		`theme="gruvbox"`,
		`theme="catppuccin-mocha"`,
		"sidebar_width=12",
		`sidebar_sections=["home", "pinned"]`,
		`default_directory="/tmp/a=b,c"`,
		`border_top=" "`,
		`open_with={txt="nvim"}`,
	}

	require.NoError(t, ApplyConfigOverrides(&c, overrides))
	assert.Equal(t, ConfigType{
		Theme:            "catppuccin-mocha",
		Editor:           "vim",
		Debug:            true,
		SidebarWidth:     12,
		SidebarSections:  []string{"home", "pinned"},
		DefaultDirectory: "/tmp/a=b,c",
		BorderTop:        " ",
		OpenWith:         map[string]string{"txt": "nvim"},
	}, c)
}

func TestApplyConfigOverridesEmpty(t *testing.T) {
	c := ConfigType{Theme: "original", SidebarSections: []string{"home"}}
	want := c
	require.NoError(t, ApplyConfigOverrides(&c, nil))
	assert.Equal(t, want, c)
	require.NoError(t, ApplyConfigOverrides(&c, []string{}))
	assert.Equal(t, want, c)
}

func TestApplyConfigOverridesErrors(t *testing.T) {
	for _, override := range []string{
		"bogus_key=true",
		"debug",
		`debug="true"`,
		`sidebar_width="wide"`,
		`sidebar_sections=["home"`,
	} {
		t.Run(override, func(t *testing.T) {
			var c ConfigType
			err := ApplyConfigOverrides(&c, []string{override})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid config override")
		})
	}
}
