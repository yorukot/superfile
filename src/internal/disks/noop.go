package disks

import "context"

type noopManager struct{}

func NewNoop() Manager {
	return noopManager{}
}

func (noopManager) List(ctx context.Context) ([]Disk, error) {
	return nil, nil
}

func (noopManager) Close() error {
	return nil
}
