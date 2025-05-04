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
		s.logger.Error("invalid organization ID",
			"method", "CreateCellsGroup",
			"org_id", orgID)
		return nil, ErrInvalidOrganization
	}
	if storageGroupID == uuid.Nil {
		s.logger.Error("invalid storage group ID",
			"method", "CreateCellsGroup",
			"org_id", orgID,
			"storage_group_id", storageGroupID)
		return nil, ErrInvalidStorageGroup
	}
	if err := s.validateStorageGroupData(name, alias); err != nil {
		s.logger.Error("validation failed",
			"method", "CreateCellsGroup",
			"org_id", orgID,
			"storage_group_id", storageGroupID,
			"name", name,
			"alias", alias,
			"error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	group, err := s.queries.CreateCellsGroup(ctx, database.CreateCellsGroupParams{
		OrgID:          utils.PgUUID(orgID),
		StorageGroupID: utils.PgUUID(storageGroupID),
		Name:           name,
		Alias:          alias,
	})
	if err != nil {
		s.logger.Error("failed to create cells group",
			"method", "CreateCellsGroup",
			"org_id", orgID,
			"storage_group_id", storageGroupID,
			"name", name,
			"alias", alias,
			"error", err)
		return nil, fmt.Errorf("failed to create cells group: %w", err)
	}

	s.logger.Info("cells group created successfully",
		"method", "CreateCellsGroup",
		"org_id", orgID,
		"storage_group_id", storageGroupID,
		"group_id", group.ID,
		"name", name,
		"alias", alias)

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
		s.logger.Error("invalid cells group: nil",
			"method", "UpdateCellsGroup")
		return nil, ErrInvalidCellsGroup
	}
	if group.ID == uuid.Nil {
		s.logger.Error("invalid cells group ID",
			"method", "UpdateCellsGroup")
		return nil, ErrInvalidCellsGroup
	}
	if group.OrgID == uuid.Nil {
		s.logger.Error("invalid organization ID",
			"method", "UpdateCellsGroup",
			"group_id", group.ID)
		return nil, ErrInvalidOrganization
	}

	if err := s.validateStorageGroupData(group.Name, group.Alias); err != nil {
		s.logger.Error("validation failed",
			"method", "UpdateCellsGroup",
			"group_id", group.ID,
			"org_id", group.OrgID,
			"name", group.Name,
			"alias", group.Alias,
			"error", err)
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updatedGroup, err := s.queries.UpdateCellsGroup(ctx, database.UpdateCellsGroupParams{
		ID:    utils.PgUUID(group.ID),
		Name:  group.Name,
		Alias: group.Alias,
	})
	if err != nil {
		s.logger.Error("failed to update cells group",
			"method", "UpdateCellsGroup",
			"group_id", group.ID,
			"org_id", group.OrgID,
			"error", err)
		return nil, fmt.Errorf("failed to update cells group: %w", err)
	}

	s.logger.Info("cells group updated successfully",
		"method", "UpdateCellsGroup",
		"group_id", group.ID,
		"org_id", group.OrgID,
		"name", group.Name,
		"alias", group.Alias)

	return toCellsGroup(updatedGroup)
}

func (s *StorageService) DeleteCellsGroup(ctx context.Context, orgID uuid.UUID, id uuid.UUID) error {
	if orgID == uuid.Nil {
		s.logger.Error("invalid organization ID",
			"method", "DeleteCellsGroup",
			"org_id", orgID)
		return ErrInvalidOrganization
	}
	if id == uuid.Nil {
		s.logger.Error("invalid cells group ID",
			"method", "DeleteCellsGroup",
			"org_id", orgID,
			"group_id", id)
		return ErrInvalidCellsGroup
	}

	err := s.queries.DeleteCellsGroup(ctx, database.DeleteCellsGroupParams{
		OrgID: utils.PgUUID(orgID),
		ID:    utils.PgUUID(id),
	})
	if err != nil {
		s.logger.Error("failed to delete cells group",
			"method", "DeleteCellsGroup",
			"group_id", id,
			"org_id", orgID,
			"error", err)
		return fmt.Errorf("failed to delete cells group: %w", err)
	}

	s.logger.Info("cells group deleted successfully",
		"method", "DeleteCellsGroup",
		"group_id", id,
		"org_id", orgID)

	return nil
}
