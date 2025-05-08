package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/services/storage"
	"github.com/let-store-it/backend/internal/usecases"
	"github.com/let-store-it/backend/internal/utils"
)

type StorageUseCase struct {
	service     *storage.StorageService
	orgService  *organization.OrganizationService
	authService *auth.AuthService
}

type StorageUseCaseConfig struct {
	Service     *storage.StorageService
	OrgService  *organization.OrganizationService
	AuthService *auth.AuthService
}

func New(config StorageUseCaseConfig) *StorageUseCase {
	return &StorageUseCase{
		authService: config.AuthService,
		service:     config.Service,
		orgService:  config.OrgService,
	}
}

func (uc *StorageUseCase) Create(ctx context.Context, unitID uuid.UUID, parentID *uuid.UUID, name string, alias string) (*models.StorageGroup, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.CreateStorageGroup(ctx, validateResult.OrgID, unitID, parentID, name, alias)
}

func (uc *StorageUseCase) GetAll(ctx context.Context) ([]*models.StorageGroup, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetAllStorageGroups(ctx, validateResult.OrgID)
}

func (uc *StorageUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.StorageGroup, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	group, err := uc.service.GetStorageGroupByID(ctx, validateResult.OrgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get storage group: %w", err)
	}

	return group, nil
}

func (uc *StorageUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return usecases.ErrNotAuthorized
	}

	return uc.service.DeleteStorageGroup(ctx, validateResult.OrgID, id)
}

func (uc *StorageUseCase) Update(ctx context.Context, group *models.StorageGroup) (*models.StorageGroup, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.UpdateStoragrGroup(ctx, group)
}

// CellsGroups

func (uc *StorageUseCase) GetCellsGroups(ctx context.Context) ([]*models.CellsGroup, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetCellsGroups(ctx, validateResult.OrgID)
}

func (uc *StorageUseCase) CreateCellsGroup(ctx context.Context, unitID uuid.UUID, storageGroupID *uuid.UUID, name string, alias string) (*models.CellsGroup, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	cellsGroup := &models.CellsGroup{
		OrgID:          validateResult.OrgID,
		UnitID:         unitID,
		StorageGroupID: storageGroupID,
		Name:           name,
		Alias:          alias,
	}

	return uc.service.CreateCellsGroup(ctx, cellsGroup, name, alias)
}

func (uc *StorageUseCase) GetCellsGroupByID(ctx context.Context, id uuid.UUID) (*models.CellsGroup, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetCellsGroup(ctx, validateResult.OrgID, id)
}

func (uc *StorageUseCase) UpdateCellsGroup(ctx context.Context, cellGroup *models.CellsGroup) (*models.CellsGroup, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	cellGroup.OrgID = validateResult.OrgID

	return uc.service.UpdateCellsGroup(ctx, cellGroup)
}

func (uc *StorageUseCase) DeleteCellsGroup(ctx context.Context, id uuid.UUID) error {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return usecases.ErrNotAuthorized
	}

	return uc.service.DeleteCellsGroup(ctx, validateResult.OrgID, id)
}

// Cells

func (uc *StorageUseCase) GetCells(ctx context.Context, cellsGroupID uuid.UUID) ([]*models.Cell, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetCells(ctx, validateResult.OrgID, cellsGroupID)
}

func (uc *StorageUseCase) CreateCell(ctx context.Context, cellsGroupID uuid.UUID, alias string, row int, level int, position int) (*models.Cell, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.CreateCell(ctx, validateResult.OrgID, cellsGroupID, alias, row, level, position)
}

func (uc *StorageUseCase) GetCellByID(ctx context.Context, id uuid.UUID) (*models.Cell, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetCellByID(ctx, validateResult.OrgID, id)
}

func (uc *StorageUseCase) DeleteCell(ctx context.Context, cellsGroupID uuid.UUID, id uuid.UUID) error {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return usecases.ErrNotAuthorized
	}

	return uc.service.DeleteCell(ctx, validateResult.OrgID, cellsGroupID, id)
}

func (uc *StorageUseCase) UpdateCell(ctx context.Context, cellsGroupID uuid.UUID, cell *models.Cell) (*models.Cell, error) {
	validateResult, err := utils.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	cell.OrgID = validateResult.OrgID
	cell.CellsGroupID = cellsGroupID

	return uc.service.UpdateCell(ctx, cell)
}
