package audit

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/utils"
)

type AuditUseCase struct {
	authService  *auth.AuthService
	auditService *audit.AuditService
}

type AuditUseCaseConfig struct {
	AuthService  *auth.AuthService
	AuditService *audit.AuditService
}

func New(config *AuditUseCaseConfig) *AuditUseCase {
	return &AuditUseCase{
		authService:  config.AuthService,
		auditService: config.AuditService,
	}
}

func (uc *AuditUseCase) CreateObjectChange(ctx context.Context, objectChange *models.ObjectChange) error {
	orgID, err := utils.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get org id: %w", err)
	}

	userID, err := utils.GetUserIdFromContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get user id: %w", err)
	}

	objectChange.UserID = userID
	objectChange.OrgID = orgID

	if err := uc.auditService.CreateObjectChange(ctx, objectChange); err != nil {
		return fmt.Errorf("failed to create object change: %w", err)
	}

	return nil
}

func (uc *AuditUseCase) GetObjectChanges(ctx context.Context, targetObjectTypeId models.ObjectTypeId, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	orgID, err := utils.GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get org id: %w", err)
	}

	changes, err := uc.auditService.GetObjectChanges(ctx, orgID, targetObjectTypeId, targetObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object changes: %w", err)
	}

	return changes, nil
}
