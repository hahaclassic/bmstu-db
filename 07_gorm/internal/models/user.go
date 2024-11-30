package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `fake:"-"`
	Name              string    `fake:"{username}"`
	RegistrationDate  time.Time `fake:"-"`
	BirthDate         time.Time `fake:"-"`
	Premium           bool      `fake:"{bool}"`
	PremiumExpiration time.Time `fake:"-"`
}
