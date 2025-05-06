package main

import (
	"fmt"
	"sync"
	"time"
)

// CacheEntry represents a cached response
type CacheEntry struct {
	Response    []byte
	ContentType string
	StatusCode  int
	Expiry      time.Time
}

// InMemoryCache is a simple in-memory cache
type InMemoryCache struct {
	cache map[string]CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
}

// NewInMemoryCache creates a new InMemoryCache with a default TTL
func NewInMemoryCache(ttl time.Duration) *InMemoryCache {
	inMemoryCache := &InMemoryCache{
		cache: make(map[string]CacheEntry),
		ttl:   ttl,
	}
	// Starts a cleanup timer to remove expired items
	go inMemoryCache.startCleanupTimer()
	return inMemoryCache
}

// Set adds a response to the cache
func (c *InMemoryCache) Set(key string, entry CacheEntry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = entry
}

// Get retrieves a cached response by key
func (c *InMemoryCache) Get(key string) (CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, found := c.cache[key]
	if !found {
		return CacheEntry{}, false
	}

	if time.Now().After(entry.Expiry) {
		// We could check if a cache item has expired and remove it from the cache here.
		// But it's generally better to avoid modifications under a read lock.
		// Instead, we just return nil and let the caller handle it.
		return CacheEntry{}, false
	}

	return entry, found
}

// cleanupExpired removes expired items from the cache.
func (c *InMemoryCache) cleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, item := range c.cache {
		if now.After(item.Expiry) {
			delete(c.cache, key)
			fmt.Printf("Cache cleanup: Removed expired key %s\n", key)
		}
	}
}

func (c *InMemoryCache) startCleanupTimer() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for range ticker.C {
		c.cleanupExpired()
	}
}

// Clear clears all items from the cache.
func (c *InMemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create a new map to replace the old one.
	c.cache = make(map[string]CacheEntry)
	fmt.Println("Cache cleared.")
}
