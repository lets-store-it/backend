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
	OrgID    uuid.UUID  `json:"org_id"`
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

	Variants  *[]ItemVariant  `json:"variants"`
	Instances *[]ItemInstance `json:"instances"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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
	ID             uuid.UUID  `json:"id"`
	OrgID          uuid.UUID  `json:"org_id"`
	UnitID         uuid.UUID  `json:"unit_id"`
	StorageGroupID *uuid.UUID `json:"storage_group_id"`
	Name           string     `json:"name"`
	Alias          string     `json:"alias"`
}

type CellPathObjectType string

const (
	CellPathObjectTypeCellsGroup   CellPathObjectType = "cells_group"
	CellPathObjectTypeStorageGroup CellPathObjectType = "storage_group"
	CellPathObjectTypeUnit         CellPathObjectType = "unit"
)

type CellPathSegment struct {
	ID         uuid.UUID          `json:"id"`
	Name       string             `json:"name"`
	ObjectType CellPathObjectType `json:"object_type"`
	Alias      string             `json:"alias"`
}

type Cell struct {
	ID           uuid.UUID `json:"id"`
	OrgID        uuid.UUID `json:"org_id"`
	CellsGroupID uuid.UUID `json:"cells_group_id"`
	Alias        string    `json:"alias"`
	Row          int       `json:"row"`
	Level        int       `json:"level"`
	Position     int       `json:"position"`

	Path *[]CellPathSegment `json:"path"`
}

type ObjectChangeAction string

const (
	ObjectChangeActionCreate ObjectChangeAction = "create"
	ObjectChangeActionUpdate ObjectChangeAction = "update"
	ObjectChangeActionDelete ObjectChangeAction = "delete"
)

type ObjectTypeId int

const (
	ObjectTypeOrganization ObjectTypeId = 1
	ObjectTypeUnit         ObjectTypeId = 2
	ObjectTypeStorageGroup ObjectTypeId = 3
	ObjectTypeCellsGroup   ObjectTypeId = 4
	ObjectTypeCell         ObjectTypeId = 5
	ObjectTypeItem         ObjectTypeId = 6
	ObjectTypeItemInstance ObjectTypeId = 7
	ObjectTypeUserRoles    ObjectTypeId = 8
)

type ObjectType struct {
	ID    ObjectTypeId `json:"id"`
	Group string       `json:"group"`
	Name  string       `json:"name"`
}

type ObjectChange struct {
	ID                 uuid.UUID          `json:"id"`
	OrgID              uuid.UUID          `json:"org_id"`
	UserID             uuid.UUID          `json:"user_id"`
	Action             ObjectChangeAction `json:"action"`
	TargetObjectTypeId ObjectTypeId       `json:"target_object_type_id"`
	TargetObjectID     uuid.UUID          `json:"target_object_id"`
	PrechangeState     json.RawMessage    `json:"prechange_state"`
	PostchangeState    json.RawMessage    `json:"postchange_state"`
	Timestamp          time.Time          `json:"timestamp"`

	Employee   *Employee   `json:"employee"`
	ObjectType *ObjectType `json:"object_type"`
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

	Cell    *Cell        `json:"cell"`
	Variant *ItemVariant `json:"variant"`
}

type ApiToken struct {
	ID        uuid.UUID  `json:"id"`
	OrgID     uuid.UUID  `json:"org_id"`
	Name      string     `json:"name"`
	Token     string     `json:"token"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}

type TaskType string

const (
	TaskTypePickItem TaskType = "pick"
	TaskTypeMovement TaskType = "movement"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

type TaskItem struct {
	ID           uuid.UUID `json:"id"`
	InstanceID   uuid.UUID `json:"instance_id"`
	TargetCellID uuid.UUID `json:"target_cell_id"`

	Instance   *ItemInstance `json:"instance"`
	TargetCell *Cell         `json:"target_cell"`
}

type Task struct {
	ID          uuid.UUID  `json:"id"`
	OrgID       uuid.UUID  `json:"org_id"`
	UnitID      uuid.UUID  `json:"unit_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      TaskStatus `json:"status"`
	Type        TaskType   `json:"type"`
	Items       []TaskItem `json:"items"`
}

type RoleName string

const (
	RoleOwner   RoleName = "org_owner"
	RoleAdmin   RoleName = "org_admin"
	RoleManager RoleName = "org_manager"
	RoleWorker  RoleName = "org_worker"
)

type Role struct {
	ID          int
	Name        RoleName
	DisplayName string
	Description string
}

type Employee struct {
	UserID     uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	MiddleName *string   `json:"middle_name"`
	RoleID     int       `json:"role_id"`

	Role *Role `json:"role"`
}
