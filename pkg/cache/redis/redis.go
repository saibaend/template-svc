package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrNotFound    = errors.New("key not found in cache")
	ErrEmptyAddrs  = errors.New("empty redis addrs")
	ErrEmptyMaster = errors.New("empty redis master")
	ErrConnect     = errors.New("failed to connect")
)

type InitParams struct {
	Addrs      []string
	Master     string
	DB         int
	DefaultTTL time.Duration
}

type Cache struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func New(ctx context.Context, params InitParams) (*Cache, error) {
	if len(params.Addrs) == 0 {
		return nil, ErrEmptyAddrs
	}

	if params.Master == "" {
		return nil, ErrEmptyMaster
	}

	client, err := connect(ctx, params.Addrs, params.Master, params.DB)
	if err != nil {
		return nil, errors.Join(ErrConnect, err)
	}

	return &Cache{
		client:     client,
		defaultTTL: params.DefaultTTL,
	}, nil
}

func (c Cache) Get(ctx context.Context, key string) (string, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNotFound
		}

		return "", err
	}

	return result, nil
}

func (c Cache) Set(ctx context.Context, key, value string) error {
	return c.client.Set(ctx, key, value, c.defaultTTL).Err()
}

func (c Cache) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c Cache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c Cache) Close() error {
	return c.client.Close()
}
