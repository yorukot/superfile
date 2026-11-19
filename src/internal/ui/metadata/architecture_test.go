package metadata

import (
	"bytes"
	"debug/elf"
	"os"
	"path/filepath"
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

// Regression test for #1550: files that are not ELF/PE/Mach-O must be rejected
// by the cheap magic-byte gate without invoking the debug/* parsers.
// Previously, pe.Open accepted MZ-less files as raw COFF objects and could
// read and allocate memory proportional to the file size (multi-GB videos)
// before failing.
func TestGetBinaryArchitecture_MediaLikeFileIsGated(t *testing.T) {
	dir := t.TempDir()

	// An MP4-like header: size-prefixed "ftyp" box, followed by bytes that
	// would decode as huge COFF section/symbol counts if misparsed.
	path := filepath.Join(dir, "movie.mp4")
	content := append([]byte{0x00, 0x00, 0x00, 0x20, 'f', 't', 'y', 'p', 'i', 's', 'o', 'm'},
		bytes.Repeat([]byte{0xff}, 4096)...)
	require.NoError(t, os.WriteFile(path, content, 0o644))

	_, err := GetBinaryArchitecture(path)
	require.ErrorIs(t, err, errNotBinary)

	// The format detector itself must classify it as unknown.
	require.Equal(t, formatUnknown, detectBinaryFormat(path))
}

func TestDetectBinaryFormat(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, data []byte) string {
		p := filepath.Join(dir, name)
		require.NoError(t, os.WriteFile(p, data, 0o644))
		return p
	}

	testCases := []struct {
		name     string
		data     []byte
		expected binaryFormat
	}{
		{"elf", []byte{0x7f, 'E', 'L', 'F', 2, 1, 1, 0}, formatELF},
		{"pe_mz", []byte{'M', 'Z', 0x90, 0x00}, formatPE},
		{"coff_amd64", []byte{0x64, 0x86, 0x01, 0x00}, formatPE},
		{"macho_64", []byte{0xfe, 0xed, 0xfa, 0xcf}, formatMachO},
		{"macho_fat", []byte{0xca, 0xfe, 0xba, 0xbe}, formatMachO},
		{"macho_64_swapped", []byte{0xcf, 0xfa, 0xed, 0xfe}, formatMachO},
		{"mp4", []byte{0x00, 0x00, 0x00, 0x20, 'f', 't', 'y', 'p'}, formatUnknown},
		{"text", []byte("hello world"), formatUnknown},
		{"tiny", []byte{0x7f}, formatUnknown},
		{"empty", nil, formatUnknown},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, detectBinaryFormat(write(tt.name, tt.data)))
		})
	}
}
