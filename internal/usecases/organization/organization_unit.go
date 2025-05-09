package organization

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/usecases"
)

func (uc *OrganizationUseCase) CreateUnit(ctx context.Context, name string, alias string, address string) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, common.ErrNotAuthorized
	}

	return uc.service.CreateUnit(ctx, validateResult.OrgID, name, alias, address)
}

func (uc *OrganizationUseCase) GetAllUnits(ctx context.Context) ([]*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, common.ErrNotAuthorized
	}

	return uc.service.GetAllUnits(ctx, validateResult.OrgID)
}

func (uc *OrganizationUseCase) GetUnitByID(ctx context.Context, id uuid.UUID) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, common.ErrNotAuthorized
	}

	unit, err := uc.service.GetUnitByID(ctx, validateResult.OrgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	return unit, nil
}

func (uc *OrganizationUseCase) DeleteUnit(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return common.ErrNotAuthorized
	}

	return uc.service.DeleteUnit(ctx, validateResult.OrgID, id)
}

func (uc *OrganizationUseCase) UpdateUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, common.ErrNotAuthorized
	}

	unit.OrgID = validateResult.OrgID

	return uc.service.UpdateUnit(ctx, unit)
}

func (uc *OrganizationUseCase) PatchUnit(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateOrgAndUserAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, common.ErrNotAuthorized
	}

	unit, err := uc.service.GetUnitByID(ctx, validateResult.OrgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		unit.Name = name
	}
	if alias, ok := updates["alias"].(string); ok {
		unit.Alias = alias
	}
	if address, ok := updates["address"].(string); ok {
		unit.Address = &address
	}

	return uc.service.UpdateUnit(ctx, unit)
}
