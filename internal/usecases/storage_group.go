package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/storage"
)

type StorageUseCase struct {
	service     *storage.StorageService
	orgService  *organization.OrganizationService
	authService *auth.AuthService
}

func NewStorageUseCase(service *storage.StorageService, orgService *organization.OrganizationService, authService *auth.AuthService) *StorageUseCase {
	return &StorageUseCase{
		authService: authService,
		service:     service,
		orgService:  orgService,
	}
}

func (uc *StorageUseCase) validateOrganizationAccess(ctx context.Context) (uuid.UUID, error) {
	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get organization ID: %w", err)
	}
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	_, err = uc.authService.GetUserRole(ctx, userID, orgID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	// if _, ok := roles[auth.RoleOwner]; !ok {
	// 	return uuid.Nil, fmt.Errorf("user is not an owner of the organization")
	// }

	return orgID, nil
}

func (uc *StorageUseCase) Create(ctx context.Context, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.CreateStorageGroup(ctx, orgID, unitID, parentID, name, alias)
}

func (uc *StorageUseCase) GetAll(ctx context.Context) ([]*models.StorageGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.GetAllStorageGroups(ctx, orgID)
}

func (uc *StorageUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.StorageGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	group, err := uc.service.GetStorageGroupByID(ctx, orgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage group: %w", err)
	}

	return group, nil
}

func (uc *StorageUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return err
	}

	return uc.service.DeleteStorageGroup(ctx, orgID, id)
}

func (uc *StorageUseCase) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	_, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.UpdateStoragrGroup(ctx, group)
}

// CellsGroups

func (uc *StorageUseCase) GetCellsGroups(ctx context.Context) ([]*models.CellsGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.GetCellsGroups(ctx, orgID)
}

func (uc *StorageUseCase) CreateCellsGroup(ctx context.Context, unitID uuid.UUID, storageGroupID *uuid.UUID, name string, alias string) (*models.CellsGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	cellsGroup := &models.CellsGroup{
		OrgID:          orgID,
		UnitID:         unitID,
		StorageGroupID: storageGroupID,
		Name:           name,
		Alias:          alias,
	}

	return uc.service.CreateCellsGroup(ctx, cellsGroup, name, alias)
}

func (uc *StorageUseCase) GetCellsGroupByID(ctx context.Context, id uuid.UUID) (*models.CellsGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.GetCellsGroup(ctx, orgID, id)
}

func (uc *StorageUseCase) UpdateCellsGroup(ctx context.Context, cellGroup *models.CellsGroup) (*models.CellsGroup, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	cellGroup.OrgID = orgID

	return uc.service.UpdateCellsGroup(ctx, cellGroup)
}

func (uc *StorageUseCase) DeleteCellsGroup(ctx context.Context, id uuid.UUID) error {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return err
	}

	return uc.service.DeleteCellsGroup(ctx, orgID, id)
}

// Cells

func (uc *StorageUseCase) GetCells(ctx context.Context, cellsGroupID uuid.UUID) ([]*models.Cell, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.GetCells(ctx, orgID, cellsGroupID)
}

func (uc *StorageUseCase) CreateCell(ctx context.Context, cellsGroupID uuid.UUID, alias string, row int, level int, position int) (*models.Cell, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.CreateCell(ctx, orgID, cellsGroupID, alias, row, level, position)
}

func (uc *StorageUseCase) GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.GetCellByID(ctx, orgID, id)
}

func (uc *StorageUseCase) DeleteCell(ctx context.Context, cellsGroupID uuid.UUID, id uuid.UUID) error {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return err
	}

	return uc.service.DeleteCell(ctx, orgID, cellsGroupID, id)
}

func (uc *StorageUseCase) UpdateCell(ctx context.Context, cellsGroupID uuid.UUID, cell *models.Cell) (*models.Cell, error) {
	orgID, err := uc.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	cell.OrgID = orgID
	cell.CellsGroupID = cellsGroupID

	return uc.service.UpdateCell(ctx, cell)
}
