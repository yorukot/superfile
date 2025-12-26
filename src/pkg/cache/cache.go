package cache

import (
	"sync"
	"time"
)

type cacheItemInternal[T any] struct {
	obj       T
	Timestamp time.Time
}

type Cache[T any] struct {
	cache      map[string]cacheItemInternal[T]
	mutex      sync.RWMutex
	maxEntries int
	expiration time.Duration
}

func New[T any](maxEntries int, expiration time.Duration) *Cache[T] {
	cache := &Cache[T]{
		cache:      make(map[string]cacheItemInternal[T]),
		maxEntries: maxEntries,
		expiration: expiration,
	}

	// Start a cleanup goroutine
	go cache.periodicCleanup()

	return cache
}

// periodicCleanup removes expired entries periodically
func (c *Cache[T]) periodicCleanup() {
	//nolint:mnd // half of expiration for cleanup interval
	ticker := time.NewTicker(c.expiration / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpired()
	}
}

// cleanupExpired removes expired cache entries
func (c *Cache[T]) cleanupExpired() {
	now := time.Now()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, entry := range c.cache {
		if now.Sub(entry.Timestamp) > c.expiration {
			delete(c.cache, key)
		}
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if entry, exists := c.cache[key]; exists {
		return entry.obj, true
	}
	var res T
	return res, false
}

func (c *Cache[T]) Set(key string, obj T) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if we need to evict entries
	if len(c.cache) >= c.maxEntries {
		c.evictOldest()
	}

	c.cache[key] = cacheItemInternal[T]{
		obj:       obj,
		Timestamp: time.Now(),
	}
}

// evictOldest removes the oldest entry from the cache
func (c *Cache[T]) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	// Find the oldest entry
	for key, entry := range c.cache {
		if oldestKey == "" || entry.Timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.Timestamp
		}
	}

	// Remove the oldest entry
	if oldestKey != "" {
		delete(c.cache, oldestKey)
	}
}
