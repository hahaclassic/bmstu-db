package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID     `json:"id"`
	Name              string        `json:"name"`
	RegistrationDate  time.Time     `json:"registration_date"`
	BirthDate         JsonBirthDate `json:"birth_date"`
	Premium           bool          `json:"premium"`
	PremiumExpiration time.Time     `json:"premium_expiration"`
}

type JsonBirthDate time.Time

func (j *JsonBirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonBirthDate(t)
	return nil
}

func (j JsonBirthDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func (j JsonBirthDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}
