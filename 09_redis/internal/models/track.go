package models

import (
	"github.com/google/uuid"
)

type Track struct {
	ID          uuid.UUID `fake:"-"`
	Name        string    `fake:"{sentence:3}"`
	Album       string    `fake:"{sentence:2}"`
	Explicit    bool      `fake:"{bool}"`
	Duration    int       `fake:"{number:180,300}"`
	Genre       string    `fake:"-"`
	StreamCount int       `fake:"{number:0,5000000}"`
}
