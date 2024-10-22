package models

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID           uuid.UUID `fake:"-"`
	Name         string    `fake:"{sentence:3}"`
	OrderInAlbum int       `fake:"-"`
	AlbumID      uuid.UUID `fake:"-"`
	Explicit     bool      `fake:"{bool}"`
	Duration     int       `fake:"{number:180,300}"` // duration in seconds
	Genre        string    `fake:"-"`
	StreamCount  int64     `fake:"{number:0,5000000}"`
}

type PlaylistTrack struct {
	ID         uuid.UUID
	PlaylistID uuid.UUID // It can be uuid of playlist or album
	TrackOrder int       // If TrackOrder == -1, Track goes to last position (working only postgresql)
	DateAdded  time.Time
}
