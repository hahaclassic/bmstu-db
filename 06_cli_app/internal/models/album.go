package models

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID `fake:"-"`
	Title       string    `fake:"{sentence:3}"`
	ReleaseDate time.Time `fake:"-"`
	Label       string    `fake:"{sentence:1}"`
	Genre       string    `fake:"-"`
}
