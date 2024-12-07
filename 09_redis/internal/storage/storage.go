package storage

import (
	"context"
	"errors"

	"github.com/hahaclassic/databases/09_redis/internal/models"
)

var (
	ErrTop10MostStreamedTracks  = errors.New("failed to get top10 most streamed tracks")
	ErrAddTracks                = errors.New("failed to add track")
	ErrDeleteRandomTrack        = errors.New("failed to delete random track")
	ErrUpdateRandomTrackStreams = errors.New("failed to update random track streams")

	ErrTableAlreadyExists = errors.New("table already exists")
	ErrStorageConnection  = errors.New("storage: can't connect to the database")
	ErrNoRowsAffected     = errors.New("no rows affected")
	ErrNotFound           = errors.New("not found")
)

type Storage interface {
	Top10MostStreamedTracks(ctx context.Context) ([]*models.Track, error)
	AddTrack(ctx context.Context, track *models.Track) error
	DeleteRandomTrack(ctx context.Context) error
	UpdateRandomTrackStreams(ctx context.Context) error
}
