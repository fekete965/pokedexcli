package internal

import (
	"sync"
	"time"
)

type Cache struct {
	mu    sync.Mutex
	entry map[string]cacheEntry
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	newCache := Cache{
		mu:    sync.Mutex{},
		entry: make(map[string]cacheEntry),
	}

	go newCache.reapLoop(interval)

	return &newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, hasEntry := c.entry[key]; hasEntry {
		return entry.val, true
	}

	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.deleteExpiresEntries(interval)
	}
}

func (c *Cache) deleteExpiresEntries(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.entry {
		if time.Since(entry.createdAt) > interval {
			delete(c.entry, key)
		}
	}
}
