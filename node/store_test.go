package node

import (
	"testing"
)

func TestKeyValueStore(t *testing.T) {
	kv := NewKeyValueStore()

	kv.Put("key1", "value1")
	value, exists := kv.Get("key1")
	if !exists || value != "value1" {
		t.Errorf("Expected 'value1' for 'key1', got '%s'", value)
	}

	// Test Get with a non-existing key
	_, exists = kv.Get("nonexistent")
	if exists {
		t.Error("Expected 'false' for 'nonexistent' key, but it exists")
	}

	deleted := kv.Delete("key1")
	if !deleted {
		t.Error("Expected true for deleting 'key1', but it was not deleted")
	}

	// Test Delete on non-existing key
	deleted = kv.Delete("nonexistent")
	if deleted {
		t.Error("Expected false for deleting 'nonexistent', but it was deleted")
	}

}
