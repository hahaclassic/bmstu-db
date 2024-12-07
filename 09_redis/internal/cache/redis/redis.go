package cacheredis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/hahaclassic/databases/09_redis/config"
	"github.com/hahaclassic/databases/09_redis/internal/cache"
	"github.com/hahaclassic/databases/09_redis/internal/models"
	"github.com/redis/go-redis/v9"
)

type CacheRedis struct {
	config *config.RedisConfig
	client *redis.Client
}

func New(ctx context.Context, cfg *config.RedisConfig) (*CacheRedis, error) {
	cli := redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(cfg.Host, cfg.Port),
	})

	if err := cli.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", cache.ErrConnection, err)
	}

	return &CacheRedis{
		client: cli,
		config: cfg,
	}, nil
}

func (c *CacheRedis) Get(ctx context.Context, key string) ([]*models.Track, error) {
	var tracks []*models.Track
	data, err := c.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("%w: %w", cache.ErrSetData, cache.ErrCacheMiss)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", cache.ErrGetData, err)
	}

	err = json.Unmarshal(data, &tracks)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", cache.ErrGetData, err)
	}

	return tracks, nil
}

func (c *CacheRedis) Set(ctx context.Context, key string, tracks []*models.Track) error {
	data, err := json.Marshal(tracks)
	if err != nil {
		return fmt.Errorf("%w: %w", cache.ErrSetData, err)
	}

	err = c.client.Set(ctx, key, data, c.config.Expiration).Err()
	if err != nil {
		return fmt.Errorf("%w: %w", cache.ErrSetData, err)
	}

	return nil
}

func (c *CacheRedis) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("%w: %w", cache.ErrDeleteData, err)
	}

	return nil
}
