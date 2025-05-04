package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

func (s *StorageService) GetCells(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID) ([]*models.Cell, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if cellsGroupID == uuid.Nil {
		return nil, ErrInvalidCellsGroup
	}

	cells, err := s.queries.GetCells(ctx, database.GetCellsParams{
		OrgID:        utils.PgUUID(orgID),
		CellsGroupID: utils.PgUUID(cellsGroupID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cells: %w", err)
	}

	result := make([]*models.Cell, len(cells))
	for i, cell := range cells {
		result[i], err = toCell(cell)
		if err != nil {
			return nil, fmt.Errorf("failed to convert cell: %w", err)
		}
	}
	return result, nil
}

func (s *StorageService) GetCellByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.Cell, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return nil, ErrInvalidCell
	}

	cell, err := s.queries.GetCell(ctx, database.GetCellParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cell: %w", err)
	}
	return toCell(cell)
}

func (s *StorageService) CreateCell(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID, alias string, row int, level int, position int) (*models.Cell, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if cellsGroupID == uuid.Nil {
		return nil, ErrInvalidCellsGroup
	}
	if err := s.validateAlias(alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	cell, err := s.queries.CreateCell(ctx, database.CreateCellParams{
		OrgID:        utils.PgUUID(orgID),
		CellsGroupID: utils.PgUUID(cellsGroupID),
		Alias:        alias,
		Row:          int32(row),
		Level:        int32(level),
		Position:     int32(position),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cell: %w", err)
	}
	return toCell(cell)
}

func (s *StorageService) UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	if cell == nil {
		return nil, ErrInvalidCell
	}
	if cell.ID == uuid.Nil {
		return nil, ErrInvalidCell
	}
	if cell.OrgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if cell.CellsGroupID == uuid.Nil {
		return nil, ErrInvalidCellsGroup
	}
	if err := s.validateAlias(cell.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updatedCell, err := s.queries.UpdateCell(ctx, database.UpdateCellParams{
		ID:           utils.PgUUID(cell.ID),
		OrgID:        utils.PgUUID(cell.OrgID),
		CellsGroupID: utils.PgUUID(cell.CellsGroupID),
		Alias:        cell.Alias,
		Row:          int32(cell.Row),
		Level:        int32(cell.Level),
		Position:     int32(cell.Position),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update cell: %w", err)
	}
	return toCell(updatedCell)
}

func (s *StorageService) DeleteCell(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID, id uuid.UUID) error {
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if cellsGroupID == uuid.Nil {
		return ErrInvalidCellsGroup
	}
	if id == uuid.Nil {
		return ErrInvalidCell
	}

	err := s.queries.DeleteCell(ctx, database.DeleteCellParams{
		ID:           utils.PgUUID(id),
		OrgID:        utils.PgUUID(orgID),
		CellsGroupID: utils.PgUUID(cellsGroupID),
	})
	if err != nil {
		return fmt.Errorf("failed to delete cell: %w", err)
	}
	return nil
}

func (s *StorageService) GetCellPath(ctx context.Context, orgID uuid.UUID, cellID uuid.UUID) ([]models.CellPathSegment, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if cellID == uuid.Nil {
		return nil, ErrInvalidCell
	}

	segments, err := s.queries.GetCellPath(ctx, database.GetCellPathParams{
		ID:    utils.PgUUID(cellID),
		OrgID: utils.PgUUID(orgID),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get cell path: %w", err)
	}

	result := make([]models.CellPathSegment, len(segments))
	for i, segment := range segments {
		result[i] = models.CellPathSegment{
			ID:         *utils.UuidFromPgx(segment.ID),
			Name:       segment.Name,
			ObjectType: models.CellPathObjectType(segment.Type),
			Alias:      segment.Alias,
		}
	}
	return result, nil
}
