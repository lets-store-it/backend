package organization

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/usecases"
)

func (uc *OrganizationUseCase) CreateUnit(ctx context.Context, name string, alias string, address string) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAuthorized {
		return nil, usecases.ErrNotAuthorized
	}

	createdUnit, err := uc.service.CreateUnit(ctx, validateResult.OrgID, name, alias, address)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(createdUnit)
	if err != nil {
		return nil, err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionCreate,
		TargetObjectTypeId: models.ObjectTypeUnit,
		TargetObjectID:     createdUnit.ID,
		PrechangeState:     nil,
		PostchangeState:    postchangeState,
	})
	if err != nil {
		return nil, err
	}

	return createdUnit, nil
}

func (uc *OrganizationUseCase) GetAllUnits(ctx context.Context) ([]*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAuthorized {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetAllUnits(ctx, validateResult.OrgID)
}

func (uc *OrganizationUseCase) GetUnitByID(ctx context.Context, id uuid.UUID) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAuthorized {
		return nil, usecases.ErrNotAuthorized
	}

	unit, err := uc.service.GetUnitByID(ctx, validateResult.OrgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	return unit, nil
}

func (uc *OrganizationUseCase) DeleteUnit(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return err
	}

	if !validateResult.IsAuthorized {
		return usecases.ErrNotAuthorized
	}

	return uc.service.DeleteUnit(ctx, validateResult.OrgID, id)
}

func (uc *OrganizationUseCase) UpdateUnit(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAuthorized {
		return nil, usecases.ErrNotAuthorized
	}

	beforeUpdate, err := uc.service.GetUnitByID(ctx, validateResult.OrgID, unit.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	prechangeState, err := json.Marshal(beforeUpdate)
	if err != nil {
		return nil, err
	}

	unit.OrgID = validateResult.OrgID

	updatedUnit, err := uc.service.UpdateUnit(ctx, unit)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(updatedUnit)
	if err != nil {
		return nil, err
	}

	err = uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		OrgID:              validateResult.OrgID,
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionUpdate,
		TargetObjectTypeId: models.ObjectTypeUnit,
		TargetObjectID:     unit.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    postchangeState,
	})
	if err != nil {
		return nil, err
	}

	return updatedUnit, nil
}

func (uc *OrganizationUseCase) PatchUnit(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationUnit, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAuthorized {
		return nil, usecases.ErrNotAuthorized
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
