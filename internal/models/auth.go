package models

import (
	"time"

	"github.com/google/uuid"
)

type UserSession struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Token  string    `json:"secret"`

	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type ApiToken struct {
	ID    uuid.UUID `json:"id"`
	OrgID uuid.UUID `json:"org_id"`
	Name  string    `json:"name"`
	Token string    `json:"token"`

	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}
