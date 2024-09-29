package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `faker:"-"`
	Name              string    `faker:"name"`
	RegistrationDate  time.Time `faker:"-"`
	BirthDate         time.Time `faker:"-"`
	Premium           bool      `faker:"boolean"`
	PremiumExpiration time.Time `faker:"-"`
}
