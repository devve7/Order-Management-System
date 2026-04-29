package cache

import (
	"Order-Management-System/internal/application/ports"
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type LoggingCache struct {
	next   ports.Cache
	logger *logrus.Logger
}

func NewLoggingCache(next ports.Cache, logger *logrus.Logger) *LoggingCache {
	return &LoggingCache{
		next:   next,
		logger: logger,
	}
}

func (c *LoggingCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.next.Get(ctx, key)
	if err != nil {
		status := "error"

		if errors.Is(err, redis.Nil) {
			status = "miss"
		}

		c.logger.WithFields(logrus.Fields{
			"cache_key": key,
			"cache":     status,
			"error":     err.Error(),
		}).Info("cache get")

		return "", err
	}
	c.logger.WithFields(logrus.Fields{
		"cache_key": key,
		"cache":     "hit",
	}).Info("cache get")

	return value, nil
}

func (c *LoggingCache) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	err := c.next.Set(ctx, key, value, ttlSeconds)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"cache_key": key,
			"cache":     "set_error",
			"error":     err.Error(),
		}).Info("cache set")
		return err
	}

	c.logger.WithFields(logrus.Fields{
		"cache_key": key,
		"cache":     "set",
		"ttl":       ttlSeconds,
	}).Info("cache set")

	return nil
}

func (c *LoggingCache) Delete(ctx context.Context, key string) error {
	err := c.next.Delete(ctx, key)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"cache_key": key,
			"cache":     "delete_error",
			"error":     err.Error(),
		}).Info("cache delete")
		return err
	}

	c.logger.WithFields(logrus.Fields{
		"cache_key": key,
		"cache":     "delete",
	}).Info("cache delete")

	return nil
}
