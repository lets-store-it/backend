package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/usecases"
)

type StorageUseCase struct {
	storageService *storage.StorageService
	orgService     *organization.OrganizationService
	authService    *auth.AuthService
}

type StorageUseCaseConfig struct {
	StorageService *storage.StorageService
	OrgService     *organization.OrganizationService
	AuthService    *auth.AuthService
}

func New(config StorageUseCaseConfig) *StorageUseCase {
	if config.AuthService == nil || config.StorageService == nil || config.OrgService == nil {
		panic("AuthService, StorageService and OrgService are required")
	}

	return &StorageUseCase{
		authService:    config.AuthService,
		storageService: config.StorageService,
		orgService:     config.OrgService,
	}
}

func (uc *StorageUseCase) CreateStorageGroup(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	group.OrgID = validateResult.OrgID

	createdGroup, err := uc.storageService.CreateStorageGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return createdGroup, nil
}

func (uc *StorageUseCase) GetAllStorageGroup(ctx context.Context) ([]*models.StorageGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.storageService.GetAllStorageGroups(ctx, validateResult.OrgID)
}

func (uc *StorageUseCase) GetStorageGroupByID(ctx context.Context, id uuid.UUID) (*models.StorageGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	group, err := uc.storageService.GetStorageGroupByID(ctx, validateResult.OrgID, id)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (uc *StorageUseCase) DeleteStorageGroup(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.storageService.DeleteStorageGroup(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *StorageUseCase) UpdateStorageGroup(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	group.OrgID = validateResult.OrgID

	updatedGroup, err := uc.storageService.UpdateStorageGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return updatedGroup, nil
}

// CellsGroups

func (uc *StorageUseCase) GetCellsGroups(ctx context.Context) ([]*models.CellsGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.storageService.GetCellsGroups(ctx, validateResult.OrgID)
}

func (uc *StorageUseCase) CreateCellsGroup(ctx context.Context, group *models.CellsGroup) (*models.CellsGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	group.OrgID = validateResult.OrgID

	createdGroup, err := uc.storageService.CreateCellsGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return createdGroup, nil
}

func (uc *StorageUseCase) GetCellsGroupByID(ctx context.Context, id uuid.UUID) (*models.CellsGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.storageService.GetCellsGroup(ctx, validateResult.OrgID, id)
}

func (uc *StorageUseCase) UpdateCellsGroup(ctx context.Context, cellGroup *models.CellsGroup) (*models.CellsGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	cellGroup.OrgID = validateResult.OrgID

	updatedGroup, err := uc.storageService.UpdateCellsGroup(ctx, cellGroup)
	if err != nil {
		return nil, err
	}

	return updatedGroup, nil
}

func (uc *StorageUseCase) DeleteCellsGroup(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.storageService.DeleteCellsGroup(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	return nil
}

// Cells

func (uc *StorageUseCase) GetCells(ctx context.Context, cellsGroupID uuid.UUID) ([]*models.Cell, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.storageService.GetCells(ctx, validateResult.OrgID, cellsGroupID)
}

func (uc *StorageUseCase) CreateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	cell.OrgID = validateResult.OrgID

	createdCell, err := uc.storageService.CreateCell(ctx, cell)
	if err != nil {
		return nil, err
	}

	return createdCell, nil
}

func (uc *StorageUseCase) GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.storageService.GetCellByID(ctx, validateResult.OrgID, id)
}

func (uc *StorageUseCase) DeleteCell(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.storageService.DeleteCell(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *StorageUseCase) UpdateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	cell.OrgID = validateResult.OrgID

	updatedCell, err := uc.storageService.UpdateCell(ctx, cell)
	if err != nil {
		return nil, err
	}

	return updatedCell, nil
}
