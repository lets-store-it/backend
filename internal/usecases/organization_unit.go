package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/organization"
)

type OrganizationUnitUseCase struct {
	orgService  *organization.OrganizationService
	authUseCase *AuthUseCase
}

func NewOrganizationUnitUseCase(orgService *organization.OrganizationService, authUseCase *AuthUseCase) *OrganizationUnitUseCase {
	return &OrganizationUnitUseCase{
		authUseCase: authUseCase,
		orgService:  orgService,
	}
}

func (uc *OrganizationUnitUseCase) Create(ctx context.Context, name string, alias string, address string) (*models.OrganizationUnit, error) {
	orgID, err := uc.authUseCase.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.orgService.CreateUnit(ctx, orgID, name, alias, address)
}

func (uc *OrganizationUnitUseCase) GetAll(ctx context.Context) ([]*models.OrganizationUnit, error) {
	orgID, err := uc.authUseCase.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	return uc.orgService.GetAllUnits(ctx, orgID)
}

func (uc *OrganizationUnitUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationUnit, error) {
	orgID, err := uc.authUseCase.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	unit, err := uc.orgService.GetUnitByID(ctx, orgID, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	return unit, nil
}

func (uc *OrganizationUnitUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	orgID, err := uc.authUseCase.validateOrganizationAccess(ctx)
	if err != nil {
		return err
	}

	return uc.orgService.DeleteUnit(ctx, orgID, id)
}

func (uc *OrganizationUnitUseCase) Update(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	orgID, err := uc.authUseCase.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}
	unit.OrgID = orgID

	return uc.orgService.UpdateUnit(ctx, unit)
}

func (uc *OrganizationUnitUseCase) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationUnit, error) {
	orgID, err := uc.authUseCase.validateOrganizationAccess(ctx)
	if err != nil {
		return nil, err
	}

	unit, err := uc.orgService.GetUnitByID(ctx, orgID, id)
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

	return uc.orgService.UpdateUnit(ctx, unit)
}
