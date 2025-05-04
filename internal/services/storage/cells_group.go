package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

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
		OrgID:          utils.PgUUID(orgID),
		StorageGroupID: utils.PgUUID(storageGroupID),
		Name:           name,
		Alias:          alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cells group: %w", err)
	}
	return toCellsGroup(group)
}

func (s *StorageService) GetAllCellsGroups(ctx context.Context, orgID uuid.UUID) ([]*models.CellsGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	groups, err := s.queries.GetCellsGroups(ctx, utils.PgUUID(orgID))
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
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cells group: %w", err)
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
		ID:    utils.PgUUID(group.ID),
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
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		return fmt.Errorf("failed to delete cells group: %w", err)
	}
	return nil
}
