package models

import (
	"github.com/google/uuid"
)

type Artist struct {
	ID        uuid.UUID `fake:"-"`
	Name      string    `fake:"-"`
	Genre     string    `fake:"-"`
	Country   string    `fake:"-"`
	DebutYear int       `fake:"-"`
}
