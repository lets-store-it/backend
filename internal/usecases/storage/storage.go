package storage

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/usecases"
)

type StorageUseCase struct {
	service      *storage.StorageService
	orgService   *organization.OrganizationService
	authService  *auth.AuthService
	auditService *audit.AuditService
}

type StorageUseCaseConfig struct {
	Service      *storage.StorageService
	OrgService   *organization.OrganizationService
	AuthService  *auth.AuthService
	AuditService *audit.AuditService
}

func New(config StorageUseCaseConfig) *StorageUseCase {
	return &StorageUseCase{
		authService:  config.AuthService,
		service:      config.Service,
		orgService:   config.OrgService,
		auditService: config.AuditService,
	}
}

func (uc *StorageUseCase) Create(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	group.OrgID = validateResult.OrgID

	createdGroup, err := uc.service.CreateStorageGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return createdGroup, nil
}

func (uc *StorageUseCase) GetAll(ctx context.Context) ([]*models.StorageGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetAllStorageGroups(ctx, validateResult.OrgID)
}

func (uc *StorageUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.StorageGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	group, err := uc.service.GetStorageGroupByID(ctx, validateResult.OrgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage group: %w", err)
	}

	return group, nil
}

func (uc *StorageUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.service.DeleteStorageGroup(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *StorageUseCase) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	group.OrgID = validateResult.OrgID

	updatedGroup, err := uc.service.UpdateStorageGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return updatedGroup, nil
}

// CellsGroups

func (uc *StorageUseCase) GetCellsGroups(ctx context.Context) ([]*models.CellsGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetCellsGroups(ctx, validateResult.OrgID)
}

func (uc *StorageUseCase) CreateCellsGroup(ctx context.Context, group *models.CellsGroup) (*models.CellsGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	group.OrgID = validateResult.OrgID

	createdGroup, err := uc.service.CreateCellsGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return createdGroup, nil
}

func (uc *StorageUseCase) GetCellsGroupByID(ctx context.Context, id uuid.UUID) (*models.CellsGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetCellsGroup(ctx, validateResult.OrgID, id)
}

func (uc *StorageUseCase) UpdateCellsGroup(ctx context.Context, cellGroup *models.CellsGroup) (*models.CellsGroup, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	cellGroup.OrgID = validateResult.OrgID

	updatedGroup, err := uc.service.UpdateCellsGroup(ctx, cellGroup)
	if err != nil {
		return nil, err
	}

	return updatedGroup, nil
}

func (uc *StorageUseCase) DeleteCellsGroup(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.service.DeleteCellsGroup(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	return nil
}

// Cells

func (uc *StorageUseCase) GetCells(ctx context.Context, cellsGroupID uuid.UUID) ([]*models.Cell, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetCells(ctx, validateResult.OrgID, cellsGroupID)
}

func (uc *StorageUseCase) CreateCell(ctx context.Context, cell *models.Cell) (*models.Cell, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelManager, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	cell.OrgID = validateResult.OrgID

	createdCell, err := uc.service.CreateCell(ctx, cell)
	if err != nil {
		return nil, err
	}

	return createdCell, nil
}

func (uc *StorageUseCase) GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetCellByID(ctx, validateResult.OrgID, id)
}

func (uc *StorageUseCase) DeleteCell(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.service.DeleteCell(ctx, validateResult.OrgID, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *StorageUseCase) UpdateCell(ctx context.Context, cellsGroupID uuid.UUID, cell *models.Cell) (*models.Cell, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}

	cell.OrgID = validateResult.OrgID
	cell.CellsGroupID = cellsGroupID

	updatedCell, err := uc.service.UpdateCell(ctx, cell)
	if err != nil {
		return nil, err
	}

	return updatedCell, nil
}
