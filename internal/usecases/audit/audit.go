package audit

import (
	"context"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
	"github.com/let-store-it/backend/internal/usecases"
)

type AuditUseCase struct {
	authService  *auth.AuthService
	auditService *audit.AuditService
}

type AuditUseCaseConfig struct {
	AuthService  *auth.AuthService
	AuditService *audit.AuditService
}

func New(config AuditUseCaseConfig) *AuditUseCase {
	return &AuditUseCase{
		authService:  config.AuthService,
		auditService: config.AuditService,
	}
}

func (uc *AuditUseCase) GetObjectChanges(ctx context.Context, targetObjectTypeId models.ObjectTypeId, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	validateResult, err := usecases.ValidateAccessWithOptionalApiToken(ctx, uc.authService, models.AccessLevelWorker, true)
	if err != nil {
		return nil, err
	}

	if !validateResult.IsAllowed {
		return nil, usecases.ErrNotAuthorized
	}
	changes, err := uc.auditService.GetObjectChanges(ctx, validateResult.OrgID, targetObjectTypeId, targetObjectID)
	if err != nil {
		return nil, err
	}

	return changes, nil
}
