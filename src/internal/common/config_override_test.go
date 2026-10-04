package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyConfigOverrides(t *testing.T) {
	t.Run("scalar overrides", func(t *testing.T) {
		cfg := ConfigType{Theme: "default", Debug: false, SidebarWidth: 10}

		err := ApplyConfigOverrides(&cfg, []string{"debug=true", "theme=catppuccin", "sidebar_width=15"})
		require.NoError(t, err)

		assert.True(t, cfg.Debug)
		assert.Equal(t, "catppuccin", cfg.Theme)
		assert.Equal(t, 15, cfg.SidebarWidth)
	})

	t.Run("keys are case insensitive", func(t *testing.T) {
		cfg := ConfigType{}

		err := ApplyConfigOverrides(&cfg, []string{"Debug=TRUE"})
		require.NoError(t, err)

		assert.True(t, cfg.Debug)
	})

	t.Run("value may contain an equals sign", func(t *testing.T) {
		cfg := ConfigType{}

		err := ApplyConfigOverrides(&cfg, []string{"theme=my=theme"})
		require.NoError(t, err)

		assert.Equal(t, "my=theme", cfg.Theme)
	})

	t.Run("unknown key is an error", func(t *testing.T) {
		cfg := ConfigType{}

		err := ApplyConfigOverrides(&cfg, []string{"not_a_config_key=true"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not_a_config_key")
	})

	t.Run("composite keys are not supported", func(t *testing.T) {
		cfg := ConfigType{}

		err := ApplyConfigOverrides(&cfg, []string{"open_with=go"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not supported")

		err = ApplyConfigOverrides(&cfg, []string{"sidebar_sections=home"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not supported")
	})

	t.Run("invalid bool value is an error", func(t *testing.T) {
		cfg := ConfigType{}

		err := ApplyConfigOverrides(&cfg, []string{"debug=maybe"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "boolean")
	})

	t.Run("invalid int value is an error", func(t *testing.T) {
		cfg := ConfigType{}

		err := ApplyConfigOverrides(&cfg, []string{"sidebar_width=ten"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "integer")
	})

	t.Run("missing equals sign is an error", func(t *testing.T) {
		cfg := ConfigType{}

		err := ApplyConfigOverrides(&cfg, []string{"debug"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "KEY=VALUE")
	})

	t.Run("no overrides is a no-op", func(t *testing.T) {
		cfg := ConfigType{Theme: "default"}

		err := ApplyConfigOverrides(&cfg, nil)
		require.NoError(t, err)
		assert.Equal(t, "default", cfg.Theme)
	})

	t.Run("overridden values go through config validation", func(t *testing.T) {
		cfg := ConfigType{}

		// SidebarWidth must be 5-20, or 0. Applying an out of range
		// override and then validating must fail, which is why overrides
		// are applied before ValidateConfig in LoadConfigFile().
		err := ApplyConfigOverrides(&cfg, []string{"sidebar_width=3"})
		require.NoError(t, err)

		assert.Error(t, ValidateConfig(&cfg))
	})
}
