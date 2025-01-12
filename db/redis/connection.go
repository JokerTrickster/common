package redis

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/JokerTrickster/common/aws"
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
func (r *RedisService) Initialize(ctx context.Context, isLocal bool, ssmKeys []string) error {
	var err error
	r.initOnce.Do(func() {
		connectionString, e := r.getConnectionString(ctx, isLocal, ssmKeys)
		if e != nil {
			err = e
			return
		}

		opt, e := redis.ParseURL(connectionString)
		if e != nil {
			err = fmt.Errorf("failed to parse Redis URL: %w", e)
			return
		}

		r.client = redis.NewClient(opt)

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

// getConnectionString generates the Redis connection string
func (r *RedisService) getConnectionString(ctx context.Context, isLocal bool, ssmKeys []string) (string, error) {
	if isLocal {
		user := getEnvOrFallback("REDIS_USER", "")
		password := getEnvOrFallback("REDIS_PASSWORD", "")
		return fmt.Sprintf("redis://%s:%s@localhost:6379/0", user, password), nil
	}

	ssmServiceClient := aws.SSMService{}
	dbInfos, err := ssmServiceClient.AwsSsmGetParams(ctx, ssmKeys)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Redis SSM parameters: %w", err)
	}

	// SSM Keys order: [host, port, db, user, password]
	return fmt.Sprintf("redis://%s:%s@%s:%s/%s",
		dbInfos[3], // user
		dbInfos[4], // password
		dbInfos[0], // host
		dbInfos[1], // port
		dbInfos[2], // db
	), nil
}
