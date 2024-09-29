package models

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID `faker:"-"`
	Title       string    `faker:"sentence:3"`
	Description string    `faker:"sentence:5"`
	Private     bool      `faker:"boolean"`
	Rating      int       `faker:"number:1,5"`
	LastUpdated time.Time
}

type UserPlaylist struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	IsFavorite  bool
	AccessLevel int
}
