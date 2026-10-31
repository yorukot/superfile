package disks

import (
	"testing"

	"github.com/google/glazier/go/storage"
	"github.com/stretchr/testify/assert"
)

func TestGetVolumeName(t *testing.T) {
	tests := []struct {
		name   string
		volume *storage.Volume
		want   string
	}{
		{name: "FileSystemLabel", volume: &storage.Volume{
			FileSystemLabel: "movies",
		}, want: "movies"},
		{name: "FileSystemLabel and DriveLetter", volume: &storage.Volume{
			DriveLetter:     "D",
			FileSystemLabel: "movies",
		}, want: "movies (D)"},
		{name: "FileSystemLabel and invalid DriveLetter", volume: &storage.Volume{
			DriveLetter:     "\x00",
			Path:            `D:\`,
			FileSystemLabel: "movies",
		}, want: "movies"},
		{name: "DriveLetter only", volume: &storage.Volume{
			DriveLetter:     "D",
			Path:            `D:\`,
			FileSystemLabel: "",
		}, want: "D:"},
		{name: "invalid DriveLetter, empty FileSystemLabel", volume: &storage.Volume{
			DriveLetter:     "\x00",
			Path:            `D:\`,
			FileSystemLabel: "",
		}, want: `D:\`},
		{name: "invalid DriveLetter, empty FileSystemLabel, GUID path", volume: &storage.Volume{
			DriveLetter:     "\x00",
			Path:            `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`,
			FileSystemLabel: "",
		}, want: `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`},
		{name: "FileSystemLabel and all other fields valid", volume: &storage.Volume{
			DriveLetter:     "F",
			Path:            `F:\`,
			FileSystemLabel: "sandisk",
		}, want: "sandisk (F)"},
		{name: "no usable volume identifiers", volume: &storage.Volume{
			DriveLetter:     "\x00",
			Path:            "",
			FileSystemLabel: "",
		}, want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getVolumeName(tc.volume)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestGetVolumeMountPoint(t *testing.T) {
	tests := []struct {
		name   string
		volume *storage.Volume
		want   string
	}{
		{name: "DriveLetter and GUID Path", volume: &storage.Volume{
			DriveLetter: "D",
			Path:        `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`,
		}, want: `D:\`},
		{name: "invalid DriveLetter, GUID path", volume: &storage.Volume{
			DriveLetter: "\x00",
			Path:        `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`,
		}, want: `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`},
		{name: "empty DriveLetter, GUID path", volume: &storage.Volume{
			DriveLetter: "",
			Path:        `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`,
		}, want: `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`},
		{name: "DriveLetter based path, empty DriveLetter", volume: &storage.Volume{
			DriveLetter: "",
			Path:        `D:\`,
		}, want: `D:\`},
		{name: "empty DriveLetter and Path", volume: &storage.Volume{
			DriveLetter: "",
			Path:        "",
		}, want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := getVolumeMountPoint(tc.volume)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestShouldHideVolume(t *testing.T) {
	tests := []struct {
		name       string
		partitions []storage.Partition
		volume     storage.Volume
		want       bool
	}{
		{name: "all fields pass",
			partitions: []storage.Partition{
				{AccessPaths: []string{
					`D:\`,
					`\\?\Volume{12345678-1234-1234-1234-123456789abc}\`},
				},
			},
			volume: storage.Volume{
				Path: `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`,
				Size: 2 * (1 << 40),
			},
			want: false},
		{name: "volume size is 0",
			partitions: []storage.Partition{
				{AccessPaths: []string{
					`D:\`,
					`\\?\Volume{12345678-1234-1234-1234-123456789abc}\`},
				},
			},
			volume: storage.Volume{
				Path: `D:\`,
				Size: 0,
			},
			want: true},
		{name: "AccessPaths only has GUID path",
			partitions: []storage.Partition{
				{AccessPaths: []string{
					`\\?\Volume{12345678-1234-1234-1234-123456789abc}\`},
				},
			},
			volume: storage.Volume{
				Path: `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`,
				Size: 16 * (1 << 30),
			},
			want: false},
		{name: "AccessPaths only has DriveLetter based path",
			partitions: []storage.Partition{
				{AccessPaths: []string{
					`E:\`},
				},
			},
			volume: storage.Volume{
				Path: `E:\`,
				Size: 16 * (1 << 30),
			},
			want: false},
		{name: "nil AccessPaths",
			partitions: []storage.Partition{
				{AccessPaths: []string{}},
			},
			volume: storage.Volume{
				Path: `\\?\Volume{12345678-1234-1234-1234-123456789abc}\`,
				Size: 500 * (1 << 30),
			},
			want: false},
		{name: "Path in AccessPaths, partition is hidden",
			partitions: []storage.Partition{
				{AccessPaths: []string{
					`D:\`,
					`\\?\Volume{12345678-1234-1234-1234-123456789abc}\`},
					IsHidden: true,
				},
			},
			volume: storage.Volume{
				Path: `D:\`,
				Size: 500 * (1 << 30),
			},
			want: true},
		{name: "Path in AccessPaths, partition is not hidden",
			partitions: []storage.Partition{
				{AccessPaths: []string{
					`D:\`,
					`\\?\Volume{12345678-1234-1234-1234-123456789abc}\`},
				},
			},
			volume: storage.Volume{
				Path: `D:\`,
				Size: 500 * (1 << 30),
			},
			want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldHideVolume(tc.partitions, tc.volume)
			assert.Equal(t, tc.want, got)
		})
	}
}
