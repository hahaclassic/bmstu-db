package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/databases/06_cli_app/internal/models"
)

var (
	ErrCountArtists             = errors.New("failed to count artists")
	ErrTracksAdditionalInfo     = errors.New("failed to get additional info")
	ErrRankedTracks             = errors.New("failed to get ranked tracks")
	ErrTablesNames              = errors.New("failed to get tables names")
	ErrCountArtistTracks        = errors.New("failed to count artist tracks")
	ErrAlbumsInfo               = errors.New("failed to get artist albums")
	ErrFreePremiumForPensioners = errors.New("failed to set premium")
	ErrCurrentUser              = errors.New("failed to get current user")
	ErrCreateTable              = errors.New("failed to create table")
	ErrAddReview                = errors.New("failed to add review")

	ErrTableAlreadyExists = errors.New("table already exists")
	ErrStorageConnection  = errors.New("storage: can't connect to the database")
	ErrNoRowsAffected     = errors.New("no rows affected")
	ErrNotFound           = errors.New("not found")
)

type MusicServiceStorage interface {
	CountArtists(ctx context.Context) (int, error)
	TracksAdditionalInfo(ctx context.Context, limit int) (tracks []*models.TracksAdditionalInfo, err error)
	BestTracks(ctx context.Context, limit int) (tracks []*models.RankedTrack, err error)
	TablesNames(ctx context.Context) (names []string, err error)
	CountArtistTracks(ctx context.Context, artistID uuid.UUID) (int, error)
	ArtistAlbums(ctx context.Context, artistID uuid.UUID) (tracks []*models.Album, err error)
	FreePremiumForPensioners(ctx context.Context) error
	CurrentUser(ctx context.Context) (string, error)
	CreateReviewsTable(ctx context.Context) error
	AddReview(ctx context.Context, review *models.Review) error

	Close()
}
