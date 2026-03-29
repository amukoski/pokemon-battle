package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/pokemon-battle/backend/internal/model"
	"github.com/redis/go-redis/v9"
)

const (
	pokemonKeyPrefix = "pokemon:"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisCache(addr string, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{Addr: addr})
	return &RedisCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RedisCache) Get(ctx context.Context, name string) (*model.Pokemon, error) {
	key := pokemonKeyPrefix + name

	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		slog.WarnContext(ctx, "redis GET failed for key", "key", key, "error", err)
		return nil, nil
	}

	var p model.Pokemon
	if err = json.Unmarshal(data, &p); err != nil {
		slog.WarnContext(ctx, "redis unmarshal failed for key", "key", key, "error", err)
		return nil, nil
	}

	return &p, nil
}

func (c *RedisCache) Set(ctx context.Context, name string, p *model.Pokemon) error {
	key := pokemonKeyPrefix + name

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshalling pokemon for cache: %w", err)
	}

	if err = c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		slog.WarnContext(ctx, "redis SET failed for key", "key", key, "error", err)
		return nil
	}

	return nil
}

func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}
