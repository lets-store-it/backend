package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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
	ObjectTypeEmployee     ObjectTypeId = 8
	ObjectTypeTask         ObjectTypeId = 9
	ObjectTypeItemVariant  ObjectTypeId = 10
	ObjectTypeApiToken     ObjectTypeId = 11
)

type ObjectType struct {
	ID    ObjectTypeId `json:"id"`
	Group string       `json:"group"`
	Name  string       `json:"name"`
}

type ObjectChangeCreate struct {
	Action ObjectChangeAction `json:"action"`

	TargetObjectType ObjectTypeId `json:"target_object_type_id"`
	TargetObjectID   uuid.UUID    `json:"target_object_id"`

	PrechangeState  any `json:"prechange_state"`
	PostchangeState any `json:"postchange_state"`
}

type ObjectChange struct {
	ID     uuid.UUID  `json:"id"`
	OrgID  uuid.UUID  `json:"org_id"`
	UserID *uuid.UUID `json:"user_id"`

	Action ObjectChangeAction `json:"action"`

	TargetObjectType ObjectTypeId `json:"target_object_type_id"`
	TargetObjectID   uuid.UUID    `json:"target_object_id"`

	PrechangeState  json.RawMessage `json:"prechange_state"`
	PostchangeState json.RawMessage `json:"postchange_state"`

	Timestamp time.Time `json:"timestamp"`

	Employee   *Employee   `json:"employee"`
	ObjectType *ObjectType `json:"object_type"`
}
