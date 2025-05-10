package organization

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/services/organization"
	"github.com/let-store-it/backend/internal/usecases"
)

type OrganizationUseCase struct {
	service      *organization.OrganizationService
	authService  *auth.AuthService
	auditService *audit.AuditService
}

type OrganizationUseCaseConfig struct {
	Service      *organization.OrganizationService
	AuthService  *auth.AuthService
	AuditService *audit.AuditService
}

func New(config OrganizationUseCaseConfig) *OrganizationUseCase {
	return &OrganizationUseCase{
		service:      config.Service,
		authService:  config.AuthService,
		auditService: config.AuditService,
	}
}

func (uc *OrganizationUseCase) Create(ctx context.Context, name string, subdomain string) (*models.Organization, error) {
	userId, err := usecases.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	org, err := uc.service.Create(ctx, name, subdomain)
	if err != nil {
		return nil, err
	}

	err = uc.authService.SetUserRole(ctx, org.ID, userId, models.RoleOwnerID)
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
		UserID:             &userId,
		Action:             models.ObjectChangeActionCreate,
		TargetObjectTypeId: models.ObjectTypeOrganization,
		TargetObjectID:     org.ID,
		PrechangeState:     nil,
		PostchangeState:    postchangeState,
	})

	return org, nil
}

func (uc *OrganizationUseCase) GetUsersOrgs(ctx context.Context) ([]*models.Organization, error) {
	userId, err := usecases.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return uc.service.GetUsersOrgs(ctx, userId)
}

func (uc *OrganizationUseCase) GetByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	validateResult, err := usecases.ValidateAccess(ctx, uc.authService, models.AccessLevelAdmin)

	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
	}

	return uc.service.GetByID(ctx, validateResult.OrgID)
}

func (uc *OrganizationUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	validateResult, err := usecases.ValidateAccess(ctx, uc.authService, models.AccessLevelOwner)
	if err != nil {
		return err
	}

	if !validateResult.HasAccess {
		return usecases.ErrNotAuthorized
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
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionDelete,
		TargetObjectTypeId: models.ObjectTypeOrganization,
		TargetObjectID:     org.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    nil,
	})
	return nil
}

func (uc *OrganizationUseCase) Update(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	validateResult, err := usecases.ValidateAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
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
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionUpdate,
		TargetObjectTypeId: models.ObjectTypeOrganization,
		TargetObjectID:     org.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    postchangeState,
	})

	return orgUpdated, nil
}

func (uc *OrganizationUseCase) Patch(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Organization, error) {
	validateResult, err := usecases.ValidateAccess(ctx, uc.authService, models.AccessLevelAdmin)
	if err != nil {
		return nil, err
	}

	if !validateResult.HasAccess {
		return nil, usecases.ErrNotAuthorized
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
		UserID:             validateResult.UserID,
		Action:             models.ObjectChangeActionUpdate,
		TargetObjectTypeId: models.ObjectTypeOrganization,
		TargetObjectID:     org.ID,
		PrechangeState:     prechangeState,
		PostchangeState:    postchangeState,
	})

	return orgUpdated, nil
}
