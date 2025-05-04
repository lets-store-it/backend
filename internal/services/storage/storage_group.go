package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/generated/database"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/utils"
)

func (s *StorageService) validateStorageGroupData(name, alias string) error {
	if err := s.validateName(name); err != nil {
		return err
	}
	return s.validateAlias(alias)
}

func (s *StorageService) CreateStorageGroup(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if unitID == uuid.Nil {
		return nil, ErrInvalidUnit
	}
	if err := s.validateStorageGroupData(name, alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	group, err := s.queries.CreateStorageGroup(ctx, database.CreateStorageGroupParams{
		OrgID:    utils.PgUUID(orgID),
		UnitID:   utils.PgUUID(unitID),
		ParentID: utils.NullablePgUUID(parentID),
		Name:     name,
		Alias:    alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create storage group: %w", err)
	}

	return toStorageGroup(group)
}

func (s *StorageService) GetAllStorageGroups(ctx context.Context, orgID uuid.UUID) ([]*models.StorageGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}

	groups, err := s.queries.GetStorageGroups(ctx, utils.PgUUID(orgID))
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

func (s *StorageService) GetStorageGroupByID(ctx context.Context, orgID uuid.UUID, id uuid.UUID) (*models.StorageGroup, error) {
	if orgID == uuid.Nil {
		return nil, ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return nil, ErrInvalidStorageGroup
	}

	group, err := s.queries.GetStorageGroup(ctx, database.GetStorageGroupParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
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

func (s *StorageService) DeleteStorageGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if orgID == uuid.Nil {
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		return ErrInvalidStorageGroup
	}

	err := s.queries.DeleteStorageGroup(ctx, database.DeleteStorageGroupParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		return fmt.Errorf("failed to delete storage group: %w", err)
	}
	return nil
}

func (s *StorageService) UpdateStoragrGroup(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
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
		ID:    utils.PgUUID(group.ID),
		Name:  group.Name,
		Alias: group.Alias,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update storage group: %w", err)
	}
	return toStorageGroup(updatedGroup)
}
