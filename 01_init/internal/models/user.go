package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID
	Name              string
	RegistrationDate  time.Time
	BirthDate         time.Time
	Premium           bool
	PremiumExpiration time.Time
}
