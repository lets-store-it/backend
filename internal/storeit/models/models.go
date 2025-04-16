package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        uuid.UUID
	Name      string
	Subdomain string

	CreatedAt time.Time
	DeletedAt *time.Time
}

type OrganizationUnit struct {
	ID      uuid.UUID
	OrgID   uuid.UUID
	Name    string
	Alias   string
	Address *string

	CreatedAt time.Time
	DeletedAt *time.Time
}

type StorageGroup struct {
	ID       uuid.UUID
	UnitID   uuid.UUID
	ParentID *uuid.UUID
	Name     string
	Alias    string

	CreatedAt time.Time
	DeletedAt *time.Time
}

type Item struct {
	ID          uuid.UUID
	OrgID       uuid.UUID
	Name        string
	Description *string

	Variants  *[]ItemVariant
	CreatedAt time.Time
	DeletedAt *time.Time
}

type ItemVariant struct {
	ID     uuid.UUID
	ItemID uuid.UUID
	OrgID  uuid.UUID

	Name string

	Article *string
	EAN13   *int

	CreatedAt time.Time
	DeletedAt *time.Time
}
