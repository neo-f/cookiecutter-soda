package utils

import (
	"sync"
	"time"
)

type Cache[T any] struct {
	storage map[string]T
	ttl     map[string]time.Time
	mu      sync.RWMutex
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{
		storage: make(map[string]T),
		ttl:     make(map[string]time.Time),
		mu:      sync.RWMutex{},
	}
}

func (c *Cache[T]) Get(key string, ttl time.Duration, fn func() (T, error)) (T, error) {
	c.mu.RLock()
	if v, ok := c.storage[key]; ok {
		if time.Now().Before(c.ttl[key]) {
			c.mu.RUnlock()
			return v, nil
		}
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.storage[key]; ok {
		if time.Now().Before(c.ttl[key]) {
			return v, nil
		}
		delete(c.storage, key)
		delete(c.ttl, key)
	}

	v, err := fn()
	if err != nil {
		return v, err
	}
	c.storage[key] = v
	c.ttl[key] = time.Now().Add(ttl)
	return v, nil
}
