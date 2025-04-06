// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type AppUser struct {
	ID       pgtype.UUID
	Fullname pgtype.Text
	Email    string
}

type Cell struct {
	ID        pgtype.UUID
	ShortName string
	GroupID   pgtype.UUID
	KindID    pgtype.UUID
	Rack      string
	Level     int32
	Position  int32
}

type CellGroup struct {
	ID        pgtype.UUID
	SpaceID   pgtype.UUID
	Name      string
	ShortName pgtype.Text
}

type CellKind struct {
	ID        pgtype.UUID
	OrgID     pgtype.UUID
	Name      string
	Height    int32
	Width     int32
	Depth     int32
	MaxWeight int32
}

type Employee struct {
	UserID pgtype.UUID
	OrgID  pgtype.UUID
}

type Item struct {
	ID    pgtype.UUID
	OrgID pgtype.UUID
	Name  string
	Ean   string
}

type ItemInstance struct {
	ID                    pgtype.UUID
	TrackingID            string
	ItemID                pgtype.UUID
	VariantID             pgtype.UUID
	CellID                pgtype.UUID
	Status                string
	AffectedByOperationID pgtype.UUID
}

type ItemProperty struct {
	ID     pgtype.UUID
	ItemID pgtype.UUID
	Name   string
	Type   string
	Value  string
}

type ItemVariant struct {
	ID     pgtype.UUID
	ItemID pgtype.UUID
	Name   string
	Width  pgtype.Int4
	Depth  pgtype.Int4
	Height pgtype.Int4
	Weight pgtype.Int4
}

type ItemVariantProperty struct {
	VariantID  pgtype.UUID
	PropertyID pgtype.UUID
	Value      string
}

type Operation struct {
	ID         pgtype.UUID
	Type       string
	AssignedTo pgtype.UUID
	CreatedAt  pgtype.Timestamp
}

type OperationItem struct {
	OperationID       pgtype.UUID
	ItemInstanceID    pgtype.UUID
	Status            string
	OriginCellID      pgtype.UUID
	DestinationCellID pgtype.UUID
}

type Org struct {
	ID        pgtype.UUID
	Name      string
	Subdomain string
	IsDeleted bool
	CreatedAt pgtype.Timestamp
	CreatedBy pgtype.UUID
	UpdatedAt pgtype.Timestamp
	UpdatedBy pgtype.UUID
}

type OrgUnit struct {
	ID      pgtype.UUID
	OrgID   pgtype.UUID
	Name    string
	Address pgtype.Text
}

type Role struct {
	ID          pgtype.UUID
	Name        string
	DisplayName string
}

type RoleBinding struct {
	ID         pgtype.UUID
	RoleID     pgtype.UUID
	EmployeeID pgtype.UUID
}

type RolePermission struct {
	RoleID     pgtype.UUID
	Permission string
}

type StorageSpace struct {
	ID        pgtype.UUID
	UnitID    pgtype.UUID
	Name      string
	ShortName pgtype.Text
}
