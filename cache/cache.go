package goutil

import (
	"sync"
	"time"

	"github.com/AspieSoft/goutil/v7"
)

type CacheMap[K goutil.Hashable, V any] struct {
	value map[K]V
	err map[K]error
	lastUse map[K]time.Time
	exp time.Duration
	mu sync.Mutex
	null V
}

// NewCache creates a new cache map
func NewCache[K goutil.Hashable, V any](exp time.Duration) *CacheMap[K, V] {
	cache := CacheMap[K, V]{
		value: map[K]V{},
		err: map[K]error{},
		lastUse: map[K]time.Time{},
		exp: exp,
	}

	go func(){
		for {
			time.Sleep(10 * time.Minute)

			// default: remove cache items have not been accessed in over 2 hours
			cacheTime := cache.exp

			// SysFreeMemory returns the total free system memory in megabytes
			mb := goutil.SysFreeMemory()
			if mb < 200 && mb != 0 {
				// low memory: remove cache items have not been accessed in over 10 minutes
				cacheTime = 10 * time.Minute
			}else if mb < 500 && mb != 0 {
				// low memory: remove cache items have not been accessed in over 30 minutes
				cacheTime = 30 * time.Minute
			}else if mb < 2000 && mb != 0 {
				// low memory: remove cache items have not been accessed in over 1 hour
				cacheTime = 1 * time.Hour
			}else if mb > 64000 {
				// high memory: remove cache items have not been accessed in over 12 hour
				cacheTime = 12 * time.Hour
			}else if mb > 32000 {
				// high memory: remove cache items have not been accessed in over 6 hour
				cacheTime = 6 * time.Hour
			}else if mb > 16000 {
				// high memory: remove cache items have not been accessed in over 3 hour
				cacheTime = 3 * time.Hour
			}

			if cacheTime == 0 {
				continue
			}

			cache.DelOld(cacheTime)

			time.Sleep(10 * time.Second)

			// clear cache if were still critically low on available memory
			if mb := goutil.SysFreeMemory(); mb < 10 && mb != 0 {
				cache.DelOld(0)
			}
		}
	}()

	return &cache
}

// NewCacheCB is just like the NewCache method,
// but it returns the loop in a callback function, to avoid creating more goroutines
func NewCacheCB[K goutil.Hashable, V any](exp time.Duration) (*CacheMap[K, V], func()) {
	cache := CacheMap[K, V]{
		value: map[K]V{},
		err: map[K]error{},
		lastUse: map[K]time.Time{},
		exp: exp,
	}

	clearFunc := func(){
		// default: remove cache items have not been accessed in over 2 hours
		cacheTime := cache.exp

		// SysFreeMemory returns the total free system memory in megabytes
		mb := goutil.SysFreeMemory()
		if mb < 200 && mb != 0 {
			// low memory: remove cache items have not been accessed in over 10 minutes
			cacheTime = 10 * time.Minute
		}else if mb < 500 && mb != 0 {
			// low memory: remove cache items have not been accessed in over 30 minutes
			cacheTime = 30 * time.Minute
		}else if mb < 2000 && mb != 0 {
			// low memory: remove cache items have not been accessed in over 1 hour
			cacheTime = 1 * time.Hour
		}else if mb > 64000 {
			// high memory: remove cache items have not been accessed in over 12 hour
			cacheTime = 12 * time.Hour
		}else if mb > 32000 {
			// high memory: remove cache items have not been accessed in over 6 hour
			cacheTime = 6 * time.Hour
		}else if mb > 16000 {
			// high memory: remove cache items have not been accessed in over 3 hour
			cacheTime = 3 * time.Hour
		}

		if cacheTime == 0 {
			return
		}

		cache.DelOld(cacheTime)

		time.Sleep(10 * time.Second)

		// clear cache if were still critically low on available memory
		if mb := goutil.SysFreeMemory(); mb < 10 && mb != 0 {
			cache.DelOld(0)
		}
	}

	return &cache, clearFunc
}

// Get returns a value or an error if it exists
//
// if the object key does not exist, it will return both a nil/zero value (of the relevant type) and nil error
func (cache *CacheMap[K, V]) Get(key K) (V, error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if err, ok := cache.err[key]; ok {
		cache.lastUse[key] = time.Now()
		return cache.null, err
	}else if val, ok := cache.value[key]; ok {
		cache.lastUse[key] = time.Now()
		return val, nil
	}

	return cache.null, nil
}

// Set sets or adds a new key with either a value, or an error
func (cache *CacheMap[K, V]) Set(key K, value V, err error) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if err != nil {
		cache.err[key] = err
		delete(cache.value, key)
		cache.lastUse[key] = time.Now()
	}else{
		cache.value[key] = value
		delete(cache.err, key)
		cache.lastUse[key] = time.Now()
	}
}

// Del removes a cache item by key
func (cache *CacheMap[K, V]) Del(key K){
	cache.mu.Lock()
	defer cache.mu.Unlock()

	delete(cache.value, key)
	delete(cache.err, key)
	delete(cache.lastUse, key)
}

// DelOld removes old cache items
func (cache *CacheMap[K, V]) DelOld(cacheTime time.Duration){
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if cacheTime == 0 {
		for key := range cache.lastUse {
			delete(cache.value, key)
			delete(cache.err, key)
			delete(cache.lastUse, key)
		}
		return
	}

	now := time.Now().UnixNano()

	for key, lastUse := range cache.lastUse {
		if now - lastUse.UnixNano() > int64(cacheTime) {
			delete(cache.value, key)
			delete(cache.err, key)
			delete(cache.lastUse, key)
		}
	}
}

// Has returns true if a key value exists and is not an error
func (cache *CacheMap[K, V]) Has(key K) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	if _, ok := cache.err[key]; ok {
		cache.lastUse[key] = time.Now()
		return false
	}else if _, ok := cache.value[key]; ok {
		cache.lastUse[key] = time.Now()
		return true
	}

	return false
}

// Expire sets the ttl for all cache items
func (cache *CacheMap[K, V]) Expire(exp time.Duration) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.exp = exp

	return false
}

// Touch resets a cache items expiration so it will stay in the cache longer
func (cache *CacheMap[K, V]) Touch(key K) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.lastUse[key] = time.Now()

	return false
}

// ForEach runs a callback function for each cache item that has not expired
//
// in the callback, return true to continue, and false to break the loop
func (cache *CacheMap[K, V]) ForEach(cb func(key K, value V) bool, touch ...bool){
	cache.mu.Lock()
	keyList := []K{}
	for key := range cache.value {
		keyList = append(keyList, key)
	}
	cache.mu.Unlock()

	now := time.Now()
	for _, key := range keyList {
		cache.mu.Lock()
		if _, ok := cache.err[key]; ok {
			cache.lastUse[key] = now
			cache.mu.Unlock()
			continue
		}

		if now.UnixNano() - cache.lastUse[key].UnixNano() > int64(cache.exp) {
			delete(cache.value, key)
			delete(cache.err, key)
			delete(cache.lastUse, key)
			cache.mu.Unlock()
			continue
		}

		var val V
		if v, ok := cache.value[key]; ok {
			val = v
		}

		cache.mu.Unlock()

		if !cb(key, val) {
			break
		}
	}
}
