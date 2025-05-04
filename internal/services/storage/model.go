package storage

import (
	"errors"

	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

func toStorageGroup(group database.StorageGroup) (*models.StorageGroup, error) {
	id := utils.UuidFromPgx(group.ID)
	if id == nil {
		return nil, errors.New("failed to convert storage group")
	}
	unitID := utils.UuidFromPgx(group.UnitID)
	if unitID == nil {
		return nil, errors.New("failed to convert storage group")
	}

	return &models.StorageGroup{
		ID:       *id,
		UnitID:   *unitID,
		ParentID: utils.UuidFromPgx(group.ParentID),
		Name:     group.Name,
		Alias:    group.Alias,
	}, nil
}

func toCellsGroup(group database.CellsGroup) (*models.CellsGroup, error) {
	id := utils.UuidFromPgx(group.ID)
	if id == nil {
		return nil, errors.New("failed to convert cells group")
	}
	storageGroupID := utils.UuidFromPgx(group.StorageGroupID)
	if storageGroupID == nil {
		return nil, errors.New("failed to convert cells group")
	}
	orgID := utils.UuidFromPgx(group.OrgID)
	if orgID == nil {
		return nil, errors.New("failed to convert cells group")
	}

	return &models.CellsGroup{
		ID:             *id,
		OrgID:          *orgID,
		StorageGroupID: *storageGroupID,
		Name:           group.Name,
		Alias:          group.Alias,
	}, nil
}

func toCell(cell database.Cell) (*models.Cell, error) {
	id := utils.UuidFromPgx(cell.ID)
	if id == nil {
		return nil, errors.New("failed to convert cell")
	}
	cellsGroupID := utils.UuidFromPgx(cell.CellsGroupID)
	if cellsGroupID == nil {
		return nil, errors.New("failed to convert cell")
	}

	return &models.Cell{
		ID:           *id,
		CellsGroupID: *cellsGroupID,
		Alias:        cell.Alias,
		Row:          int(cell.Row),
		Level:        int(cell.Level),
		Position:     int(cell.Position),
	}, nil
}
