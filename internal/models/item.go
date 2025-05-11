package models

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID `json:"id"`
	OrgID       uuid.UUID `json:"org_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`

	Variants  []*ItemVariant  `json:"variants"`
	Instances []*ItemInstance `json:"instances"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type ItemVariant struct {
	ID     uuid.UUID `json:"id"`
	ItemID uuid.UUID `json:"item_id"`
	OrgID  uuid.UUID `json:"org_id"`

	Name string `json:"name"`

	Article *string `json:"article"`
	EAN13   *int32  `json:"ean13"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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
	CellID                *uuid.UUID         `json:"cell_id"`
	Status                ItemInstanceStatus `json:"status"`
	AffectedByOperationID *uuid.UUID          `json:"affected_by_operation_id"`

	Item    *Item        `json:"item"`
	Cell    *Cell        `json:"cell"`
	Variant *ItemVariant `json:"variant"`
}
