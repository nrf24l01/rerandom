package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/nrf24l01/rerandom/backend/core"
	"github.com/redis/go-redis/v9"
)

// RedisClient обёртка для Redis
type RedisClient struct {
    client       *redis.Client
    ctx          context.Context
    keysSetName  string
}

func CreateRedisFromCFG(cfg *core.Config) (*RedisClient, error) {
	return NewRedisClient(
		fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		cfg.RedisPassword,
		cfg.RedisDB,
		cfg.KeysSetName,
	), nil
}

// NewRedisClient создаёт подключение к Redis
func NewRedisClient(addr string, password string, db int, keysSetName string) *RedisClient {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    ctx := context.Background()

    // Проверка соединения
    if err := rdb.Ping(ctx).Err(); err != nil {
        log.Fatalf("Не удалось подключиться к Redis: %v", err)
    }

    return &RedisClient{
        client:      rdb,
        ctx:         ctx,
        keysSetName: keysSetName,
    }
}

func (r *RedisClient) Close() error {
    return r.client.Close()
}

func (r *RedisClient) GetFirstFromSet() (string, error) {
    val, err := r.client.SPop(r.ctx, r.keysSetName).Result()
    if err == redis.Nil {
        return "", fmt.Errorf("множество %s пусто", r.keysSetName)
    }
    return val, err
}

func (r *RedisClient) AddToSet(value string) error {
    return r.client.SAdd(r.ctx, r.keysSetName, value).Err()
}