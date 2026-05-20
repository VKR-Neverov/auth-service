package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fidesy-pay/auth-service/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct {
	db *redis.Client
}

func NewRedis(ctx context.Context) (*Redis, error) {
	c := &Redis{}

	cli := redis.NewClient(&redis.Options{
		Addr:     config.Get(config.RedisHost).(string),
		Password: config.Get(config.RedisPassword).(string),
		DB:       0,
	})

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis.Ping: %w", err)
	}

	c.db = cli

	return c, nil
}

func (c *Redis) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = c.db.Set(ctx, key, bytes, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis.Set: %w", err)
	}

	return nil
}

func (c *Redis) Get(ctx context.Context, key string, dst interface{}) (bool, error) {
	result := c.db.Get(ctx, key)
	if err := result.Err(); err != nil {
		// not found
		if result.Err() == redis.Nil {
			return false, nil
		}

		return false, fmt.Errorf("redis.Get: %w", err)
	}

	bytes, err := result.Bytes()
	if err != nil {
		return false, fmt.Errorf("result.Bytes: %w", err)
	}

	err = json.Unmarshal(bytes, &dst)
	if err != nil {
		return false, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return true, nil
}

func (c *Redis) Delete(ctx context.Context, key string) error {
	result := c.db.Del(ctx, key)
	if result.Err() != nil {
		return fmt.Errorf("db.Del: %w", result.Err())
	}

	return nil
}
