package preview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildBatArgs(t *testing.T) {
	tests := []struct {
		name         string
		batFlags     []string
		expectedArgs []string
	}{
		{
			name:     "default flags",
			batFlags: []string{"--plain", "--force-colorization"},
			expectedArgs: []string{
				"example.go",
				"--plain",
				"--force-colorization",
				"--line-range",
				":24",
			},
		},
		{
			name:     "custom flags",
			batFlags: []string{"--style=numbers", "--color=always"},
			expectedArgs: []string{
				"example.go",
				"--style=numbers",
				"--color=always",
				"--line-range",
				":24",
			},
		},
		{
			name:     "no flags",
			batFlags: nil,
			expectedArgs: []string{
				"example.go",
				"--line-range",
				":24",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedArgs, buildBatArgs("example.go", 24, test.batFlags))
		})
	}
}
