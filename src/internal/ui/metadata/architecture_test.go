package metadata

import (
	"debug/elf"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBinaryArchitecture_NonBinaryFile(t *testing.T) {
	tmpFile := t.TempDir() + "/test.txt"

	err := os.WriteFile(tmpFile, []byte("This is not a binary file"), 0o644)
	require.NoError(t, err)

	arch, err := GetBinaryArchitecture(tmpFile)
	require.Error(t, err)
	assert.Empty(t, arch)
}

func TestGetBinaryArchitecture_NonExistentFile(t *testing.T) {
	arch, err := GetBinaryArchitecture("/nonexistent/file/path")
	require.Error(t, err)
	assert.Empty(t, arch)
}

func TestGetBinaryArchitecture_CurrentBinary(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Skip("Could not get current executable path")
	}

	arch, err := GetBinaryArchitecture(executable)
	require.NoError(t, err)
	assert.NotEmpty(t, arch)

	hasValidPrefix := strings.HasPrefix(arch, "ELF") ||
		strings.HasPrefix(arch, "PE") ||
		strings.HasPrefix(arch, "Mach-O")
	assert.True(t, hasValidPrefix,
		"Architecture should start with a known format prefix, got: %s", arch)
}

func TestElfMachineToString(t *testing.T) {
	tests := []struct {
		name     string
		input    uint16
		expected string
	}{
		{"x86-64", 0x3E, archX8664},
		{"i386", 0x03, archI386},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, elfMachineToString(elf.Machine(tt.input)))
		})
	}
}

func TestPeArchitectureToString(t *testing.T) {
	assert.Equal(t, archI386, peArchitectureToString(0x14c))
	assert.Equal(t, archX8664, peArchitectureToString(0x8664))
	assert.Equal(t, archARM, peArchitectureToString(0x1c0))
	assert.Equal(t, archARM64, peArchitectureToString(0xaa64))
	assert.Contains(t, peArchitectureToString(0x9999), "Unknown")
}
