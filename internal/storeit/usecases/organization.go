package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/storeit/models"
	"github.com/let-store-it/backend/internal/storeit/services"
)

type OrganizationUseCase struct {
	service     *services.OrganizationService
	authService *services.AuthService
}

func NewOrganizationUseCase(service *services.OrganizationService, authService *services.AuthService) *OrganizationUseCase {
	return &OrganizationUseCase{
		service:     service,
		authService: authService,
	}
}

func (uc *OrganizationUseCase) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	org, err := uc.service.Create(ctx, name, subdomain)
	if err != nil {
		return nil, err
	}

	err = uc.authService.AssignRoleToUser(ctx, org.ID, userID, services.RoleOwner)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (uc *OrganizationUseCase) GetUsersOrgs(ctx context.Context) ([]*models.Organization, error) {
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.GetUsersOrgs(ctx, userID)
}

func (uc *OrganizationUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return uc.service.GetByID(ctx, id)
}

func (uc *OrganizationUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return err
	}

	roles, err := uc.authService.GetUserRoles(ctx, userID, id)
	if err != nil {
		return err
	}

	if len(roles) == 0 {
		return errors.New("no permissions to delete organization")
	}

	if _, ok := roles[services.RoleOwner]; !ok {
		return errors.New("no permissions to delete organization")
	}

	return uc.service.Delete(ctx, id)
}

func (uc *OrganizationUseCase) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := uc.authService.GetUserRoles(ctx, userID, org.ID)
	if err != nil {
		return nil, err
	}

	if _, ok := roles[services.RoleOwner]; !ok {
		return nil, errors.New("no permissions to update organization")
	}

	return uc.service.Update(ctx, org)
}

func (uc *OrganizationUseCase) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Organization, error) {
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	roles, err := uc.authService.GetUserRoles(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	if _, ok := roles[services.RoleOwner]; !ok {
		return nil, errors.New("no permissions to update organization")
	}

	org, err := uc.service.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Apply updates
	if name, ok := updates["name"].(string); ok {
		org.Name = name
	}
	if subdomain, ok := updates["subdomain"].(string); ok {
		org.Subdomain = subdomain
	}

	return uc.service.Update(ctx, org)
}

func (uc *OrganizationUseCase) IsOrganizationExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return uc.service.IsOrganizationExists(ctx, id)
}
