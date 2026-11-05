package common

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validBaseConfig() ConfigType {
	return ConfigType{
		Theme:                   "default",
		FilePanelNamePercent:    50,
		FilePreviewWidth:        0,
		SidebarWidth:            10,
		SidebarSections:         []string{"home", "pinned", "disks"},
		DefaultSortType:         0,
		BorderTop:               "─",
		BorderBottom:            "─",
		BorderLeft:              "│",
		BorderRight:             "│",
		BorderTopLeft:           "┌",
		BorderTopRight:          "┐",
		BorderBottomLeft:        "└",
		BorderBottomRight:       "┘",
		BorderMiddleLeft:        "├",
		BorderMiddleRight:       "┤",
		ZoxideMaxVisibleResults: DefaultZoxideMaxVisibleResults,
		ZoxideModalWidth:        DefaultZoxideModalWidthPercent,
		ZoxideModalHeight:       DefaultZoxideModalHeightPercent,
	}
}

func TestZoxideConfigParsing(t *testing.T) {
	tomlContent := `
theme = "default"
zoxide_support = true
zoxide_max_visible_results = 12
zoxide_modal_width = 75
zoxide_modal_height = 65
`
	var cfg ConfigType
	err := toml.Unmarshal([]byte(tomlContent), &cfg)
	require.NoError(t, err)

	assert.True(t, cfg.ZoxideSupport)
	assert.Equal(t, 12, cfg.ZoxideMaxVisibleResults)
	assert.Equal(t, 75, cfg.ZoxideModalWidth)
	assert.Equal(t, 65, cfg.ZoxideModalHeight)
}

func TestValidateConfigZoxideBounds(t *testing.T) {
	tests := []struct {
		name        string
		modify      func(c *ConfigType)
		expectError bool
		errorField  string
	}{
		{
			name:        "valid default zoxide settings",
			modify:      func(c *ConfigType) {},
			expectError: false,
		},
		{
			name: "valid custom zoxide settings",
			modify: func(c *ConfigType) {
				c.ZoxideMaxVisibleResults = 15
				c.ZoxideModalWidth = 70
				c.ZoxideModalHeight = 60
			},
			expectError: false,
		},
		{
			name: "zoxide max visible results zero is invalid",
			modify: func(c *ConfigType) {
				c.ZoxideMaxVisibleResults = 0
			},
			expectError: true,
			errorField:  "zoxide_max_visible_results",
		},
		{
			name: "zoxide max visible results negative is invalid",
			modify: func(c *ConfigType) {
				c.ZoxideMaxVisibleResults = -1
			},
			expectError: true,
			errorField:  "zoxide_max_visible_results",
		},
		{
			name: "zoxide max visible results above 100 is invalid",
			modify: func(c *ConfigType) {
				c.ZoxideMaxVisibleResults = 101
			},
			expectError: true,
			errorField:  "zoxide_max_visible_results",
		},
		{
			name: "zoxide modal width below 10 is invalid",
			modify: func(c *ConfigType) {
				c.ZoxideModalWidth = 9
			},
			expectError: true,
			errorField:  "zoxide_modal_width",
		},
		{
			name: "zoxide modal width above 100 is invalid",
			modify: func(c *ConfigType) {
				c.ZoxideModalWidth = 101
			},
			expectError: true,
			errorField:  "zoxide_modal_width",
		},
		{
			name: "zoxide modal height below 10 is invalid",
			modify: func(c *ConfigType) {
				c.ZoxideModalHeight = 5
			},
			expectError: true,
			errorField:  "zoxide_modal_height",
		},
		{
			name: "zoxide modal height above 100 is invalid",
			modify: func(c *ConfigType) {
				c.ZoxideModalHeight = 105
			},
			expectError: true,
			errorField:  "zoxide_modal_height",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validBaseConfig()
			tt.modify(&cfg)
			err := ValidateConfig(&cfg)
			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorField)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
