package storage

import (
	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/sqlc"
	"github.com/let-store-it/backend/internal/database"
	"github.com/let-store-it/backend/internal/models"
)

func toStorageGroupModel(group sqlc.StorageGroup) *models.StorageGroup {

	var parentID *uuid.UUID
	if group.ParentID.Valid {
		id := database.UUIDFromPgx(group.ParentID)
		parentID = &id
	}

	return &models.StorageGroup{
		ID:       database.UUIDFromPgx(group.ID),
		UnitID:   database.UUIDFromPgx(group.UnitID),
		ParentID: parentID,
		Name:     group.Name,
		Alias:    group.Alias,
		OrgID:    group.OrgID.Bytes,
	}
}

func toCellsGroupModel(group sqlc.CellsGroup) *models.CellsGroup {
	id := database.UUIDFromPgx(group.ID)
	orgID := database.UUIDFromPgx(group.OrgID)
	unitID := database.UUIDFromPgx(group.UnitID)
	storageGroupID := database.UUIDPtrFromPgx(group.StorageGroupID)

	return &models.CellsGroup{
		ID:             id,
		OrgID:          orgID,
		UnitID:         unitID,
		StorageGroupID: storageGroupID,
		Name:           group.Name,
		Alias:          group.Alias,
		CreatedAt:      group.CreatedAt.Time,
	}
}

func toCellModel(cell sqlc.Cell) *models.Cell {
	return &models.Cell{
		ID:           database.UUIDFromPgx(cell.ID),
		CellsGroupID: database.UUIDFromPgx(cell.CellsGroupID),
		Alias:        cell.Alias,
		Row:          int(cell.Row),
		Level:        int(cell.Level),
		Position:     int(cell.Position),
	}
}

func toCellPathModel(segments []sqlc.GetCellPathRow) []models.CellPathSegment {
	result := make([]models.CellPathSegment, len(segments))
	for i, segment := range segments {
		result[i] = models.CellPathSegment{
			ID:         database.UUIDFromPgx(segment.ID),
			Name:       segment.Name,
			ObjectType: models.CellPathObjectType(segment.Type),
			Alias:      segment.Alias,
		}
	}

	return result
}
