package redis

import (
	"context"
	"fmt"
	"time"
)

// SetKey sets a key-value pair in Redis
func (r *RedisService) SetKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if r.client == nil {
		return fmt.Errorf("Redis client is not initialized")
	}
	return r.client.Set(ctx, key, value, expiration).Err()
}

// GetKey retrieves a value by key from Redis
func (r *RedisService) GetKey(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", fmt.Errorf("Redis client is not initialized")
	}
	return r.client.Get(ctx, key).Result()
}

// DeleteKey deletes a key from Redis
func (r *RedisService) DeleteKey(ctx context.Context, key string) error {
	if r.client == nil {
		return fmt.Errorf("Redis client is not initialized")
	}
	return r.client.Del(ctx, key).Err()
}
