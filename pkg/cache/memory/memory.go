package memory

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"
)

var ErrNotFound = errors.New("key not found in cache")

type Cache struct {
	cache      *cache.Cache
	defaultTTL time.Duration
}

func New(cleanupInterval, defaultTTL time.Duration) *Cache {
	return &Cache{
		cache:      cache.New(defaultTTL, cleanupInterval),
		defaultTTL: defaultTTL,
	}
}

func (c *Cache) Get(_ context.Context, key string) (string, error) {
	value, ok := c.cache.Get(key)
	if !ok {
		return "", ErrNotFound
	}

	str, ok := value.(string)
	if !ok {
		return "", ErrNotFound
	}

	return str, nil
}

func (c *Cache) Set(_ context.Context, key string, value string) error {
	c.cache.Set(key, value, c.defaultTTL)
	return nil
}

func (c *Cache) SetWithTTL(_ context.Context, key string, value string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	c.cache.Set(key, value, ttl)
	return nil
}

func (c *Cache) Delete(_ context.Context, key string) error {
	c.cache.Delete(key)
	return nil
}

func (c *Cache) Close() error {
	c.cache.Flush()
	return nil
}
