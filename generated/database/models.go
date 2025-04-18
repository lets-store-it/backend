// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type AppApiToken struct {
	ID        pgtype.UUID
	OrgID     pgtype.UUID
	Name      string
	Token     string
	CreatedAt pgtype.Timestamp
	RevokedAt pgtype.Timestamp
}

type AppObjectChange struct {
	ID               pgtype.UUID
	OrgID            pgtype.UUID
	UserID           pgtype.UUID
	Action           string
	Time             pgtype.Timestamp
	TargetObjectType int32
	TargetObjectID   pgtype.UUID
	PrechangeState   []byte
	PostchangeState  []byte
}

type AppRole struct {
	ID          int32
	Name        string
	DisplayName string
	Description pgtype.Text
}

type AppRoleBinding struct {
	ID     pgtype.UUID
	OrgID  pgtype.UUID
	RoleID int32
	UserID pgtype.UUID
}

type AppUser struct {
	ID         pgtype.UUID
	Email      string
	FirstName  string
	LastName   string
	MiddleName pgtype.Text
	YandexID   pgtype.Text
	CreatedAt  pgtype.Timestamp
}

type AppUserSession struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	Token     string
	CreatedAt pgtype.Timestamp
	ExpiresAt pgtype.Timestamp
	RevokedAt pgtype.Timestamp
}

type Cell struct {
	ID           pgtype.UUID
	OrgID        pgtype.UUID
	CellsGroupID pgtype.UUID
	Alias        string
	Row          int32
	Level        int32
	Position     int32
	CreatedAt    pgtype.Timestamp
	DeletedAt    pgtype.Timestamp
}

type CellsGroup struct {
	ID             pgtype.UUID
	OrgID          pgtype.UUID
	StorageGroupID pgtype.UUID
	Name           string
	Alias          string
	CreatedAt      pgtype.Timestamp
	DeletedAt      pgtype.Timestamp
}

type Item struct {
	ID          pgtype.UUID
	OrgID       pgtype.UUID
	Name        string
	Description pgtype.Text
	Width       pgtype.Int4
	Depth       pgtype.Int4
	Height      pgtype.Int4
	Weight      pgtype.Int4
	CreatedAt   pgtype.Timestamp
	DeletedAt   pgtype.Timestamp
}

type ItemInstance struct {
	ID                    pgtype.UUID
	OrgID                 pgtype.UUID
	ItemID                pgtype.UUID
	VariantID             pgtype.UUID
	CellID                pgtype.UUID
	Status                string
	AffectedByOperationID pgtype.UUID
	CreatedAt             pgtype.Timestamp
	DeletedAt             pgtype.Timestamp
}

type ItemVariant struct {
	ID        pgtype.UUID
	OrgID     pgtype.UUID
	ItemID    pgtype.UUID
	Name      string
	Article   pgtype.Text
	Ean13     pgtype.Int4
	CreatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type ObjectType struct {
	ID          int32
	ObjectGroup string
	ObjectName  string
}

type Org struct {
	ID        pgtype.UUID
	Name      string
	Subdomain string
	CreatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type OrgUnit struct {
	ID        pgtype.UUID
	OrgID     pgtype.UUID
	Name      string
	Alias     string
	Address   pgtype.Text
	CreatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type StorageGroup struct {
	ID          pgtype.UUID
	OrgID       pgtype.UUID
	UnitID      pgtype.UUID
	ParentID    pgtype.UUID
	Name        string
	Alias       string
	Description pgtype.Text
	CreatedAt   pgtype.Timestamp
	DeletedAt   pgtype.Timestamp
}
