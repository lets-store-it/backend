package storage

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
)

var (
	ErrStorageGroupNotFound = errors.New("storage group not found")
	ErrCellsGroupNotFound   = errors.New("cells group not found")
	ErrCellNotFound         = errors.New("cell not found")
	ErrInvalidStorageGroup  = errors.New("invalid storage group")
	ErrInvalidCellsGroup    = errors.New("invalid cells group")
	ErrInvalidCell          = errors.New("invalid cell")
	ErrInvalidName          = errors.New("invalid name")
	ErrInvalidAlias         = errors.New("invalid alias")
	ErrInvalidOrganization  = errors.New("invalid organization")
	ErrInvalidUnit          = errors.New("invalid unit")
)

func uuidFromPgx(id pgtype.UUID) *uuid.UUID {
	if !id.Valid {
		return nil
	}
	result := uuid.UUID(id.Bytes)
	return &result
}

type StorageService struct {
	queries *database.Queries
}

func New(queries *database.Queries) *StorageService {
	return &StorageService{
		queries: queries,
	}
}

func toStorageGroup(group database.StorageGroup) (*models.StorageGroup, error) {
	id := uuidFromPgx(group.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert storage group: %w", ErrInvalidStorageGroup)
	}
	unitID := uuidFromPgx(group.UnitID)
	if unitID == nil {
		return nil, fmt.Errorf("failed to convert storage group: %w", ErrInvalidUnit)
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
		return fmt.Errorf("%w: name cannot be empty", ErrInvalidName)
	}
	if len(name) > 100 {
		return fmt.Errorf("%w: name is too long (max 100 characters)", ErrInvalidName)
	}

	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("%w: alias cannot be empty", ErrInvalidAlias)
	}
	if len(alias) > 100 {
		return fmt.Errorf("%w: alias is too long (max 100 characters)", ErrInvalidAlias)
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("%w: alias can only contain letters, numbers, and hyphens (no spaces)", ErrInvalidAlias)
	}
	return nil
}

func (s *StorageService) Create(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if unitID == uuid.Nil {
		return nil, ErrInvalidUnit
	}
	if err := s.validateStorageGroupData(name, alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	var parentIDPgx pgtype.UUID
	if parentID != nil {
		if *parentID == uuid.Nil {
			return nil, ErrInvalidStorageGroup
		}
		parentIDPgx = pgtype.UUID{Bytes: *parentID, Valid: true}
	}

	group, err := s.queries.CreateStorageGroup(ctx, database.CreateStorageGroupParams{
		OrgID:    pgtype.UUID{Bytes: orgID, Valid: true},
		UnitID:   pgtype.UUID{Bytes: unitID, Valid: true},
		ParentID: parentIDPgx,
		Name:     name,
		Alias:    alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create storage group: %w", err)
	}

	return toStorageGroup(group)
}

func (s *StorageService) GetAll(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	groups, err := s.queries.GetStorageGroups(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get storage groups: %w", err)
	}

	result := make([]*models.StorageGroup, len(groups))
	for i, group := range groups {
		result[i], err = toStorageGroup(group)
		if err != nil {
			return nil, fmt.Errorf("failed to convert storage group: %w", err)
		}
	}

	return result, nil
}

func (s *StorageService) GetByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.StorageGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return nil, ErrInvalidStorageGroup
	}

	group, err := s.queries.GetStorageGroup(ctx, database.GetStorageGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get storage group: %w", err)
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
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return ErrInvalidStorageGroup
	}

	err := s.queries.DeleteStorageGroup(ctx, database.DeleteStorageGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to delete storage group: %w", err)
	}
	return nil
}

func (s *StorageService) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	if group == nil {
		return nil, ErrInvalidStorageGroup
	}
	if group.ID == uuid.Nil {
		return nil, ErrInvalidStorageGroup
	}

	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updatedGroup, err := s.queries.UpdateStorageGroup(ctx, database.UpdateStorageGroupParams{
		ID:    pgtype.UUID{Bytes: group.ID, Valid: true},
		Name:  group.Name,
		Alias: group.Alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update storage group: %w", err)
	}
	return toStorageGroup(updatedGroup)
}

func toCellsGroup(group database.CellsGroup) (*models.CellsGroup, error) {
	id := uuidFromPgx(group.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert cells group: %w", ErrInvalidCellsGroup)
	}
	storageGroupID := uuidFromPgx(group.StorageGroupID)
	if storageGroupID == nil {
		return nil, fmt.Errorf("failed to convert cells group: %w", ErrInvalidStorageGroup)
	}
	orgID := uuidFromPgx(group.OrgID)
	if orgID == nil {
		return nil, fmt.Errorf("failed to convert cells group: %w", ErrInvalidOrganization)
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
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	groups, err := s.queries.GetCellsGroups(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get cells groups: %w", err)
	}

	result := make([]*models.CellsGroup, len(groups))
	for i, group := range groups {
		result[i], err = toCellsGroup(group)
		if err != nil {
			return nil, fmt.Errorf("failed to convert cells group: %w", err)
		}
	}

	return result, nil
}

func (s *StorageService) GetCellsGroupByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.CellsGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return nil, ErrInvalidCellsGroup
	}

	group, err := s.queries.GetCellsGroup(ctx, database.GetCellsGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cells group: %w", err)
	}
	return toCellsGroup(group)
}

func (s *StorageService) CreateCellsGroup(ctx context.Context, orgID uuid.UUID, storageGroupID uuid.UUID, name string, alias string) (*models.CellsGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if storageGroupID == uuid.Nil {
		return nil, ErrInvalidStorageGroup
	}
	if err := s.validateStorageGroupData(name, alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	group, err := s.queries.CreateCellsGroup(ctx, database.CreateCellsGroupParams{
		OrgID:          pgtype.UUID{Bytes: orgID, Valid: true},
		StorageGroupID: pgtype.UUID{Bytes: storageGroupID, Valid: true},
		Name:           name,
		Alias:          alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cells group: %w", err)
	}
	return toCellsGroup(group)
}

func (s *StorageService) UpdateCellsGroup(ctx context.Context, group *models.CellsGroup) (*models.CellsGroup, error) {
	if group == nil {
		return nil, ErrInvalidCellsGroup
	}
	if group.ID == uuid.Nil {
		return nil, ErrInvalidCellsGroup
	}
	if group.OrgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updatedGroup, err := s.queries.UpdateCellsGroup(ctx, database.UpdateCellsGroupParams{
		ID:    pgtype.UUID{Bytes: group.ID, Valid: true},
		Name:  group.Name,
		Alias: group.Alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update cells group: %w", err)
	}
	return toCellsGroup(updatedGroup)
}

func (s *StorageService) DeleteCellsGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return ErrInvalidCellsGroup
	}

	err := s.queries.DeleteCellsGroup(ctx, database.DeleteCellsGroupParams{
		OrgID: pgtype.UUID{Bytes: orgID, Valid: true},
		ID:    pgtype.UUID{Bytes: id, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to delete cells group: %w", err)
	}
	return nil
}

func toCell(cell database.Cell) (*models.Cell, error) {
	id := uuidFromPgx(cell.ID)
	if id == nil {
		return nil, fmt.Errorf("failed to convert cell: %w", ErrInvalidCell)
	}
	cellsGroupID := uuidFromPgx(cell.CellsGroupID)
	if cellsGroupID == nil {
		return nil, fmt.Errorf("failed to convert cell: %w", ErrInvalidCellsGroup)
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
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if cellsGroupID == uuid.Nil {
		return nil, ErrInvalidCellsGroup
	}

	cells, err := s.queries.GetCells(ctx, database.GetCellsParams{
		OrgID:        pgtype.UUID{Bytes: orgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cellsGroupID, Valid: true},
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

func (s *StorageService) GetCellByID(ctx context.Context, orgID uuid.UUID, cellsGroupID uuid.UUID, id uuid.UUID) (*models.Cell, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if cellsGroupID == uuid.Nil {
		return nil, ErrInvalidCellsGroup
	}
	if id == uuid.Nil {
		return nil, ErrInvalidCell
	}

	cell, err := s.queries.GetCell(ctx, database.GetCellParams{
		OrgID:        pgtype.UUID{Bytes: orgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cellsGroupID, Valid: true},
		ID:           pgtype.UUID{Bytes: id, Valid: true},
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
	if strings.TrimSpace(alias) == "" {
		return nil, fmt.Errorf("%w: alias cannot be empty", ErrInvalidAlias)
	}

	cell, err := s.queries.CreateCell(ctx, database.CreateCellParams{
		OrgID:        pgtype.UUID{Bytes: orgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cellsGroupID, Valid: true},
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
	if strings.TrimSpace(cell.Alias) == "" {
		return nil, fmt.Errorf("%w: alias cannot be empty", ErrInvalidAlias)
	}

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
		ID:           pgtype.UUID{Bytes: id, Valid: true},
		OrgID:        pgtype.UUID{Bytes: orgID, Valid: true},
		CellsGroupID: pgtype.UUID{Bytes: cellsGroupID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to delete cell: %w", err)
	}
	return nil
}
