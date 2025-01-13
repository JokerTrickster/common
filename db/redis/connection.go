package redis

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client   *redis.Client
	initOnce sync.Once
}

var instance *RedisService
var once sync.Once

// GetRedisService returns a singleton RedisService instance
func GetRedisService() *RedisService {
	once.Do(func() {
		instance = &RedisService{}
	})
	return instance
}

// Initialize initializes the Redis client
func (r *RedisService) Initialize(ctx context.Context, connectionString string) error {
	var err error
	r.initOnce.Do(func() {
		// Parse the Redis connection string
		opt, e := redis.ParseURL(connectionString)
		if e != nil {
			err = fmt.Errorf("failed to parse Redis URL: %w", e)
			return
		}

		// Create a new Redis client
		r.client = redis.NewClient(opt)

		// Test the connection
		_, e = r.client.Ping(ctx).Result()
		if e != nil {
			err = fmt.Errorf("failed to connect to Redis: %w", e)
			return
		}

		log.Println("Connected to Redis!")
	})
	return err
}

// GetClient returns the Redis client
func (r *RedisService) GetClient(ctx context.Context) (*redis.Client, error) {
	if r.client == nil {
		return nil, fmt.Errorf("Redis client is not initialized. Call Initialize first")
	}
	return r.client, nil
}
