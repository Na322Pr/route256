package cache

import (
	"sync"
	"time"
)

type CacheClient[K comparable, V any] struct {
	ttl  time.Duration
	lock sync.RWMutex
	data map[K]*Cached[V]
}

func NewCacheClient[K comparable, V any](ttl time.Duration) *CacheClient[K, V] {
	return &CacheClient[K, V]{
		ttl:  ttl,
		data: make(map[K]*Cached[V]),
	}
}

func (c *CacheClient[K, V]) Get(key K) (V, bool) {
	c.lock.RLock()
	v, ok := c.data[key]
	c.lock.RUnlock()

	if ok && !v.Expired(time.Now()) {
		return v.Value(), true
	}

	return (&Cached[V]{}).Value(), false
}

func (c *CacheClient[K, V]) Set(key K, value V, now time.Time) {
	wrapped := NewCached(time.Now().Add(c.ttl), value)

	c.lock.Lock()
	c.data[key] = wrapped
	c.lock.Unlock()
}
