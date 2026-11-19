package disks

import (
	"context"
	"io"

	"github.com/google/deck"
	"github.com/google/deck/backends/logger"
	"github.com/google/glazier/go/storage"
)

type manager struct {
	Service *storage.Service
}

func New() (Manager, error) {
	// This discards the "no backends configured" log when calling m.Service.GetVolumes().
	// The entire deck package is not used anywhere else besides for here to discard
	// the log, there is no other way to do so - to the best of my knowledge.
	deck.Add(logger.Init(io.Discard, 0))

	svc, err := storage.Connect()
	if err != nil {
		return nil, err
	}
	return &manager{Service: &svc}, nil
}

func (m *manager) List(ctx context.Context) ([]Disk, error) {
	var disks []Disk

	if errP := getPartitionlessDisks(m.Service, &disks); errP != nil {
		return nil, errP
	}

	if errV := getVolumes(m.Service, &disks); errV != nil {
		return disks, errV
	}

	return disks, nil
}

func (m *manager) Close() error {
	m.Service.Close()
	return nil
}
