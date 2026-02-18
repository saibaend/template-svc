package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultTimeout = 10 * time.Second
)

func connect(ctx context.Context, addrs []string, masterName string, db int) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	//nolint:exhaustruct
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    masterName,
		SentinelAddrs: addrs,
		DB:            db,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
