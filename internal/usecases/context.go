package usecases

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/let-store-it/backend/internal/models"
)

func GetOrganizationIDFromContext(ctx context.Context) (uuid.UUID, error) {
	orgID, ok := ctx.Value(models.OrganizationIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("organization ID not found in context or invalid type")
	}
	if orgID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("organization ID not found in context")
	}
	return orgID, nil
}

func GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(models.UserIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user ID not found in context or invalid type")
	}
	if userID == uuid.Nil {
		return uuid.Nil, fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

type AuthService interface {
	CheckUserAccess(ctx context.Context, orgID uuid.UUID, userID uuid.UUID, accessLevel models.AccessLevel) (bool, error)
}

type ValidateOrgAndUserAccessResult struct {
	HasAccess bool
	OrgID     uuid.UUID
	UserID    *uuid.UUID
}

func ValidateOrgAndUserAccess(ctx context.Context, service AuthService, accessLevel models.AccessLevel) (ValidateOrgAndUserAccessResult, error) {
	orgID, err := GetOrganizationIDFromContext(ctx)
	if err != nil {
		return ValidateOrgAndUserAccessResult{}, err
	}
	userID, err := GetUserIdFromContext(ctx)
	if err != nil {
		return ValidateOrgAndUserAccessResult{}, err
	}

	ok, err := service.CheckUserAccess(ctx, orgID, userID, accessLevel)
	if err != nil {
		return ValidateOrgAndUserAccessResult{}, err
	}

	return ValidateOrgAndUserAccessResult{
		HasAccess: ok,
		OrgID:     orgID,
		UserID:    &userID,
	}, nil
}
