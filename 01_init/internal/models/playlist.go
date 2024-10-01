package models

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID `fake:"-"`
	Title       string    `fake:"{sentence:3}"`
	Description string    `fake:"{sentence:6}"`
	Private     bool      `fake:"-"`
	Rating      int       `fake:"{number:1,25}"`
	LastUpdated time.Time `faker:"-"`
}

type UserPlaylist struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	IsFavorite  bool
	AccessLevel AccessLevel
}
