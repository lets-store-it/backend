package storage

import (
	"errors"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toStorageGroup(group sqlc.StorageGroup) (*models.StorageGroup, error) {
	if !group.ID.Valid {
		return nil, errors.New("failed to convert storage group: invalid ID")
	}
	if !group.UnitID.Valid {
		return nil, errors.New("failed to convert storage group: invalid unit ID")
	}

	var parentID *uuid.UUID
	if group.ParentID.Valid {
		id := database.UuidFromPgx(group.ParentID)
		parentID = &id
	}

	return &models.StorageGroup{
		ID:       database.UuidFromPgx(group.ID),
		UnitID:   database.UuidFromPgx(group.UnitID),
		ParentID: parentID,
		Name:     group.Name,
		Alias:    group.Alias,
		OrgID:    group.OrgID.Bytes,
	}, nil
}

func toCellsGroup(group sqlc.CellsGroup) (*models.CellsGroup, error) {
	if !group.ID.Valid {
		return nil, errors.New("failed to convert cells group: invalid ID")
	}
	if !group.OrgID.Valid {
		return nil, errors.New("failed to convert cells group: invalid org ID")
	}
	if !group.UnitID.Valid {
		return nil, errors.New("failed to convert cells group: invalid unit ID")
	}

	var storageGroupID *uuid.UUID
	if group.StorageGroupID.Valid {
		id := database.UuidFromPgx(group.StorageGroupID)
		storageGroupID = &id
	}

	return &models.CellsGroup{
		ID:             database.UuidFromPgx(group.ID),
		OrgID:          database.UuidFromPgx(group.OrgID),
		UnitID:         database.UuidFromPgx(group.UnitID),
		StorageGroupID: storageGroupID,
		Name:           group.Name,
		Alias:          group.Alias,
	}, nil
}

func toCell(cell sqlc.Cell) (*models.Cell, error) {
	if !cell.ID.Valid {
		return nil, errors.New("failed to convert cell: invalid ID")
	}
	if !cell.CellsGroupID.Valid {
		return nil, errors.New("failed to convert cell: invalid cells group ID")
	}

	return &models.Cell{
		ID:           database.UuidFromPgx(cell.ID),
		CellsGroupID: database.UuidFromPgx(cell.CellsGroupID),
		Alias:        cell.Alias,
		Row:          int(cell.Row),
		Level:        int(cell.Level),
		Position:     int(cell.Position),
	}, nil
}
