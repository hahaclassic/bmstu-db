package models

import (
	"github.com/google/uuid"
)

type Playlist struct {
	ID          uuid.UUID
	Title       string
	Description string
	Private     bool
	Rating      int
}
