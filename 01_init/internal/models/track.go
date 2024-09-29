package models

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID           uuid.UUID `faker:"-"`
	Name         string    `faker:"sentence:3"`
	Explicit     bool      `faker:"boolean"`
	Duration     int       `faker:"number:180,300"` // duration in seconds
	Genre        string    `faker:"-"`
	StreamCount  int64     `faker:"number:1000,50000"`
	OrderInAlbum int       `faker:"-"`
	AlbumID      uuid.UUID `faker:"-"`
}

type PlaylistTrack struct {
	ID         uuid.UUID
	PlaylistID uuid.UUID // It can be uuid of playlist or album
	TrackOrder int       // If TrackOrder == -1, Track goes to last position (working only postgresql)
	DateAdded  time.Time
}
