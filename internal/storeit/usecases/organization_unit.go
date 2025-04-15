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

type OrganizationUnitUseCase struct {
	service    *services.OrganizationUnitService
	orgService *services.OrganizationService
}

func NewOrganizationUnitUseCase(service *services.OrganizationUnitService, orgService *services.OrganizationService) *OrganizationUnitUseCase {
	return &OrganizationUnitUseCase{
		service:    service,
		orgService: orgService,
	}
}

func (uc *OrganizationUnitUseCase) validateOrganizationUnitData(name string, alias string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("organization unit name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("organization unit name is too long (max 100 characters)")
	}

	if strings.TrimSpace(alias) == "" {
		return fmt.Errorf("organization unit alias cannot be empty")
	}
	if len(alias) > 100 {
		return fmt.Errorf("organization unit alias is too long (max 100 characters)")
	}
	matched, _ := regexp.MatchString("^[\\w-]+$", alias)
	if !matched {
		return fmt.Errorf("organization unit alias can only contain letters, numbers, and hyphens (no spaces)")
	}
	return nil
}

func (uc *OrganizationUnitUseCase) checkUnitBelongsToOrganization(ctx context.Context, orgID uuid.UUID, unitID uuid.UUID) error {
	exists, err := uc.service.IsOrganizationUnitExistsForOrganization(ctx, orgID, unitID)
	if err != nil {
		return fmt.Errorf("failed to check unit ownership: %w", err)
	}
	if !exists {
		return services.ErrOrganizationUnitNotFound
	}
	return nil
}

func (uc *OrganizationUnitUseCase) validateOrganizationAccess(ctx context.Context, unitID uuid.UUID) (uuid.UUID, error) {
	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get organization ID: %w", err)
	}

	if unitID != uuid.Nil {
		if err := uc.checkUnitBelongsToOrganization(ctx, orgID, unitID); err != nil {
			return uuid.Nil, err
		}
	}

	return orgID, nil
}

func (uc *OrganizationUnitUseCase) Create(ctx context.Context, name string, alias string, address string) (*models.OrganizationUnit, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, uuid.Nil)
	if err != nil {
		return nil, err
	}

	if err := uc.validateOrganizationUnitData(name, alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return uc.service.Create(ctx, orgID, name, alias, address)
}

func (uc *OrganizationUnitUseCase) GetAll(ctx context.Context) ([]*models.OrganizationUnit, error) {
	orgID, err := uc.validateOrganizationAccess(ctx, uuid.Nil)
	if err != nil {
		return nil, err
	}

	return uc.service.GetAll(ctx, orgID)
}

func (uc *OrganizationUnitUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.OrganizationUnit, error) {
	_, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return nil, err
	}

	unit, err := uc.service.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization unit: %w", err)
	}

	return unit, nil
}

func (uc *OrganizationUnitUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return err
	}

	return uc.service.Delete(ctx, id)
}

func (uc *OrganizationUnitUseCase) Update(ctx context.Context, unit *models.OrganizationUnit) (*models.OrganizationUnit, error) {
	_, err := uc.validateOrganizationAccess(ctx, unit.ID)
	if err != nil {
		return nil, err
	}

	if err := uc.validateOrganizationUnitData(unit.Name, unit.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return uc.service.Update(ctx, unit)
}

func (uc *OrganizationUnitUseCase) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.OrganizationUnit, error) {
	_, err := uc.validateOrganizationAccess(ctx, id)
	if err != nil {
		return nil, err
	}

	unit, err := uc.service.GetByID(ctx, id)
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

	if err := uc.validateOrganizationUnitData(unit.Name, unit.Alias); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return uc.service.Update(ctx, unit)
}
