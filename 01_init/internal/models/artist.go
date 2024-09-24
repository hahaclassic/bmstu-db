package models

import (
	"github.com/google/uuid"
)

type Artist struct {
	ID        uuid.UUID `faker:"-"`
	Name      string    `faker:"name"`
	Genre     string    `faker:"word"`
	Country   string    `faker:"country"`
	DebutYear int       `faker:"year"`
}
