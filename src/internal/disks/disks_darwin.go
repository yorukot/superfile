package disks

import (
	"context"

	"github.com/shirou/gopsutil/v4/disk"
)

type manager struct{}

func New() (Manager, error) {
	return &manager{}, nil
}

func (m *manager) List(ctx context.Context) ([]Disk, error) {
	parts, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []Disk
	for _, part := range parts {
		if shouldListDisk(part.Mountpoint) {
			disks = append(disks, Disk{
				Name:       diskName(part.Mountpoint),
				MountPoint: part.Mountpoint,
				Mounted:    part.Mountpoint != "",
			})
		}
	}

	return disks, nil
}

func (m *manager) Close() error {
	return nil
}
