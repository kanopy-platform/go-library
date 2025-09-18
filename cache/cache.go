package cache

import (
	"maps"
	"sync"
)

// Cache represents a thread-safe in-memory cache with read/write mutex protection
// K is the type for keys, V is the type for values
type Cache[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// New creates and returns a new Cache instance
func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		data: make(map[K]V),
	}
}

// Read retrieves a value from the cache by key.
// Returns the value and a boolean indicating whether the key was found.
func (c *Cache[K, V]) Read(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.data[key]
	return value, exists
}

// Write stores a key-value pair in the cache.
func (c *Cache[K, V]) Write(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// Delete removes a key-value pair from the cache.
// Returns true if the key existed and was deleted, false otherwise.
func (c *Cache[K, V]) Delete(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, exists := c.data[key]
	if exists {
		delete(c.data, key)
	}
	return exists
}

// Size returns the number of items in the cache.
func (c *Cache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}

// Clear removes all items from the cache.
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[K]V)
}

// Keys returns all keys in the cache.
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}
	return keys
}

// Replace replaces the entire cache data with the provided map.
// The provided map is copied to ensure thread safety.
func (c *Cache[K, V]) Replace(data map[K]V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Create a new map and copy the data to avoid sharing references
	c.data = make(map[K]V, len(data))
	maps.Copy(c.data, data)
}
