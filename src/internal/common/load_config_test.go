package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateConfigImageRenderer(t *testing.T) {
	for _, value := range []string{"", "auto", "kitty", "sixel", "ansi", "Sixel"} {
		c := ConfigType{ImageRenderer: value}
		err := ValidateConfig(&c)
		// Other fields of the zero config may be invalid; only image_renderer matters here.
		if err != nil {
			assert.NotContains(t, err.Error(), "image_renderer", "value %q should be valid", value)
		}
	}

	c := ConfigType{ImageRenderer: "iterm"}
	err := ValidateConfig(&c)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "image_renderer")
}
