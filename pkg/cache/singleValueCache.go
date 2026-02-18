package cache

import (
	"sync"
	"time"
)

type SingleValueCache struct {
	mu        sync.RWMutex
	value     []byte
	expiresAt time.Time
}

func (c *SingleValueCache) Set(value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.value = value
	c.expiresAt = time.Now().Add(ttl)
}

func (c *SingleValueCache) Get() ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Now().After(c.expiresAt) {
		return nil, false
	}

	return c.value, true
}
