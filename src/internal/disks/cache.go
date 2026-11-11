package disks

import (
	"context"
	"time"

	"github.com/yorukot/superfile/src/pkg/cache"
)

const (
	cacheSize       = 1 // cacheSize is 1 because the Disk slice is stored under one key
	cacheExpiration = 2 * time.Second
	cacheKey        = "disks"
)

type Cache struct {
	manager Manager
	cache   *cache.Cache[listResult]
}

// listResult holds the outcome of one inner List() call
type listResult struct {
	disks []Disk
	err   error
}

func NewCache(manager Manager) *Cache {
	return &Cache{
		manager: manager,
		cache:   cache.New[listResult](cacheSize, cacheExpiration),
	}
}

func (c *Cache) List(ctx context.Context) ([]Disk, error) {
	if res, ok := c.cache.Get(cacheKey); ok {
		return res.disks, res.err
	}

	list, err := c.manager.List(ctx)
	c.cache.Set(cacheKey, listResult{disks: list, err: err})
	return list, err
}

func (c *Cache) Close() error {
	return c.manager.Close()
}
