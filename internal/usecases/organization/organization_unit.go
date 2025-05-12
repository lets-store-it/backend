package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/usecases"
)

func (uc *OrganizationUseCase) CreateUnit(ctx context.Context, name string, alias string, address string) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	createdUnit, err := uc.service.CreateUnit(ctx, validateResult.OrgID, name, alias, address)
	if err != nil {
		return nil, err
	}

	return createdUnit, nil
}

func (uc *OrganizationUseCase) GetAllUnits(ctx context.Context) ([]*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetAllUnits(ctx, validateResult.OrgID)
}

func (uc *OrganizationUseCase) GetUnitByID(ctx context.Context, id uuid.UUID) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	unit, err := uc.service.GetUnitByID(ctx, validateResult.OrgID, id)
	if err != nil {
		return nil, err
	}

	return unit, nil
}

func (uc *OrganizationUseCase) DeleteUnit(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	return uc.service.DeleteUnit(ctx, validateResult.OrgID, id)
}

func (uc *OrganizationUseCase) UpdateUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	unit.OrgID = validateResult.OrgID

	updatedUnit, err := uc.service.UpdateUnit(ctx, unit)
	if err != nil {
		return nil, err
	}

	return updatedUnit, nil
}
