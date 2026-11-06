package disks

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

const dbusCallTimeout = 2 * time.Second

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
	ctx, cancel := context.WithTimeout(ctx, dbusCallTimeout)
	defer cancel()

	var objs map[dbus.ObjectPath]map[string]map[string]dbus.Variant

	obj := m.conn.Object("org.freedesktop.UDisks2", "/org/freedesktop/UDisks2")
	err := obj.CallWithContext(ctx, "org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&objs)
	if err != nil {
		return nil, err
	}

	var disks []Disk
	for path, ifaces := range objs {
		if disk, ok := diskFromObject(path, ifaces); ok {
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
	return m.conn.Close()
}
