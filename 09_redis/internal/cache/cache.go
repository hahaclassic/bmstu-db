package cache

import (
	"context"
	"errors"

	"github.com/hahaclassic/databases/09_redis/internal/models"
)

var (
	ErrGetData    = errors.New("cache: failed to get data")
	ErrSetData    = errors.New("cache: failed to set data")
	ErrDeleteData = errors.New("cache: failed to delete data")
	ErrCacheMiss  = errors.New("cache: cache miss")

	ErrConnection = errors.New("cache: failed to connect")
)

const (
	Key = "top_10_most_streamed_tracks"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]*models.Track, error)
	Set(ctx context.Context, key string, tracks []*models.Track) error
	Delete(ctx context.Context, key string) error
}
