package memory

import (
	"context"
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration time.Time
}

type MemoryCache struct {
	mu         sync.RWMutex
	items      map[string]item
	maxSize    int
	defaultTTL time.Duration
}

func New(maxSize int, defaultTTL time.Duration) *MemoryCache {
	return &MemoryCache{
		items:      make(map[string]item),
		maxSize:    maxSize,
		defaultTTL: defaultTTL,
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ttl == 0 {
		ttl = c.defaultTTL
	}

	if len(c.items) >= c.maxSize {
		// We could opt for random eviction, double linked lists or whatever if we wanted to optimize the current O(n) approach
		c.evictOldest()
	}

	c.items[key] = item{
		value:      value,
		expiration: time.Now().Add(ttl),
	}

	return nil
}

func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

func (c *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestTime.IsZero() || item.expiration.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.expiration
		}
	}

	if oldestKey != "" {
		delete(c.items, oldestKey)
	}
}

