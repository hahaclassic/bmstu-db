package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/databases/01_init/internal/models"
)

var (
	ErrStorageConnection = errors.New("storage: can't connect to the database")
	ErrNoRowsAffected    = errors.New("no rows affected")
	ErrNotFound          = errors.New("not found")

	ErrCreateArtist   = errors.New("failed to create artist")
	ErrCreateAlbum    = errors.New("failed to create album")
	ErrCreateTrack    = errors.New("failed to create track")
	ErrCreatePlaylist = errors.New("failed to create playlist")
	ErrCreateUser     = errors.New("failed to create user")

	ErrAddPlaylist        = errors.New("failed to add playlist to user")
	ErrAddTrackToPlaylist = errors.New("failed to add track to playlist")
	ErrAddTrackToAlbum    = errors.New("failed to add track to album")
	ErrAddArtistTrack     = errors.New("failed to add artist track")
	ErrAddArtistAlbum     = errors.New("failed to add artist album")

	ErrDeleteAll = errors.New("failed to delete all records")
)

type MusicServiceStorage interface {
	CreateArtist(ctx context.Context, artist *models.Artist) error
	CreateAlbum(ctx context.Context, album *models.Album) error
	CreateTrack(ctx context.Context, track *models.Track) error
	CreatePlaylist(ctx context.Context, playlist *models.Playlist) error
	CreateUser(ctx context.Context, user *models.User) error

	AddPlaylist(ctx context.Context, userPlaylist *models.UserPlaylist) error
	AddTrackToPlaylist(ctx context.Context, track *models.PlaylistTrack) error
	AddArtistTrack(ctx context.Context, trackID uuid.UUID, artistID uuid.UUID) error
	AddArtistAlbum(ctx context.Context, albumID uuid.UUID, artistID uuid.UUID) error

	DeleteAll(ctx context.Context) error

	Close()
}
