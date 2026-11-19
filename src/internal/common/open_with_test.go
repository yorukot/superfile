package common

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveOpenWith(t *testing.T) {
	t.Run("legacy command", func(t *testing.T) {
		applications, err := ResolveOpenWith(map[string]any{"pdf": "zathura"}, "/tmp/report.PDF")

		require.NoError(t, err)
		require.Len(t, applications, 1)
		assert.Equal(t, OpenWithApplication{Name: "zathura", Command: "zathura"}, applications[0])
	})

	t.Run("multiple structured applications", func(t *testing.T) {
		var config ConfigType
		require.NoError(t, toml.Unmarshal([]byte(`
[open_with]
md = [
  { name = "Neovim", command = "nvim", args = ["--clean"], terminal = true },
  { name = "Code", command = "code" }
]
`), &config))

		applications, err := ResolveOpenWith(config.OpenWith, "/tmp/read me.md")

		require.NoError(t, err)
		assert.Equal(t, []OpenWithApplication{
			{Name: "Neovim", Command: "nvim", Args: []string{"--clean"}, Terminal: true},
			{Name: "Code", Command: "code"},
		}, applications)
	})

	t.Run("no configured extension", func(t *testing.T) {
		applications, err := ResolveOpenWith(map[string]any{"pdf": "zathura"}, "/tmp/report.txt")

		require.NoError(t, err)
		assert.Empty(t, applications)
	})

	t.Run("invalid application", func(t *testing.T) {
		_, err := ResolveOpenWith(map[string]any{
			"pdf": []any{map[string]any{"name": "Missing command"}},
		}, "/tmp/report.pdf")

		assert.ErrorIs(t, err, ErrInvalidOpenWithApplication)
	})
}

func TestValidateOpenWithRejectsInvalidArguments(t *testing.T) {
	err := ValidateOpenWith(map[string]any{
		"txt": map[string]any{"command": "less", "args": []any{"-R", 1}},
	})

	assert.ErrorIs(t, err, ErrInvalidOpenWithApplication)
}
