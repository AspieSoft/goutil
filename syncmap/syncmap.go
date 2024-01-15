package goutil

import (
	"sync"

	"github.com/AspieSoft/goutil/v7"
)

type SyncMap[K goutil.Hashable, V any] struct {
	value map[K]V
	hasVal map[K]bool
	mu sync.Mutex
	null V
}

func NewMap[K goutil.Hashable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		value: map[K]V{},
		hasVal: map[K]bool{},
	}
}

// Get returns a value or an error if it exists
func (syncmap *SyncMap[K, V]) Get(key K) (V, bool) {
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	if hasVal, ok := syncmap.hasVal[key]; !ok || !hasVal {
		return syncmap.null, false
	}else if val, ok := syncmap.value[key]; ok {
		return val, true
	}

	return syncmap.null, false
}

// Set sets or adds a new key with a value
func (syncmap *SyncMap[K, V]) Set(key K, value V) {
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	syncmap.value[key] = value
	syncmap.hasVal[key] = true
}

// Del removes an item by key
func (syncmap *SyncMap[K, V]) Del(key K){
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	delete(syncmap.value, key)
	delete(syncmap.hasVal, key)
}

// Has returns true if a key value exists in the list
func (syncmap *SyncMap[K, V]) Has(key K) bool {
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	if hasVal, ok := syncmap.hasVal[key]; !ok || !hasVal {
		return false
	}else if _, ok := syncmap.value[key]; ok {
		return true
	}

	return false
}

// ForEach runs a callback function for each key value pair
//
// in the callback, return true to continue, and false to break the loop
func (syncmap *SyncMap[K, V]) ForEach(cb func(key K, value V) bool){
	syncmap.mu.Lock()
	defer syncmap.mu.Unlock()

	for k, v := range syncmap.value {
		if !cb(k, v) {
			break
		}
	}
}
