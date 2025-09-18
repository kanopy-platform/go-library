package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheReadWrite(t *testing.T) {
	cache := New[string, string]()

	// Test writing and reading a value
	key := "test-key"
	value := "test-value"
	cache.Write(key, value)

	result, exists := cache.Read(key)
	if !exists {
		t.Errorf("Expected key '%s' to exist in cache", key)
	}
	if result != value {
		t.Errorf("Expected value '%s', got '%s'", value, result)
	}
}

func TestCacheReadNonExistent(t *testing.T) {
	cache := New[string, string]()

	// Test reading a non-existent key
	_, exists := cache.Read("non-existent-key")
	if exists {
		t.Error("Expected non-existent key to return false")
	}
}

func TestCacheWithDifferentTypes(t *testing.T) {
	// Test string cache
	stringCache := New[string, string]()
	stringCache.Write("string", "test")
	if val, exists := stringCache.Read("string"); !exists || val != "test" {
		t.Errorf("String test failed: exists=%v, val=%v", exists, val)
	}

	// Test slice cache
	sliceCache := New[string, []int]()
	sliceCache.Write("slice", []int{1, 2, 3})
	if val, exists := sliceCache.Read("slice"); !exists {
		t.Error("Slice test failed: key not found")
	} else if len(val) != 3 {
		t.Errorf("Slice test failed: val=%v", val)
	}

	// Test map cache
	mapCache := New[string, map[string]int]()
	mapCache.Write("map", map[string]int{"a": 1, "b": 2})
	if val, exists := mapCache.Read("map"); !exists {
		t.Error("Map test failed: key not found")
	} else if val["a"] != 1 {
		t.Errorf("Map test failed: val=%v", val)
	}
}

func TestCacheDelete(t *testing.T) {
	cache := New[string, string]()

	// Add a key-value pair
	cache.Write("key", "value")

	// Verify it exists
	_, exists := cache.Read("key")
	if !exists {
		t.Error("Key should exist before deletion")
	}

	// Delete the key
	deleted := cache.Delete("key")
	if !deleted {
		t.Error("Delete should return true for existing key")
	}

	// Verify it no longer exists
	_, exists = cache.Read("key")
	if exists {
		t.Error("Key should not exist after deletion")
	}

	// Try to delete non-existent key
	deleted = cache.Delete("non-existent")
	if deleted {
		t.Error("Delete should return false for non-existent key")
	}
}

func TestCacheSize(t *testing.T) {
	cache := New[string, string]()

	// Initial size should be 0
	if cache.Size() != 0 {
		t.Errorf("Expected initial size to be 0, got %d", cache.Size())
	}

	// Add some items
	cache.Write("key1", "value1")
	cache.Write("key2", "value2")
	cache.Write("key3", "value3")

	if cache.Size() != 3 {
		t.Errorf("Expected size to be 3, got %d", cache.Size())
	}

	// Delete one item
	cache.Delete("key2")

	if cache.Size() != 2 {
		t.Errorf("Expected size to be 2 after deletion, got %d", cache.Size())
	}
}

func TestCacheClear(t *testing.T) {
	cache := New[string, string]()

	// Add some items
	cache.Write("key1", "value1")
	cache.Write("key2", "value2")
	cache.Write("key3", "value3")

	// Verify items exist
	if cache.Size() != 3 {
		t.Errorf("Expected size to be 3 before clear, got %d", cache.Size())
	}

	// Clear the cache
	cache.Clear()

	// Verify cache is empty
	if cache.Size() != 0 {
		t.Errorf("Expected size to be 0 after clear, got %d", cache.Size())
	}

	// Verify items no longer exist
	_, exists := cache.Read("key1")
	if exists {
		t.Error("Key should not exist after clear")
	}
}

func TestCacheKeys(t *testing.T) {
	cache := New[string, string]()

	// Add some items
	cache.Write("key1", "value1")
	cache.Write("key2", "value2")
	cache.Write("key3", "value3")

	keys := cache.Keys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}

	// Verify all keys are present (order doesn't matter)
	keyMap := make(map[string]bool)
	for _, key := range keys {
		keyMap[key] = true
	}

	expectedKeys := []string{"key1", "key2", "key3"}
	for _, expectedKey := range expectedKeys {
		if !keyMap[expectedKey] {
			t.Errorf("Expected key '%s' not found in keys list", expectedKey)
		}
	}
}

func TestCacheConcurrency(t *testing.T) {
	cache := New[string, string]()
	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup

	// Start multiple goroutines performing writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := fmt.Sprintf("value-%d-%d", id, j)
				cache.Write(key, value)
			}
		}(i)
	}

	// Start multiple goroutines performing reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				// Wait a bit to allow some writes to complete
				time.Sleep(time.Microsecond)
				cache.Read(key)
			}
		}(i)
	}

	// Start some goroutines performing mixed operations
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations/2; j++ {
				key := fmt.Sprintf("mixed-key-%d-%d", id, j)
				cache.Write(key, "mixed-value")
				cache.Read(key)
				cache.Delete(key)
				_ = cache.Size()
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is still functional after concurrent operations
	cache.Write("test-after-concurrency", "test-value")
	value, exists := cache.Read("test-after-concurrency")
	if !exists || value != "test-value" {
		t.Error("Cache not functional after concurrent operations")
	}
}

func TestCacheReplace(t *testing.T) {
	cache := New[string, string]()

	// Add some initial data
	cache.Write("key1", "value1")

	// Verify initial state
	if cache.Size() != 1 {
		t.Errorf("Expected initial size to be 1, got %d", cache.Size())
	}

	// Create new data to replace with
	newData := map[string]string{
		"newKey1": "newValue1",
		"newKey2": "newValue2",
	}

	// Replace the cache data
	cache.Replace(newData)

	// Verify new size
	if cache.Size() != 2 {
		t.Errorf("Expected size to be 2 after replace, got %d", cache.Size())
	}

	// Verify old keys no longer exist
	if _, exists := cache.Read("key1"); exists {
		t.Error("Old key1 should not exist after replace")
	}

	// Verify new keys exist with correct values
	for key, expectedValue := range newData {
		if value, exists := cache.Read(key); !exists {
			t.Errorf("New key '%s' should exist after replace", key)
		} else if value != expectedValue {
			t.Errorf("Expected value '%s' for key '%s', got '%s'", expectedValue, key, value)
		}
	}
}

func TestCacheReplaceWithNilMap(t *testing.T) {
	cache := New[string, string]()

	// Add some initial data
	cache.Write("key1", "value1")

	// Replace with nil map
	cache.Replace(nil)

	// Verify cache is now empty
	if cache.Size() != 0 {
		t.Errorf("Expected size to be 0 after replace with nil map, got %d", cache.Size())
	}

	// Verify cache is still functional
	cache.Write("test", "value")
	if value, exists := cache.Read("test"); !exists || value != "value" {
		t.Error("Cache should be functional after replace with nil map")
	}
}

func TestCacheReplaceDataIndependence(t *testing.T) {
	cache := New[string, string]()

	// Create source data
	sourceData := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	// Replace cache data
	cache.Replace(sourceData)

	// Modify the source data after replace
	sourceData["key1"] = "modified"
	sourceData["key3"] = "new"

	// Verify cache data is not affected by modifications to source
	if value, exists := cache.Read("key1"); !exists || value != "value1" {
		t.Errorf("Cache should not be affected by source modifications: exists=%v, value=%s", exists, value)
	}
	if _, exists := cache.Read("key3"); exists {
		t.Error("Cache should not have key3 added to source after replace")
	}

	// Verify cache has only the original replaced data
	if cache.Size() != 2 {
		t.Errorf("Expected cache size to remain 2, got %d", cache.Size())
	}
}

func TestCacheReplaceConcurrency(t *testing.T) {
	cache := New[string, string]()
	const numGoroutines = 50
	const numOperations = 50

	var wg sync.WaitGroup

	// Start goroutines that perform replace operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				replaceData := map[string]string{
					fmt.Sprintf("replace-key-%d-%d", id, j): fmt.Sprintf("replace-value-%d-%d", id, j),
				}
				cache.Replace(replaceData)
			}
		}(i)
	}

	// Start goroutines that perform read operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				_ = cache.Size()
				cache.Keys()
				// Try to read various keys
				cache.Read(fmt.Sprintf("replace-key-%d-%d", id%10, j%10))
			}
		}(i)
	}

	// Start goroutines that perform mixed operations
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations/2; j++ {
				// Write some data
				cache.Write(fmt.Sprintf("mixed-key-%d-%d", id, j), "mixed-value")

				// Replace with different data
				replaceData := map[string]string{
					fmt.Sprintf("concurrent-key-%d", id): fmt.Sprintf("concurrent-value-%d", j),
				}
				cache.Replace(replaceData)

				// Read and check size
				_ = cache.Size()
				cache.Read(fmt.Sprintf("concurrent-key-%d", id))
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is still functional after concurrent operations
	testData := map[string]string{
		"final-test-key": "final-test-value",
	}
	cache.Replace(testData)

	if cache.Size() != 1 {
		t.Errorf("Expected size 1 after final replace, got %d", cache.Size())
	}

	if value, exists := cache.Read("final-test-key"); !exists || value != "final-test-value" {
		t.Error("Cache not functional after concurrent replace operations")
	}
}
