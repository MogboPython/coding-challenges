package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewInMemoryCache(t *testing.T) {
	ttl := 5 * time.Minute
	cache := NewInMemoryCache(ttl)

	assert.NotNil(t, cache)
	assert.NotNil(t, cache.cache)
	assert.Equal(t, ttl, cache.ttl)
}

func TestSetAndGet(t *testing.T) {
	cache := NewInMemoryCache(5 * time.Minute)

	key := "test-key"
	entry := CacheEntry{
		Response:    []byte("test response"),
		ContentType: "text/plain",
		StatusCode:  200,
		Expiry:      time.Now().Add(time.Minute),
	}

	cache.Set(key, entry)

	retrieved, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, entry.Response, retrieved.Response)
	assert.Equal(t, entry.ContentType, retrieved.ContentType)
	assert.Equal(t, entry.StatusCode, retrieved.StatusCode)
	assert.Equal(t, entry.Expiry.Unix(), retrieved.Expiry.Unix())
}

func TestGetNonExistent(t *testing.T) {
	cache := NewInMemoryCache(5 * time.Minute)

	_, found := cache.Get("non-existent-key")
	assert.False(t, found)
}

func TestGetExpired(t *testing.T) {
	cache := NewInMemoryCache(5 * time.Minute)

	// Create an already expired entry
	key := "expired-key"
	entry := CacheEntry{
		Response:    []byte("expired data"),
		ContentType: "text/plain",
		StatusCode:  200,
		Expiry:      time.Now().Add(-time.Minute), // Expired 1 minute ago
	}

	cache.Set(key, entry)

	_, found := cache.Get(key)
	assert.False(t, found)
}

func TestCleanupExpired(t *testing.T) {
	cache := NewInMemoryCache(5 * time.Minute)

	cache.Set("valid", CacheEntry{
		Response:    []byte("valid"),
		ContentType: "text/plain",
		StatusCode:  200,
		Expiry:      time.Now().Add(time.Minute),
	})

	cache.Set("expired1", CacheEntry{
		Response:    []byte("expired1"),
		ContentType: "text/plain",
		StatusCode:  200,
		Expiry:      time.Now().Add(-time.Second),
	})

	cache.Set("expired2", CacheEntry{
		Response:    []byte("expired2"),
		ContentType: "text/plain",
		StatusCode:  200,
		Expiry:      time.Now().Add(-2 * time.Second),
	})

	// Run the cleanup
	cache.cleanupExpired()

	_, found1 := cache.Get("expired1")
	assert.False(t, found1)
	_, found2 := cache.Get("expired2")
	assert.False(t, found2)

	// Check that valid item is still there (direct cache access to bypass expiry check in Get)
	cache.mu.RLock()
	_, validFound := cache.cache["valid"]
	cache.mu.RUnlock()
	assert.True(t, validFound)
}

func TestStartCleanupTimer(t *testing.T) {
	// Use a short TTL for testing
	shortTTL := 100 * time.Millisecond
	cache := NewInMemoryCache(shortTTL)

	cache.Set("soon-to-expire", CacheEntry{
		Response:    []byte("soon gone"),
		ContentType: "text/plain",
		StatusCode:  200,
		Expiry:      time.Now().Add(50 * time.Millisecond),
	})

	cache.Set("staying-valid", CacheEntry{
		Response:    []byte("staying"),
		ContentType: "text/plain",
		StatusCode:  200,
		Expiry:      time.Now().Add(500 * time.Millisecond),
	})

	// Wait for cleanup to occur
	time.Sleep(200 * time.Millisecond)

	cache.mu.RLock()
	_, expiredFound := cache.cache["soon-to-expire"]
	cache.mu.RUnlock()
	assert.False(t, expiredFound)

	// Check that valid entry is still there
	cache.mu.RLock()
	_, validFound := cache.cache["staying-valid"]
	cache.mu.RUnlock()
	assert.True(t, validFound)
}
