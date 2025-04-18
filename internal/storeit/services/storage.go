package services

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/storeit/models"
)

var (
	ErrStorageGroupNotFound = errors.New("storage group not found")
)

type StorageService struct {
	queries *database.Queries
}

func NewStorageService(queries *database.Queries) *StorageService {
	return &StorageService{
		queries: queries,
	}
}

func toStorageGroup(group database.StorageGroup) (*models.StorageGroup, error) {
	id := uuidFromPgx(group.ID)
	if id == nil {
		return nil, errors.New("id is nil")
	}
	unitID := uuidFromPgx(group.UnitID)
	if unitID == nil {
		return nil, errors.New("unit_id is nil")
	}

	return &models.StorageGroup{
		ID:       *id,
		UnitID:   *unitID,
		ParentID: uuidFromPgx(group.ParentID),
		Name:     group.Name,
		Alias:    group.Alias,
	}, nil
}

func (s *StorageService) validateStorageGroupData(name string, alias string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("storage group name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("storage group name is too long (max 100 characters)")
	}

	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("storage group alias cannot be empty")
	}
	if len(alias) > 100 {
		return fmt.Errorf("storage group alias is too long (max 100 characters)")
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("storage group alias can only contain letters, numbers, and hyphens (no spaces)")
	}
	return nil
}

func (s *StorageService) Create(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	var parentIDPgx pgtype.UUID
	if parentID != nil {
		parentIDPgx = pgtype.UUID{Bytes: *parentID, Valid: true}
	} else {
		parentIDPgx = pgtype.UUID{Valid: false}
	}

	group, err := s.queries.CreateStorageGroup(ctx, database.CreateStorageGroupParams{
		OrgID:    pgtype.UUID{Bytes: orgID, Valid: true},
		UnitID:   pgtype.UUID{Bytes: unitID, Valid: true},
		ParentID: parentIDPgx,
		Name:     name,
		Alias:    alias,
	})
	if err != nil {
		return nil, err
	}

	return toStorageGroup(group)
}

func (s *StorageService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	groups, err := s.queries.GetActiveStorageGroups(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]*models.StorageGroup, len(groups))
	for i, group := range groups {
		result[i], err = toStorageGroup(group)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *StorageService) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.StorageGroup, error) {
	group, err := s.queries.GetStorageGroup(ctx, database.GetStorageGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	model, err := toStorageGroup(group)

	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, ErrStorageGroupNotFound
	}
	return model, nil
}

func (s *StorageService) Delete(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.queries.DeleteStorageGroup(ctx, database.DeleteStorageGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
}

func (s *StorageService) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	updatedGroup, err := s.queries.UpdateStorageGroup(ctx, database.UpdateStorageGroupParams{
		ID:    pgtype.UUID{Bytes: group.ID, Valid: true},
		Name:  group.Name,
		Alias: group.Alias,
	})
	if err != nil {
		return nil, err
	}
	return toStorageGroup(updatedGroup)
}

func (s *StorageService) IsStorageGroupExists(ctx context.Context, orgID uuid.UUID, groupID uuid.UUID) (bool, error) {
	return s.queries.IsStorageGroupExists(ctx, database.IsStorageGroupExistsParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: groupID, Valid: true},
	})
}

// CellsGroups

func toCellsGroup(group database.CellsGroup) (*models.CellsGroup, error) {
	id := uuidFromPgx(group.ID)
	if id == nil {
		return nil, errors.New("id is nil")
	}
	storageGroupID := uuidFromPgx(group.StorageGroupID)
	if storageGroupID == nil {
		return nil, errors.New("storage_group_id is nil")
	}
	orgID := uuidFromPgx(group.OrgID)
	if orgID == nil {
		return nil, errors.New("org_id is nil")
	}

	return &models.CellsGroup{
		ID:             *id,
		OrgID:          *orgID,
		StorageGroupID: *storageGroupID,
		Name:           group.Name,
		Alias:          group.Alias,
	}, nil
}

func (s *StorageService) GetCellsGroups(ctx context.Context, orgID uuid.UUID) ([]*models.CellsGroup, error) {
	groups, err := s.queries.GetCellsGroups(ctx, pgtype.UUID{Bytes: orgID, Valid: true})

	result := make([]*models.CellsGroup, len(groups))
	for i, group := range groups {
		result[i], err = toCellsGroup(group)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *StorageService) GetCellsGroupByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.CellsGroup, error) {
	group, err := s.queries.GetCellsGroup(ctx, database.GetCellsGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return toCellsGroup(group)
}

func (s *StorageService) CreateCellsGroup(ctx context.Context, orgID uuid.UUID, storageGroupID uuid.UUID, name string, alias string) (*models.CellsGroup, error) {
	group, err := s.queries.CreateCellsGroup(ctx, database.CreateCellsGroupParams{
		OrgID:          pgtype.UUID{Bytes: orgID, Valid: true},
		StorageGroupID: pgtype.UUID{Bytes: storageGroupID, Valid: true},
		Name:           name,
		Alias:          alias,
	})
	if err != nil {
		return nil, err
	}
	return toCellsGroup(group)
}

func (s *StorageService) UpdateCellsGroup(ctx context.Context, group *models.CellsGroup) (*models.CellsGroup, error) {
	updatedGroup, err := s.queries.UpdateCellsGroup(ctx, database.UpdateCellsGroupParams{
		ID:    pgtype.UUID{Bytes: group.ID, Valid: true},
		Name:  group.Name,
		Alias: group.Alias,
	})
	if err != nil {
		return nil, err
	}
	return toCellsGroup(updatedGroup)
}

func (s *StorageService) DeleteCellsGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	return s.queries.DeleteCellsGroup(ctx, database.DeleteCellsGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
}

// Cells

func toCell(cell database.Cell) (*models.Cell, error) {
	id := uuidFromPgx(cell.ID)
	if id == nil {
		return nil, errors.New("id is nil")
	}
	cellsGroupID := uuidFromPgx(cell.CellsGroupID)
	if cellsGroupID == nil {
		return nil, errors.New("cells_group_id is nil")
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

func (s *StorageService) GetCells(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID) ([]*models.Cell, error) {
	cells, err := s.queries.GetCells(ctx, database.GetCellsParams{
		OrgID:        pgtype.UUID{Bytes: orgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cellsGroupID, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	result := make([]*models.Cell, len(cells))
	for i, cell := range cells {
		result[i], err = toCell(cell)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *StorageService) GetCellByID(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID, id uuid.UUID) (*models.Cell, error) {
	cell, err := s.queries.GetCell(ctx, database.GetCellParams{
		OrgID:        pgtype.UUID{Bytes: orgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cellsGroupID, Valid: true},
		ID:           pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, err
	}
	return toCell(cell)
}

func (s *StorageService) CreateCell(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID, alias string, row int, level int, position int) (*models.Cell, error) {
	cell, err := s.queries.CreateCell(ctx, database.CreateCellParams{
		OrgID:        pgtype.UUID{Bytes: orgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cellsGroupID, Valid: true},
		Alias:        alias,
		Row:          int32(row),
		Level:        int32(level),
		Position:     int32(position),
	})
	if err != nil {
		return nil, err
	}
	return toCell(cell)
}

func (s *StorageService) UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	updatedCell, err := s.queries.UpdateCell(ctx, database.UpdateCellParams{
		ID:           pgtype.UUID{Bytes: cell.ID, Valid: true},
		OrgID:        pgtype.UUID{Bytes: cell.OrgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cell.CellsGroupID, Valid: true},
		Alias:        cell.Alias,
		Row:          int32(cell.Row),
		Level:        int32(cell.Level),
		Position:     int32(cell.Position),
	})
	if err != nil {
		return nil, err
	}
	return toCell(updatedCell)
}

func (s *StorageService) DeleteCell(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID, id uuid.UUID) error {
	return s.queries.DeleteCell(ctx, database.DeleteCellParams{
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		OrgID:        pgtype.UUID{Bytes: orgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cellsGroupID, Valid: true},
	})
}
