package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `faker:"-"`
	Name              string    `faker:"name"`
	RegistrationDate  time.Time
	BirthDate         time.Time
	Premium           bool `faker:"boolean"`
	PremiumExpiration time.Time
}
