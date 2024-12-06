package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	RegistrationDate  time.Time `json:"registration_date"`
	BirthDate         time.Time `json:"birth_date"`
	Premium           bool      `json:"premium"`
	PremiumExpiration time.Time `json:"premium_expiration"`
}
