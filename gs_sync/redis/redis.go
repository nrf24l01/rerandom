package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nrf24l01/rerandom/gs_sync/config"
	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the redis client
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
	saveKey string
}

func InitRedisFromCFG(cfg *config.Config) *RedisClient {
	return InitRedis(cfg.REDIS_HOST, cfg.REDIS_PASSWORD, cfg.REDIS_DB, cfg.REDIS_KEY)
}

// InitRedis initializes the Redis client
func InitRedis(addr string, password string, db int, saveKey string) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		Protocol: 2, // Use RESP2 protocol for better compatibility with older Redis servers
	})
	ctx := context.Background()

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return &RedisClient{
		client: rdb,
		ctx:    ctx,
		saveKey:  saveKey,
	}
}

// SaveStruct saves any struct as JSON in Redis
func (r *RedisClient) SaveStruct(value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal struct: %w", err)
	}

	err = r.client.Set(r.ctx, r.saveKey, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to save data to redis: %w", err)
	}

	return nil
}

// LoadStruct loads a JSON object from Redis into the given struct
func (r *RedisClient) LoadStruct(dest any) error {
	data, err := r.client.Get(r.ctx, r.saveKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key does not exist")
		}
		return fmt.Errorf("failed to get data from redis: %w", err)
	}

	err = json.Unmarshal(data, dest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}