package disks

import (
	"testing"

	"github.com/godbus/dbus/v5"

	"github.com/stretchr/testify/assert"
)

func TestDiskFromObject(t *testing.T) {
	tests := []struct {
		name     string
		path     dbus.ObjectPath
		ifaces   map[string]map[string]dbus.Variant
		wantDisk Disk
		wantBool bool
	}{
		{
			name: "no Block interface",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sda"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Filesystem": {
					"MountPoints": dbus.MakeVariant([][]byte{}),
				},
			},
			wantDisk: Disk{},
			wantBool: false,
		},
		{
			name: "no Filesystem interface",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/nvme0n1p4"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Block": {
					"Size":       dbus.MakeVariant(uint64(100_000_000_000)),
					"IdLabel":    dbus.MakeVariant("sandisk"),
					"HintSystem": dbus.MakeVariant(false),
					"HintIgnore": dbus.MakeVariant(false),
				},
			},
			wantDisk: Disk{},
			wantBool: false,
		},
		{
			name: "HintIgnore true",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sda"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Block": {
					"Size":       dbus.MakeVariant(uint64(100_000_000_000)),
					"IdLabel":    dbus.MakeVariant("sandisk"),
					"HintSystem": dbus.MakeVariant(true),
					"HintIgnore": dbus.MakeVariant(true),
				},
				"org.freedesktop.UDisks2.Filesystem": {
					"MountPoints": dbus.MakeVariant([][]byte{
						[]byte("/run/media/user/sandisk"),
					}),
				},
			},
			wantDisk: Disk{},
			wantBool: false,
		},
		{
			name: "Required property missing - HintSystem",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sda1"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Block": {
					"Size":       dbus.MakeVariant(uint64(100_000_000_000)),
					"IdLabel":    dbus.MakeVariant("sandisk"),
					"HintIgnore": dbus.MakeVariant(false),
				},
				"org.freedesktop.UDisks2.Filesystem": {
					"MountPoints": dbus.MakeVariant([][]byte{
						[]byte("/run/media/user/sandisk"),
					}),
				},
			},
			wantDisk: Disk{},
			wantBool: false,
		},
		{
			name: "Wrong type",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sdb"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Block": {
					"Size":       dbus.MakeVariant(uint64(100_000_000_000)),
					"IdLabel":    dbus.MakeVariant("sandisk"),
					"HintSystem": dbus.MakeVariant("indeed I am true"),
					"HintIgnore": dbus.MakeVariant(false),
				},
				"org.freedesktop.UDisks2.Filesystem": {
					"MountPoints": dbus.MakeVariant([][]byte{
						[]byte("/run/media/user/movies-backup"),
					}),
				},
			},
			wantDisk: Disk{},
			wantBool: false,
		},
		{
			name: "unmounted external disk",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sdb1"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Block": {
					"Size":       dbus.MakeVariant(uint64(100_000_000_000)),
					"IdLabel":    dbus.MakeVariant("work"),
					"HintSystem": dbus.MakeVariant(false),
					"HintIgnore": dbus.MakeVariant(false),
				},
				"org.freedesktop.UDisks2.Filesystem": {
					"MountPoints": dbus.MakeVariant([][]byte{}),
				},
			},
			wantDisk: Disk{
				Name:       "work",
				MountPoint: "",
				Mounted:    false,
				Internal:   false,
				ObjectPath: "/org/freedesktop/UDisks2/block_devices/sdb1",
			},
			wantBool: true,
		},
		{
			name: "mounted external disk",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/sdb1"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Block": {
					"Size":       dbus.MakeVariant(uint64(100_000_000_000)),
					"IdLabel":    dbus.MakeVariant("work"),
					"HintSystem": dbus.MakeVariant(false),
					"HintIgnore": dbus.MakeVariant(false),
				},
				"org.freedesktop.UDisks2.Filesystem": {
					"MountPoints": dbus.MakeVariant([][]byte{
						[]byte("/run/media/user/work\x00"), // the trailing NUL byte is optional
					}),
				},
			},
			wantDisk: Disk{
				Name:       "work",
				MountPoint: "/run/media/user/work",
				Mounted:    true,
				Internal:   false,
				ObjectPath: "/org/freedesktop/UDisks2/block_devices/sdb1",
			},
			wantBool: true,
		},
		{
			name: "unmounted internal disk",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/nvme0n1"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Block": {
					"Size":       dbus.MakeVariant(uint64(1000_000_000_000)),
					"IdLabel":    dbus.MakeVariant(""),
					"HintSystem": dbus.MakeVariant(true),
					"HintIgnore": dbus.MakeVariant(false),
				},
				"org.freedesktop.UDisks2.Filesystem": {
					"MountPoints": dbus.MakeVariant([][]byte{}),
				},
			},
			wantDisk: Disk{
				Name:       "1.0 TB Volume",
				MountPoint: "",
				Mounted:    false,
				Internal:   true,
				ObjectPath: "/org/freedesktop/UDisks2/block_devices/nvme0n1",
			},
			wantBool: true,
		},
		{
			name: "mounted internal disk",
			path: dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/nvme0n1"),
			ifaces: map[string]map[string]dbus.Variant{
				"org.freedesktop.UDisks2.Block": {
					"Size":       dbus.MakeVariant(uint64(2000_000_000_000)),
					"IdLabel":    dbus.MakeVariant("work"),
					"HintSystem": dbus.MakeVariant(true),
					"HintIgnore": dbus.MakeVariant(true),
				},
				"org.freedesktop.UDisks2.Filesystem": {
					"MountPoints": dbus.MakeVariant([][]byte{
						[]byte("/"),
						[]byte("/home"),
						[]byte("/var/log"),
					}),
				},
			},
			wantDisk: Disk{},
			wantBool: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotDisk, gotBool := diskFromObject(tc.path, tc.ifaces)
			assert.Equal(t, tc.wantDisk, gotDisk, "disks mismatch")
			assert.Equal(t, tc.wantBool, gotBool)
		})
	}
}

func TestDiskName(t *testing.T) {
	tests := []struct {
		name    string
		idLabel string
		size    uint64
		want    string
	}{
		{name: "label wins over size", idLabel: "movies backup", size: 16_000_000_000, want: "movies backup"},
		{name: "label wins even when size is zero", idLabel: "ventoy", size: 0, want: "ventoy"},

		{name: "no label, zero size", idLabel: "", size: 0, want: "0 MB Volume"},
		{name: "no label, MB range", idLabel: "", size: 500_000_000, want: "500 MB Volume"},
		{name: "no label, just below 1 GB", idLabel: "", size: 999_999_999, want: "999 MB Volume"},

		{name: "no label, exactly 1 GB", idLabel: "", size: 1_000_000_000, want: "1 GB Volume"},
		{name: "no label, GB range", idLabel: "", size: 10_000_000_000, want: "10 GB Volume"},
		{name: "no label, GB truncates", idLabel: "", size: 339_900_000_000, want: "339 GB Volume"},
		{name: "no label, just below 1 TB", idLabel: "", size: 999_999_999_999, want: "999 GB Volume"},

		{name: "no label, exactly 1 TB", idLabel: "", size: 1_000_000_000_000, want: "1.0 TB Volume"},
		{name: "no label, TB with decimal", idLabel: "", size: 1_900_000_000_000, want: "1.9 TB Volume"},
		{name: "no label, TB rounds up", idLabel: "", size: 1_960_000_000_000, want: "2.0 TB Volume"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := diskName(tc.idLabel, tc.size)
			assert.Equal(t, tc.want, got, "wrong name")
		})
	}
}

func TestDiskMountPoint(t *testing.T) {
	tests := []struct {
		name        string
		mountPoints [][]byte
		want        string
	}{
		{name: "empty slice, empty string", mountPoints: [][]byte{}, want: ""},
		{name: "nil slice, empty string", mountPoints: nil, want: ""},
		{
			name:        "single entry with trailing NUL byte",
			mountPoints: [][]byte{[]byte("/run/media/user/movies\x00")},
			want:        "/run/media/user/movies",
		},
		{
			name:        "no trailing NUL byte",
			mountPoints: [][]byte{[]byte("/run/media/user/movies")},
			want:        "/run/media/user/movies",
		},
		{
			name: "multiple entries returns only the first",
			mountPoints: [][]byte{
				[]byte("/run/media/user/movies\x00"),
				[]byte("/mnt/other\x00"),
			},
			want: "/run/media/user/movies",
		},
		{
			name:        "entry of only a NUL byte, empty string",
			mountPoints: [][]byte{[]byte("\x00")},
			want:        "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := diskMountPoint(tc.mountPoints)
			assert.Equal(t, tc.want, got, "invalid mount point")
		})
	}
}
