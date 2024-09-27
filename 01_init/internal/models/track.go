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
