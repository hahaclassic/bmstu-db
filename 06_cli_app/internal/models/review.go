package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	AlbumID    uuid.UUID
	Rating     int
	Comment    string
	ReviewDate time.Time
}
