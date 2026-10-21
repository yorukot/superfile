package disks

import (
	"context"
)

type Disk struct {
	Name       string
	MountPoint string
	Mounted    bool
	// If the disk is an internal disk parition
	// that is not running the current OS
	Internal bool

	// Linux only field to store the ObjectPath
	// of a block device for sorting. See
	// disks_linux.go List()
	ObjectPath string
}

type Manager interface {
	List(ctx context.Context) ([]Disk, error)
	Close() error
}
