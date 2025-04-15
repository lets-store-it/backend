package usecases

import (
	"context"
	"fmt"
	"regexp"
	"strings"

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

func (uc *StorageGroupUseCase) validateStorageGroupData(name string, alias string) error {
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

func (uc *StorageGroupUseCase) checkGroupBelongsToOrganization(ctx context.Context, orgID uuid.UUID, groupID uuid.UUID) error {
	exists, err := uc.service.IsStorageGroupExistsForOrganization(ctx, orgID, groupID)
	if err != nil {
		return fmt.Errorf("failed to check group ownership: %w", err)
	}
	if !exists {
		return services.ErrStorageGroupNotFound
	}
	return nil
}

func (uc *StorageGroupUseCase) validateOrganizationAccess(ctx context.Context, groupID uuid.UUID) (uuid.UUID, error) {
	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get organization ID: %w", err)
	}

	if groupID != uuid.Nil {
		if err := uc.checkGroupBelongsToOrganization(ctx, orgID, groupID); err != nil {
			return uuid.Nil, err
		}
	}

	return orgID, nil
}

func (uc *StorageGroupUseCase) Create(ctx context.Context, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, uuid.Nil)
	if err != nil {
		return nil, err
	}

	if err := uc.validateStorageGroupData(name, alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
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
	_, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return nil, err
	}

	group, err := uc.service.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage group: %w", err)
	}

	return group, nil
}

func (uc *StorageGroupUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return err
	}

	return uc.service.Delete(ctx, id)
}

func (uc *StorageGroupUseCase) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	_, err := uc.validateOrganizationAccess(ctx, group.ID)
	if err != nil {
		return nil, err
	}

	if err := uc.validateStorageGroupData(group.Name, group.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return uc.service.Update(ctx, group)
}

func (uc *StorageGroupUseCase) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.StorageGroup, error) {
	_, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return nil, err
	}

	group, err := uc.service.GetByID(ctx, id)
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

	if err := uc.validateStorageGroupData(group.Name, group.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return uc.service.Update(ctx, group)
}
