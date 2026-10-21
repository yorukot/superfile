package disks

import (
	"context"
)

type manager struct{}

func New() (Manager, error) {
	return &manager{}, nil
}

func (m *manager) List(ctx context.Context) ([]Disk, error) {
	return nil, nil
}

func (m *manager) Close() error {
	// TODO ...
	return nil
}
