package health

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisChecker struct {
	redis *redis.Client
}

func NewRedisChecker(redis *redis.Client) *RedisChecker {
	return &RedisChecker{
		redis: redis,
	}
}

func (c *RedisChecker) Ping(ctx context.Context) error {
	return c.redis.Ping(ctx).Err()
}
