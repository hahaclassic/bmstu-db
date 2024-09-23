package models

import (
	"time"

	"github.com/google/uuid"
)

type Track struct {
	ID          uuid.UUID
	Name        string
	Explicit    bool
	Duration    time.Duration
	Genre       string
	StreamCount int64
}
