package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
	"github.com/let-store-it/backend/internal/services/audit"
	"github.com/let-store-it/backend/internal/services/auth"
)

// AuditUseCase handles audit-related operations with authentication checks
type AuditUseCase struct {
	authService  *auth.AuthService
	auditService *audit.AuditService
}

// NewAuditUseCase creates a new instance of AuditUseCase
func NewAuditUseCase(authService *auth.AuthService, auditService *audit.AuditService) *AuditUseCase {
	return &AuditUseCase{
		authService:  authService,
		auditService: auditService,
	}
}

// CreateObjectChange creates an audit record for an object change after verifying user permissions
func (uc *AuditUseCase) CreateObjectChange(ctx context.Context, userID uuid.UUID, orgID uuid.UUID, objectChange *models.ObjectChange) error {
	// Verify user has access to the organization
	if _, err := uc.authService.GetUserRole(ctx, userID, orgID); err != nil {
		return fmt.Errorf("failed to get user role: %w", err)
	}

	// Set the user and organization IDs in the object change
	objectChange.UserID = userID
	objectChange.OrgID = orgID

	// Create the audit record
	if err := uc.auditService.CreateObjectChange(ctx, objectChange); err != nil {
		return fmt.Errorf("failed to create object change: %w", err)
	}

	return nil
}

// GetObjectChanges retrieves audit records for a specific object after verifying user permissions
func (uc *AuditUseCase) GetObjectChanges(ctx context.Context, targetObjectTypeId models.ObjectTypeId, targetObjectID uuid.UUID) ([]*models.ObjectChange, error) {
	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get org id: %w", err)
	}

	// Get the audit records
	changes, err := uc.auditService.GetObjectChanges(ctx, orgID, targetObjectTypeId, targetObjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get object changes: %w", err)
	}

	return changes, nil
}
