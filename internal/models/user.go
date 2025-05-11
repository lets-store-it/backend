package models

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `json:"id"`

	Email      string  `json:"email"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	MiddleName *string `json:"middle_name"`

	YandexID *string `json:"yandex_id"`
}
