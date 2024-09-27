package models

import (
	"github.com/google/uuid"
)

type Artist struct {
	ID        uuid.UUID `faker:"-"`
	Name      string    `faker:"word"`
	Genre     string    `faker:"-"`
	Country   string    `faker:"country"`
	DebutYear int       `faker:"year"`
}
