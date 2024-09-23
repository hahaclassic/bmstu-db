package storage

import (
	"context"
	"errors"
)

var (
	ErrStorageConnection = errors.New("storage: can't connect to the database")
	ErrNoRowsAffected    = errors.New("no rows affected")
	ErrNotFound          = errors.New("not found")
)

// TODO
type MusicServiceStorage interface {
	AddArtist(ctx context.Context)
	AddAlbum(ctx context.Context)
	AddTrack(ctx context.Context)
	AddPlaylist(ctx context.Context)
	AddUser(ctx context.Context)

	DeleteAll()
}
