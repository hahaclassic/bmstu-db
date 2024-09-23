package models

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID
	Title       string
	ReleaseDate time.Time
	Label       string
	Genre       string
}
