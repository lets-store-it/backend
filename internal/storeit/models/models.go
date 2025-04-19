package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	Secret string    `json:"secret"`
}

type User struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName *string   `json:"middle_name"`
	YandexID   *string   `json:"yandex_id"`
}

type Organization struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Subdomain string    `json:"subdomain"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type OrganizationUnit struct {
	ID      uuid.UUID `json:"id"`
	OrgID   uuid.UUID `json:"org_id"`
	Name    string    `json:"name"`
	Alias   string    `json:"alias"`
	Address *string   `json:"address"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type StorageGroup struct {
	ID       uuid.UUID  `json:"id"`
	UnitID   uuid.UUID  `json:"unit_id"`
	ParentID *uuid.UUID `json:"parent_id"`
	Name     string     `json:"name"`
	Alias    string     `json:"alias"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type Item struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	Variants  *[]ItemVariant `json:"variants"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
}

type ItemVariant struct {
	ID     uuid.UUID `json:"id"`
	ItemID uuid.UUID `json:"item_id"`
	OrgID  uuid.UUID `json:"org_id"`

	Name string `json:"name"`

	Article *string `json:"article"`
	EAN13   *int    `json:"ean13"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CellsGroup struct {
	ID             uuid.UUID `json:"id"`
	OrgID          uuid.UUID `json:"org_id"`
	StorageGroupID uuid.UUID `json:"storage_group_id"`
	Name           string    `json:"name"`
	Alias          string    `json:"alias"`
}

type Cell struct {
	ID           uuid.UUID `json:"id"`
	OrgID        uuid.UUID `json:"org_id"`
	CellsGroupID uuid.UUID `json:"cells_group_id"`
	Alias        string    `json:"alias"`
	Row          int       `json:"row"`
	Level        int       `json:"level"`
	Position     int       `json:"position"`
}

type ObjectChangeAction string

const (
	ObjectChangeActionCreate ObjectChangeAction = "create"
	ObjectChangeActionUpdate ObjectChangeAction = "update"
	ObjectChangeActionDelete ObjectChangeAction = "delete"
)

type ObjectType int

const (
	ObjectTypeOrganization ObjectType = 1
	ObjectTypeUnit         ObjectType = 2
	ObjectTypeStorageGroup ObjectType = 3
	ObjectTypeCellsGroup   ObjectType = 4
	ObjectTypeCell         ObjectType = 5
	ObjectTypeItem         ObjectType = 6
	ObjectTypeItemInstance ObjectType = 7
	ObjectTypeUserRoles    ObjectType = 8
)

type ObjectChange struct {
	ID               uuid.UUID          `json:"id"`
	OrgID            uuid.UUID          `json:"org_id"`
	UserID           uuid.UUID          `json:"user_id"`
	Action           ObjectChangeAction `json:"action"`
	TargetObjectType ObjectType         `json:"target_object_type"`
	TargetObjectID   uuid.UUID          `json:"target_object_id"`
	PrechangeState   json.RawMessage    `json:"prechange_state"`
	PostchangeState  json.RawMessage    `json:"postchange_state"`
}

type ItemInstanceStatus string

const (
	ItemInstanceStatusAvailable ItemInstanceStatus = "available"
	ItemInstanceStatusReserved  ItemInstanceStatus = "reserved"
	ItemInstanceStatusConsumed  ItemInstanceStatus = "consumed"
)

type ItemInstance struct {
	ID                    uuid.UUID          `json:"id"`
	OrgID                 uuid.UUID          `json:"org_id"`
	ItemID                uuid.UUID          `json:"item_id"`
	VariantID             uuid.UUID          `json:"variant_id"`
	CellID                uuid.UUID          `json:"cell_id"`
	Status                ItemInstanceStatus `json:"status"`
	AffectedByOperationID uuid.UUID          `json:"affected_by_operation_id"`
}

type ApiToken struct {
	ID        uuid.UUID  `json:"id"`
	OrgID     uuid.UUID  `json:"org_id"`
	Name      string     `json:"name"`
	Token     string     `json:"token"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}
