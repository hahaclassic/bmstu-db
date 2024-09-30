package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `fake:"-"`
	Name              string    `fake:"{firstname}"`
	RegistrationDate  time.Time `fake:"-"`
	BirthDate         time.Time `fake:"-"`
	Premium           bool      `fake:"-"`
	PremiumExpiration time.Time `fake:"-"`
}
