package models

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID        uuid.UUID
	Name      string
	Genre     string
	Country   string
	DebutYear time.Time
}
