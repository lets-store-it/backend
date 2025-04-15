package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
	"github.com/let-store-it/backend/internal/storeit/services"
)

type StorageGroupUseCase struct {
	service    *services.StorageGroupService
	orgService *services.OrganizationService
}

func NewStorageGroupUseCase(service *services.StorageGroupService, orgService *services.OrganizationService) *StorageGroupUseCase {
	return &StorageGroupUseCase{
		service:    service,
		orgService: orgService,
	}
}

func (uc *StorageGroupUseCase) validateOrganizationAccess(ctx context.Context, groupID uuid.UUID) (uuid.UUID, error) {
	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get organization ID: %w", err)
	}

	if groupID != uuid.Nil {
		exists, err := uc.service.IsStorageGroupExists(ctx, orgID, groupID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("failed to check group ownership: %w", err)
		}
		if !exists {
			return uuid.Nil, services.ErrStorageGroupNotFound
		}
	}

	return orgID, nil
}

func (uc *StorageGroupUseCase) Create(ctx context.Context, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, uuid.Nil)
	if err != nil {
		return nil, err
	}

	return uc.service.Create(ctx, orgID, unitID, parentID, name, alias)
}

func (uc *StorageGroupUseCase) GetAll(ctx context.Context) ([]*models.StorageGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, uuid.Nil)
	if err != nil {
		return nil, err
	}

	return uc.service.GetAll(ctx, orgID)
}

func (uc *StorageGroupUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.StorageGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return nil, err
	}

	group, err := uc.service.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage group: %w", err)
	}

	return group, nil
}

func (uc *StorageGroupUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	orgID, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return err
	}

	return uc.service.Delete(ctx, orgID, id)
}

func (uc *StorageGroupUseCase) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	_, err := uc.validateOrganizationAccess(ctx, group.ID)
	if err != nil {
		return nil, err
	}

	return uc.service.Update(ctx, group)
}

func (uc *StorageGroupUseCase) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.StorageGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return nil, err
	}

	group, err := uc.service.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage group: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		group.Name = name
	}
	if alias, ok := updates["alias"].(string); ok {
		group.Alias = alias
	}
	if parentID, ok := updates["parent_id"].(uuid.UUID); ok {
		group.ParentID = parentID
	}

	return uc.service.Update(ctx, group)
}
