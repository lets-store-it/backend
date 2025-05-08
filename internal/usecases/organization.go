package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
)

type OrganizationUseCase struct {
	service      *organization.OrganizationService
	authService  *auth.AuthService
	auditService *audit.AuditService
}

func NewOrganizationUseCase(service *organization.OrganizationService, authService *auth.AuthService, auditService *audit.AuditService) *OrganizationUseCase {
	return &OrganizationUseCase{
		service:      service,
		authService:  authService,
		auditService: auditService,
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

	err = uc.authService.SetUserRole(ctx, org.ID, userID, 1) // owner
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		ID:                 uuid.New(),
		OrgID:              org.ID,
		UserID:             userID,
		Action:             models.ObjectChangeActionCreate,
		TargetObjectTypeId: models.ObjectTypeOrganization,
		TargetObjectID:     org.ID,
		PrechangeState:     nil,
		PostchangeState:    postchangeState,
	})

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

	_, err = uc.authService.GetUserRole(ctx, userID, id)
	if err != nil {
		return err
	}

	org, err := uc.service.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = uc.service.Delete(ctx, id)
	if err != nil {
		return err
	}

	prechangeState, err := json.Marshal(org)
	if err != nil {
		return err
	}

	uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		ID:                 uuid.New(),
		OrgID:              org.ID,
		UserID:             userID,
		Action:             models.ObjectChangeActionDelete,
		TargetObjectTypeId: models.ObjectTypeOrganization,
		TargetObjectID:     org.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    nil,
	})
	return nil
}

func (uc *OrganizationUseCase) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = uc.authService.GetUserRole(ctx, userID, org.ID)
	if err != nil {
		return nil, err
	}

	orgUpdated, err := uc.service.Update(ctx, org)
	if err != nil {
		return nil, err
	}

	prechangeState, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(orgUpdated)
	if err != nil {
		return nil, err
	}

	uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		ID:                 uuid.New(),
		OrgID:              org.ID,
		UserID:             userID,
		Action:             models.ObjectChangeActionUpdate,
		TargetObjectTypeId: models.ObjectTypeOrganization,
		TargetObjectID:     org.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    postchangeState,
	})

	return orgUpdated, nil
}

func (uc *OrganizationUseCase) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Organization, error) {
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	_, err = uc.authService.GetUserRole(ctx, userID, id)
	if err != nil {
		return nil, err
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

	orgUpdated, err := uc.service.Update(ctx, org)
	if err != nil {
		return nil, err
	}

	prechangeState, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	postchangeState, err := json.Marshal(orgUpdated)
	if err != nil {
		return nil, err
	}

	uc.auditService.CreateObjectChange(ctx, &models.ObjectChange{
		ID:                 uuid.New(),
		OrgID:              org.ID,
		UserID:             userID,
		Action:             models.ObjectChangeActionUpdate,
		TargetObjectTypeId: models.ObjectTypeOrganization,
		TargetObjectID:     org.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    postchangeState,
	})

	return orgUpdated, nil
}

func (uc *OrganizationUseCase) IsOrganizationExists(ctx context.Context, id uuid.UUID) (bool, error) {
	return uc.service.IsOrganizationExists(ctx, id)
}
