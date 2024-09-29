package models

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID `faker:"-"`
	Title       string    `faker:"sentence:3"`
	ReleaseDate time.Time `faker:"-"`
	Label       string    `faker:"word"`
	Genre       string    `faker:"-"`
}
