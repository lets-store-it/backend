package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID uuid.UUID `json:"id"`

	Name      string `json:"name"`
	Subdomain string `json:"subdomain"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type OrganizationUnit struct {
	ID    uuid.UUID `json:"id"`
	OrgID uuid.UUID `json:"org_id"`

	Name    string  `json:"name"`
	Alias   string  `json:"alias"`
	Address *string `json:"address"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
