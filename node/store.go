package node

import "sync"

// KeyValueStore represents an in-memory key-value store
type KeyValueStore struct {
	mu    sync.RWMutex
	store map[string]string
}

// NewKeyValueStore initializes a new key-value store
func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		store: make(map[string]string),
	}
}

func (kv *KeyValueStore) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	val, ok := kv.store[key]
	return val, ok
}

func (kv *KeyValueStore) Put(key, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[key] = value
}

func (kv *KeyValueStore) Delete(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	_, ok := kv.store[key]
	if ok {
		delete(kv.store, key)
	}
	return ok
}
