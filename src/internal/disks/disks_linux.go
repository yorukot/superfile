package disks

import (
	"cmp"
	"context"
	"slices"
	"strings"

	"github.com/godbus/dbus/v5"
)

type manager struct {
	conn *dbus.Conn
}

func New() (Manager, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return &manager{conn: conn}, nil
}

func (m *manager) List(ctx context.Context) ([]Disk, error) {
	var objs map[dbus.ObjectPath]map[string]map[string]dbus.Variant

	obj := m.conn.Object("org.freedesktop.UDisks2", "/org/freedesktop/UDisks2")
	err := obj.CallWithContext(
		ctx, "org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0,
	).Store(&objs)
	if err != nil {
		return nil, err
	}

	var disks []Disk
	for path, ifaces := range objs {
		if disk, ok := linuxDiskFromObject(path, ifaces); ok {
			disks = append(disks, disk)
		}
	}

	// Sort external disks before internal ones, then by object path
	slices.SortFunc(disks, func(a, b Disk) int {
		return cmp.Or(
			cmp.Compare(boolToInt(a.Internal), boolToInt(b.Internal)),
			strings.Compare(a.ObjectPath, b.ObjectPath),
		)
	})

	return disks, nil
}

func (m *manager) Close() error {
	if err := m.conn.Close(); err != nil {
		return err
	}
	return nil
}
