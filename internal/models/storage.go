package models

import (
	"time"

	"github.com/google/uuid"
)

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

type CellsGroup struct {
	ID             uuid.UUID  `json:"id"`
	OrgID          uuid.UUID  `json:"org_id"`
	UnitID         uuid.UUID  `json:"unit_id"`
	StorageGroupID *uuid.UUID `json:"storage_group_id"`
	Name           string     `json:"name"`
	Alias          string     `json:"alias"`

	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
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
