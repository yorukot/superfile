package filepreview

import (
	"sync"
	"time"
)

// ImagePreviewCache stores cached image previews
type ImagePreviewCache struct {
	cache      map[string]*CachedPreview
	mutex      sync.RWMutex
	maxEntries int
	expiration time.Duration
}

// CachedPreview represents a cached image preview
type CachedPreview struct {
	Preview    string
	Timestamp  time.Time
	Renderer   ImageRenderer
	Dimensions string // "width,height,bgColor,sideAreaWidth"
}

// NewImagePreviewCache creates a new image preview cache
func NewImagePreviewCache(maxEntries int, expiration time.Duration) *ImagePreviewCache {
	cache := &ImagePreviewCache{
		cache:      make(map[string]*CachedPreview),
		maxEntries: maxEntries,
		expiration: expiration,
	}

	// Start a cleanup goroutine
	go cache.periodicCleanup()

	return cache
}

// periodicCleanup removes expired entries periodically
func (c *ImagePreviewCache) periodicCleanup() {
	//nolint:mnd // half of expiration for cleanup interval
	ticker := time.NewTicker(c.expiration / 2)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpired()
	}
}

// cleanupExpired removes expired cache entries
func (c *ImagePreviewCache) cleanupExpired() {
	now := time.Now()
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for key, entry := range c.cache {
		if now.Sub(entry.Timestamp) > c.expiration {
			delete(c.cache, key)
		}
	}
}

// Get retrieves a cached preview if available
func (c *ImagePreviewCache) Get(path, dimensions string, renderer ImageRenderer) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	cacheKey := path + ":" + dimensions

	if entry, exists := c.cache[cacheKey]; exists {
		if entry.Renderer == renderer && time.Since(entry.Timestamp) < c.expiration {
			return entry.Preview, true
		}
	}

	return "", false
}

// Set stores a preview in the cache
func (c *ImagePreviewCache) Set(path, dimensions, preview string, renderer ImageRenderer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Check if we need to evict entries
	if len(c.cache) >= c.maxEntries {
		c.evictOldest()
	}

	cacheKey := path + ":" + dimensions
	c.cache[cacheKey] = &CachedPreview{
		Preview:    preview,
		Timestamp:  time.Now(),
		Renderer:   renderer,
		Dimensions: dimensions,
	}
}

// evictOldest removes the oldest entry from the cache
func (c *ImagePreviewCache) evictOldest() {
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
