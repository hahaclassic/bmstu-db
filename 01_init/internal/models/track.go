package models

import (
	"github.com/google/uuid"
)

type Track struct {
	ID          uuid.UUID `faker:"-"`
	Name        string    `faker:"sentence:3"`
	Explicit    bool      `faker:"boolean"`
	Duration    int       `faker:"number:180,300"` // duration in seconds
	Genre       string    `faker:"-"`
	StreamCount int64     `faker:"number:1000,50000"`
}

type ListTrack struct {
	ID         uuid.UUID
	ListID     uuid.UUID // It can be uuid of playlist or album
	TrackOrder int       // If TrackOrder == -1, Track goes to last position (working only postgresql)
}
