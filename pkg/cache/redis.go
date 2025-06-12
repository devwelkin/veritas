package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("could not ping redis: %w", err)
	}

	return client, nil
}
