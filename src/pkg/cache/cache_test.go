package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newTestCache returns a cache without exercising the periodic-cleanup
// goroutine's expiration math; tests control eviction directly.
func newTestCache(maxEntries int) *Cache[string] {
	return New[string](maxEntries, time.Hour)
}

func TestCacheSetAndGet(t *testing.T) {
	c := newTestCache(4)
	_, ok := c.Get("missing")
	assert.False(t, ok, "Get on a missing key should report miss")

	c.Set("a", "a1")
	v, ok := c.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "a1", v)

	// Overwrite updates the value in place.
	c.Set("a", "a2")
	v, ok = c.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "a2", v)
}

func TestCacheRemove(t *testing.T) {
	c := newTestCache(4)
	c.Set("a", "a1")
	c.Remove("a")
	_, ok := c.Get("a")
	assert.False(t, ok, "Get after Remove should miss")

	// Removing a missing key is a no-op (no panic).
	c.Remove("never-set")
}

func TestCacheEvictsOldestWhenFull(t *testing.T) {
	c := newTestCache(2)
	c.Set("a", "a1")
	time.Sleep(2 * time.Millisecond)
	c.Set("b", "b1")
	// A brand-new third key must evict the oldest entry ("a").
	c.Set("c", "c1")

	_, aOk := c.Get("a")
	assert.False(t, aOk, "oldest key should be evicted once capacity is reached")
	_, bOk := c.Get("b")
	assert.True(t, bOk)
	_, cOk := c.Get("c")
	assert.True(t, cOk)
}

// TestCacheSetOverwriteDoesNotEvictSibling is the regression case: refreshing
// an existing key at capacity must not evict an unrelated entry. Before the
// fix, Set ran evictOldest whenever len>=maxEntries, so overwriting the newest
// key dropped its oldest sibling.
func TestCacheSetOverwriteDoesNotEvictSibling(t *testing.T) {
	c := newTestCache(2)
	c.Set("a", "a1") // oldest
	time.Sleep(2 * time.Millisecond)
	c.Set("b", "b1") // newest

	// Refresh the newest entry at capacity. The cache size is unchanged, so
	// nothing should be evicted.
	time.Sleep(2 * time.Millisecond)
	c.Set("b", "b2")

	v, ok := c.Get("b")
	assert.True(t, ok)
	assert.Equal(t, "b2", v)
	_, aOk := c.Get("a")
	assert.True(t, aOk, "overwriting an existing key must not evict a different entry")
}

// TestCacheSetOverwriteOldestKeepsCapacity guards the other branch of the fix:
// when the key being refreshed is itself the oldest one, the cache must still
// stay within maxEntries and keep the unrelated sibling.
func TestCacheSetOverwriteOldestKeepsCapacity(t *testing.T) {
	c := newTestCache(2)
	c.Set("a", "a1")
	time.Sleep(2 * time.Millisecond)
	c.Set("b", "b1")

	// Refresh the oldest entry.
	c.Set("a", "a2")
	assert.Len(t, c.cache, 2)

	v, ok := c.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "a2", v)
	_, bOk := c.Get("b")
	assert.True(t, bOk)
}
