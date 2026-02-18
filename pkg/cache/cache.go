package cache

import (
	"context"
	"errors"
	"time"

	"github.com/saibaend/template-svc/pkg/cache/memory"
	"github.com/saibaend/template-svc/pkg/cache/redis"
)

type cacheType string

const (
	TypeRedis    cacheType = "redis"
	TypeInMemory cacheType = "inMemory"
)

var ErrUnknownType = errors.New("unknown type")

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	SetWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Close() error
}

type Config struct {
	Type       cacheType
	DefaultTTL time.Duration

	RedisAddrs  []string
	RedisMaster string
	RedisDB     int

	InMemCleanupInterval time.Duration
}

func New(ctx context.Context, config Config) (Cache, error) {
	switch config.Type {
	case TypeRedis:
		return redis.New(ctx, redis.InitParams{
			Addrs:      config.RedisAddrs,
			Master:     config.RedisMaster,
			DB:         config.RedisDB,
			DefaultTTL: config.DefaultTTL,
		})
	case TypeInMemory:
		return memory.New(config.InMemCleanupInterval, config.DefaultTTL), nil
	default:
		return nil, ErrUnknownType
	}
}
