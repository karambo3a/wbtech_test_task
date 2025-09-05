package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/karambo3a/wbtech_test_task/internal/model"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=redis_cache.go -destination=../../test/mocks/redis_cache_mock.go

type RedisCache interface {
	Init(orders []*model.Order) error
	Get(ctx context.Context, key string) (*model.Order, error)
	Set(ctx context.Context, key string, value *model.Order, expiration time.Duration) error
}

type RedisCacheImpl struct {
	client  *redis.Client
	maxSize int64
}

func NewRedisCache(maxMemoryMB int) *RedisCacheImpl {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS"),
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	maxMemory := fmt.Sprintf("%dmb", maxMemoryMB)
	client.ConfigSet(ctx, "maxmemory", maxMemory)
	client.ConfigSet(ctx, "maxmemory-policy", "allkeys-lru")

	return &RedisCacheImpl{
		client:  client,
		maxSize: int64(maxMemoryMB * 1024 * 1024),
	}
}

func (rc *RedisCacheImpl) Init(orders []*model.Order) error {
	for _, order := range orders {
		if err := rc.Set(context.TODO(), order.OrderUID, order, 24*time.Hour); err != nil {
			continue
		}
	}

	log.Println("cache initialized")
	return nil
}

func (rc *RedisCacheImpl) Get(ctx context.Context, key string) (*model.Order, error) {
	bytes, err := rc.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return &model.Order{}, fmt.Errorf("cache miss for key=%s: %w", key, err)
	} else if err != nil {
		return &model.Order{}, fmt.Errorf("failed to get value by key=%s: %w", key, err)
	}

	var order model.Order
	if err = json.Unmarshal(bytes, &order); err != nil {
		return &model.Order{}, fmt.Errorf("failed to parse json: %w", err)
	}

	log.Printf("got order_uid=%s from cache\n", key)
	return &order, nil
}

func (rc *RedisCacheImpl) Set(ctx context.Context, key string, value *model.Order, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to create json: %w", err)
	}
	if err = rc.client.Set(ctx, key, bytes, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set data: %w", err)
	}

	log.Printf("set order_uid=%s in cache\n", key)
	return nil
}
