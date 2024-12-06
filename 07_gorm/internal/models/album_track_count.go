package models

import "github.com/google/uuid"

type AlbumTrackCount struct {
	AlbumID    uuid.UUID
	Title      string
	TrackCount int
}
