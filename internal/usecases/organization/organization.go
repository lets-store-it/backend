package organization

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/common"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/usecases"
)

type OrganizationUseCase struct {
	service     *organization.OrganizationService
	authService *auth.AuthService
}

type OrganizationUseCaseConfig struct {
	Service     *organization.OrganizationService
	AuthService *auth.AuthService
}

func New(config OrganizationUseCaseConfig) *OrganizationUseCase {
	if config.Service == nil || config.AuthService == nil {
		panic("Service and AuthService are required")
	}
	return &OrganizationUseCase{
		service:     config.Service,
		authService: config.AuthService,
	}
}

func (uc *OrganizationUseCase) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	userId, err := common.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	org, err := uc.service.CreateOrganization(ctx, name, subdomain)
	if err != nil {
		return nil, err
	}

	err = uc.authService.SetUserRole(ctx, org.ID, userId, models.RoleOwnerID)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (uc *OrganizationUseCase) GetUsersOrgs(ctx context.Context) ([]*models.Organization, error) {
	userId, err := common.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.GetUsersOrganization(ctx, userId)
}

func (uc *OrganizationUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)

	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	return uc.service.GetOrganizationByID(ctx, validateResult.OrgID)
}

func (uc *OrganizationUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelOwner, false)
	if err != nil {
		return err
	}

	if !validateResult.IsAllowed {
		return usecases.ErrNotAuthorized
	}

	err = uc.service.DeleteOrganization(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *OrganizationUseCase) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelAdmin, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrForbidden
	}

	org.ID = validateResult.OrgID

	orgUpdated, err := uc.service.UpdateOrganization(ctx, org)
	if err != nil {
		return nil, err
	}

	return orgUpdated, nil
}
