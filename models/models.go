package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	FullName string
	Email    string
}

type Organization struct {
	ID        uuid.UUID
	Name      string
	Subdomain string
}

type OrganizationUnit struct {
	ID      uuid.UUID
	OrgID   uuid.UUID
	Name    string
	Address string
}

type StorageGroup struct {
	ID         uuid.UUID
	UnitID     uuid.UUID
	ParentID   uuid.UUID
	Name       string
	ShortAlias string
}

type CellKind struct {
	ID                   uuid.UUID
	OrgID                uuid.UUID
	Name                 string
	Height, Width, Depth int
	MaxWeight            int
}

type CellGroup struct {
	ID         uuid.UUID
	SpaceID    uuid.UUID
	Name       string
	ShortAlias string
}

type Cell struct {
	ID              uuid.UUID
	ShortAlias      string
	GroupID         uuid.UUID
	KindID          uuid.UUID
	Rack            string
	Level, Position int
}
