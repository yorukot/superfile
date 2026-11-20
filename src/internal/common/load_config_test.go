package common

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	variable "github.com/yorukot/superfile/src/config"
	"github.com/yorukot/superfile/src/pkg/utils"
)

type repairMessageTestConfig struct {
	Existing string `toml:"existing"`
	Added    string `toml:"added"`
}

func TestFormatHotkeysLoadError(t *testing.T) {
	originalFixHotkeys := variable.FixHotkeys
	originalLipglossError := LipglossError
	t.Cleanup(func() {
		variable.FixHotkeys = originalFixHotkeys
		LipglossError = originalLipglossError
	})

	LipglossError = "Error | "

	t.Run("successful repair is not labelled as an error", func(t *testing.T) {
		defaultData := "existing = 'default'\nadded = 'default'\n"
		hotkeysFile := filepath.Join(t.TempDir(), "hotkeys.toml")
		require.NoError(t, os.WriteFile(hotkeysFile, []byte("existing = 'custom'\n"), utils.ConfigFilePerm))

		var target repairMessageTestConfig
		err := utils.LoadTomlFile(hotkeysFile, defaultData, &target, true, false)
		require.Error(t, err)

		variable.FixHotkeys = true
		message, toExit := formatHotkeysLoadError(err)

		assert.False(t, toExit)
		assert.NotContains(t, message, LipglossError)
		assert.True(t, strings.HasPrefix(message, "config file had issues. It's fixed successfully."))
	})

	t.Run("fatal errors retain the error label", func(t *testing.T) {
		variable.FixHotkeys = true

		message, toExit := formatHotkeysLoadError(errors.New("could not load hotkeys"))

		assert.True(t, toExit)
		assert.Equal(t, LipglossError+"could not load hotkeys", message)
	})

	t.Run("unresolved missing fields retain the error label", func(t *testing.T) {
		defaultData := "existing = 'default'\nadded = 'default'\n"
		hotkeysFile := filepath.Join(t.TempDir(), "hotkeys.toml")
		require.NoError(t, os.WriteFile(hotkeysFile, []byte("existing = 'custom'\n"), utils.ConfigFilePerm))

		var target repairMessageTestConfig
		err := utils.LoadTomlFile(hotkeysFile, defaultData, &target, false, false)
		require.Error(t, err)

		variable.FixHotkeys = false
		message, toExit := formatHotkeysLoadError(err)

		assert.False(t, toExit)
		assert.True(t, strings.HasPrefix(message, LipglossError))
		assert.Contains(t, message, "--fix-hotkeys")
	})
}
