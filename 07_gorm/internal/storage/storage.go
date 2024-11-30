package storage

import (
	"context"
	"errors"

	"github.com/hahaclassic/databases/07_gorm/internal/models"
)

var (
	ErrExplicitTracks               = errors.New("failed to get explicit tracks")
	ErrCountTracksByGenre           = errors.New("failed to count tracks by genres")
	ErrAlbumsWithMaxTracks          = errors.New("failed to get albums with min num of tracks")
	ErrArtistsWithReleasedAlbumYear = errors.New("failed to get artists")
	ErrUsersOlderThan               = errors.New("failed to get users older than specified")

	ErrTableAlreadyExists = errors.New("table already exists")
	ErrStorageConnection  = errors.New("storage: can't connect to the database")
	ErrNoRowsAffected     = errors.New("no rows affected")
	ErrNotFound           = errors.New("not found")
)

type Storage interface {
	BestExplicitTracks(ctx context.Context, limit int) ([]*models.Track, error)
	CountTracksByGenre(ctx context.Context) ([]*models.GenreCount, error)
	AlbumsWithMaxTracks(ctx context.Context, minNumOfTracks int) ([]*models.Album, error)
	ArtistsWithReleasedAlbumYear(ctx context.Context, year int) ([]*models.Artist, error)
	UsersOlderThan(ctx context.Context, age int) ([]*models.User, error)

	GetTracksByGenre(ctx context.Context, genre string) ([]*models.Track, error)
	GetAlbumsWithTrackCounts(ctx context.Context) ([]*models.AlbumTrackCount, error)

	AddUser(ctx context.Context, user *models.User) error
	UpdateUserName(ctx context.Context, userID int, newName string) error
	DeleteUser(ctx context.Context, userID int) error

	AlbumsByArtist(ctx context.Context, artistID string) ([]*models.Album, error)
}
