package models

import (
	"time"

	"github.com/google/uuid"
)

type TvBoard struct {
	ID     uuid.UUID `json:"id"`
	OrgID  uuid.UUID `json:"org_id"`
	Name   string    `json:"name"`
	UnitID uuid.UUID `json:"unit_id"`
	Token  string    `json:"token"`

	Unit   *OrganizationUnit `json:"unit"`
	
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
